package components

import (
	time "time"

	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
)

type Timer interface {
	Init() tea.Cmd
	Update(tea.Msg) (tea.Model, tea.Cmd)
	View() string
	ID() int
	GetTimeout() time.Duration
	SetTimeout(time.Duration)
	Stop() tea.Cmd
	Toggle() tea.Cmd
}

type timerImpl struct {
	model timer.Model
}

func NewTimerWithInterval(timeout, interval time.Duration) Timer {
	return &timerImpl{
		model: timer.NewWithInterval(timeout, interval),
	}
}

func (t *timerImpl) Init() tea.Cmd {
	return t.model.Init()
}

func (t *timerImpl) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m, cmd := t.model.Update(msg)
	t.model = m
	return t, cmd
}

func (t *timerImpl) View() string {
	return t.model.View()
}

func (t *timerImpl) GetTimeout() time.Duration {
	return t.model.Timeout
}

func (t *timerImpl) ID() int {
	return t.model.ID()
}

func (t *timerImpl) SetTimeout(timeout time.Duration) {
	t.model.Timeout = timeout
}

func (t *timerImpl) Stop() tea.Cmd {
	return t.model.Stop()
}

func (t *timerImpl) Toggle() tea.Cmd {
	return t.model.Toggle()
}
