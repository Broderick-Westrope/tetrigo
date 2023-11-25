package main

import (
	"fmt"
	"os"

	"github.com/Broderick-Westrope/tetrigo/internal/marathon"
	"github.com/Broderick-Westrope/tetrigo/internal/menu"
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

func main() {
	ctx := kong.Parse(&cli)
	switch ctx.Command() {
	case "menu":
		startTeaModel(menu.InitialModel(cli.Menu.Fullscreen))
	case "marathon":
		var fullscreen *marathon.FullscreenInfo
		if cli.Marathon.Fullscreen {
			fullscreen = &marathon.FullscreenInfo{
				Width:  0,
				Height: 0,
			}
		}

		startTeaModel(marathon.InitialModel(cli.Marathon.Level, fullscreen))
	default:
		panic(ctx.Command())
	}
}

func startTeaModel(m tea.Model) {
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
