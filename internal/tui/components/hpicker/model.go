package hpicker

import (
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var _ tea.Model = &Model{}

// Model is the model for the horizontal picker component.
type Model struct {
	// cursor is the index of the currently selected option.
	selected int
	// options is a list of the possible options for this component.
	options []KeyValuePair
	// keymap encodes the keybindings recognized by the component.
	keymap *KeyMap
	styles Styles
}

type KeyValuePair struct {
	Key   string
	Value any
}

type Option func(*Model)

func NewModel(options []KeyValuePair, opts ...Option) *Model {
	m := &Model{
		options: options,
		keymap:  defaultKeyMap(),
		styles:  DefaultStyles(),
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func WithRange(minValue, maxValue int) Option {
	return func(m *Model) {
		m.options = make([]KeyValuePair, (maxValue-minValue)+1)
		for i := minValue - 1; i < maxValue; i++ {
			m.options[i].Key = strconv.Itoa(i + 1)
			m.options[i].Value = i + 1
		}
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

// Update is the Tea update function which binds keystrokes to pagination.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(msg, m.keymap.Next):
			m.Next()
		case key.Matches(msg, m.keymap.Prev):
			m.Prev()
		}
	}

	return m, nil
}

// View renders the cursor to a string.
func (m *Model) View() string {
	var prev lipgloss.Style
	if m.isFirst() {
		prev = m.styles.PrevDisabledStyle
	} else {
		prev = m.styles.PrevStyle
	}

	var next lipgloss.Style
	if m.isLast() {
		next = m.styles.NextDisabledStyle
	} else {
		next = m.styles.NextStyle
	}

	return lipgloss.JoinHorizontal(lipgloss.Center,
		prev.Render(m.styles.PrevIndicator),
		m.styles.SelectionStyle.Render(m.options[m.selected].Key),
		next.Render(m.styles.NextIndicator),
	)
}

// Prev is a helper function for navigating one option backward.
// It will not go page beyond the first option (i.e. option 0).
func (m *Model) Prev() {
	if !m.isFirst() {
		m.selected--
	}
}

// Next is a helper function for navigating one option forward.
// It will not go beyond the last option (i.e. len(options) - 1).
func (m *Model) Next() {
	if !m.isLast() {
		m.selected++
	}
}

func (m *Model) isFirst() bool {
	return m.selected == 0
}

func (m *Model) isLast() bool {
	return m.selected == len(m.options)-1
}

func (m *Model) GetSelection() KeyValuePair {
	return m.options[m.selected]
}
