package textinput

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type TextInputModel struct {
	Child textinput.Model
}

func NewTextInputModel(placeholder string, width int, charLimit int) TextInputModel {
	c := textinput.New()
	c.Placeholder = placeholder
	c.Width = width
	c.CharLimit = charLimit
	c.Focus()

	return TextInputModel{
		Child: c,
	}
}

func (m TextInputModel) Init() tea.Cmd {
	return nil
}

func (m TextInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Child, cmd = m.Child.Update(msg)
	return m, cmd
}

func (m TextInputModel) View() string {
	return m.Child.View()
}
