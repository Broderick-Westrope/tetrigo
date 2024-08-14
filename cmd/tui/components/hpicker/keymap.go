package hpicker

import "github.com/charmbracelet/bubbles/key"

// KeyMap is the key bindings for different actions within the component.
type KeyMap struct {
	Prev key.Binding
	Next key.Binding
}

// DefaultKeyMap is the default set of key bindings for navigating and acting  upon the component.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Prev: key.NewBinding(key.WithKeys("left", "h")),
		Next: key.NewBinding(key.WithKeys("right", "l")),
	}
}
