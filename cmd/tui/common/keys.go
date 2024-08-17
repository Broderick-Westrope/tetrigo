package common

import "github.com/charmbracelet/bubbles/key"

type Keys struct {
	ForceQuit   []string
	Exit        []string
	Help        []string
	Submit      []string
	Up          []string
	Down        []string
	Left        []string
	Right       []string
	LeftRotate  []string
	RightRotate []string
}

func DefaultKeys() *Keys {
	return &Keys{
		ForceQuit:   []string{"ctrl+c"},
		Exit:        []string{"esc"},
		Help:        []string{"?"},
		Submit:      []string{" ", "enter"},
		Up:          []string{"i"},
		Down:        []string{"k"},
		Left:        []string{"j"},
		Right:       []string{"l"},
		LeftRotate:  []string{"u"},
		RightRotate: []string{"o"},
	}
}

func ConstructKeyBinding(keys []string, desc string) key.Binding {
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
