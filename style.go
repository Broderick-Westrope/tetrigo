package main

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	Program lipgloss.Style
}

func DefaultStyles() *Styles {
	return &Styles{
		Program: lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0),
	}
}
