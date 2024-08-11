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
	mode     common.Mode
	switchIn common.SwitchModeInput
	db       *sql.DB
}

func NewInput(mode common.Mode, switchIn common.SwitchModeInput, db *sql.DB) *Input {
	return &Input{
		mode:     mode,
		switchIn: switchIn,
		db:       db,
	}
}

var _ tea.Model = &Model{}

type Model struct {
	child tea.Model
	db    *sql.DB
	keys  *common.Keys
}

func NewModel(in *Input) (*Model, error) {
	m := &Model{
		db:   in.db,
		keys: common.DefaultKeys(),
	}

	err := m.setChild(in.mode, in.switchIn)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Model) Init() tea.Cmd {
	return m.child.Init()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case common.SwitchModeMsg:
		err := m.setChild(msg.Target, msg.Input)
		if err != nil {
			panic(err)
		}
		return m, m.child.Init()
	}

	var cmd tea.Cmd
	m.child, cmd = m.child.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	return m.child.View()
}

func (m *Model) setChild(mode common.Mode, switchIn common.SwitchModeInput) error {
	if rv := reflect.ValueOf(switchIn); !rv.IsValid() || rv.IsNil() {
		return errors.New("switchIn is not valid")
	}

	switch mode {
	case common.MODE_MENU:
		menuIn, ok := switchIn.(*common.MenuInput)
		if !ok {
			return ErrInvalidSwitchModeInput
		}
		m.child = menu.NewModel(menuIn, m.keys)
	case common.MODE_MARATHON:
		marathonIn, ok := switchIn.(*common.MarathonInput)
		if !ok {
			return ErrInvalidSwitchModeInput
		}
		child, err := marathon.NewModel(marathonIn, m.keys)
		if err != nil {
			return err
		}
		m.child = child
	case common.MODE_LEADERBOARD:
		leaderboardIn, ok := switchIn.(*common.LeaderboardInput)
		if !ok {
			return ErrInvalidSwitchModeInput
		}
		child, err := leaderboard.NewModel(leaderboardIn, m.db, m.keys)
		if err != nil {
			return err
		}
		m.child = child
	default:
		return ErrInvalidSwitchMode
	}
	return nil
}
