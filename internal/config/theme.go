package config

type Theme struct {
	Colours struct {
		TetriminoCells struct {
			I string
			O string
			T string
			S string
			Z string
			J string
			L string
		}
		EmptyCell string
		GhostCell string
	}
	Characters struct {
		Tetriminos string
		EmptyCell  string
		GhostCell  string
	}
}

func defaultTheme() *Theme {
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
