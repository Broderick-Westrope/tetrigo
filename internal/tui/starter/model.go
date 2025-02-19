package starter

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"

	"github.com/Broderick-Westrope/charmutils"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/stuttgart-things/sthings-tetris/internal/config"
	"github.com/stuttgart-things/sthings-tetris/internal/tui"
	"github.com/stuttgart-things/sthings-tetris/internal/tui/views"
)

type Input struct {
	mode     tui.Mode
	switchIn tui.SwitchModeInput
	db       *sql.DB
	cfg      *config.Config
}

func NewInput(mode tui.Mode, switchIn tui.SwitchModeInput, db *sql.DB, cfg *config.Config) *Input {
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
	cfg          *config.Config
	forceQuitKey key.Binding

	width  int
	height int

	ExitError error
}

func NewModel(in *Input) (*Model, error) {
	m := &Model{
		db:           in.db,
		cfg:          in.cfg,
		forceQuitKey: key.NewBinding(key.WithKeys(in.cfg.Keys.ForceQuit...)),
	}

	err := m.setChild(in.mode, in.switchIn)
	if err != nil {
		return nil, fmt.Errorf("setting child model: %w", err)
	}
	return m, nil
}

func (m *Model) Init() tea.Cmd {
	return m.initChild()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tui.FatalErrorMsg:
		m.ExitError = msg
		return m, tea.Quit

	case tea.KeyMsg:
		if key.Matches(msg, m.forceQuitKey) {
			return m, tea.Quit
		}

	case tui.SwitchModeMsg:
		err := m.setChild(msg.Target, msg.Input)
		if err != nil {
			return m, tui.FatalErrorCmd(fmt.Errorf("setting child model: %w", err))
		}
		cmd := m.initChild()
		return m, cmd

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	var cmd tea.Cmd
	m.child, cmd = m.child.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	return m.child.View()
}

func (m *Model) setChild(mode tui.Mode, switchIn tui.SwitchModeInput) error {
	if rv := reflect.ValueOf(switchIn); !rv.IsValid() || rv.IsNil() {
		return errors.New("switchIn is not valid")
	}

	switch mode {
	case tui.ModeMenu:
		menuIn, ok := switchIn.(*tui.MenuInput)
		if !ok {
			return fmt.Errorf("switchIn is not a MenuInput: %w", charmutils.ErrInvalidTypeAssertion)
		}
		m.child = views.NewMenuModel(menuIn)

	case tui.ModeMarathon, tui.ModeSprint, tui.ModeUltra:
		singleIn, ok := switchIn.(*tui.SingleInput)
		if !ok {
			return fmt.Errorf("switchIn is not a SingleInput: %w", charmutils.ErrInvalidTypeAssertion)
		}
		child, err := views.NewSingleModel(singleIn, m.cfg)
		if err != nil {
			return fmt.Errorf("creating single model: %w", err)
		}
		m.child = child

	case tui.ModeLeaderboard:
		leaderboardIn, ok := switchIn.(*tui.LeaderboardInput)
		if !ok {
			return fmt.Errorf("switchIn is not a LeaderboardInput: %w", charmutils.ErrInvalidTypeAssertion)
		}
		child, err := views.NewLeaderboardModel(leaderboardIn, m.db)
		if err != nil {
			return fmt.Errorf("creating leaderboard model: %w", err)
		}
		m.child = child

	default:
		return errors.New("invalid Mode")
	}
	return nil
}

func (m *Model) initChild() tea.Cmd {
	var cmds []tea.Cmd
	cmd := m.child.Init()
	cmds = append(cmds, cmd)
	m.child, cmd = m.child.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
	cmds = append(cmds, cmd)
	return tea.Batch(cmds...)
}
