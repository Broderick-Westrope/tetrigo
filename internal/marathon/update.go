package marathon

import (
	"fmt"
	"time"

	"github.com/Broderick-Westrope/tetrigo/internal"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
)

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
	case internal.SwitchModeMsg:
		if msg.Target == internal.MODE_MENU {
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
			fallInterval := m.game.ToggleSoftDrop()
			m.fallStopwatch.Interval = fallInterval
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
		lockedDown, err := m.game.Lower()
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
