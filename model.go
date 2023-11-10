package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	playfield  Playfield
	styles     *Styles
	help       help.Model
	keys       *KeyMap
	currentTet *Tetrimino
	fall       Fall
}

type Fall struct {
	stopwatch    stopwatch.Model
	defaultTime  time.Duration
	softDropTime time.Duration
	currentMode  uint8
}

func (f *Fall) toggleSoftDrop() {
	if f.currentMode == 0 {
		f.currentMode = 1
		f.stopwatch.Interval = f.softDropTime
		return
	}
	f.currentMode = 0
	f.stopwatch.Interval = f.defaultTime
}

func InitialModel() *Model {
	m := &Model{
		playfield: Playfield{},
		styles:    DefaultStyles(),
		help:      help.New(),
		keys:      DefaultKeyMap(),
		fall: Fall{
			defaultTime:  time.Millisecond * 300,
			softDropTime: time.Millisecond * 100,
		},
	}
	m.fall.stopwatch = stopwatch.NewWithInterval(m.fall.defaultTime)
	m.currentTet = m.playfield.NewTetrimino()
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
			var err error
			for {
				if !m.currentTet.canMoveDown(m.playfield) {
					m.playfield.removeCompletedLines(m.currentTet)
					m.currentTet = m.playfield.NewTetrimino()
					break
				}

				err = m.currentTet.MoveDown(&m.playfield)
				if err != nil {
					panic(fmt.Errorf("failed to move tetrimino down with hard drop: %w", err))
				}
			}
		case key.Matches(msg, m.keys.SoftDrop):
			m.fall.toggleSoftDrop()
		}
	case stopwatch.TickMsg:
		if !m.currentTet.canMoveDown(m.playfield) {
			m.playfield.removeCompletedLines(m.currentTet)
			m.currentTet = m.playfield.NewTetrimino()
			break
		}
		err := m.currentTet.MoveDown(&m.playfield)
		if err != nil {
			panic(fmt.Errorf("failed to move tetrimino down: %w", err))
		}
	}

	var cmd tea.Cmd
	m.fall.stopwatch, cmd = m.fall.stopwatch.Update(msg)

	return m, cmd
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
