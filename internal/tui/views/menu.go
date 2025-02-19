package views

import (
	"errors"
	"fmt"

	"github.com/Broderick-Westrope/charmutils"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	"github.com/stuttgart-things/sthings-tetris/internal/tui"
)

const (
	titleStr = `
    ______________________  ______________
   /_  __/ ____/_  __/ __ \/  _/ ____/ __ \
    / / / __/   / / / /_/ // // / __/ / / /
   / / / /___  / / / _, _// // /_/ / /_/ /
  /_/ /_____/ /_/ /_/ |_/___/\____/\____/
`
)

var _ tea.Model = &MenuModel{}

type MenuModel struct {
	form                   *huh.Form
	hasAnnouncedCompletion bool
	keys                   *menuKeyMap
	formData               *MenuFormData

	width  int
	height int
}

type MenuFormData struct {
	Username string
	GameMode tui.Mode
	Level    int
}

func NewMenuModel(_ *tui.MenuInput) *MenuModel {
	formData := new(MenuFormData)
	keys := defaultMenuKeyMap()

	return &MenuModel{
		formData: formData,
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewInput().Value(&formData.Username).
					Title("Username:").CharLimit(100).
					Validate(func(s string) error {
						if len(s) == 0 {
							return errors.New("empty username not allowed")
						}
						return nil
					}),
				huh.NewSelect[tui.Mode]().Value(&formData.GameMode).
					Title("Game Mode:").
					Options(
						huh.NewOption("Marathon", tui.ModeMarathon),
						huh.NewOption("Sprint (40 Lines)", tui.ModeSprint),
						huh.NewOption("Ultra (Time Trial)", tui.ModeUltra),
					),
				huh.NewSelect[int]().Value(&formData.Level).
					Title("Starting Level:").
					Options(charmutils.HuhIntRangeOptions(1, 15)...),
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
		if key.Matches(msg, m.keys.Exit) {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		formWidth := msg.Width / 2
		formWidth = min(formWidth, lipgloss.Width(titleStr))
		m.form = m.form.WithWidth(formWidth)
		return m, nil
	}

	var cmds []tea.Cmd
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmds = append(cmds, cmd)
	}

	if m.form.State == huh.StateCompleted && !m.hasAnnouncedCompletion {
		cmds = append(cmds, m.announceCompletion())
	}

	return m, tea.Batch(cmds...)
}

func (m *MenuModel) announceCompletion() tea.Cmd {
	m.hasAnnouncedCompletion = true

	switch m.formData.GameMode {
	case tui.ModeMarathon, tui.ModeSprint, tui.ModeUltra:
		in := tui.NewSingleInput(m.formData.GameMode, m.formData.Level, m.formData.Username)
		return tui.SwitchModeCmd(m.formData.GameMode, in)

	case tui.ModeMenu, tui.ModeLeaderboard:
		fallthrough
	default:
		return tui.FatalErrorCmd(fmt.Errorf("invalid mode for starting game %q", m.formData.GameMode))
	}
}

func (m *MenuModel) View() string {
	output := lipgloss.JoinVertical(lipgloss.Center,
		titleStr+"\n",
		m.form.View(),
	)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, output)
}
