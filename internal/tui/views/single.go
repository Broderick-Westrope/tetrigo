package views

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"
	"time"

	"github.com/Broderick-Westrope/tetrigo/internal/tui/components"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/stopwatch"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Broderick-Westrope/charmutils"

	"github.com/Broderick-Westrope/tetrigo/internal/config"
	"github.com/Broderick-Westrope/tetrigo/internal/data"
	"github.com/Broderick-Westrope/tetrigo/internal/tui"
	"github.com/Broderick-Westrope/tetrigo/pkg/tetris"
	"github.com/Broderick-Westrope/tetrigo/pkg/tetris/modes/single"
)

const (
	pausedMessage = `
    ____                            __  
   / __ \____ ___  __________  ____/ /  
  / /_/ / __ ^/ / / / ___/ _ \/ __  /   
 / ____/ /_/ / /_/ (__  )  __/ /_/ /    
/_/    \__,_/\__,_/____/\___/\__,_/     
Press PAUSE to continue or HOLD to exit.
`
	gameOverMessage = `
   ______                        ____                 
  / ____/___ _____ ___  ___     / __ \_   _____  _____
 / / __/ __ ^/ __ ^__ \/ _ \   / / / / | / / _ \/ ___/
/ /_/ / /_/ / / / / / /  __/  / /_/ /| |/ /  __/ /    
\____/\__,_/_/ /_/ /_/\___/   \____/ |___/\___/_/     

            Press EXIT or HOLD to continue.           
`
	timerUpdateInterval = time.Millisecond * 13
)

var _ tea.Model = &SingleModel{}

type SingleModel struct {
	username        string
	game            *single.Game
	nextQueueLength int
	fallStopwatch   components.Stopwatch
	mode            tui.Mode

	gameTimer     components.Timer
	gameStopwatch components.Stopwatch

	styles   *components.GameStyles
	help     help.Model
	keys     *components.GameKeyMap
	isPaused bool
	rand     *rand.Rand

	width  int
	height int
}

func NewSingleModel(
	in *tui.SingleInput,
	cfg *config.Config,
	opts ...func(*SingleModel),
) (*SingleModel, error) {
	// Setup initial model
	m := &SingleModel{
		username:        in.Username,
		styles:          components.CreateGameStyles(cfg.Theme),
		help:            help.New(),
		keys:            components.ConstructGameKeyMap(cfg.Keys),
		isPaused:        false,
		nextQueueLength: cfg.NextQueueLength,
		mode:            in.Mode,
		//nolint:gosec // This random source is not for any security-related tasks.
		rand: rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64())),
	}

	for _, opt := range opts {
		opt(m)
	}

	// Get game input
	var gameIn *single.Input
	switch in.Mode {
	case tui.ModeMarathon:
		gameIn = &single.Input{
			Level:         in.Level,
			MaxLevel:      cfg.MaxLevel,
			IncreaseLevel: true,
			EndOnMaxLevel: cfg.EndOnMaxLevel,

			GhostEnabled: cfg.GhostEnabled,
		}
		m.gameStopwatch = components.NewStopwatchWithInterval(timerUpdateInterval)

	case tui.ModeSprint:
		gameIn = &single.Input{
			Level:         in.Level,
			MaxLevel:      cfg.MaxLevel,
			IncreaseLevel: true,

			MaxLines:      40,
			EndOnMaxLines: true,

			GhostEnabled: cfg.GhostEnabled,
		}
		m.gameStopwatch = components.NewStopwatchWithInterval(timerUpdateInterval)

	case tui.ModeUltra:
		gameIn = &single.Input{
			Level:        in.Level,
			GhostEnabled: cfg.GhostEnabled,
		}
		m.gameTimer = components.NewTimerWithInterval(time.Minute*2, timerUpdateInterval)

	case tui.ModeMenu, tui.ModeLeaderboard:
		fallthrough
	default:
		return nil, fmt.Errorf("invalid single player game mode: %v", in.Mode)
	}
	gameIn.Rand = m.rand

	// Create game
	var err error
	m.game, err = single.NewGame(gameIn)
	if err != nil {
		return nil, fmt.Errorf("creating single player game: %w", err)
	}

	// Setup game dependents
	m.fallStopwatch = components.NewStopwatchWithInterval(m.game.GetDefaultFallInterval())

	return m, nil
}

