package main

import (
	"log"

	"github.com/adrg/xdg"
	"github.com/alecthomas/kong"
)

type CLI struct {
	GlobalVars

	Menu        MenuCmd        `cmd:"" help:"Start in the menu" default:"1"`
	Play        PlayCmd        `cmd:"" help:"Play a specific game mode"`
	Leaderboard LeaderboardCmd `cmd:"" help:"Start on the leaderboard"`
}

type GlobalVars struct {
	Config string `help:"Path to config file" default:"config.toml" type:"path"`
	DB     string `help:"Path to database file. Empty value will use XDG data directory." default:""`
}

func main() {
	cli := CLI{}
	ctx := kong.Parse(&cli,
		kong.Name("tetrigo"),
		kong.Description("A tetris TUI written in Go"),
		kong.UsageOnError(),
	)

	if cli.GlobalVars.DB == "" {
		var err error
		cli.GlobalVars.DB, err = xdg.DataFile("./tetrigo/tetrigo.db")
		if err != nil {
			log.Fatal(err)
		}
	}

	// Call the Run() method of the selected parsed command.
	err := ctx.Run(&cli.GlobalVars)
	ctx.FatalIfErrorf(err)
}
