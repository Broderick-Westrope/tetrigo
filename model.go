package main

import (
	"fmt"

	"github.com/Broderick-Westrope/tetrigo/tetris"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	playfield  tetris.Playfield
	styles     *Styles
	help       help.Model
	keys       *KeyMap
	currentTet *tetris.Tetrimino
	holdTet    *tetris.Tetrimino
	canHold    bool
	fall       *Fall
	scoring    *tetris.Scoring
	bag        *tetris.Bag
}

func InitialModel() *Model {
	m := &Model{
		playfield: tetris.Playfield{},
		styles:    DefaultStyles(),
		help:      help.New(),
		keys:      DefaultKeyMap(),
		scoring:   tetris.NewScoring(1),
		holdTet: &tetris.Tetrimino{
			Cells: [][]bool{
				{false, false, false},
				{false, false, false},
				{false, false, false},
			},
			Value: 0,
		},
		canHold: true,
	}
	m.bag = tetris.NewBag(len(m.playfield))
	m.fall = defaultFall(m.scoring.Level())
	m.currentTet = m.bag.Next()
	err := m.playfield.AddTetrimino(m.currentTet)
	if err != nil {
		panic(fmt.Errorf("failed to add tetrimino to playfield: %w", err))
	}
	return m
}

func (m Model) Init() tea.Cmd {
	return m.fall.stopwatch.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Left):
			err := m.currentTet.MoveLeft(&m.playfield)
			if err != nil {
				panic(fmt.Errorf("failed to move tetrimino left: %w", err))
			}
		case key.Matches(msg, m.keys.Right):
			err := m.currentTet.MoveRight(&m.playfield)
			if err != nil {
				panic(fmt.Errorf("failed to move tetrimino right: %w", err))
			}
		case key.Matches(msg, m.keys.Clockwise):
			err := m.currentTet.Rotate(&m.playfield, true)
			if err != nil {
				panic(fmt.Errorf("failed to rotate tetrimino clockwise: %w", err))
			}
		case key.Matches(msg, m.keys.CounterClockwise):
			err := m.currentTet.Rotate(&m.playfield, false)
			if err != nil {
				panic(fmt.Errorf("failed to rotate tetrimino counter-clockwise: %w", err))
			}
		case key.Matches(msg, m.keys.HardDrop):
			for {
				finished, err := m.lowerTetrimino()
				if err != nil {
					panic(fmt.Errorf("failed to lower tetrimino (hard drop): %w", err))
				}
				if finished {
					break
				}
			}
		case key.Matches(msg, m.keys.SoftDrop):
			m.fall.toggleSoftDrop()
		case key.Matches(msg, m.keys.Hold):
			err := m.holdTetrimino()
			if err != nil {
				panic(fmt.Errorf("failed to hold tetrimino: %w", err))
			}
		}
	case stopwatch.TickMsg:
		_, err := m.lowerTetrimino()
		if err != nil {
			panic(fmt.Errorf("failed to lower tetrimino (tick): %w", err))
		}
	}

	var cmd tea.Cmd
	m.fall.stopwatch, cmd = m.fall.stopwatch.Update(msg)

	return m, cmd
}

func (m Model) View() string {
	var output = lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.JoinVertical(lipgloss.Right, m.holdView(), m.informationView()),
		m.playfieldView(),
		m.bagView(),
	)

	return output + "\n" + m.help.View(m.keys)
}

func (m *Model) playfieldView() string {
	var output string
	for row := (len(m.playfield) - 20); row < len(m.playfield); row++ {
		for col := range m.playfield[row] {
			output += m.renderCell(m.playfield[row][col])
		}
		if row < len(m.playfield)-1 {
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
	var output string
	output += fmt.Sprintln("Score: ", m.scoring.Total())
	output += fmt.Sprintln("Level: ", m.scoring.Level())
	return m.styles.Information.Render(output)
}

func (m *Model) holdView() string {
	output := "Hold:\n" + m.renderTetrimino(m.holdTet)
	return m.styles.Hold.Render(output)
}

func (m *Model) bagView() string {
	output := "Next:\n"
	for i, t := range m.bag.Elements {
		if i > 5 {
			break
		}
		output += "\n" + m.renderTetrimino(&t)
	}
	return m.styles.Bag.Render(output)
}

func (m *Model) renderTetrimino(t *tetris.Tetrimino) string {
	var output string
	for row := range t.Cells {
		for col := range t.Cells[row] {
			if t.Cells[row][col] {
				output += m.renderCell(t.Value)
			} else {
				output += m.renderCell(0)
			}
		}
		output += "\n"
	}
	return output
}

func (m *Model) renderCell(cell byte) string {
	switch cell {
	case 0:
		return m.styles.ColIndicator.Render("▕ ")
	case 'G':
		return "░░"
	default:
		cellStyle, ok := m.styles.TetriminoStyles[cell]
		if ok {
			return cellStyle.Render("██")
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

	m.playfield.RemoveTetrimino(m.holdTet)

	// Reset the position of the hold tetrimino
	var found bool
	for _, t := range tetris.Tetriminos {
		if t.Value == m.holdTet.Value {
			m.holdTet.Pos = t.Pos
			m.holdTet.Pos.Y += (len(m.playfield) - 20)
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("failed to find tetrimino with value '%v'", m.currentTet.Value)
	}

	// Add the current tetrimino to the playfield
	err := m.playfield.AddTetrimino(m.currentTet)
	if err != nil {
		return fmt.Errorf("failed to add tetrimino to playfield: %w", err)
	}

	m.canHold = false
	return nil
}

func (m *Model) lowerTetrimino() (bool, error) {
	if !m.currentTet.CanMoveDown(m.playfield) {
		action := m.playfield.RemoveCompletedLines(m.currentTet)
		m.scoring.ProcessAction(action)
		m.currentTet = m.bag.Next()
		err := m.playfield.AddTetrimino(m.currentTet)
		if err != nil {
			return false, fmt.Errorf("failed to add tetrimino to playfield: %w", err)
		}
		m.canHold = true
		return true, nil
	}

	err := m.currentTet.MoveDown(&m.playfield)
	if err != nil {
		return false, fmt.Errorf("failed to move tetrimino down: %w", err)
	}

	return false, nil
}
