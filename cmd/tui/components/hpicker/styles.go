package hpicker

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	NextIndicator     string
	NextStyle         lipgloss.Style
	NextDisabledStyle lipgloss.Style

	PrevIndicator     string
	PrevStyle         lipgloss.Style
	PrevDisabledStyle lipgloss.Style

	SelectionStyle lipgloss.Style
}

func DefaultStyles() Styles {
	return Styles{
		NextIndicator:     " >",
		NextStyle:         lipgloss.NewStyle().Foreground(lipgloss.Color("#FF00FF")),
		NextDisabledStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")),

		PrevIndicator:     "< ",
		PrevStyle:         lipgloss.NewStyle().Foreground(lipgloss.Color("#FF00FF")),
		PrevDisabledStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")),

		SelectionStyle: lipgloss.NewStyle(),
	}
}
