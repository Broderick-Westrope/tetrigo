package main

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Quit             key.Binding
	Help             key.Binding
	Left             key.Binding
	Right            key.Binding
	Clockwise        key.Binding
	CounterClockwise key.Binding
	SoftDrop         key.Binding
	HardDrop         key.Binding
}

func DefaultKeyMap() *KeyMap {
	return &KeyMap{
		Quit:             key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
		Help:             key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
		Left:             key.NewBinding(key.WithKeys("left", "j", "a"), key.WithHelp("a, j, left", "move left")),
		Right:            key.NewBinding(key.WithKeys("right", "l", "d"), key.WithHelp("d, l, right", "move right")),
		Clockwise:        key.NewBinding(key.WithKeys(">", "e", "o"), key.WithHelp("e, o, >", "rotate clockwise")),
		CounterClockwise: key.NewBinding(key.WithKeys("<", "q", "u"), key.WithHelp("q, u, <", "rotate counter-clockwise")),
		SoftDrop:         key.NewBinding(key.WithKeys("down", "s", "k"), key.WithHelp("s, k, down", "soft drop")),
		HardDrop:         key.NewBinding(key.WithKeys("up", "w", "i"), key.WithHelp("w, i, up", "hard drop")),
	}
}

func (k *KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Quit,
		k.Help,
	}
}

func (k *KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			k.Quit,
			k.Help,
		},
		{
			k.Left,
			k.Right,
		},
		{
			k.Clockwise,
			k.CounterClockwise,
		},
		{
			k.SoftDrop,
			k.HardDrop,
		},
	}
}
