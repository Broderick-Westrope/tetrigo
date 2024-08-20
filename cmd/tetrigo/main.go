package main

import (
	"fmt"

	"github.com/Broderick-Westrope/tetrigo/internal/config"
	"github.com/Broderick-Westrope/tetrigo/internal/data"
	"github.com/Broderick-Westrope/tetrigo/internal/tui/common"
	"github.com/Broderick-Westrope/tetrigo/internal/tui/starter"
	"github.com/alecthomas/kong"

	tea "github.com/charmbracelet/bubbletea"
)

type CLI struct {
	GlobalVars

	Menu        MenuCmd        `cmd:"" help:"Start in the menu" default:"1"`
	Play        PlayCmd        `cmd:"" help:"Play a specific game mode"`
	Leaderboard LeaderboardCmd `cmd:"" help:"Start on the leaderboard"`
}

type GlobalVars struct {
	Config string `help:"Path to config file" default:"config.toml" type:"path"`
	DB     string `help:"Path to database file" default:"tetrigo.db"`
}

type MenuCmd struct{}

func (c *MenuCmd) Run(globals *GlobalVars) error {
	return launchStarter(globals, common.ModeMenu, common.NewMenuInput())
}

type PlayCmd struct {
	GameMode string `arg:"" help:"Game mode to play" default:"marathon"`
	Level    uint   `help:"Level to start at" short:"l" default:"1"`
	Name     string `help:"Name of the player" short:"n" default:"Anonymous"`
}

func (c *PlayCmd) Run(globals *GlobalVars) error {
	switch c.GameMode {
	case "marathon":
		return launchStarter(globals, common.ModeMarathon, common.NewMarathonInput(c.Level, c.Name))
	default:
		return fmt.Errorf("invalid game mode: %s", c.GameMode)
	}
}

type LeaderboardCmd struct {
	GameMode string `arg:"" help:"Game mode to display" default:"marathon"`
}

func (c *LeaderboardCmd) Run(globals *GlobalVars) error {
	return launchStarter(globals, common.ModeLeaderboard, common.NewLeaderboardInput(c.GameMode))
}

func main() {
	cli := CLI{}
	ctx := kong.Parse(&cli,
		kong.Name("tetrigo"),
		kong.Description("A tetris TUI written in Go"),
		kong.UsageOnError(),
	)
	// Call the Run() method of the selected parsed command.
	err := ctx.Run(&cli.GlobalVars)
	ctx.FatalIfErrorf(err)
}

func launchStarter(globals *GlobalVars, starterMode common.Mode, switchIn common.SwitchModeInput) error {
	db, err := data.NewDB(globals.DB)
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}

	cfg, err := config.GetConfig(globals.Config)
	if err != nil {
		return fmt.Errorf("error getting config: %w", err)
	}

	model, err := starter.NewModel(
		starter.NewInput(starterMode, switchIn, db, cfg),
	)
	if err != nil {
		return fmt.Errorf("error creating starter model: %w", err)
	}

	if _, err = tea.NewProgram(model, tea.WithAltScreen()).Run(); err != nil {
		return fmt.Errorf("error running tea program: %w", err)
	}

	return nil
}
