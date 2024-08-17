package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Broderick-Westrope/tetrigo/cmd/tui/common"
	"github.com/Broderick-Westrope/tetrigo/cmd/tui/starter"
	"github.com/Broderick-Westrope/tetrigo/internal/config"
	"github.com/Broderick-Westrope/tetrigo/internal/data"
	"github.com/alecthomas/kong"
	tea "github.com/charmbracelet/bubbletea"
)

var cli struct {
	Menu struct {
	} `cmd:"" help:"Play the game" default:"1"`

	Marathon struct {
		Level uint   `help:"Level to start at" short:"l" default:"1"`
		Name  string `help:"Name of the player" default:"Anonymous"`
	} `cmd:"" help:"Play marathon mode"`

	Leaderboard struct {
		GameMode string `help:"Game mode to display" default:"marathon"`
	} `cmd:"" help:"View the leaderboard for a given mode"`
}

var subcommandToStarterMode = map[string]common.Mode{
	"menu":        common.MODE_MENU,
	"marathon":    common.MODE_MARATHON,
	"leaderboard": common.MODE_LEADERBOARD,
}

func main() {
	ctx := kong.Parse(&cli)

	starterMode, ok := subcommandToStarterMode[ctx.Command()]
	if !ok {
		fmt.Printf("Invalid command: %s\n", ctx.Command())
	}

	db, err := data.NewDB("data.db")
	if err != nil {
		log.Printf("error opening database: %v", err)
		os.Exit(1)
	}

	cfg, err := config.GetConfig("config.toml")
	if err != nil {
		log.Printf("error getting config: %v", err)
		os.Exit(1)
	}

	switchIn, err := getSwitchModeInput(starterMode)
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
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func getSwitchModeInput(starterMode common.Mode) (common.SwitchModeInput, error) {
	switch starterMode {
	case common.MODE_MENU:
		return common.NewMenuInput(), nil
	case common.MODE_MARATHON:
		return common.NewMarathonInput(cli.Marathon.Level, cli.Marathon.Name), nil
	case common.MODE_LEADERBOARD:
		return common.NewLeaderboardInput(cli.Leaderboard.GameMode), nil
	default:
		return nil, fmt.Errorf("invalid starter mode: %v", starterMode)
	}
}
