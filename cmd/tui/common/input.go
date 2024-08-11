package common

import "github.com/Broderick-Westrope/tetrigo/internal/data"

type MarathonInput struct {
	IsFullscreen bool
	Level        uint
	MaxLevel     uint
}

func NewMarathonInput(isFullscreen bool, level, maxLevel uint) *MarathonInput {
	return &MarathonInput{
		IsFullscreen: isFullscreen,
		Level:        level,
		MaxLevel:     maxLevel,
	}
}

func (in *MarathonInput) isSwitchModeInput() {}

type MenuInput struct {
	IsFullscreen bool
}

func NewMenuInput(isFullscreen bool) *MenuInput {
	return &MenuInput{
		IsFullscreen: isFullscreen,
	}
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
