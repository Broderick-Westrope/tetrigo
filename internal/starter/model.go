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

func NewModel(in *Input) (*Model, error) {
	child, err := in.getChild(in.mode)
	if err != nil {
		return nil, err
	}
	return &Model{
		input: in,
		child: child,
	}, nil
}

func (m *Model) Init() tea.Cmd {
	return m.child.Init()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case internal.SwitchModeMsg:
		model, err := m.input.getChild(msg.Target)
		if err != nil {
			panic(err)
		}
		m.child = model
		return m, m.child.Init()
	}

	var cmd tea.Cmd
	m.child, cmd = m.child.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	return m.child.View()
}
