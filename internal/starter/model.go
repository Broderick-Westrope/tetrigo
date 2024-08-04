package starter

import (
	"errors"

	"github.com/Broderick-Westrope/tetrigo/internal"
	"github.com/Broderick-Westrope/tetrigo/internal/marathon"
	"github.com/Broderick-Westrope/tetrigo/internal/menu"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	input *Input
	child tea.Model
}

func InitialModel(in *Input) (*Model, error) {
	var child tea.Model
	switch in.mode {
	case internal.MODE_MENU:
		child = menu.InitialModel(in.ToMenuInput())
	case internal.MODE_MARATHON:
		child = marathon.InitialModel(in.ToMarathonInput())
	default:
		return nil, errors.New("invalid starter mode")
	}

	return &Model{
		input: in,
		child: child,
	}, nil
}

func (m Model) Init() tea.Cmd {
	return m.child.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case internal.SwitchModeMsg:
		switch msg.Target {
		case internal.MODE_MENU:
			m.child = menu.InitialModel(m.input.ToMenuInput())
		case internal.MODE_MARATHON:
			m.child = marathon.InitialModel(m.input.ToMarathonInput())
		}
	}
	return m.child.Update(msg)
}

func (m Model) View() string {
	return m.child.View()
}
