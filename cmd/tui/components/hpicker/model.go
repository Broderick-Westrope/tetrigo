package hpicker

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// KeyMap is the key bindings for different actions within the component.
type KeyMap struct {
	Prev key.Binding
	Next key.Binding
}

// DefaultKeyMap is the default set of key bindings for navigating and acting  upon the component.
var DefaultKeyMap = KeyMap{
	Prev: key.NewBinding(key.WithKeys("left", "h")),
	Next: key.NewBinding(key.WithKeys("right", "l")),
}

var _ tea.Model = &Model{}

// Model is the model for the horizontal picker component.
type Model struct {
	// Selection is the current selection number.
	Selection      int
	SelectionStyle lipgloss.Style
	// Options is a list of the possible options for this component.
	Options []string
	// KeyMap encodes the keybindings recognized by the component.
	KeyMap KeyMap

	NextIndicator     string
	NextStyle         lipgloss.Style
	NextDisabledStyle lipgloss.Style

	PrevIndicator     string
	PrevStyle         lipgloss.Style
	PrevDisabledStyle lipgloss.Style
}

func (m *Model) Init() tea.Cmd {
	return nil
}

// Update is the Tea update function which binds keystrokes to pagination.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.Next):
			m.Next()
		case key.Matches(msg, m.KeyMap.Prev):
			m.Prev()
		}
	}

	return m, nil
}

// View renders the selection to a string.
func (m *Model) View() string {
	var prev lipgloss.Style
	if m.isFirstPage() {
		prev = m.PrevDisabledStyle
	} else {
		prev = m.PrevStyle
	}

	var next lipgloss.Style
	if m.isLastPage() {
		next = m.NextDisabledStyle
	} else {
		next = m.NextStyle
	}

	return lipgloss.JoinHorizontal(lipgloss.Center,
		prev.Render(m.PrevIndicator),
		m.SelectionStyle.Render(m.Options[m.Selection]),
		next.Render(m.NextIndicator),
	)
}

// Prev is a helper function for navigating one option backward. It will not go page beyond the first option (i.e. option 0).
func (m *Model) Prev() {
	if !m.isFirstPage() {
		m.Selection--
	}
}

// Next is a helper function for navigating one option forward. It will not go beyond the last option (i.e. len(options) - 1).
func (m *Model) Next() {
	if !m.isLastPage() {
		m.Selection++
	}
}

func (m *Model) isFirstPage() bool {
	return m.Selection == 0
}

func (m *Model) isLastPage() bool {
	return m.Selection == len(m.Options)-1
}
