package views

import (
	"github.com/charmbracelet/bubbles/key"
)

type menuKeyMap struct {
	Exit  key.Binding
	Help  key.Binding
	Left  key.Binding
	Right key.Binding
	Up    key.Binding
	Down  key.Binding
	Start key.Binding
}

func defaultMenuKeyMap() *menuKeyMap {
	return &menuKeyMap{
		Exit:  key.NewBinding(key.WithKeys("esc"), key.WithHelp("escape", "exit")),
		Help:  key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
		Left:  key.NewBinding(key.WithKeys("left"), key.WithHelp("left arrow", "move left")),
		Right: key.NewBinding(key.WithKeys("right"), key.WithHelp("right arrow", "move right")),
		Up:    key.NewBinding(key.WithKeys("up"), key.WithHelp("up arrow", "move up")),
		Down:  key.NewBinding(key.WithKeys("down"), key.WithHelp("down arrow", "move down")),
		Start: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "start")),
	}
}

func (k *menuKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Start,
		k.Exit,
		k.Help,
	}
}

func (k *menuKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			k.Exit,
			k.Help,
			k.Left,
			k.Right,
		},
		{
			k.Up,
			k.Down,
			k.Start,
		},
	}
}
