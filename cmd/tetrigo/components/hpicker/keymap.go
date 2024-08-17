package hpicker

import (
	"github.com/charmbracelet/bubbles/key"
)

// KeyMap is the key bindings for different actions within the component.
type KeyMap struct {
	Prev key.Binding
	Next key.Binding
}

func defaultKeyMap() *KeyMap {
	return &KeyMap{
		Prev: key.NewBinding(
			key.WithKeys("left"),
			key.WithHelp("<-", "move left"),
		),
		Next: key.NewBinding(
			key.WithKeys("right"),
			key.WithHelp("->", "move right"),
		),
	}
}
