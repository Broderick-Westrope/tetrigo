package menu

import (
	"fmt"

	"github.com/Broderick-Westrope/tetrigo/internal/marathon"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	modeMenu = iota
	modeGame
)

type option interface{}

type setting struct {
	name    string
	options []option
	index   int
}

type Model struct {
	settings     []setting
	settingIndex int
	game         tea.Model
	mode         int

	keys   *keyMap
	styles *styles
	help   help.Model

	isFullscreen bool
}

func InitialModel(isFullscreen bool) *Model {
	m := Model{
		settings: []setting{
			{
				name:    "Level",
				options: []option{1, 5, 10, 15},
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
		mode:         modeMenu,
		help:         help.New(),
		isFullscreen: isFullscreen,
	}
	return &m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.mode == modeGame {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.keys.Quit):
				m.mode = modeMenu
				m.game = nil
				return m, nil
			}
		}
		var cmd tea.Cmd
		m.game, cmd = m.game.Update(msg)
		return m, cmd
	}

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
	if m.mode == modeGame {
		return m.game.View()
	}

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
	output := `    ______________________  ______________ 
   /_  __/ ____/_  __/ __ \/  _/ ____/ __ \
    / / / __/   / / / /_/ // // / __/ / / /
   / / / /___  / / / _, _// // /_/ / /_/ / 
  /_/ /_____/ /_/ /_/ |_/___/\____/\____/  
`
	return output
}

func (m *Model) renderSetting(index int, isSelected bool) string {
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

func (m *Model) startGame() (tea.Cmd, error) {
	var level uint
	var mode string
	// var players uint
	for _, setting := range m.settings {
		switch setting.name {
		case "Level":
			intLevel := setting.options[setting.index].(int)
			level = uint(intLevel)
		// case "Players":
		// 	players = setting.options[setting.index].(uint)
		case "Mode":
			mode = setting.options[setting.index].(string)
		}
	}

	switch mode {
	case "Marathon":
		m.mode = modeGame

		var fullscreen *marathon.FullscreenInfo
		if m.isFullscreen {
			fullscreen = &marathon.FullscreenInfo{
				Width:  m.styles.programFullscreen.GetWidth(),
				Height: m.styles.programFullscreen.GetHeight(),
			}
		}

		m.game = marathon.InitialModel(level, fullscreen)
		return m.game.Init(), nil
	}
	return nil, fmt.Errorf("invalid mode: %v", mode)
}
