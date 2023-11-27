package marathon

import (
	"fmt"
	"time"

	"github.com/Broderick-Westrope/tetrigo/internal/config"
	"github.com/Broderick-Westrope/tetrigo/tetris"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FullscreenInfo struct {
	Width  int
	Height int
}

type Model struct {
	matrix     tetris.Matrix
	styles     *Styles
	help       help.Model
	keys       *keyMap
	currentTet *tetris.Tetrimino
	holdTet    *tetris.Tetrimino
	canHold    bool
	fall       *fall
	scoring    *tetris.Scoring
	bag        *tetris.Bag
	timer      stopwatch.Model
	cfg        *config.Config
	fullscreen *FullscreenInfo
	paused     bool
	startLine  int
}

func InitialModel(level uint, fullscreen *FullscreenInfo) *Model {
	m := &Model{
		matrix:  tetris.Matrix{},
		styles:  defaultStyles(),
		help:    help.New(),
		keys:    defaultKeyMap(),
		scoring: tetris.NewScoring(level),
		holdTet: &tetris.Tetrimino{
			Cells: [][]bool{
				{false, false, false},
				{false, false, false},
			},
			Value: 0,
		},
		canHold: true,
		timer:   stopwatch.NewWithInterval(time.Millisecond),
		paused:  false,
	}
	m.bag = tetris.NewBag(len(m.matrix))
	m.fall = defaultFall(level)
	m.currentTet = m.bag.Next()
	// TODO: Check if the game is over at the starting position
	m.currentTet.Pos.Y++
	err := m.matrix.AddTetrimino(m.currentTet)
	if err != nil {
		panic(fmt.Errorf("failed to add tetrimino to matrix: %w", err))
	}

	cfg, err := config.GetConfig("config.toml")
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}
	m.styles = CreateStyles(&cfg.Theme)
	m.cfg = cfg

	if fullscreen != nil {
		m.fullscreen = fullscreen
		m.styles.ProgramFullscreen.Width(fullscreen.Width).Height(fullscreen.Height)
	}

	m.startLine = len(m.matrix)

	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.fall.stopwatch.Init(), m.timer.Init())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.timer, cmd = m.timer.Update(msg)
	cmds = append(cmds, cmd)

	m.fall.stopwatch, cmd = m.fall.stopwatch.Update(msg)
	cmds = append(cmds, cmd)

	// Operations that can be performed all the time
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Pause):
			m.paused = !m.paused
			cmds = append(cmds, m.timer.Toggle())
			cmds = append(cmds, m.fall.stopwatch.Toggle())
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		}
	case tea.WindowSizeMsg:
		m.styles.ProgramFullscreen.Width(msg.Width).Height(msg.Height)
	}

	if m.paused {
		return m, tea.Batch(cmds...)
	}

	// Operations that can be performed when the game is not paused
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Left):
			err := m.currentTet.MoveLeft(&m.matrix)
			if err != nil {
				panic(fmt.Errorf("failed to move tetrimino left: %w", err))
			}
		case key.Matches(msg, m.keys.Right):
			err := m.currentTet.MoveRight(&m.matrix)
			if err != nil {
				panic(fmt.Errorf("failed to move tetrimino right: %w", err))
			}
		case key.Matches(msg, m.keys.Clockwise):
			err := m.currentTet.Rotate(&m.matrix, true)
			if err != nil {
				panic(fmt.Errorf("failed to rotate tetrimino clockwise: %w", err))
			}
		case key.Matches(msg, m.keys.CounterClockwise):
			err := m.currentTet.Rotate(&m.matrix, false)
			if err != nil {
				panic(fmt.Errorf("failed to rotate tetrimino counter-clockwise: %w", err))
			}
		case key.Matches(msg, m.keys.HardDrop):
			m.hardDrop()
			cmds = append(cmds, m.fall.stopwatch.Reset())
		case key.Matches(msg, m.keys.SoftDrop):
			m.toggleSoftDrop()
		case key.Matches(msg, m.keys.Hold):
			err := m.holdTetrimino()
			if err != nil {
				panic(fmt.Errorf("failed to hold tetrimino: %w", err))
			}
		}
	case stopwatch.TickMsg:
		if m.fall.stopwatch.ID() != msg.ID {
			break
		}
		finished, err := m.lowerTetrimino()
		if err != nil {
			panic(fmt.Errorf("failed to lower tetrimino (tick): %w", err))
		}
		if finished {
			if m.fall.isSoftDrop {
				linesCleared := m.currentTet.Pos.Y - m.startLine
				if linesCleared > 0 {
					m.scoring.AddSoftDrop(uint(linesCleared))
				}
			}

			err := m.nextTetrimino()
			if err != nil {
				panic(fmt.Errorf("failed to get next tetrimino (tick): %w", err))
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var output = lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.JoinVertical(lipgloss.Right, m.holdView(), m.informationView()),
		m.matrixView(),
		m.bagView(),
	)

	output = lipgloss.JoinVertical(lipgloss.Left, output, m.help.View(m.keys))

	if m.fullscreen != nil {
		return m.styles.ProgramFullscreen.Render(output)
	}
	return output
}

