package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	// The number of tetriminos to display in the Next Queue.
	NextQueueLength int `toml:"next_queue_length"`

	// Whether a ghost piece will be displayed beneath the current tetrimino.
	GhostEnabled bool `toml:"ghost_enabled"`

	// TODO: What mode to use when locking down a tetrimino.
	LockDownMode string `toml:"lock_down_mode"`

	// The maximum level to reach before the game ends or the level stops increasing.
	MaxLevel int `toml:"max_level"`

	// Whether the game ends when the max level is reached.
	EndOnMaxLevel bool `toml:"end_on_max_level"`

	// The styling for the game in all modes
	Theme *Theme `toml:"theme"`

	// The keybindings for the game
	Keys *Keys `toml:"keys"`
}

func GetConfig(path string) (*Config, error) {
	c := Config{
		NextQueueLength: 5,
		GhostEnabled:    true,
		LockDownMode:    "Extended",
		MaxLevel:        15,
		EndOnMaxLevel:   false,

		Theme: defaultTheme(),
		Keys:  defaultKeys(),
	}

	_, err := toml.DecodeFile(path, &c)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &c, nil
		}
		return nil, err
	}

	err = c.validate()
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &c, nil
}

func (c *Config) validate() error {
	if c.NextQueueLength < 0 || c.NextQueueLength > 7 {
		return fmt.Errorf("NextQueueLength '%d' must be between 0 and 7", c.NextQueueLength)
	}
	if c.LockDownMode != "Extended" && c.LockDownMode != "Infinite" && c.LockDownMode != "Classic" {
		return fmt.Errorf("LockDownMode '%s' must be one of 'Extended', 'Infinite', or 'Classic'", c.LockDownMode)
	}
	return nil
}
