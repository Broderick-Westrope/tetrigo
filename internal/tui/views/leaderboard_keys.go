package views

import (
	"github.com/charmbracelet/bubbles/key"
)

type leaderboardKeyMap struct {
	Exit  key.Binding
	Help  key.Binding
	Left  key.Binding
	Right key.Binding
	Up    key.Binding
	Down  key.Binding
}

func defaultLeaderboardKeyMap() *leaderboardKeyMap {
	return &leaderboardKeyMap{
		Exit:  key.NewBinding(key.WithKeys("esc"), key.WithHelp("escape", "exit")),
		Help:  key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
		Left:  key.NewBinding(key.WithKeys("left"), key.WithHelp("left arrow", "move left")),
		Right: key.NewBinding(key.WithKeys("right"), key.WithHelp("right arrow", "move right")),
		Up:    key.NewBinding(key.WithKeys("up"), key.WithHelp("up arrow", "move up")),
		Down:  key.NewBinding(key.WithKeys("down"), key.WithHelp("down arrow", "move down")),
	}
}

func (k *leaderboardKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Exit,
		k.Help,
	}
}

func (k *leaderboardKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			k.Exit,
			k.Help,
		},
		{
			k.Left,
			k.Right,
		},
		{
			k.Up,
			k.Down,
		},
	}
}
