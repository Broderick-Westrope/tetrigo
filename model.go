package main

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Playfield [40][10]byte

type Model struct {
	playfield Playfield
	styles    *Styles
	help      help.Model
	keys      *KeyMap
}

func InitialModel() *Model {
	return &Model{
		playfield: Playfield{
			{0, 'I', 'O', 'T', 'S', 'Z', 'J', 'L', 'G'},
		},
		styles: DefaultStyles(),
		help:   help.New(),
		keys:   DefaultKeyMap(),
	}
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
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		}
	}

	return m, nil
}

func (m Model) View() string {
	var output string

	for row := (len(m.playfield) - 20); row < len(m.playfield); row++ {
		for col := range m.playfield[row] {
			switch m.playfield[row][col] {
			case 0:
				output += m.styles.ColIndicator.Render("▕ ")
			case 'G':
				output += "░░"
			default:
				cellStyle, ok := m.styles.TetriminoStyles[m.playfield[row][col]]
				if ok {
					output += cellStyle.Render("██")
				} else {
					output += "? "
				}
			}
		}
		if row < 19 {
			output += "\n"
		}
	}

	return m.styles.Program.Render(output) + "\n" + m.help.View(m.keys)
}
