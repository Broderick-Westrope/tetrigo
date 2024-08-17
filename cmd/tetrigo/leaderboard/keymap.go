package leaderboard

import (
	"github.com/Broderick-Westrope/tetrigo/cmd/tetrigo/common"
	"github.com/charmbracelet/bubbles/key"
)

type keyMap struct {
	ForceQuit key.Binding
	Exit      key.Binding
	Help      key.Binding
	Left      key.Binding
	Right     key.Binding
	Up        key.Binding
	Down      key.Binding
}

func constructKeyMap(keys *common.Keys) *keyMap {
	return &keyMap{
		ForceQuit: common.ConstructKeyBinding(keys.ForceQuit, "force quit"),
		Exit:      common.ConstructKeyBinding(keys.Exit, "exit"),
		Help:      common.ConstructKeyBinding(keys.Help, "help"),
		Up:        common.ConstructKeyBinding(keys.Up, "move up"),
		Down:      common.ConstructKeyBinding(keys.Down, "move down"),
		Left:      common.ConstructKeyBinding(keys.Left, "move left"),
		Right:     common.ConstructKeyBinding(keys.Right, "move right"),
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
