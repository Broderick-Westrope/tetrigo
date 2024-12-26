package config

type Theme struct {
	Colours struct {
		TetriminoCells struct {
			I string `toml:"I"`
			O string `toml:"O"`
			T string `toml:"T"`
			S string `toml:"S"`
			Z string `toml:"Z"`
			J string `toml:"J"`
			L string `toml:"L"`
		} `toml:"tetrimino_cells"`
		EmptyCell string `toml:"empty_cell"`
		GhostCell string `toml:"ghost_cell"`
	} `toml:"colours"`
	Characters struct {
		Tetriminos string `toml:"tetriminos"`
		EmptyCell  string `toml:"empty_cell"`
		GhostCell  string `toml:"ghost_cell"`
	} `toml:"characters"`
}

func DefaultTheme() *Theme {
	theme := new(Theme)

	theme.Colours.TetriminoCells.I = "#64C4EB"
	theme.Colours.TetriminoCells.O = "#F1D448"
	theme.Colours.TetriminoCells.T = "#A15398"
	theme.Colours.TetriminoCells.S = "#64B452"
	theme.Colours.TetriminoCells.Z = "#DC3A35"
	theme.Colours.TetriminoCells.J = "#5C65A8"
	theme.Colours.TetriminoCells.L = "#E07F3A"
	theme.Colours.EmptyCell = "#303040"
	theme.Colours.GhostCell = "white"

	theme.Characters.Tetriminos = "██"
	theme.Characters.EmptyCell = "▕ "
	theme.Characters.GhostCell = "░░"

	return theme
}
