package menu

import (
	"github.com/Broderick-Westrope/tetrigo/cmd/tui/common"
	"github.com/charmbracelet/bubbles/key"
)

type keyMap struct {
	Exit  key.Binding
	Help  key.Binding
	Left  key.Binding
	Right key.Binding
	Up    key.Binding
	Down  key.Binding
	Start key.Binding
}

func constructKeyMap(keys *common.Keys) *keyMap {
	exitKeys := append(keys.Exit, keys.ForceQuit...)
	return &keyMap{
		Exit:  common.ConstructKeyBinding(exitKeys, "exit"),
		Help:  common.ConstructKeyBinding(keys.Help, "help"),
		Left:  common.ConstructKeyBinding(keys.Left, "move left"),
		Right: common.ConstructKeyBinding(keys.Right, "move right"),
		Up:    common.ConstructKeyBinding(keys.Down, "toggle soft drop"),
		Down:  common.ConstructKeyBinding(keys.Up, "hard drop"),
		Start: common.ConstructKeyBinding(keys.Submit, "hold"),
	}
}

func (k *keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Exit,
		k.Help,
	}
}

func (k *keyMap) FullHelp() [][]key.Binding {
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
