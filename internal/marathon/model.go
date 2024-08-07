package marathon

import (
	"fmt"
	"time"

	"github.com/Broderick-Westrope/tetrigo/internal/config"
	"github.com/Broderick-Westrope/tetrigo/pkg/tetris"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Input struct {
	isFullscreen bool
	level        uint
}

func NewInput(isFullscreen bool, level uint) *Input {
	return &Input{
		isFullscreen: isFullscreen,
		level:        level,
	}
}

var _ tea.Model = &Model{}

type Model struct {
	matrix            tetris.Matrix
	styles            *Styles
	help              help.Model
	keys              *keyMap
	currentTet        *tetris.Tetrimino
	holdTet           *tetris.Tetrimino
	canHold           bool
	fall              *fall
	scoring           *tetris.Scoring
	bag               *tetris.Bag
	timer             stopwatch.Model
	cfg               *config.Config
	isFullscreen      bool
	paused            bool
	startLine         int
	gameOver          bool
	gameOverStopwatch stopwatch.Model
}

func NewModel(in *Input) *Model {
	m := &Model{
		matrix:  tetris.Matrix{},
		styles:  defaultStyles(),
		help:    help.New(),
		keys:    defaultKeyMap(),
		scoring: tetris.NewScoring(in.level),
		holdTet: &tetris.Tetrimino{
			Cells: [][]bool{
				{false, false, false},
				{false, false, false},
			},
			Value: 0,
		},
		canHold:      true,
		timer:        stopwatch.NewWithInterval(time.Millisecond),
		paused:       false,
		gameOver:     false,
		isFullscreen: in.isFullscreen,
	}
	m.bag = tetris.NewBag(len(m.matrix))
	m.fall = defaultFall(in.level)
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

	if in.isFullscreen {
		m.styles.ProgramFullscreen.Width(0).Height(0)
	}

	m.startLine = len(m.matrix)

	return m
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(m.fall.stopwatch.Init(), m.timer.Init())
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
	var header string
	headerStyle := lipgloss.NewStyle().Width(m.styles.Information.GetWidth()).AlignHorizontal(lipgloss.Center).Bold(true).Underline(true)
	if m.gameOver {
		header = headerStyle.Render("GAME OVER")
	} else if m.paused {
		header = headerStyle.Render("PAUSED")
	} else {
		header = headerStyle.Render("RUNNING")
	}

	var output string
	output += fmt.Sprintln("Score: ", m.scoring.Total())
	output += fmt.Sprintln("level: ", m.scoring.Level())
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
