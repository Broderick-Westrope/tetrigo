package textinput

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Child textinput.Model
}

func NewModel(placeholder string, width int, charLimit int) Model {
	c := textinput.New()
	c.Placeholder = placeholder
	c.Width = width
	c.CharLimit = charLimit
	c.Focus()

	return Model{
		Child: c,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Child, cmd = m.Child.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.Child.View()
}
