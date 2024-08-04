package marathon

import (
	"fmt"
	"time"

	"github.com/Broderick-Westrope/tetrigo/internal"
	"github.com/Broderick-Westrope/tetrigo/pkg/tetris"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
)

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
			if m.gameOver {
				break
			}
			m.paused = !m.paused
			cmds = append(cmds, m.timer.Toggle())
			cmds = append(cmds, m.fall.stopwatch.Toggle())
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

	if m.paused || m.gameOver {
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
			newModel, err := m.holdTetrimino()
			m = newModel
			if err != nil {
				panic(fmt.Errorf("failed to hold tetrimino: %w", err))
			}
		}
	case stopwatch.TickMsg:
		if m.fall.stopwatch.ID() != msg.ID {
			break
		}
		newModel, finished, err := m.lowerTetrimino()
		m = newModel
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

			m, gameOver, err := m.nextTetrimino()
			if err != nil {
				panic(fmt.Errorf("failed to get next tetrimino (tick): %w", err))
			}

			if gameOver {
				m.gameOver = true
				cmds = append(cmds, m.timer.Toggle())
				cmds = append(cmds, m.fall.stopwatch.Toggle())
				m.gameOverStopwatch = stopwatch.NewWithInterval(time.Second * 5)
				cmds = append(cmds, m.gameOverStopwatch.Start())
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) holdTetrimino() (Model, error) {
	if !m.canHold {
		return m, nil
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
		return m, fmt.Errorf("failed to find tetrimino with value '%v'", m.currentTet.Value)
	}

	// Add the current tetrimino to the matrix
	err := m.matrix.AddTetrimino(m.currentTet)
	if err != nil {
		return m, fmt.Errorf("failed to add tetrimino to matrix: %w", err)
	}

	m.canHold = false
	return m, nil
}

func (m Model) lowerTetrimino() (Model, bool, error) {
	if !m.currentTet.CanMoveDown(m.matrix) {
		action := m.matrix.RemoveCompletedLines(m.currentTet)
		m.scoring.ProcessAction(action, m.cfg.MaxLevel)
		m.fall.calculateFallSpeeds(m.scoring.Level())
		return m, true, nil
	}

	err := m.currentTet.MoveDown(&m.matrix)
	if err != nil {
		return m, false, fmt.Errorf("failed to move tetrimino down: %w", err)
	}

	return m, false, nil
}

func (m Model) nextTetrimino() (Model, bool, error) {
	m.currentTet = m.bag.Next()

	// Block Out
	if m.currentTet.IsOverlapping(&m.matrix) {
		return m, true, nil
	}

	if m.currentTet.CanMoveDown(m.matrix) {
		m.currentTet.Pos.Y++
	} else {
		// Lock Out
		if m.currentTet.IsAbovePlayfield(len(m.matrix)) {
			return m, true, nil
		}
	}

	if err := m.matrix.AddTetrimino(m.currentTet); err != nil {
		return m, false, fmt.Errorf("failed to add tetrimino to matrix: %w", err)
	}
	m.canHold = true

	if m.fall.isSoftDrop {
		m.startLine = m.currentTet.Pos.Y
	}

	return m, false, nil
}

func (m Model) hardDrop() (Model, bool) {
	m.startLine = m.currentTet.Pos.Y
	for {
		newModel, finished, err := m.lowerTetrimino()
		m = newModel
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

	m, gameOver, err := m.nextTetrimino()
	if err != nil {
		panic(fmt.Errorf("failed to get next tetrimino (hard drop): %w", err))
	}
	return m, gameOver
}

func (m Model) toggleSoftDrop() {
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
