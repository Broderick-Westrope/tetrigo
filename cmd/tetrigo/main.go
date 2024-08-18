package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Broderick-Westrope/tetrigo/cmd/tetrigo/common"
	"github.com/Broderick-Westrope/tetrigo/cmd/tetrigo/starter"
	"github.com/Broderick-Westrope/tetrigo/internal/config"
	"github.com/Broderick-Westrope/tetrigo/internal/data"
	"github.com/alecthomas/kong"

	tea "github.com/charmbracelet/bubbletea"
)

type CLI struct {
	Config string `help:"Path to config file" default:"config.toml" type:"path"`
	DB     string `help:"Path to database file" default:"tetrigo.db"`

	Menu        MenuCmd        `cmd:"" help:"Start in the menu" default:"1"`
	Marathon    MarathonCmd    `cmd:"" help:"Play in marathon mode"`
	Leaderboard LeaderboardCmd `cmd:"" help:"Start on the leaderboard"`
}

type MenuCmd struct{}

type MarathonCmd struct {
	Level uint   `help:"Level to start at" short:"l" default:"1"`
	Name  string `help:"Name of the player" default:"Anonymous"`
}

type LeaderboardCmd struct {
	GameMode string `help:"Game mode to display" default:"marathon"`
}

var subcommandToStarterMode = map[string]common.Mode{
	"menu":        common.ModeMenu,
	"marathon":    common.ModeMarathon,
	"leaderboard": common.ModeLeaderboard,
}

func main() {
	cli := CLI{}
	ctx := kong.Parse(&cli,
		kong.Name("tetrigo"),
		kong.Description("A tetris TUI written in Go"),
		kong.UsageOnError(),
	)

	starterMode, ok := subcommandToStarterMode[ctx.Command()]
	if !ok {
		log.Printf("Invalid command: %s\n", ctx.Command())
	}

	db, err := data.NewDB(cli.DB)
	if err != nil {
		log.Printf("error opening database: %v", err)
		os.Exit(1)
	}

	cfg, err := config.GetConfig(cli.Config)
	if err != nil {
		log.Printf("error getting config: %v", err)
		os.Exit(1)
	}

	switchIn, err := cli.getSwitchModeInput(starterMode)
	if err != nil {
		log.Printf("error getting switch mode input: %v", err)
		os.Exit(1)
	}

	model, err := starter.NewModel(
		starter.NewInput(starterMode, switchIn, db, cfg),
	)
	if err != nil {
		log.Printf("error creating starter model: %v", err)
		os.Exit(1)
	}

	if _, err = tea.NewProgram(model, tea.WithAltScreen()).Run(); err != nil {
		log.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func (cli CLI) getSwitchModeInput(starterMode common.Mode) (common.SwitchModeInput, error) {
	switch starterMode {
	case common.ModeMenu:
		return common.NewMenuInput(), nil
	case common.ModeMarathon:
		return common.NewMarathonInput(cli.Marathon.Level, cli.Marathon.Name), nil
	case common.ModeLeaderboard:
		return common.NewLeaderboardInput(cli.Leaderboard.GameMode), nil
	default:
		return nil, fmt.Errorf("invalid starter mode: %v", starterMode)
	}
}
