package starter

import (
	"github.com/Broderick-Westrope/tetrigo/internal"
	tea "github.com/charmbracelet/bubbletea"
)

var _ tea.Model = &Model{}

type Model struct {
	input *Input
	child tea.Model
}

func InitialModel(in *Input) (*Model, error) {
	return &Model{
		input: in,
		child: in.getChild(in.mode),
	}, nil
}

func (m *Model) Init() tea.Cmd {
	return m.child.Init()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case internal.SwitchModeMsg:
		m.child = m.input.getChild(msg.Target)
		return m, m.child.Init()
	}

	var cmd tea.Cmd
	m.child, cmd = m.child.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	return m.child.View()
}
