package views

import (
	"fmt"

	"github.com/Broderick-Westrope/tetrigo/internal/tui"
	"github.com/Broderick-Westrope/tetrigo/internal/tui/components/hpicker"
	"github.com/Broderick-Westrope/tetrigo/internal/tui/components/textinput"
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

var _ tea.Model = MenuModel{}

type MenuModel struct {
	items    []menuItem
	selected int

	keys   *menuKeyMap
	styles *menuStyles
	help   help.Model
}

type menuItem struct {
	label     string
	model     tea.Model
	hideLabel bool
}

func NewMenuModel(_ *tui.MenuInput) *MenuModel {
	nameInput := textinput.NewTextInputModel("Enter your name", 20, 20)
	modePicker := hpicker.NewHPickerModel([]hpicker.KeyValuePair{
		{Key: "Marathon", Value: tui.ModeMarathon},
		{Key: "Sprint (40 Lines)", Value: tui.ModeSprint},
		{Key: "Ultra (Time Trial)", Value: tui.ModeUltra},
	})
	levelPicker := hpicker.NewHPickerModel(nil, hpicker.WithRange(1, 15))

	return &MenuModel{
		items: []menuItem{
			{label: "Name", model: nameInput, hideLabel: true},
			{label: "Mode", model: modePicker},
			{label: "Starting Level", model: levelPicker},
		},
		selected: 0,

		keys:   defaultMenuKeyMap(),
		styles: defaultMenuStyles(),
		help:   help.New(),
	}
}

func (m MenuModel) Init() tea.Cmd {
	return nil
}

func (m MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m MenuModel) View() string {
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

func (m MenuModel) renderItem(index int) string {
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

func (m MenuModel) startGame() (tea.Cmd, error) {
	var level int
	var mode tui.Mode
	var playerName string

	errInvalidModel := fmt.Errorf("invalid model for item %q", m.items[m.selected].label)
	errInvalidValue := fmt.Errorf("invalid value for model of item %q", m.items[m.selected].label)

	for _, i := range m.items {
		switch i.label {
		case "Starting Level":
			picker, ok := i.model.(*hpicker.HPickerModel)
			if !ok {
				return nil, errInvalidModel
			}
			level, ok = picker.GetSelection().Value.(int)
			if !ok {
				return nil, errInvalidValue
			}
		case "Mode":
			picker, ok := i.model.(*hpicker.HPickerModel)
			if !ok {
				return nil, errInvalidModel
			}
			mode, ok = picker.GetSelection().Value.(tui.Mode)
			if !ok {
				return nil, errInvalidValue
			}
		case "Name":
			playerName = i.model.(textinput.TextInputModel).Child.Value()
		default:
			return nil, fmt.Errorf("invalid item label: %q", i.label)
		}
	}

	switch mode {
	case tui.ModeMarathon, tui.ModeSprint, tui.ModeUltra:
		in := tui.NewSingleInput(mode, level, playerName)
		return tui.SwitchModeCmd(mode, in), nil
	case tui.ModeMenu, tui.ModeLeaderboard:
		return nil, fmt.Errorf("invalid mode for starting game: %q", mode)
	default:
		return nil, fmt.Errorf("invalid mode from menu: %q", mode)
	}
}
