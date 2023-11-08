package main

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	playfield [40][20]byte
	styles    *Styles
	keys      *KeyMap
}

func InitialModel() *Model {
	return &Model{
		playfield: [40][20]byte{},
		styles:    DefaultStyles(),
		keys:      DefaultKeyMap(),
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
		}
	}

	return m, nil
}

func (m Model) View() string {
	var output string

	for row := 0; row < 20; row++ {
		for col := range m.playfield[row] {
			switch m.playfield[row][col] {
			case 0:
				output += " "
			case 'I':
			case 'O':
			case 'T':
			case 'S':
			case 'Z':
			case 'J':
			case 'L':
			default:
				output += "?"
			}
		}
		if row < 19 {
			output += "\n"
		}
	}

	return m.styles.Program.Render(output)
}