func WithRandSource(r *rand.Rand) func(*SingleModel) {
	return func(m *SingleModel) {
		m.rand = r
	}
}

func (m *SingleModel) Init() tea.Cmd {
	var cmd tea.Cmd
	if m.gameTimer != nil {
		cmd = m.gameTimer.Init()
	} else {
		cmd = m.gameStopwatch.Init()
	}

	return tea.Batch(m.fallStopwatch.Init(), cmd)
}

func (m *SingleModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// Dependencies
	m, cmd = m.dependenciesUpdate(msg)
	cmds = append(cmds, cmd)

	// Operations that can be performed all the time
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
			return m, tea.Batch(cmds...)
		case key.Matches(msg, m.keys.ForceQuit):
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, tea.Batch(cmds...)
	}

	// Game Over
	if m.game.IsGameOver() {
		m, cmd = m.gameOverUpdate(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}

	// Paused
	if m.isPaused {
		m, cmd = m.pausedUpdate(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}

	// Playing
	m, cmd = m.playingUpdate(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *SingleModel) dependenciesUpdate(msg tea.Msg) (*SingleModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	var err error

	if m.gameTimer != nil {
		cmd, err = charmutils.UpdateTypedModel(&m.gameTimer, msg)
		if err != nil {
			cmds = append(cmds, tui.FatalErrorCmd(err))
		}
	} else {
		cmd, err = charmutils.UpdateTypedModel(&m.gameStopwatch, msg)
		if err != nil {
			cmds = append(cmds, tui.FatalErrorCmd(err))
		}
	}
	cmds = append(cmds, cmd)

	cmd, err = charmutils.UpdateTypedModel(&m.fallStopwatch, msg)
	if err != nil {
		cmds = append(cmds, tui.FatalErrorCmd(err))
	}
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *SingleModel) gameOverUpdate(msg tea.Msg) (*SingleModel, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		if key.Matches(msg, m.keys.Exit, m.keys.Hold) {
			modeStr := m.mode.String()
			newEntry := &data.Score{
				GameMode: modeStr,
				Name:     m.username,
				Score:    m.game.GetTotalScore(),
				Lines:    m.game.GetLinesCleared(),
				Level:    m.game.GetLevel(),
			}

			return m, tui.SwitchModeCmd(tui.ModeLeaderboard,
				tui.NewLeaderboardInput(modeStr, tui.WithNewEntry(newEntry)),
			)
		}
	}

	return m, nil
}

func (m *SingleModel) pausedUpdate(msg tea.Msg) (*SingleModel, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(msg, m.keys.Exit):
			return m, m.togglePause()
		case key.Matches(msg, m.keys.Hold):
			return m, tui.SwitchModeCmd(tui.ModeMenu, tui.NewMenuInput())
		}
	}

	return m, nil
}

func (m *SingleModel) playingUpdate(msg tea.Msg) (*SingleModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.playingKeyMsgUpdate(msg)

	case stopwatch.TickMsg:
		if msg.ID != m.fallStopwatch.ID() {
			break
		}
		return m, m.fallStopwatchTick()

	case timer.TimeoutMsg:
		if msg.ID != m.gameTimer.ID() {
			break
		}
		return m, m.triggerGameOver()
	}

	return m, nil
}

