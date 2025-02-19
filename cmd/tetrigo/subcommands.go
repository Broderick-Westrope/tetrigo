package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/stuttgart-things/sthings-tetris/internal/config"
	"github.com/stuttgart-things/sthings-tetris/internal/data"
	"github.com/stuttgart-things/sthings-tetris/internal/tui"
	"github.com/stuttgart-things/sthings-tetris/internal/tui/starter"
)

type MenuCmd struct{}

func (c *MenuCmd) Run(globals *GlobalVars) error {
	return launchStarter(globals, tui.ModeMenu, tui.NewMenuInput())
}

type PlayCmd struct {
	GameMode string `arg:"" help:"Game mode to play" default:"marathon"`
	Level    int    `help:"Level to start at" short:"l" default:"1"`
	Name     string `help:"Name of the player" short:"n" default:"Anonymous"`
}

func (c *PlayCmd) Run(globals *GlobalVars) error {
	singlePlayerModes := map[string]tui.Mode{
		"marathon": tui.ModeMarathon,
		"sprint":   tui.ModeSprint,
		"ultra":    tui.ModeUltra,
	}

	mode, ok := singlePlayerModes[c.GameMode]
	if !ok {
		return fmt.Errorf("invalid game mode: %s", c.GameMode)
	}

	return launchStarter(globals, mode, tui.NewSingleInput(tui.ModeMarathon, c.Level, c.Name))
}

type LeaderboardCmd struct {
	GameMode string `arg:"" help:"Game mode to display" default:"marathon"`
}

func (c *LeaderboardCmd) Run(globals *GlobalVars) error {
	return launchStarter(globals, tui.ModeLeaderboard, tui.NewLeaderboardInput(c.GameMode))
}

func launchStarter(globals *GlobalVars, starterMode tui.Mode, switchIn tui.SwitchModeInput) error {
	db, err := data.NewDB(globals.DB)
	if err != nil {
		return fmt.Errorf("opening database: %w", err)
	}

	cfg, err := config.GetConfig(globals.Config)
	if err != nil {
		return fmt.Errorf("getting config: %w", err)
	}

	model, err := starter.NewModel(
		starter.NewInput(starterMode, switchIn, db, cfg),
	)
	if err != nil {
		return fmt.Errorf("creating starter model: %w", err)
	}

	exitModel, err := tea.NewProgram(model, tea.WithAltScreen()).Run()
	if err != nil {
		return fmt.Errorf("failed to run program: %w", err)
	}

	typedExitModel, ok := exitModel.(*starter.Model)
	if !ok {
		return fmt.Errorf("failed to assert exit model type: %w", err)
	}

	if err = typedExitModel.ExitError; err != nil {
		return fmt.Errorf("starter model exited with an error: %w", err)
	}

	return nil
}
