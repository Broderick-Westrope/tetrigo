package common

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
}

func NewLeaderboardInput(gameMode string) *LeaderboardInput {
	return &LeaderboardInput{
		GameMode: gameMode,
	}
}

func (in *LeaderboardInput) isSwitchModeInput() {}
