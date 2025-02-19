package components

import (
	"github.com/stuttgart-things/sthings-tetris/internal/config"

	"github.com/charmbracelet/lipgloss"
)

type GameStyles struct {
	Playfield           lipgloss.Style
	EmptyCell           lipgloss.Style
	TetriminoCellStyles map[byte]lipgloss.Style
	GhostCell           lipgloss.Style
	Hold                holdStyles
	Information         lipgloss.Style
	RowIndicator        lipgloss.Style
	Bag                 lipgloss.Style
	CellChar            cellCharacters
}

type holdStyles struct {
	View  lipgloss.Style
	Label lipgloss.Style
	Item  lipgloss.Style
}

type cellCharacters struct {
	Empty      string
	Ghost      string
	Tetriminos string
}

func CreateGameStyles(theme *config.Theme) *GameStyles {
	s := GameStyles{
		Playfield: lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0),
		EmptyCell: lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Colours.EmptyCell)),
		TetriminoCellStyles: map[byte]lipgloss.Style{
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
			View: lipgloss.NewStyle().Width(10).Height(5).
				Border(lipgloss.RoundedBorder(), true, false, true, true).
				Align(lipgloss.Center, lipgloss.Center),
			Label: lipgloss.NewStyle().Width(10).PaddingLeft(1).PaddingBottom(1),
			Item:  lipgloss.NewStyle().Width(10).Height(2).Align(lipgloss.Center, lipgloss.Center),
		},
		Information: lipgloss.NewStyle().Width(13).Align(lipgloss.Left, lipgloss.Top),
		RowIndicator: lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Characters.EmptyCell)).
			Align(lipgloss.Left).Padding(0, 1, 0),
		Bag: lipgloss.NewStyle().PaddingTop(1),
		CellChar: cellCharacters{
			Empty:      theme.Characters.EmptyCell,
			Ghost:      theme.Characters.GhostCell,
			Tetriminos: theme.Characters.Tetriminos,
		},
	}
	return &s
}
