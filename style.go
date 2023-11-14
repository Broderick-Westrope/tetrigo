package main

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	Playfield       lipgloss.Style
	ColIndicator    lipgloss.Style
	TetriminoStyles map[byte]lipgloss.Style
	Hold            lipgloss.Style
	Information     lipgloss.Style
	RowIndicator    lipgloss.Style
}

func DefaultStyles() *Styles {
	s := Styles{
		Playfield:    lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0),
		ColIndicator: lipgloss.NewStyle().Foreground(lipgloss.Color("#303040")),
		TetriminoStyles: map[byte]lipgloss.Style{
			'I': lipgloss.NewStyle().Foreground(lipgloss.Color("#64C4EB")),
			'O': lipgloss.NewStyle().Foreground(lipgloss.Color("#F1D448")),
			'T': lipgloss.NewStyle().Foreground(lipgloss.Color("#A15398")),
			'S': lipgloss.NewStyle().Foreground(lipgloss.Color("#64B452")),
			'Z': lipgloss.NewStyle().Foreground(lipgloss.Color("#DC3A35")),
			'J': lipgloss.NewStyle().Foreground(lipgloss.Color("#5C65A8")),
			'L': lipgloss.NewStyle().Foreground(lipgloss.Color("#E07F3A")),
		},
		Hold:         lipgloss.NewStyle().Width(10).Height(5).Border(lipgloss.RoundedBorder(), true, false, true, true).Align(lipgloss.Center, lipgloss.Center),
		Information:  lipgloss.NewStyle().Width(13).Align(lipgloss.Left, lipgloss.Top),
		RowIndicator: lipgloss.NewStyle().Foreground(lipgloss.Color("#444049")).Align(lipgloss.Left).Padding(0, 1, 0),
	}
	return &s
}
