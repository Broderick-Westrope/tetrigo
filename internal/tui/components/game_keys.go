package components

import (
	"github.com/Broderick-Westrope/charmutils"
	"github.com/charmbracelet/bubbles/key"

	"github.com/stuttgart-things/sthings-tetris/internal/config"
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

func ConstructGameKeyMap(keys *config.Keys) *GameKeyMap {
	return &GameKeyMap{
		ForceQuit:        charmutils.ConstructKeyBinding(keys.ForceQuit, "force quit"),
		Exit:             charmutils.ConstructKeyBinding(keys.Exit, "exit"),
		Help:             charmutils.ConstructKeyBinding(keys.Help, "help"),
		Left:             charmutils.ConstructKeyBinding(keys.Left, "move left"),
		Right:            charmutils.ConstructKeyBinding(keys.Right, "move right"),
		Clockwise:        charmutils.ConstructKeyBinding(keys.RotateClockwise, "rotate clockwise"),
		CounterClockwise: charmutils.ConstructKeyBinding(keys.RotateCounterClockwise, "rotate counter-clockwise"),
		SoftDrop:         charmutils.ConstructKeyBinding(keys.Down, "toggle soft drop"),
		HardDrop:         charmutils.ConstructKeyBinding(keys.Up, "hard drop"),
		Hold:             charmutils.ConstructKeyBinding(keys.Submit, "hold"),
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
