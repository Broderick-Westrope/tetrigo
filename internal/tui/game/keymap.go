package game

import (
	"github.com/Broderick-Westrope/tetrigo/internal/config"
	"github.com/charmbracelet/bubbles/key"
)

type GameKeyMap struct {
	ForceQuit        key.Binding
	Exit             key.Binding
	Help             key.Binding
	Left             key.Binding
	Right            key.Binding
	Clockwise        key.Binding
	CounterClockwise key.Binding
	SoftDrop         key.Binding
	HardDrop         key.Binding
	Hold             key.Binding
}

func constructKeyBinding(keys []string, desc string) key.Binding {
	buildHelpKeys := func(keys []string) string {
		helpKeys := ""
		for _, key := range keys {
			if key == " " {
				key = "space"
			}
			helpKeys += key + ", "
		}
		return helpKeys[:len(helpKeys)-2]
	}

	return key.NewBinding(key.WithKeys(keys...), key.WithHelp(buildHelpKeys(keys), desc))
}

func ConstructGameKeyMap(keys *config.Keys) *GameKeyMap {
	return &GameKeyMap{
		ForceQuit:        constructKeyBinding(keys.ForceQuit, "force quit"),
		Exit:             constructKeyBinding(keys.Exit, "exit"),
		Help:             constructKeyBinding(keys.Help, "help"),
		Left:             constructKeyBinding(keys.Left, "move left"),
		Right:            constructKeyBinding(keys.Right, "move right"),
		Clockwise:        constructKeyBinding(keys.RotateClockwise, "rotate clockwise"),
		CounterClockwise: constructKeyBinding(keys.RotateCounterClockwise, "rotate counter-clockwise"),
		SoftDrop:         constructKeyBinding(keys.Down, "toggle soft drop"),
		HardDrop:         constructKeyBinding(keys.Up, "hard drop"),
		Hold:             constructKeyBinding(keys.Submit, "hold"),
	}
}

func (k *GameKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Exit,
		k.Help,
	}
}

func (k *GameKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			k.Exit,
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
