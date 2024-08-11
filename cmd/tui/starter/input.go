package starter

import (
	"database/sql"
	"errors"
	"reflect"

	"github.com/Broderick-Westrope/tetrigo/cmd/tui/common"
	"github.com/Broderick-Westrope/tetrigo/cmd/tui/leaderboard"
	"github.com/Broderick-Westrope/tetrigo/cmd/tui/marathon"
	"github.com/Broderick-Westrope/tetrigo/cmd/tui/menu"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	ErrInvalidSwitchModeInput = errors.New("invalid SwitchModeInput")
	ErrInvalidSwitchMode      = errors.New("invalid SwitchMode")
)

type Input struct {
	isFullscreen bool
	level        uint
	mode         common.Mode
	db           *sql.DB
}

func NewInput(mode common.Mode, isFullscreen bool, level uint, db *sql.DB) *Input {
	return &Input{
		mode:         mode,
		isFullscreen: isFullscreen,
		level:        level,
		db:           db,
	}
}

func (in *Input) getChild(mode common.Mode, switchIn common.SwitchModeInput) (tea.Model, error) {
	// If the switchIn is nil, then we use the default input for the mode.
	if rv := reflect.ValueOf(switchIn); !rv.IsValid() || rv.IsNil() {
		switch mode {
		case common.MODE_MENU:
			return menu.NewModel(
				common.NewMenuInput(in.isFullscreen),
			), nil
		case common.MODE_MARATHON:
			return marathon.NewModel(
				common.NewMarathonInput(in.isFullscreen, in.level, 15),
			)
		case common.MODE_LEADERBOARD:
			return leaderboard.NewModel(
				common.NewLeaderboardInput(in.db, "marathon"),
			)
		default:
			return nil, ErrInvalidSwitchMode
		}
	}

	switch mode {
	case common.MODE_MENU:
		menuIn, ok := switchIn.(*common.MenuInput)
		if !ok {
			return nil, ErrInvalidSwitchModeInput
		}
		return menu.NewModel(menuIn), nil
	case common.MODE_MARATHON:
		marathonIn, ok := switchIn.(*common.MarathonInput)
		if !ok {
			return nil, ErrInvalidSwitchModeInput
		}
		return marathon.NewModel(marathonIn)
	case common.MODE_LEADERBOARD:
		leaderboardIn, ok := switchIn.(*common.LeaderboardInput)
		if !ok {
			return nil, ErrInvalidSwitchModeInput
		}
		return leaderboard.NewModel(leaderboardIn)
	default:
		return nil, ErrInvalidSwitchMode
	}
}
