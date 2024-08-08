package starter

import (
	"github.com/Broderick-Westrope/tetrigo/internal"
	"github.com/Broderick-Westrope/tetrigo/internal/marathon"
	"github.com/Broderick-Westrope/tetrigo/internal/menu"
	tea "github.com/charmbracelet/bubbletea"
)

type Input struct {
	isFullscreen bool
	level        uint
	mode         internal.Mode
}

func NewInput(mode internal.Mode, isFullscreen bool, level uint) *Input {
	return &Input{
		mode:         mode,
		isFullscreen: isFullscreen,
		level:        level,
	}
}

func (in *Input) ToMarathonInput() *marathon.Input {
	return marathon.NewInput(in.isFullscreen, in.level, 15)
}

func (in *Input) ToMenuInput() *menu.Input {
	return menu.NewInput(in.isFullscreen)
}

func (in *Input) getChild(mode internal.Mode) (tea.Model, error) {
	switch mode {
	case internal.MODE_MENU:
		return menu.NewModel(in.ToMenuInput()), nil
	case internal.MODE_MARATHON:
		return marathon.NewModel(in.ToMarathonInput())
	default:
		panic("invalid mode")
	}
}
