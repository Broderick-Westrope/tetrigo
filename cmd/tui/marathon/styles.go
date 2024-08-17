package marathon

import (
	"github.com/Broderick-Westrope/tetrigo/internal/config"
	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	Playfield       lipgloss.Style
	ColIndicator    lipgloss.Style
	TetriminoStyles map[byte]lipgloss.Style
	GhostCell       lipgloss.Style
	Hold            holdStyles
	Information     lipgloss.Style
	RowIndicator    lipgloss.Style
	Bag             lipgloss.Style
}

type holdStyles struct {
	View  lipgloss.Style
	Label lipgloss.Style
	Item  lipgloss.Style
}

func defaultStyles() *Styles {
	s := Styles{
		Playfield:    lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0),
		ColIndicator: lipgloss.NewStyle().Foreground(lipgloss.Color("#303040")),
		TetriminoStyles: map[byte]lipgloss.Style{
			'I': lipgloss.NewStyle().Foreground(lipgloss.Color("#64C4EB")),
			'O': lipgloss.NewStyle().Foreground(lipgloss.Color("#F1D448")),
			'T': lipgloss.NewStyle().Foreground(lipgloss.Color("#A15398")),
			'S': lipgloss.NewStyle().Foreground(lipgloss.Color("#64B452")),
			'Z': lipgloss.NewStyle().Foreground(lipgloss.Color("#DC3A35")),
			'J': lipgloss.NewStyle().Foreground(lipgloss.Color("#5C65A8")),
			'L': lipgloss.NewStyle().Foreground(lipgloss.Color("#E07F3A")),
		},
		GhostCell: lipgloss.NewStyle().Foreground(lipgloss.Color("white")),
		Hold: holdStyles{
			View:  lipgloss.NewStyle().Width(10).Height(5).Border(lipgloss.RoundedBorder(), true, false, true, true).Align(lipgloss.Center, lipgloss.Center),
			Label: lipgloss.NewStyle().Width(10).PaddingLeft(1).PaddingBottom(1),
			Item:  lipgloss.NewStyle().Width(10).Height(2).Align(lipgloss.Center, lipgloss.Center),
		},
		Information:  lipgloss.NewStyle().Width(13).Align(lipgloss.Left, lipgloss.Top),
		RowIndicator: lipgloss.NewStyle().Foreground(lipgloss.Color("#444049")).Align(lipgloss.Left).Padding(0, 1, 0),
		Bag:          lipgloss.NewStyle().PaddingTop(1),
	}
	return &s
}

func CreateStyles(theme *config.Theme) *Styles {
	s := Styles{
		Playfield:    lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0),
		ColIndicator: lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Colours.EmptyCell)),
		TetriminoStyles: map[byte]lipgloss.Style{
			'I': lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Colours.TetriminoCells.I)),
			'O': lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Colours.TetriminoCells.O)),
			'T': lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Colours.TetriminoCells.T)),
			'S': lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Colours.TetriminoCells.S)),
			'Z': lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Colours.TetriminoCells.Z)),
			'J': lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Colours.TetriminoCells.J)),
			'L': lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Colours.TetriminoCells.L)),
		},
		GhostCell: lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Colours.GhostCell)),
		Hold: holdStyles{
			View:  lipgloss.NewStyle().Width(10).Height(5).Border(lipgloss.RoundedBorder(), true, false, true, true).Align(lipgloss.Center, lipgloss.Center),
			Label: lipgloss.NewStyle().Width(10).PaddingLeft(1).PaddingBottom(1),
			Item:  lipgloss.NewStyle().Width(10).Height(2).Align(lipgloss.Center, lipgloss.Center),
		},
		Information:  lipgloss.NewStyle().Width(13).Align(lipgloss.Left, lipgloss.Top),
		RowIndicator: lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Characters.EmptyCell)).Align(lipgloss.Left).Padding(0, 1, 0),
		Bag:          lipgloss.NewStyle().PaddingTop(1),
	}
	return &s

}
