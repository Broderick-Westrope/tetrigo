package views

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	"github.com/Broderick-Westrope/tetrigo/internal/tui"
)

const (
	titleStr = `
    ______________________  ______________ 
   /_  __/ ____/_  __/ __ \/  _/ ____/ __ \
    / / / __/   / / / /_/ // // / __/ / / /
   / / / /___  / / / _, _// // /_/ / /_/ / 
  /_/ /_____/ /_/ /_/ |_/___/\____/\____/  
`
	formKeyUsername = "username"
	formKeyGameMode = "game_mode"
	formKeyLevel    = "level"
)

var _ tea.Model = &MenuModel{}

type MenuModel struct {
	form                   *huh.Form
	hasAnnouncedCompletion bool
	keys                   *menuKeyMap

	width  int
	height int
}

type menuItem struct {
	label     string
	model     tea.Model
	hideLabel bool
}

func NewMenuModel(_ *tui.MenuInput) *MenuModel {
	keys := defaultMenuKeyMap()
	return &MenuModel{
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewInput().Key(formKeyUsername).
					Title("Username:").CharLimit(20).
					Validate(func(s string) error {
						if s == "" {
							return errors.New("empty username not allowed")
						}
						return nil
					}).WithWidth(20),
				huh.NewSelect[tui.Mode]().Key(formKeyGameMode).
					Options(
						huh.NewOption("Marathon", tui.ModeMarathon),
						huh.NewOption("Sprint (40 Lines)", tui.ModeSprint),
						huh.NewOption("Ultra (Time Trial)", tui.ModeUltra),
					),
				huh.NewSelect[int]().Key(formKeyLevel).
					Options(tui.NewIntRangeOptions(1, 15)...),
			),
		).WithKeyMap(keys.formKeys),
		keys: keys,
	}
}

func (m *MenuModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m *MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Exit):
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		formWidth := msg.Width / 2
		formWidth = min(formWidth, 50)
		m.form = m.form.WithWidth(formWidth)
		return m, nil
	}

	if m.form.State == huh.StateCompleted {
		switch m.hasAnnouncedCompletion {
		case false:
			return m, m.announceCompletion()
		default:
			return m, nil
		}
	}

	var cmds []tea.Cmd
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *MenuModel) announceCompletion() tea.Cmd {
	username := m.form.GetString(formKeyUsername)
	level := m.form.GetInt(formKeyLevel)
	mode, ok := m.form.Get(formKeyGameMode).(tui.Mode)
	if !ok {
		return tui.FatalErrorCmd(fmt.Errorf("retrieving form mode: %w", tui.ErrInvalidTypeAssertion))
	}

	m.hasAnnouncedCompletion = true

	switch mode {
	case tui.ModeMarathon, tui.ModeSprint, tui.ModeUltra:
		in := tui.NewSingleInput(mode, level, username)
		return tui.SwitchModeCmd(mode, in)

	case tui.ModeMenu, tui.ModeLeaderboard:
		fallthrough
	default:
		return tui.FatalErrorCmd(fmt.Errorf("invalid mode for starting game %q", mode))
	}
}

func (m *MenuModel) View() string {
	output := lipgloss.JoinVertical(lipgloss.Center,
		titleStr+"\n",
		m.form.View(),
	)
	return lipgloss.Place(m.width, (m.height/10)*9, lipgloss.Center, lipgloss.Center, output)
}
