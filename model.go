package main

import tea "github.com/charmbracelet/bubbletea"

type Model struct {
	playfield [40][20]byte
}

func InitialModel() *Model {
	return &Model{
		playfield: [40][20]byte{},
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	return ""
}
