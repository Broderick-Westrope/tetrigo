package menu

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	settingSelected   lipgloss.Style
	settingUnselected lipgloss.Style
}

func DefaultStyles() *Styles {
	s := Styles{
		settingSelected: lipgloss.NewStyle().Padding(1, 2),
	}
	s.settingUnselected = s.settingSelected.Copy().Foreground(lipgloss.Color("241"))
	return &s
}