func (m *SingleModel) playingKeyMsgUpdate(msg tea.KeyMsg) (*SingleModel, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Left):
		m.game.MoveLeft()
		return m, nil

	case key.Matches(msg, m.keys.Right):
		m.game.MoveRight()
		return m, nil

	case key.Matches(msg, m.keys.Clockwise):
		err := m.game.Rotate(true)
		if err != nil {
			return nil, tui.FatalErrorCmd(fmt.Errorf("rotating clockwise: %w", err))
		}
		return m, nil

	case key.Matches(msg, m.keys.CounterClockwise):
		err := m.game.Rotate(false)
		if err != nil {
			return nil, tui.FatalErrorCmd(fmt.Errorf("rotating counter-clockwise: %w", err))
		}
		return m, nil

	case key.Matches(msg, m.keys.HardDrop):
		gameOver, err := m.game.HardDrop()
		if err != nil {
			return nil, tui.FatalErrorCmd(fmt.Errorf("hard dropping: %w", err))
		}
		var cmds []tea.Cmd
		if gameOver {
			cmds = append(cmds, m.triggerGameOver())
		}
		cmds = append(cmds, m.fallStopwatch.Reset())
		return m, tea.Batch(cmds...)

	case key.Matches(msg, m.keys.SoftDrop):
		m.game.ToggleSoftDrop()
		return m, m.fallStopwatchTick()

	case key.Matches(msg, m.keys.Hold):
		gameOver, err := m.game.Hold()
		if err != nil {
			return nil, tui.FatalErrorCmd(fmt.Errorf("holding tetrimino: %w", err))
		}
		var cmds []tea.Cmd
		if gameOver {
			cmds = append(cmds, m.triggerGameOver())
		}
		return m, tea.Batch(cmds...)

	case key.Matches(msg, m.keys.Exit):
		return m, m.togglePause()
	}
	return m, nil
}

func (m *SingleModel) fallStopwatchTick() tea.Cmd {
	gameOver, err := m.game.TickLower()
	if err != nil {
		return tui.FatalErrorCmd(fmt.Errorf("lowering tetrimino (tick): %w", err))
	}
	if gameOver {
		return m.triggerGameOver()
	}
	m.fallStopwatch.SetInterval(m.game.GetFallInterval())
	return nil
}

func (m *SingleModel) View() string {
	matrixView, err := m.matrixView()
	if err != nil {
		return "** FAILED TO BUILD MATRIX VIEW **"
	}

	var output = lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.JoinVertical(lipgloss.Right, m.holdView(), m.informationView()),
		matrixView,
		m.bagView(),
	)

	if m.game.IsGameOver() {
		output, err = charmutils.OverlayCenter(output, gameOverMessage, false)
		if err != nil {
			return "** FAILED TO OVERLAY GAME OVER MESSAGE **"
		}
	} else if m.isPaused {
		output, err = charmutils.OverlayCenter(output,
			lipgloss.NewStyle().Margin(0, 1).
				Render(pausedMessage),
			false)
		if err != nil {
			return "** FAILED TO OVERLAY PAUSED MESSAGE **"
		}
	}

	output = lipgloss.JoinVertical(lipgloss.Left, output, m.help.View(m.keys))
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, output)
}

func (m *SingleModel) matrixView() (string, error) {
	matrix, err := m.game.GetVisibleMatrix()
	if err != nil {
		return "", fmt.Errorf("getting visible matrix: %w", err)
	}

	var output string
	for row := range matrix {
		for col := range matrix[row] {
			output += m.renderCell(matrix[row][col])
		}
		if row < len(matrix)-1 {
			output += "\n"
		}
	}

	var rowIndicator string
	for i := 1; i <= 20; i++ {
		rowIndicator += fmt.Sprintf("%d\n", i)
	}
	return lipgloss.JoinHorizontal(lipgloss.Center,
		m.styles.Playfield.Render(output),
		m.styles.RowIndicator.Render(rowIndicator),
	), nil
}

