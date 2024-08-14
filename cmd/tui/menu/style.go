package menu

import "github.com/charmbracelet/lipgloss"

type styles struct {
	settingSelected   lipgloss.Style
	settingUnselected lipgloss.Style
}

func defaultStyles() *styles {
	s := styles{
		settingSelected: lipgloss.NewStyle().Padding(1, 2),
	}
	s.settingUnselected = s.settingSelected.Copy().Foreground(lipgloss.Color("241"))
	return &s
}
