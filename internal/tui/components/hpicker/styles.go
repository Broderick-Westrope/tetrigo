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
		NextStyle:         lipgloss.NewStyle(),
		NextDisabledStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("241")),

		PrevIndicator:     "< ",
		PrevStyle:         lipgloss.NewStyle(),
		PrevDisabledStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("241")),

		SelectionStyle: lipgloss.NewStyle(),
	}
}
