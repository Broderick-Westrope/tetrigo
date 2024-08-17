package menu

import (
	"fmt"
	"strconv"

	"github.com/Broderick-Westrope/tetrigo/cmd/tetrigo/common"
	"github.com/Broderick-Westrope/tetrigo/cmd/tetrigo/components/hpicker"
	"github.com/Broderick-Westrope/tetrigo/cmd/tetrigo/components/textinput"
	"github.com/Broderick-Westrope/tetrigo/internal/config"
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
	game     tea.Model

	keys   *keyMap
	styles *styles
	help   help.Model
}

type item struct {
	label     string
	model     tea.Model
	hideLabel bool
}

func NewModel(_ *common.MenuInput, keys *config.Keys) *Model {
	nameInput := textinput.NewModel("Enter your name", 20, 20)
	modePicker := hpicker.NewModel([]string{"Marathon"}, keys)
	playersPicker := hpicker.NewModel(nil, keys, hpicker.WithRange(1, 1))
	levelPicker := hpicker.NewModel(nil, keys, hpicker.WithRange(1, 15))

	return &Model{
		items: []item{
			{label: "Name", model: nameInput, hideLabel: true},
			{label: "Mode", model: modePicker},
			{label: "Players", model: playersPicker},
			{label: "Level", model: levelPicker},
		},
		selected: 0,

		keys:   constructKeyMap(keys),
		styles: defaultStyles(),
		help:   help.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
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
		label := lipgloss.NewStyle().Width(12).AlignHorizontal(lipgloss.Left).Render(i.label + ":")
		output = lipgloss.NewStyle().Width(12).AlignHorizontal(lipgloss.Right).Render(output)
		output = lipgloss.JoinHorizontal(lipgloss.Left, label, output)
	}

	if index == m.selected {
		return m.styles.settingSelected.Render(output)
	}

	return m.styles.settingUnselected.Render(output)
}

func (m Model) startGame() (tea.Cmd, error) {
	var level uint
	var players uint
	var mode string
	var playerName string

	for _, i := range m.items {
		switch i.label {
		case "Level":
			value := i.model.(*hpicker.Model).GetSelection()
			intLevel, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("failed to convert level string %q to int", value)
			}
			level = uint(intLevel)
		case "Players":
			value := i.model.(*hpicker.Model).GetSelection()
			intPlayers, err := strconv.Atoi(value)
			if err != nil {
				return nil, fmt.Errorf("failed to convert players string %q to int", value)
			}
			players = uint(intPlayers)
		case "Mode":
			mode = i.model.(*hpicker.Model).GetSelection()
		case "Name":
			playerName = i.model.(textinput.Model).Child.Value()
		default:
			return nil, fmt.Errorf("invalid item label: %q", i.label)
		}
	}

	// TODO: use players
	_ = players

	switch mode {
	case "Marathon":
		in := common.NewMarathonInput(level, playerName)
		return common.SwitchModeCmd(common.MODE_MARATHON, in), nil
	default:
		return nil, fmt.Errorf("invalid mode: %q", mode)
	}
}
