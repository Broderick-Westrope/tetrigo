package marathon

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Broderick-Westrope/tetrigo/internal/config"
	"github.com/Broderick-Westrope/tetrigo/internal/data"
	"github.com/Broderick-Westrope/tetrigo/internal/tui/common"
	"github.com/Broderick-Westrope/tetrigo/internal/tui/game"
	"github.com/Broderick-Westrope/tetrigo/pkg/tetris"
	"github.com/Broderick-Westrope/tetrigo/pkg/tetris/modes/marathon"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var _ tea.Model = &Model{}

type Model struct {
	playerName      string
	styles          *game.Styles
	help            help.Model
	keys            *game.GameKeyMap
	timerStopwatch  stopwatch.Model
	isPaused        bool
	fallStopwatch   stopwatch.Model
	game            *marathon.Game
	isGameOver      bool
	nextQueueLength int
}

func NewModel(in *common.MarathonInput, cfg *config.Config) (*Model, error) {
	gameIn := &marathon.Input{
		Level:         in.Level,
		MaxLevel:      cfg.MaxLevel,
		EndOnMaxLevel: cfg.EndOnMaxLevel,
		GhostEnabled:  cfg.GhostEnabled,
	}

	g, err := marathon.NewGame(gameIn)
	if err != nil {
		return nil, fmt.Errorf("failed to create marathon game: %w", err)
	}

	m := &Model{
		playerName:      in.PlayerName,
		styles:          game.CreateStyles(cfg.Theme),
		help:            help.New(),
		keys:            game.ConstructGameKeyMap(cfg.Keys),
		timerStopwatch:  stopwatch.NewWithInterval(time.Millisecond * 3),
		isPaused:        false,
		game:            g,
		nextQueueLength: cfg.NextQueueLength,
	}
	m.fallStopwatch = stopwatch.NewWithInterval(m.game.GetDefaultFallInterval())

	m.styles = game.CreateStyles(cfg.Theme)

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
	if msg, ok := msg.(tea.KeyMsg); ok {
		if key.Matches(msg, m.keys.Exit, m.keys.Hold) {
			newEntry := &data.Score{
				GameMode: "marathon",
				Name:     m.playerName,
				Time:     m.timerStopwatch.Elapsed(),
				Score:    int(m.game.GetTotalScore()),
				Lines:    int(m.game.GetLinesCleared()),
				Level:    int(m.game.GetLevel()),
			}

			return m, common.SwitchModeCmd(common.ModeLeaderboard,
				common.NewLeaderboardInput("marathon", common.WithNewEntry(newEntry)),
			)
		}
	}

	return m, nil
}

func (m *Model) pausedUpdate(msg tea.Msg) (*Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(msg, m.keys.Exit):
			return m, m.togglePause()
		case key.Matches(msg, m.keys.Hold):
			return m, common.SwitchModeCmd(common.ModeMenu, common.NewMenuInput())
		}
	}

	return m, nil
}

func (m *Model) playingUpdate(msg tea.Msg) (*Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.playingKeyMsgUpdate(msg)
	case stopwatch.TickMsg:
		if msg.ID != m.fallStopwatch.ID() {
			break
		}
		m.fallStopwatch.Interval = m.game.GetDefaultFallInterval()
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

func (m *Model) playingKeyMsgUpdate(msg tea.KeyMsg) (*Model, tea.Cmd) {
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
		// cmds = append(cmds, func() tea.Msg {
		// 	return stopwatch.TickMsg{ID: m.fallStopwatch.ID()}
		// })
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
	return m, nil
}

func (m *Model) View() string {
	var output = lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.JoinVertical(lipgloss.Right, m.holdView(), m.informationView()),
		m.matrixView(),
		m.bagView(),
	)

	if m.game.IsGameOver() {
		output = common.OverlayGameOverMessage(output)
	}

	if m.isPaused {
		output = common.OverlayPausedMessage(output)
	}

	output = lipgloss.JoinVertical(lipgloss.Left, output, m.help.View(m.keys))

	return output
}

func (m *Model) matrixView() string {
	matrix, err := m.game.GetVisibleMatrix()
	if err != nil {
		panic(err)
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
	)
}

func (m *Model) informationView() string {
	width := m.styles.Information.GetWidth()

	var header string
	headerStyle := lipgloss.NewStyle().Width(width).AlignHorizontal(lipgloss.Center).Bold(true).Underline(true)

	switch {
	case m.game.IsGameOver():
		header = headerStyle.Render("GAME OVER")
	case m.isPaused:
		header = headerStyle.Render("PAUSED")
	default:
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
		if i >= m.nextQueueLength {
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
