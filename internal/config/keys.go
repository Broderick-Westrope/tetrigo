package config

type Keys struct {
	ForceQuit              []string `toml:"force_quit"`
	Exit                   []string `toml:"exit"`
	Help                   []string `toml:"help"`
	Submit                 []string `toml:"submit"`
	Up                     []string `toml:"up"`
	Down                   []string `toml:"down"`
	Left                   []string `toml:"left"`
	Right                  []string `toml:"right"`
	RotateCounterClockwise []string `toml:"rotate_counter_clockwise"`
	RotateClockwise        []string `toml:"rotate_clockwise"`
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
