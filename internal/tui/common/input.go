package common

import "github.com/Broderick-Westrope/tetrigo/internal/data"

type MarathonInput struct {
	Level      int
	PlayerName string
}

func NewMarathonInput(level int, playerName string) *MarathonInput {
	return &MarathonInput{
		Level:      level,
		PlayerName: playerName,
	}
}

func (in *MarathonInput) isSwitchModeInput() {}

type SingleInput struct {
	Mode       Mode
	Level      int
	PlayerName string
}

func NewSingleInput(mode Mode, level int, playerName string) *SingleInput {
	return &SingleInput{
		Mode:       mode,
		Level:      level,
		PlayerName: playerName,
	}
}

func (in *SingleInput) isSwitchModeInput() {}

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
