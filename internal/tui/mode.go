package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/stuttgart-things/sthings-tetris/internal/data"
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

// SwitchModeInput values --------------------------------------------------

type SingleInput struct {
	Mode     Mode
	Level    int
	Username string
}

func NewSingleInput(mode Mode, level int, username string) *SingleInput {
	return &SingleInput{
		Mode:     mode,
		Level:    level,
		Username: username,
	}
}

func (in *SingleInput) isSwitchModeInput() {}

type MenuInput struct{}

func NewMenuInput() *MenuInput {
	return &MenuInput{}
}

func (in *MenuInput) isSwitchModeInput() {}

type LeaderboardInput struct {
	GameMode string
	NewEntry *data.Score
}

func NewLeaderboardInput(gameMode string, opts ...func(input *LeaderboardInput)) *LeaderboardInput {
	in := &LeaderboardInput{
		GameMode: gameMode,
	}

	for _, opt := range opts {
		opt(in)
	}

	return in
}

func (in *LeaderboardInput) isSwitchModeInput() {}

func WithNewEntry(entry *data.Score) func(input *LeaderboardInput) {
	return func(in *LeaderboardInput) {
		in.NewEntry = entry
	}
}
