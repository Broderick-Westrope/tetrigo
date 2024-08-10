package marathon

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Broderick-Westrope/tetrigo/cmd/tui/common"
	"github.com/Broderick-Westrope/tetrigo/internal/config"
	"github.com/Broderick-Westrope/tetrigo/pkg/tetris"
	"github.com/Broderick-Westrope/tetrigo/pkg/tetris/modes/marathon"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Input struct {
	isFullscreen bool
	level        uint
	maxLevel     uint
}

func NewInput(isFullscreen bool, level, maxLevel uint) *Input {
	return &Input{
		isFullscreen: isFullscreen,
		level:        level,
		maxLevel:     maxLevel,
	}
}

var _ tea.Model = &Model{}

type Model struct {
	styles            *Styles
	help              help.Model
	keys              *keyMap
	timer             stopwatch.Model
	cfg               *config.Config
	isFullscreen      bool
	paused            bool
	fallStopwatch     stopwatch.Model
	gameOverStopwatch stopwatch.Model
	game              *marathon.Game
}

func NewModel(in *Input) (*Model, error) {
	game, err := marathon.NewGame(in.level, in.maxLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to create marathon game: %w", err)
	}

	m := &Model{
		styles:       defaultStyles(),
		help:         help.New(),
		keys:         defaultKeyMap(),
		timer:        stopwatch.NewWithInterval(time.Millisecond),
		paused:       false,
		isFullscreen: in.isFullscreen,
		game:         game,
	}
	m.fallStopwatch = stopwatch.NewWithInterval(m.game.GetDefaultFallInterval())

	cfg, err := config.GetConfig("config.toml")
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	m.styles = CreateStyles(&cfg.Theme)
	m.cfg = cfg

	if in.isFullscreen {
		m.styles.ProgramFullscreen.Width(0).Height(0)
	}

	return m, nil
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(m.fallStopwatch.Init(), m.timer.Init())
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.timer, cmd = m.timer.Update(msg)
	cmds = append(cmds, cmd)

	m.fallStopwatch, cmd = m.fallStopwatch.Update(msg)
	cmds = append(cmds, cmd)

	// Operations that can be performed all the time
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Pause):
			if m.game.IsGameOver() {
				break
			}
			m.paused = !m.paused
			cmds = append(cmds, m.timer.Toggle())
			cmds = append(cmds, m.fallStopwatch.Toggle())
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		}
	case tea.WindowSizeMsg:
		m.styles.ProgramFullscreen.Width(msg.Width).Height(msg.Height)
	case stopwatch.TickMsg:
		if m.gameOverStopwatch.ID() != msg.ID {
			break
		}
		// TODO: Redirect to game over / leaderboard screen
		return m, tea.Quit
	case common.SwitchModeMsg:
		if msg.Target == common.MODE_MENU {
			return m, tea.Quit
		}
	}

	if m.paused || m.game.IsGameOver() {
		return m, tea.Batch(cmds...)
	}

	// Operations that can be performed when the game is not paused
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Left):
			err := m.game.MoveLeft()
			if err != nil {
				panic(fmt.Errorf("failed to move tetrimino left: %w", err))
			}
		case key.Matches(msg, m.keys.Right):
			err := m.game.MoveRight()
			if err != nil {
				panic(fmt.Errorf("failed to move tetrimino right: %w", err))
			}
		case key.Matches(msg, m.keys.Clockwise):
			err := m.game.Rotate(true)
			if err != nil {
				panic(fmt.Errorf("failed to rotate tetrimino clockwise: %w", err))
			}
		case key.Matches(msg, m.keys.CounterClockwise):
			err := m.game.Rotate(false)
			if err != nil {
				panic(fmt.Errorf("failed to rotate tetrimino counter-clockwise: %w", err))
			}
		case key.Matches(msg, m.keys.HardDrop):
			// TODO: handle hard drop game over
			_, err := m.game.HardDrop()
			if err != nil {
				panic(fmt.Errorf("failed to hard drop: %w", err))
			}
			cmds = append(cmds, m.fallStopwatch.Reset())
		case key.Matches(msg, m.keys.SoftDrop):
			m.fallStopwatch.Interval = m.game.ToggleSoftDrop()
			// TODO: find a fix for "pausing" momentarily before soft drop begins
			//cmds = append(cmds, func() tea.Msg {
			//	return stopwatch.TickMsg{ID: m.fallStopwatch.ID()}
			//})
		case key.Matches(msg, m.keys.Hold):
			err := m.game.Hold()
			if err != nil {
				panic(fmt.Errorf("failed to hold tetrimino: %w", err))
			}
		}
	case stopwatch.TickMsg:
		if m.fallStopwatch.ID() != msg.ID {
			break
		}
		lockedDown, err := m.game.TickLower()
		if err != nil {
			panic(fmt.Errorf("failed to lower tetrimino (tick): %w", err))
		}
		if lockedDown && m.game.IsGameOver() {
			cmds = append(cmds, m.timer.Toggle())
			cmds = append(cmds, m.fallStopwatch.Toggle())
			m.gameOverStopwatch = stopwatch.NewWithInterval(time.Second * 5)
			cmds = append(cmds, m.gameOverStopwatch.Start())
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	var output = lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.JoinVertical(lipgloss.Right, m.holdView(), m.informationView()),
		m.matrixView(),
		m.bagView(),
	)

	output = lipgloss.JoinVertical(lipgloss.Left, output, m.help.View(m.keys))

	if m.isFullscreen {
		return m.styles.ProgramFullscreen.Render(output)
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
	} else if m.paused {
		header = headerStyle.Render("PAUSED")
	} else {
		header = headerStyle.Render("MARATHON")
	}

	toFixedWidth := func(title, value string) string {
		return fmt.Sprintf("%s%*s\n", title, width-(1+len(title)), value)
	}

	elapsed := m.timer.Elapsed().Seconds()
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
