package marathon

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Broderick-Westrope/tetrigo/cmd/tui/common"
	"github.com/Broderick-Westrope/tetrigo/cmd/tui/helpers"
	"github.com/Broderick-Westrope/tetrigo/internal/config"
	"github.com/Broderick-Westrope/tetrigo/internal/data"
	"github.com/Broderick-Westrope/tetrigo/pkg/tetris"
	"github.com/Broderick-Westrope/tetrigo/pkg/tetris/modes/marathon"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var pausedMsg = ` ______  ___  _   _ _____ ___________ 
 | ___ \/ _ \| | | /  ___|  ___|  _  \
 | |_/ / /_\ \ | | \ '--.| |__ | | | |
 |  __/|  _  | | | |'--. \  __|| | | |
 | |   | | | | |_| /\__/ / |___| |/ /
 \_|   \_| |_/\___/\____/\____/|___/
Press PAUSE to continue or HOLD to exit.`

var gameOverMsg = ` _____   ___  ___  ___ _____   _____  _   _ ___________ 
|  __ \ / _ \ |  \/  ||  ___| |  _  || | | |  ___| ___ \
| |  \// /_\ \| .  . || |__   | | | || | | | |__ | |_/ /
| | __ |  _  || |\/| ||  __|  | | | || | | |  __||    / 
| |_\ \| | | || |  | || |___  \ \_/ /\ \_/ / |___| |\ \ 
 \____/\_| |_/\_|  |_/\____/   \___/  \___/\____/\_| \_|
			 Press EXIT or HOLD to continue.`

var _ tea.Model = &Model{}

type Model struct {
	styles         *Styles
	help           help.Model
	keys           *keyMap
	timerStopwatch stopwatch.Model
	cfg            *config.Config
	isFullscreen   bool
	isPaused       bool
	fallStopwatch  stopwatch.Model
	game           *marathon.Game
	isGameOver     bool
}

func NewModel(in *common.MarathonInput, keys *common.Keys) (*Model, error) {
	game, err := marathon.NewGame(in.Level, in.MaxLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to create marathon game: %w", err)
	}

	m := &Model{
		styles:         defaultStyles(),
		help:           help.New(),
		keys:           constructKeyMap(keys),
		timerStopwatch: stopwatch.NewWithInterval(time.Millisecond * 3),
		isPaused:       false,
		isFullscreen:   in.IsFullscreen,
		game:           game,
	}
	m.fallStopwatch = stopwatch.NewWithInterval(m.game.GetDefaultFallInterval())

	cfg, err := config.GetConfig("config.toml")
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	m.styles = CreateStyles(&cfg.Theme)
	m.cfg = cfg

	if in.IsFullscreen {
		m.styles.ProgramFullscreen.Width(0).Height(0)
	}

	return m, nil
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(m.fallStopwatch.Init(), m.timerStopwatch.Init())
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		m.styles.ProgramFullscreen.Width(msg.Width).Height(msg.Height)
		return m, tea.Batch(cmds...)
	case common.SwitchModeMsg:
		panic(fmt.Errorf("unexpected/unhandled SwitchModeMsg: %v", msg.Target))
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

func (m *Model) dependenciesUpdate(msg tea.Msg) (*Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.timerStopwatch, cmd = m.timerStopwatch.Update(msg)
	cmds = append(cmds, cmd)

	m.fallStopwatch, cmd = m.fallStopwatch.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) gameOverUpdate(msg tea.Msg) (*Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Exit, m.keys.Hold):
			newEntry := &data.Score{
				GameMode: "marathon",
				Name:     "Player",
				Time:     m.timerStopwatch.Elapsed(),
				Score:    int(m.game.GetTotalScore()),
				Lines:    int(m.game.GetLinesCleared()),
				Level:    int(m.game.GetLevel()),
			}

			return m, common.SwitchModeCmd(common.MODE_LEADERBOARD,
				common.NewLeaderboardInput("marathon", common.WithNewEntry(newEntry)),
			)
		}
	}

	return m, nil
}

func (m *Model) pausedUpdate(msg tea.Msg) (*Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Exit):
			return m, m.togglePause()
		case key.Matches(msg, m.keys.Hold):
			return m, common.SwitchModeCmd(common.MODE_MENU, common.NewMenuInput(m.isFullscreen))
		}
	}

	return m, nil
}

func (m *Model) playingUpdate(msg tea.Msg) (*Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Left):
			err := m.game.MoveLeft()
			if err != nil {
				panic(fmt.Errorf("failed to move left: %w", err))
			}
			return m, nil
		case key.Matches(msg, m.keys.Right):
			err := m.game.MoveRight()
			if err != nil {
				panic(fmt.Errorf("failed to move right: %w", err))
			}
			return m, nil
		case key.Matches(msg, m.keys.Clockwise):
			err := m.game.Rotate(true)
			if err != nil {
				panic(fmt.Errorf("failed to rotate clockwise: %w", err))
			}
			return m, nil
		case key.Matches(msg, m.keys.CounterClockwise):
			err := m.game.Rotate(false)
			if err != nil {
				panic(fmt.Errorf("failed to rotate counter-clockwise: %w", err))
			}
			return m, nil
		case key.Matches(msg, m.keys.HardDrop):
			gameOver, err := m.game.HardDrop()
			if err != nil {
				panic(fmt.Errorf("failed to hard drop: %w", err))
			}
			var cmds []tea.Cmd
			if gameOver {
				cmds = append(cmds, m.triggerGameOver())
			}
			cmds = append(cmds, m.fallStopwatch.Reset())
			return m, tea.Batch(cmds...)
		case key.Matches(msg, m.keys.SoftDrop):
			m.fallStopwatch.Interval = m.game.ToggleSoftDrop()
			// TODO: find a fix for "pausing" momentarily before soft drop begins
			//cmds = append(cmds, func() tea.Msg {
			//	return stopwatch.TickMsg{ID: m.fallStopwatch.ID()}
			//})
		case key.Matches(msg, m.keys.Hold):
			gameOver, err := m.game.Hold()
			if err != nil {
				panic(fmt.Errorf("failed to hold tetrimino: %w", err))
			}
			var cmds []tea.Cmd
			if gameOver {
				cmds = append(cmds, m.triggerGameOver())
			}
			return m, tea.Batch(cmds...)
		case key.Matches(msg, m.keys.Exit):
			return m, m.togglePause()
		}
	case stopwatch.TickMsg:
		if msg.ID != m.fallStopwatch.ID() {
			break
		}
		gameOver, err := m.game.TickLower()
		if err != nil {
			panic(fmt.Errorf("failed to lower tetrimino (tick): %w", err))
		}
		var cmds []tea.Cmd
		if gameOver {
			cmds = append(cmds, m.triggerGameOver())
		}
		return m, tea.Batch(cmds...)
	}

	return m, nil
}

