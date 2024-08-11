package starter

import (
	"database/sql"

	"github.com/Broderick-Westrope/tetrigo/cmd/tui/common"
	tea "github.com/charmbracelet/bubbletea"
)

var _ tea.Model = &Model{}

type Model struct {
	input *Input
	child tea.Model
	db    *sql.DB
}

func NewModel(in *Input) (*Model, error) {
	child, err := in.getChild(in.mode, nil)
	if err != nil {
		return nil, err
	}
	return &Model{
		input: in,
		child: child,
		db:    in.db,
	}, nil
}

func (m *Model) Init() tea.Cmd {
	return m.child.Init()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case common.SwitchModeMsg:
		model, err := m.input.getChild(msg.Target, msg.Input)
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
