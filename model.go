package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
)

type Playfield [40][10]byte

type Model struct {
	playfield  Playfield
	styles     *Styles
	help       help.Model
	keys       *KeyMap
	stopwatch  stopwatch.Model
	currentTet *Tetrimino
}

func InitialModel() *Model {
	m := &Model{
		playfield: Playfield{},
		styles:    DefaultStyles(),
		help:      help.New(),
		keys:      DefaultKeyMap(),
		stopwatch: stopwatch.NewWithInterval(time.Millisecond * 300),
	}
	m.currentTet = m.playfield.NewTetrimino()
	return m
}

func (m Model) Init() tea.Cmd {
	return m.stopwatch.Init()
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
		}
	case stopwatch.TickMsg:
		newTet, err := m.currentTet.MoveDown(&m.playfield)
		if err != nil {
			panic(fmt.Errorf("failed to move tetrimino down: %w", err))
		}
		if newTet != nil {
			m.currentTet = newTet
		}
	default:
		s := fmt.Sprintf("message type: %T\n", msg)
		fmt.Print(s)
	}

	var stopwatchCmd tea.Cmd
	m.stopwatch, stopwatchCmd = m.stopwatch.Update(msg)

	return m, stopwatchCmd
}

func (m Model) View() string {
	var output string

	for row := (len(m.playfield) - 20); row < len(m.playfield); row++ {
		for col := range m.playfield[row] {
			switch m.playfield[row][col] {
			case 0:
				output += m.styles.ColIndicator.Render("▕ ")
			case 'G':
				output += "░░"
			default:
				cellStyle, ok := m.styles.TetriminoStyles[m.playfield[row][col]]
				if ok {
					output += cellStyle.Render("██")
				} else {
					output += "? "
				}
			}
		}
		if row < len(m.playfield)-1 {
			output += "\n"
		}
	}

	return m.styles.Program.Render(output) + "\n" + m.help.View(m.keys)
}
