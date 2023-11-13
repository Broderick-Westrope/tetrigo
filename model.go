package main

import (
	"fmt"
	"math"
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
	fall       *Fall
	scoring    *scoring
}

type Fall struct {
	stopwatch    stopwatch.Model
	defaultTime  time.Duration
	softDropTime time.Duration
	isSoftDrop   bool
}

func (f *Fall) calculateFallSpeeds(level uint) {
	speed := math.Pow((0.8-float64(level-1)*0.007), float64(level-1)) * 1000000

	f.defaultTime = time.Microsecond * time.Duration(speed)
	f.softDropTime = time.Microsecond * time.Duration(speed/10)
}

func (f *Fall) toggleSoftDrop() {
	f.isSoftDrop = !f.isSoftDrop
	if f.isSoftDrop {
		f.stopwatch.Interval = f.softDropTime
		return
	}
	f.stopwatch.Interval = f.defaultTime
}

func defaultFall(level uint) *Fall {
	f := Fall{}
	f.calculateFallSpeeds(level)
	f.stopwatch = stopwatch.NewWithInterval(f.defaultTime)
	return &f
}

func InitialModel() *Model {
	m := &Model{
		playfield: Playfield{},
		styles:    DefaultStyles(),
		help:      help.New(),
		keys:      DefaultKeyMap(),
		scoring:   &scoring{total: 0, backToBack: false, level: 1},
	}
	m.fall = defaultFall(m.scoring.level)
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
					action := m.playfield.removeCompletedLines(m.currentTet)
					m.scoring.processAction(action)
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
			action := m.playfield.removeCompletedLines(m.currentTet)
			m.scoring.processAction(action)
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