func (m *SingleModel) informationView() string {
	width := m.styles.Information.GetWidth()

	var header string
	headerStyle := lipgloss.NewStyle().Width(width).AlignHorizontal(lipgloss.Center).Bold(true).Underline(true)

	switch {
	case m.game.IsGameOver():
		header = headerStyle.Render("GAME OVER")
	case m.isPaused:
		header = headerStyle.Render("PAUSED")
	default:
		header = headerStyle.Render(strings.ToUpper(m.mode.String()))
	}

	toFixedWidth := func(title, value string) string {
		return fmt.Sprintf("%s%*s\n", title, width-(1+len(title)), value)
	}

	var gameTime float64
	if m.gameTimer != nil {
		gameTime = m.gameTimer.GetTimeout().Seconds()
	} else {
		gameTime = m.gameStopwatch.Elapsed().Seconds()
	}

	minutes := int(gameTime) / 60

	var timeStr string
	if minutes > 0 {
		seconds := int(gameTime) % 60
		timeStr += fmt.Sprintf("%02d:%02d", minutes, seconds)
	} else {
		timeStr += fmt.Sprintf("%06.3f", gameTime)
	}

	var output string
	output += fmt.Sprintln("Score:")
	output += fmt.Sprintf("%*d\n", width-1, m.game.GetTotalScore())
	output += fmt.Sprintln("Time:")
	output += fmt.Sprintf("%*s\n", width-1, timeStr)
	output += toFixedWidth("Lines:", strconv.Itoa(m.game.GetLinesCleared()))
	output += toFixedWidth("Level:", strconv.Itoa(m.game.GetLevel()))

	return m.styles.Information.Render(lipgloss.JoinVertical(lipgloss.Left, header, output))
}

func (m *SingleModel) holdView() string {
	label := m.styles.Hold.Label.Render("Hold:")
	item := m.styles.Hold.Item.Render(m.renderTetrimino(m.game.GetHoldTetrimino(), 1))
	output := lipgloss.JoinVertical(lipgloss.Top, label, item)
	return m.styles.Hold.View.Render(output)
}

func (m *SingleModel) bagView() string {
	output := "Next:\n"
	for i, t := range m.game.GetBagTetriminos() {
		if i >= m.nextQueueLength {
			break
		}
		output += "\n" + m.renderTetrimino(&t, 1)
	}
	return m.styles.Bag.Render(output)
}

func (m *SingleModel) renderTetrimino(t *tetris.Tetrimino, background byte) string {
	var output string
	for row := range t.Cells {
		for col := range t.Cells[row] {
			if t.Cells[row][col] {
				output += m.renderCell(t.Value)
			} else {
				output += m.renderCell(background)
			}
		}
		output += "\n"
	}
	return output
}

func (m *SingleModel) renderCell(cell byte) string {
	switch cell {
	case 0:
		return m.styles.EmptyCell.Render(m.styles.CellChar.Empty)
	case 1:
		return "  "
	case 'G':
		return m.styles.GhostCell.Render(m.styles.CellChar.Ghost)
	default:
		cellStyle, ok := m.styles.TetriminoCellStyles[cell]
		if ok {
			return cellStyle.Render(m.styles.CellChar.Tetriminos)
		}
	}
	return "??"
}

func (m *SingleModel) triggerGameOver() tea.Cmd {
	m.game.EndGame()
	m.isPaused = false

	var cmds []tea.Cmd
	if m.gameTimer != nil {
		m.gameTimer.SetTimeout(0)
		cmds = append(cmds, m.gameTimer.Stop())
	} else {
		cmds = append(cmds, m.fallStopwatch.Stop())
	}

	return tea.Batch(cmds...)
}

func (m *SingleModel) togglePause() tea.Cmd {
	m.isPaused = !m.isPaused

	var cmd tea.Cmd
	if m.gameTimer != nil {
		cmd = m.gameTimer.Toggle()
	} else {
		cmd = m.gameStopwatch.Toggle()
	}

	return tea.Batch(
		m.fallStopwatch.Toggle(),
		cmd,
	)
}
