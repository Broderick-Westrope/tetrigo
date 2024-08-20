package common

import "github.com/Broderick-Westrope/tetrigo/internal/data"

type MarathonInput struct {
	Level      uint
	PlayerName string
}

func NewMarathonInput(level uint, playerName string) *MarathonInput {
	return &MarathonInput{
		Level:      level,
		PlayerName: playerName,
	}
}

func (in *MarathonInput) isSwitchModeInput() {}

type UltraInput struct {
	Level      uint
	PlayerName string
}

func NewUltraInput(level uint, playerName string) *UltraInput {
	return &UltraInput{
		Level:      level,
		PlayerName: playerName,
	}
}

func (in *UltraInput) isSwitchModeInput() {}

type MenuInput struct {
}

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
