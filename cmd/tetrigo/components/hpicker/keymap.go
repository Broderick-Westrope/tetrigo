package hpicker

import (
	"github.com/Broderick-Westrope/tetrigo/cmd/tetrigo/common"
	"github.com/Broderick-Westrope/tetrigo/internal/config"
	"github.com/charmbracelet/bubbles/key"
)

// KeyMap is the key bindings for different actions within the component.
type KeyMap struct {
	Prev key.Binding
	Next key.Binding
}

func ConstructKeyMap(keys *config.Keys) *KeyMap {
	return &KeyMap{
		Prev: common.ConstructKeyBinding(keys.Left, "move left"),
		Next: common.ConstructKeyBinding(keys.Right, "move right"),
	}
}
