package views

import (
	"github.com/Broderick-Westrope/tetrigo/internal/config"
	"github.com/Broderick-Westrope/tetrigo/internal/tui"
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

func ConstructGameKeyMap(keys *config.Keys) *GameKeyMap {
	return &GameKeyMap{
		ForceQuit:        tui.ConstructKeyBinding(keys.ForceQuit, "force quit"),
		Exit:             tui.ConstructKeyBinding(keys.Exit, "exit"),
		Help:             tui.ConstructKeyBinding(keys.Help, "help"),
		Left:             tui.ConstructKeyBinding(keys.Left, "move left"),
		Right:            tui.ConstructKeyBinding(keys.Right, "move right"),
		Clockwise:        tui.ConstructKeyBinding(keys.RotateClockwise, "rotate clockwise"),
		CounterClockwise: tui.ConstructKeyBinding(keys.RotateCounterClockwise, "rotate counter-clockwise"),
		SoftDrop:         tui.ConstructKeyBinding(keys.Down, "toggle soft drop"),
		HardDrop:         tui.ConstructKeyBinding(keys.Up, "hard drop"),
		Hold:             tui.ConstructKeyBinding(keys.Submit, "hold"),
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
