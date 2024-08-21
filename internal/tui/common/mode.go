package common

import (
	tea "github.com/charmbracelet/bubbletea"
)

type SwitchModeMsg struct {
	Target Mode
	Input  SwitchModeInput
}

type SwitchModeInput interface {
	isSwitchModeInput()
}

func SwitchModeCmd(target Mode, in SwitchModeInput) tea.Cmd {
	return func() tea.Msg {
		return SwitchModeMsg{
			Target: target,
			Input:  in,
		}
	}
}

type Mode int

const (
	ModeMenu = Mode(iota)
	ModeMarathon
	ModeSprint
	ModeUltra
	ModeLeaderboard
)

var modeToStrMap = map[Mode]string{
	ModeMenu:        "Menu",
	ModeMarathon:    "Marathon",
	ModeSprint:      "Sprint",
	ModeUltra:       "Ultra",
	ModeLeaderboard: "Leaderboard",
}

func (m Mode) String() string {
	return modeToStrMap[m]
}
