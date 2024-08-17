package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	NextQueueLength int    // The number of tetriminos to display in the Next Queue. 0-7
	GhostEnabled    bool   // Whether a ghost piece will be displayed beneath the current tetrimino
	LockDownMode    string // TODO: What mode to use when locking down a tetrimino. Choices are Extended (default), Infinite, or Classic
	MaxLevel        uint   // The maximum level to reach before the game ends or the level stops increasing . 0+ (0 = no max level)
	EndOnMaxLevel   bool   // Whether the game ends when the max level is reached

	// The styling for the game in all modes
	Theme Theme
}

type Theme struct {
	Colours struct {
		TetriminoCells struct {
			I string
			O string
			T string
			S string
			Z string
			J string
			L string
		}
		EmptyCell string
		GhostCell string
	}
	Characters struct {
		Tetriminos string
		EmptyCell  string
		GhostCell  string
	}
}

func GetConfig(path string) (*Config, error) {
	var c Config

	c.NextQueueLength = 5
	c.GhostEnabled = true
	c.LockDownMode = "Extended"
	c.MaxLevel = 15
	c.EndOnMaxLevel = false

	c.Theme.Colours.TetriminoCells.I = "#64C4EB"
	c.Theme.Colours.TetriminoCells.O = "#F1D448"
	c.Theme.Colours.TetriminoCells.T = "#A15398"
	c.Theme.Colours.TetriminoCells.S = "#64B452"
	c.Theme.Colours.TetriminoCells.Z = "#DC3A35"
	c.Theme.Colours.TetriminoCells.J = "#5C65A8"
	c.Theme.Colours.TetriminoCells.L = "#E07F3A"
	c.Theme.Colours.EmptyCell = "#303040"
	c.Theme.Colours.GhostCell = "white"

	c.Theme.Characters.Tetriminos = "██"
	c.Theme.Characters.EmptyCell = "▕ "
	c.Theme.Characters.GhostCell = "░░"

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
