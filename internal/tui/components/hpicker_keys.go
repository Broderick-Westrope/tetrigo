package components

import (
	"github.com/charmbracelet/bubbles/key"
)

// hPickerKeyMap is the key bindings for different actions within the component.
type hPickerKeyMap struct {
	Prev key.Binding
	Next key.Binding
}

func defaultHPickerKeyMap() *hPickerKeyMap {
	return &hPickerKeyMap{
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
