package tui

import tea "github.com/charmbracelet/bubbletea"

// FatalErrorMsg encloses an error which should be set on the starter model before exiting the program.
type FatalErrorMsg error

// FatalErrorCmd returns a command for creating a new FatalErrorMsg with the given error.
func FatalErrorCmd(err error) tea.Cmd {
	return func() tea.Msg {
		return FatalErrorMsg(err)
	}
}
