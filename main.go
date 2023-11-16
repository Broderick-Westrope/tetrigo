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
	Menu     struct{} `cmd:"" help:"Play the game" default:"1"`
	Marathon struct {
		Level uint `help:"Level to start at" short:"l" default:"1"`
	} `cmd:"" help:"Play marathon mode"`
}

func main() {
	ctx := kong.Parse(&cli)
	switch ctx.Command() {
	case "menu":
		startTeaModel(menu.InitialModel())
	case "marathon":
		startTeaModel(marathon.InitialModel(cli.Marathon.Level))
	default:
		panic(ctx.Command())
	}
}

func startTeaModel(m tea.Model) {
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
