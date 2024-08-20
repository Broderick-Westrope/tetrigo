package starter

import (
	"database/sql"
	"errors"
	"reflect"

	"github.com/Broderick-Westrope/tetrigo/internal/config"
	"github.com/Broderick-Westrope/tetrigo/internal/tui/common"
	"github.com/Broderick-Westrope/tetrigo/internal/tui/leaderboard"
	"github.com/Broderick-Westrope/tetrigo/internal/tui/marathon"
	"github.com/Broderick-Westrope/tetrigo/internal/tui/menu"
	"github.com/charmbracelet/bubbles/key"
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
	cfg      *config.Config
}

func NewInput(mode common.Mode, switchIn common.SwitchModeInput, db *sql.DB, cfg *config.Config) *Input {
	return &Input{
		mode:     mode,
		switchIn: switchIn,
		db:       db,
		cfg:      cfg,
	}
}

var _ tea.Model = &Model{}

type Model struct {
	child        tea.Model
	db           *sql.DB
	styles       *styles
	cfg          *config.Config
	forceQuitKey key.Binding
}

func NewModel(in *Input) (*Model, error) {
	m := &Model{
		db:           in.db,
		cfg:          in.cfg,
		styles:       defaultStyles(),
		forceQuitKey: key.NewBinding(key.WithKeys(in.cfg.Keys.ForceQuit...)),
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
	case tea.KeyMsg:
		if key.Matches(msg, m.forceQuitKey) {
			return m, tea.Quit
		}
	case common.SwitchModeMsg:
		err := m.setChild(msg.Target, msg.Input)
		if err != nil {
			panic(err)
		}
		return m, m.child.Init()
	case tea.WindowSizeMsg:
		// NOTE: Windows does not have support for reporting when resizes occur as it does not support the SIGWINCH signal.
		m.styles.programFullscreen.Width(msg.Width).Height(msg.Height)
	}

	var cmd tea.Cmd
	m.child, cmd = m.child.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	var output string

	output = m.child.View()
	output = m.styles.programFullscreen.Render(output)

	return output
}

func (m *Model) setChild(mode common.Mode, switchIn common.SwitchModeInput) error {
	if rv := reflect.ValueOf(switchIn); !rv.IsValid() || rv.IsNil() {
		return errors.New("switchIn is not valid")
	}

	switch mode {
	case common.ModeMenu:
		menuIn, ok := switchIn.(*common.MenuInput)
		if !ok {
			return ErrInvalidSwitchModeInput
		}
		m.child = menu.NewModel(menuIn)
	case common.ModeMarathon:
		marathonIn, ok := switchIn.(*common.MarathonInput)
		if !ok {
			return ErrInvalidSwitchModeInput
		}
		child, err := marathon.NewModel(marathonIn, m.cfg)
		if err != nil {
			return err
		}
		m.child = child
	case common.ModeLeaderboard:
		leaderboardIn, ok := switchIn.(*common.LeaderboardInput)
		if !ok {
			return ErrInvalidSwitchModeInput
		}
		child, err := leaderboard.NewModel(leaderboardIn, m.db)
		if err != nil {
			return err
		}
		m.child = child
	default:
		return ErrInvalidSwitchMode
	}
	return nil
}
