package views

import "github.com/charmbracelet/lipgloss"

type menuStyles struct {
	settingSelected   lipgloss.Style
	settingUnselected lipgloss.Style
}

func defaultMenuStyles() *menuStyles {
	s := menuStyles{
		settingSelected: lipgloss.NewStyle().Padding(0, 2),
	}
	s.settingUnselected = s.settingSelected.Copy().Foreground(lipgloss.Color("241"))
	return &s
}
