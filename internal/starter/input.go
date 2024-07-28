package starter

import (
	"github.com/Broderick-Westrope/tetrigo/internal/marathon"
	"github.com/Broderick-Westrope/tetrigo/internal/menu"
)

type Input struct {
	isFullscreen bool
	level        uint
	mode         Mode
}

func NewInput(mode Mode, isFullscreen bool, level uint) *Input {
	return &Input{
		mode:         mode,
		isFullscreen: isFullscreen,
		level:        level,
	}
}

func (in *Input) ToMarathonInput() *marathon.Input {
	return marathon.NewInput(in.isFullscreen, in.level)
}

func (in *Input) ToMenuInput() *menu.Input {
	return menu.NewInput(in.isFullscreen)
}
