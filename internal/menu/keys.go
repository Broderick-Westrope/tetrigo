package menu

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Quit  key.Binding
	Help  key.Binding
	Left  key.Binding
	Right key.Binding
	Up    key.Binding
	Down  key.Binding
	Start key.Binding
}

func defaultKeyMap() *keyMap {
	return &keyMap{
		Quit:  key.NewBinding(key.WithKeys("esc", "ctrl+c"), key.WithHelp("esc", "quit")),
		Help:  key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
		Left:  key.NewBinding(key.WithKeys("j", "a", "left"), key.WithHelp("a, j, left", "move left")),
		Right: key.NewBinding(key.WithKeys("l", "d", "right"), key.WithHelp("d, l, right", "move right")),
		Up:    key.NewBinding(key.WithKeys("i", "w", "up"), key.WithHelp("w, i, right", "move up")),
		Down:  key.NewBinding(key.WithKeys("k", "s", "down"), key.WithHelp("s, k, down", "move down")),
		Start: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "start game")),
	}
}

func (k *keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Quit,
		k.Help,
	}
}

func (k *keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			k.Quit,
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
