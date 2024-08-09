package tui

import tea "github.com/charmbracelet/bubbletea"

type Model interface {
	Init() tea.Cmd
	Update(tea.Msg) (Model, tea.Cmd)
}
