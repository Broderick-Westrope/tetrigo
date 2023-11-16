package marathon

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Quit             key.Binding
	Help             key.Binding
	Left             key.Binding
	Right            key.Binding
	Clockwise        key.Binding
	CounterClockwise key.Binding
	SoftDrop         key.Binding
	HardDrop         key.Binding
	Hold             key.Binding
}

func defaultKeyMap() *keyMap {
	return &keyMap{
		Quit:             key.NewBinding(key.WithKeys("esc", "ctrl+c"), key.WithHelp("esc", "quit")),
		Help:             key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
		Left:             key.NewBinding(key.WithKeys("j", "a"), key.WithHelp("a, j", "move left")),
		Right:            key.NewBinding(key.WithKeys("l", "d"), key.WithHelp("d, l", "move right")),
		Clockwise:        key.NewBinding(key.WithKeys("e", "o"), key.WithHelp("e, o", "rotate clockwise")),
		CounterClockwise: key.NewBinding(key.WithKeys("q", "u"), key.WithHelp("q, u", "rotate counter-clockwise")),
		SoftDrop:         key.NewBinding(key.WithKeys("s", "k"), key.WithHelp("s, k", "toggle soft drop")),
		HardDrop:         key.NewBinding(key.WithKeys("w", "i"), key.WithHelp("w, i", "hard drop")),
		Hold:             key.NewBinding(key.WithKeys(" "), key.WithHelp("space", "hold")),
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
		},
		{
			k.Right,
			k.Clockwise,
			k.CounterClockwise,
		},
		{
			k.SoftDrop,
			k.HardDrop,
			k.Hold,
		},
	}
}
