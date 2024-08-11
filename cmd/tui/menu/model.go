package menu

import (
	"fmt"

	"github.com/Broderick-Westrope/tetrigo/cmd/tui/common"

	//"github.com/Broderick-Westrope/tetrigo/internal/starter"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type option interface{}

type setting struct {
	name    string
	options []option
	index   int
}

var _ tea.Model = Model{}

type Model struct {
	settings     []setting
	settingIndex int
	game         tea.Model

	keys   *keyMap
	styles *styles
	help   help.Model

	isFullscreen bool
}

func NewModel(in *common.MenuInput) *Model {
	m := Model{
		settings: []setting{
			{
				name:    "Level",
				options: []option{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
				index:   0,
			},
			{
				name:    "Players",
				options: []option{uint(1)},
				index:   0,
			},
			{
				name:    "Mode",
				options: []option{"Marathon"},
				index:   0,
			},
		},
		settingIndex: 0,
		keys:         defaultKeyMap(),
		styles:       defaultStyles(),
		help:         help.New(),
		isFullscreen: in.IsFullscreen,
	}
	return &m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Left):
			m.settingIndex--
			if m.settingIndex < 0 {
				m.settingIndex = len(m.settings) - 1
			}
		case key.Matches(msg, m.keys.Right):
			m.settingIndex++
			if m.settingIndex >= len(m.settings) {
				m.settingIndex = 0
			}
		case key.Matches(msg, m.keys.Up):
			m.settings[m.settingIndex].index--
			if m.settings[m.settingIndex].index < 0 {
				m.settings[m.settingIndex].index = len(m.settings[m.settingIndex].options) - 1
			}
		case key.Matches(msg, m.keys.Down):
			m.settings[m.settingIndex].index++
			if m.settings[m.settingIndex].index >= len(m.settings[m.settingIndex].options) {
				m.settings[m.settingIndex].index = 0
			}
		case key.Matches(msg, m.keys.Start):
			cmd, err := m.startGame()
			if err != nil {
				panic(fmt.Errorf("failed to start game: %w", err))
			}
			return m, cmd
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		}
	case tea.WindowSizeMsg:
		m.styles.programFullscreen.Width(msg.Width).Height(msg.Height)
	}

	return m, nil
}

func (m Model) View() string {
	settings := make([]string, len(m.settings))
	for i := range m.settings {
		settings[i] = m.renderSetting(i, i == m.settingIndex)
	}

	output := lipgloss.JoinVertical(lipgloss.Center,
		renderTitle(),
		lipgloss.JoinHorizontal(lipgloss.Top, settings...),
	) + "\n" + m.help.View(m.keys)

	if m.isFullscreen {
		return m.styles.programFullscreen.Render(output)
	}
	return output
}

func renderTitle() string {
	output := `
    ______________________  ______________ 
   /_  __/ ____/_  __/ __ \/  _/ ____/ __ \
    / / / __/   / / / /_/ // // / __/ / / /
   / / / /___  / / / _, _// // /_/ / /_/ / 
  /_/ /_____/ /_/ /_/ |_/___/\____/\____/  
`
	return output
}

func (m Model) renderSetting(index int, isSelected bool) string {
	output := fmt.Sprintf("%v:\n", m.settings[index].name)
	for i, option := range m.settings[index].options {
		if i == m.settings[index].index {
			output += " > "
		} else {
			output += "   "
		}
		output += fmt.Sprintf("%v", option)
		if i < len(m.settings[index].options)-1 {
			output += "\n"
		}
	}
	if isSelected {
		return m.styles.settingSelected.Render(output)
	}
	return m.styles.settingUnselected.Render(output)
}

func (m Model) startGame() (tea.Cmd, error) {
	var level uint
	var mode string
	// var players uint
	for _, s := range m.settings {
		switch s.name {
		case "Level":
			intLevel := s.options[s.index].(int)
			level = uint(intLevel)
		// case "Players":
		// 	players = setting.options[setting.index].(uint)
		case "Mode":
			mode = s.options[s.index].(string)
		}
	}

	switch mode {
	case "Marathon":
		in := &common.MarathonInput{
			Level: level,
		}
		return common.SwitchModeCmd(common.MODE_MARATHON, in), nil
	default:
		return nil, fmt.Errorf("invalid mode: %v", mode)
	}
}
