package components

import (
	"time"

	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
)

type Stopwatch interface {
	Init() tea.Cmd
	Update(tea.Msg) (tea.Model, tea.Cmd)
	View() string
	Elapsed() time.Duration
}

type stopwatchImpl struct {
	model stopwatch.Model
}

func NewStopwatchWithInterval(interval time.Duration) Stopwatch {
	return &stopwatchImpl{
		stopwatch.NewWithInterval(interval),
	}
}

// Init implements Stopwatch.
func (s *stopwatchImpl) Init() tea.Cmd {
	return s.model.Init()
}

// Update implements Stopwatch.
func (s *stopwatchImpl) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m, cmd := s.model.Update(msg)
	s.model = m
	return s, cmd
}

// View implements Stopwatch.
func (s *stopwatchImpl) View() string {
	return s.model.View()
}

// Elapsed implements Stopwatch.
func (s *stopwatchImpl) Elapsed() time.Duration {
	return s.model.Elapsed()
}
