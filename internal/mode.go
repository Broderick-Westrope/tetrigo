package internal

import tea "github.com/charmbracelet/bubbletea"

type Mode int

const (
	MODE_MENU = Mode(iota)
	MODE_MARATHON
)

type SwitchModeMsg struct {
	Target Mode
	Level  uint
}

func SwitchModeCmd(target Mode, level uint) tea.Cmd {
	return func() tea.Msg {
		return SwitchModeMsg{
			Target: target,
			Level:  level,
		}
	}
}
