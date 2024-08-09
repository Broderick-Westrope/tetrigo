package starter

import (
	"github.com/Broderick-Westrope/tetrigo/internal/tui"
	"github.com/Broderick-Westrope/tetrigo/internal/tui/marathon"
	"github.com/Broderick-Westrope/tetrigo/internal/tui/menu"
	tea "github.com/charmbracelet/bubbletea"
)

type Input struct {
	isFullscreen bool
	level        uint
	mode         tui.Mode
}

func NewInput(mode tui.Mode, isFullscreen bool, level uint) *Input {
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

func (in *Input) getChild(mode tui.Mode) (tea.Model, error) {
	switch mode {
	case tui.MODE_MENU:
		return menu.NewModel(in.ToMenuInput()), nil
	case tui.MODE_MARATHON:
		return marathon.NewModel(in.ToMarathonInput())
	default:
		panic("invalid mode")
	}
}
