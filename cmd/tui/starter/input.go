package starter

import (
	"github.com/Broderick-Westrope/tetrigo/cmd/tui/common"
	"github.com/Broderick-Westrope/tetrigo/cmd/tui/marathon"
	"github.com/Broderick-Westrope/tetrigo/cmd/tui/menu"
	tea "github.com/charmbracelet/bubbletea"
)

type Input struct {
	isFullscreen bool
	level        uint
	mode         common.Mode
}

func NewInput(mode common.Mode, isFullscreen bool, level uint) *Input {
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

func (in *Input) getChild(mode common.Mode) (tea.Model, error) {
	switch mode {
	case common.MODE_MENU:
		return menu.NewModel(in.ToMenuInput()), nil
	case common.MODE_MARATHON:
		return marathon.NewModel(in.ToMarathonInput())
	default:
		panic("invalid mode")
	}
}
