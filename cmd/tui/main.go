package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Broderick-Westrope/tetrigo/internal/tui"
	"github.com/Broderick-Westrope/tetrigo/internal/tui/starter"
	"github.com/alecthomas/kong"
	tea "github.com/charmbracelet/bubbletea"
)

var cli struct {
	Menu struct {
		Fullscreen bool `help:"Start in fullscreen?" default:"true"`
	} `cmd:"" help:"Play the game" default:"1"`
	Marathon struct {
		Level      uint `help:"Level to start at" short:"l" default:"1"`
		Fullscreen bool `help:"Start in fullscreen?" default:"true"`
	} `cmd:"" help:"Play marathon mode"`
}

var subcommandToStarterMode = map[string]tui.Mode{
	"menu":     tui.MODE_MENU,
	"marathon": tui.MODE_MARATHON,
}

func main() {
	ctx := kong.Parse(&cli)

	starterMode, ok := subcommandToStarterMode[ctx.Command()]
	if !ok {
		fmt.Printf("Invalid command: %s\n", ctx.Command())
	}

	model, err := starter.NewModel(
		starter.NewInput(starterMode, cli.Menu.Fullscreen, cli.Marathon.Level),
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
