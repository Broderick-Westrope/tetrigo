package menu

import (
	"fmt"

	"github.com/Broderick-Westrope/tetrigo/internal/tui/common"
	"github.com/Broderick-Westrope/tetrigo/internal/tui/common/components/hpicker"
	"github.com/Broderick-Westrope/tetrigo/internal/tui/common/components/textinput"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const titleStr = `
    ______________________  ______________ 
   /_  __/ ____/_  __/ __ \/  _/ ____/ __ \
    / / / __/   / / / /_/ // // / __/ / / /
   / / / /___  / / / _, _// // /_/ / /_/ / 
  /_/ /_____/ /_/ /_/ |_/___/\____/\____/  
`

var _ tea.Model = Model{}

type Model struct {
	items    []item
	selected int

	keys   *keyMap
	styles *styles
	help   help.Model
}

type item struct {
	label     string
	model     tea.Model
	hideLabel bool
}

func NewModel(_ *common.MenuInput) *Model {
	nameInput := textinput.NewModel("Enter your name", 20, 20)
	modePicker := hpicker.NewModel([]hpicker.KeyValuePair{
		{Key: "Marathon", Value: common.ModeMarathon},
		// {Key: "Sprint (40 Lines)", Value: "sprint"},
		{Key: "Ultra (Time Trial)", Value: common.ModeUltra},
	})
	levelPicker := hpicker.NewModel(nil, hpicker.WithRange(1, 15))

	return &Model{
		items: []item{
			{label: "Name", model: nameInput, hideLabel: true},
			{label: "Mode", model: modePicker},
			{label: "Starting Level", model: levelPicker},
		},
		selected: 0,

		keys:   defaultKeyMap(),
		styles: defaultStyles(),
		help:   help.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch {
		case key.Matches(msg, m.keys.Exit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Up):
			if m.selected > 0 {
				m.selected--
			}
			return m, nil
		case key.Matches(msg, m.keys.Down):
			if m.selected < len(m.items)-1 {
				m.selected++
			}
			return m, nil
		case key.Matches(msg, m.keys.Start):
			cmd, err := m.startGame()
			if err != nil {
				panic(fmt.Errorf("failed to start game: %w", err))
			}
			return m, cmd
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.items[m.selected].model, cmd = m.items[m.selected].model.Update(msg)

	return m, cmd
}

func (m Model) View() string {
	items := make([]string, len(m.items))
	for i := range m.items {
		items[i] = m.renderItem(i) + "\n"
	}
	settingsStr := lipgloss.JoinVertical(lipgloss.Center, items...)

	output := lipgloss.JoinVertical(lipgloss.Center,
		titleStr,
		"\n",
		settingsStr,
		"\n\n",
		m.help.View(m.keys),
	)

	return output
}

func (m Model) renderItem(index int) string {
	i := m.items[index]
	output := i.model.View()
	if !i.hideLabel {
		label := lipgloss.NewStyle().Width(15).AlignHorizontal(lipgloss.Left).Render(i.label + ":")
		output = lipgloss.NewStyle().Width(25).AlignHorizontal(lipgloss.Right).Render(output)
		output = lipgloss.JoinHorizontal(lipgloss.Left, label, output)
	}

	if index == m.selected {
		return m.styles.settingSelected.Render(output)
	}

	return m.styles.settingUnselected.Render(output)
}

func (m Model) startGame() (tea.Cmd, error) {
	var level int
	var mode common.Mode
	var playerName string

	errInvalidModel := fmt.Errorf("invalid model for item %q", m.items[m.selected].label)
	errInvalidValue := fmt.Errorf("invalid value for model of item %q", m.items[m.selected].label)

	for _, i := range m.items {
		switch i.label {
		case "Starting Level":
			picker, ok := i.model.(*hpicker.Model)
			if !ok {
				return nil, errInvalidModel
			}
			level, ok = picker.GetSelection().Value.(int)
			if !ok {
				return nil, errInvalidValue
			}
		case "Mode":
			picker, ok := i.model.(*hpicker.Model)
			if !ok {
				return nil, errInvalidModel
			}
			mode, ok = picker.GetSelection().Value.(common.Mode)
			if !ok {
				return nil, errInvalidValue
			}
		case "Name":
			playerName = i.model.(textinput.Model).Child.Value()
		default:
			return nil, fmt.Errorf("invalid item label: %q", i.label)
		}
	}

	switch mode {
	case common.ModeMarathon:
		in := common.NewSingleInput(mode, level, playerName)
		return common.SwitchModeCmd(mode, in), nil
	case common.ModeUltra:
		in := common.NewSingleInput(mode, level, playerName)
		return common.SwitchModeCmd(mode, in), nil
	case common.ModeMenu, common.ModeLeaderboard:
		return nil, fmt.Errorf("invalid mode for starting game: %q", mode)
	default:
		return nil, fmt.Errorf("invalid mode from menu: %q", mode)
	}
}
