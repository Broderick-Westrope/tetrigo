package config

type Keys struct {
	ForceQuit              []string
	Exit                   []string
	Help                   []string
	Submit                 []string
	Up                     []string
	Down                   []string
	Left                   []string
	Right                  []string
	RotateCounterClockwise []string
	RotateClockwise        []string
}

func defaultKeys() *Keys {
	return &Keys{
		ForceQuit:              []string{"ctrl+c"},
		Exit:                   []string{"esc"},
		Help:                   []string{"?"},
		Submit:                 []string{" ", "enter"},
		Up:                     []string{"w"},
		Down:                   []string{"s"},
		Left:                   []string{"a"},
		Right:                  []string{"d"},
		RotateCounterClockwise: []string{"q"},
		RotateClockwise:        []string{"e"},
	}
}
