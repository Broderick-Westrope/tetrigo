package main

import (
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
	Config string `help:"Path to config file. Empty value will use XDG data directory." default:""`
	DB     string `help:"Path to database file. Empty value will use XDG data directory." default:""`
}

func main() {
	cli := CLI{}
	ctx := kong.Parse(&cli,
		kong.Name("sthings-tetris"),
		kong.Description("A tetris TUI written in Go"),
		kong.UsageOnError(),
	)

	if err := handleDefaultGlobals(&cli.GlobalVars); err != nil {
		ctx.FatalIfErrorf(err)
	}

	// Call the Run() method of the selected parsed command.
	err := ctx.Run(&cli.GlobalVars)
	ctx.FatalIfErrorf(err)
}

func handleDefaultGlobals(g *GlobalVars) error {
	if g.Config == "" {
		var err error
		g.Config, err = xdg.ConfigFile("./tetrigo/config.toml")
		if err != nil {
			return err
		}
	}
	if g.DB == "" {
		var err error
		g.DB, err = xdg.DataFile("./tetrigo/tetrigo.db")
		if err != nil {
			return err
		}
	}
	return nil
}
