package starter

import "github.com/charmbracelet/lipgloss"

type styles struct {
	programFullscreen lipgloss.Style
}

func defaultStyles() *styles {
	s := styles{
		// NOTE: The width and height is set  dynamically in the update loop on any tea.WindowSizeMsg
		programFullscreen: lipgloss.NewStyle().Align(lipgloss.Center, lipgloss.Center),
	}
	return &s
}