func (m *Model) View() string {
	var output = lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.JoinVertical(lipgloss.Right, m.holdView(), m.informationView()),
		m.matrixView(),
		m.bagView(),
	)

	if m.game.IsGameOver() {
		output = helpers.PlaceOverlayCenter(gameOverMsg, output)
	}

	if m.isPaused {
		output = helpers.PlaceOverlayCenter(pausedMsg, output)
	}

	output = lipgloss.JoinVertical(lipgloss.Left, output, m.help.View(m.keys))

	if m.isFullscreen {
		output = m.styles.ProgramFullscreen.Render(output)
	}

	return output
}

func (m *Model) matrixView() string {
	matrix := m.game.GetVisibleMatrix()
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
	return lipgloss.JoinHorizontal(lipgloss.Center, m.styles.Playfield.Render(output), m.styles.RowIndicator.Render(rowIndicator))
}

func (m *Model) informationView() string {
	width := m.styles.Information.GetWidth()

	var header string
	headerStyle := lipgloss.NewStyle().Width(width).AlignHorizontal(lipgloss.Center).Bold(true).Underline(true)
	if m.game.IsGameOver() {
		header = headerStyle.Render("GAME OVER")
	} else if m.isPaused {
		header = headerStyle.Render("PAUSED")
	} else {
		header = headerStyle.Render("MARATHON")
	}

	toFixedWidth := func(title, value string) string {
		return fmt.Sprintf("%s%*s\n", title, width-(1+len(title)), value)
	}

	elapsed := m.timerStopwatch.Elapsed().Seconds()
	minutes := int(elapsed) / 60

	var timeStr string
	if minutes > 0 {
		seconds := int(elapsed) % 60
		timeStr += fmt.Sprintf("%02d:%02d", minutes, seconds)
	} else {
		timeStr += fmt.Sprintf("%06.3f", elapsed)
	}

	var output string
	output += fmt.Sprintln("Score:")
	output += fmt.Sprintf("%*d\n", width-1, m.game.GetTotalScore())
	output += fmt.Sprintln("Time:")
	output += fmt.Sprintf("%*s\n", width-1, timeStr)
	output += toFixedWidth("Lines:", strconv.Itoa(int(m.game.GetLinesCleared())))
	output += toFixedWidth("Level:", strconv.Itoa(int(m.game.GetLevel())))

	return m.styles.Information.Render(lipgloss.JoinVertical(lipgloss.Left, header, output))
}

func (m *Model) holdView() string {
	label := m.styles.Hold.Label.Render("Hold:")
	item := m.styles.Hold.Item.Render(m.renderTetrimino(m.game.GetHoldTetrimino(), 1))
	output := lipgloss.JoinVertical(lipgloss.Top, label, item)
	return m.styles.Hold.View.Render(output)
}

func (m *Model) bagView() string {
	output := "Next:\n"
	for i, t := range m.game.GetBagTetriminos() {
		if i >= m.cfg.QueueLength {
			break
		}
		output += "\n" + m.renderTetrimino(&t, 1)
	}
	return m.styles.Bag.Render(output)
}

func (m *Model) renderTetrimino(t *tetris.Tetrimino, background byte) string {
	var output string
	for row := range t.Minos {
		for col := range t.Minos[row] {
			if t.Minos[row][col] {
				output += m.renderCell(t.Value)
			} else {
				output += m.renderCell(background)
			}
		}
		output += "\n"
	}
	return output
}

func (m *Model) renderCell(cell byte) string {
	switch cell {
	case 0:
		return m.styles.ColIndicator.Render(m.cfg.Theme.Characters.EmptyCell)
	case 1:
		return "  "
	case 'G':
		return m.styles.GhostCell.Render(m.cfg.Theme.Characters.GhostCell)
	default:
		cellStyle, ok := m.styles.TetriminoStyles[cell]
		if ok {
			return cellStyle.Render(m.cfg.Theme.Characters.Tetriminos)
		}
	}
	return "??"
}

func (m *Model) triggerGameOver() tea.Cmd {
	m.isGameOver = true
	m.isPaused = false
	var cmds []tea.Cmd
	cmds = append(cmds, m.timerStopwatch.Stop())
	cmds = append(cmds, m.fallStopwatch.Stop())
	return tea.Batch(cmds...)
}

func (m *Model) togglePause() tea.Cmd {
	m.isPaused = !m.isPaused
	return tea.Batch(
		m.fallStopwatch.Toggle(),
		m.timerStopwatch.Toggle(),
	)
}
