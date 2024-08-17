package common

import "github.com/charmbracelet/bubbles/key"

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
