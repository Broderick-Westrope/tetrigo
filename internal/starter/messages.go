package starter

import tea "github.com/charmbracelet/bubbletea"

type switchModeMsg struct {
	target Mode
}

func SwitchModeCmd(target Mode) tea.Msg {
	return func() tea.Msg {
		return switchModeMsg{target}
	}
}