func (m *Model) matrixView() string {
	var output string
	for row := (len(m.matrix) - 20); row < len(m.matrix); row++ {
		for col := range m.matrix[row] {
			output += m.renderCell(m.matrix[row][col])
		}
		if row < len(m.matrix)-1 {
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
	header := ""
	if m.paused {
		header = lipgloss.NewStyle().Bold(true).Render("** PAUSED **")
	}

	var output string
	output += fmt.Sprintln("Score: ", m.scoring.Total())
	output += fmt.Sprintln("Level: ", m.scoring.Level())
	output += fmt.Sprintln("Cleared: ", m.scoring.Lines())

	elapsed := m.timer.Elapsed().Seconds()
	minutes := int(elapsed) / 60

	output += "Time: "
	if minutes > 0 {
		seconds := int(elapsed) % 60
		output += fmt.Sprintf("%02d:%02d\n", minutes, seconds)
	} else {
		output += fmt.Sprintf("%06.3f\n", elapsed)
	}

	return m.styles.Information.Render(lipgloss.JoinVertical(lipgloss.Left, header, output))
}

func (m *Model) holdView() string {
	label := m.styles.Hold.Label.Render("Hold:")
	item := m.styles.Hold.Item.Render(m.renderTetrimino(m.holdTet, 1))
	output := lipgloss.JoinVertical(lipgloss.Top, label, item)
	return m.styles.Hold.View.Render(output)
}

func (m *Model) bagView() string {
	output := "Next:\n"
	for i, t := range m.bag.Elements {
		if i >= m.cfg.QueueLength {
			break
		}
		output += "\n" + m.renderTetrimino(&t, 1)
	}
	return m.styles.Bag.Render(output)
}

func (m *Model) renderTetrimino(t *tetris.Tetrimino, background byte) string {
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

func (m *Model) holdTetrimino() error {
	if !m.canHold {
		return nil
	}

	// Swap the current tetrimino with the hold tetrimino
	if m.holdTet.Value == 0 {
		m.holdTet = m.currentTet
		m.currentTet = m.bag.Next()
	} else {
		m.holdTet, m.currentTet = m.currentTet, m.holdTet
	}

	m.matrix.RemoveTetrimino(m.holdTet)

	// Reset the position of the hold tetrimino
	var found bool
	for _, t := range tetris.Tetriminos {
		if t.Value == m.holdTet.Value {
			m.holdTet.Pos = t.Pos
			m.holdTet.Pos.Y += (len(m.matrix) - 20)
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("failed to find tetrimino with value '%v'", m.currentTet.Value)
	}

	// Add the current tetrimino to the matrix
	err := m.matrix.AddTetrimino(m.currentTet)
	if err != nil {
		return fmt.Errorf("failed to add tetrimino to matrix: %w", err)
	}

	m.canHold = false
	return nil
}

func (m *Model) lowerTetrimino() (bool, error) {
	if !m.currentTet.CanMoveDown(m.matrix) {
		action := m.matrix.RemoveCompletedLines(m.currentTet)
		m.scoring.ProcessAction(action, m.cfg.MaxLevel)
		return true, nil
	}

	err := m.currentTet.MoveDown(&m.matrix)
	if err != nil {
		return false, fmt.Errorf("failed to move tetrimino down: %w", err)
	}

	return false, nil
}

func (m *Model) nextTetrimino() error {
	m.currentTet = m.bag.Next()
	// TODO: Check if the game is over at the starting position
	m.currentTet.Pos.Y++
	err := m.matrix.AddTetrimino(m.currentTet)
	if err != nil {
		return fmt.Errorf("failed to add tetrimino to matrix: %w", err)
	}
	m.canHold = true

	if m.fall.isSoftDrop {
		m.startLine = m.currentTet.Pos.Y
	}

	return nil
}

func (m *Model) hardDrop() {
	m.startLine = m.currentTet.Pos.Y
	for {
		finished, err := m.lowerTetrimino()
		if err != nil {
			panic(fmt.Errorf("failed to lower tetrimino (hard drop): %w", err))
		}
		if finished {
			break
		}
	}
	linesCleared := m.currentTet.Pos.Y - m.startLine
	if linesCleared > 0 {
		m.scoring.AddHardDrop(uint(m.currentTet.Pos.Y - m.startLine))
	}
	m.startLine = len(m.matrix)

	err := m.nextTetrimino()
	if err != nil {
		panic(fmt.Errorf("failed to get next tetrimino (hard drop): %w", err))
	}
}

func (m *Model) toggleSoftDrop() {
	m.fall.isSoftDrop = !m.fall.isSoftDrop
	if m.fall.isSoftDrop {
		m.fall.stopwatch.Interval = m.fall.softDropTime
		m.startLine = m.currentTet.Pos.Y
		return
	}
	m.fall.stopwatch.Interval = m.fall.defaultTime
	linesCleared := m.currentTet.Pos.Y - m.startLine
	if linesCleared > 0 {
		m.scoring.AddSoftDrop(uint(linesCleared))
	}
	m.startLine = len(m.matrix)
}
