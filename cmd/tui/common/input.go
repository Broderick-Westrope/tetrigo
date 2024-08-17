package common

import "github.com/Broderick-Westrope/tetrigo/internal/data"

type MarathonInput struct {
	Level      uint
	MaxLevel   uint
	PlayerName string
}

func NewMarathonInput(level, maxLevel uint, playerName string) *MarathonInput {
	return &MarathonInput{
		Level:      level,
		MaxLevel:   maxLevel,
		PlayerName: playerName,
	}
}

func (in *MarathonInput) isSwitchModeInput() {}

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
