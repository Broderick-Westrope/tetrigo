package menu

import "github.com/charmbracelet/lipgloss"

type styles struct {
	settingSelected   lipgloss.Style
	settingUnselected lipgloss.Style
	programFullscreen lipgloss.Style
}

func defaultStyles() *styles {
	s := styles{
		settingSelected:   lipgloss.NewStyle().Padding(1, 2),
		programFullscreen: lipgloss.NewStyle().Align(lipgloss.Center, lipgloss.Center),
	}
	s.settingUnselected = s.settingSelected.Copy().Foreground(lipgloss.Color("241"))
	return &s
}
