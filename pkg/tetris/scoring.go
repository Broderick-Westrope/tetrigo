package tetris

import (
	"fmt"
)

// Scoring is a scoring system for Tetris.
// It keeps track of the current level, total score, and lines cleared.
// It also has options to increase the level, end the game on max level, and end the game on max lines.
type Scoring struct {
	level         int
	maxLevel      int
	increaseLevel bool
	endOnMaxLevel bool

	lines         int
	maxLines      int
	endOnMaxLines bool

	total      int
	backToBack bool
}

// NewScoring creates a new scoring system.
func NewScoring(
	level, maxLevel int,
	increaseLevel, endOnMaxLevel bool,
	maxLines int,
	endOnMaxLines bool) (*Scoring, error) {
	s := &Scoring{
		level:         level,
		maxLevel:      maxLevel,
		increaseLevel: increaseLevel,
		endOnMaxLevel: endOnMaxLevel,

		maxLines:      maxLines,
		endOnMaxLines: endOnMaxLines,
	}
	return s, s.validate()
}

func (s *Scoring) validate() error {
	if s.level <= 0 {
		return fmt.Errorf("invalid level '%d'", s.level)
	}
	if s.maxLevel < 0 {
		return fmt.Errorf("invalid max level '%d'", s.maxLevel)
	}
	if s.maxLines < 0 {
		return fmt.Errorf("invalid max lines '%d'", s.maxLines)
	}
	if s.total < 0 {
		return fmt.Errorf("invalid total '%d'", s.total)
	}
	return nil
}

// Level returns the current level.
func (s *Scoring) Level() int {
	return s.level
}

// Total returns the total score.
func (s *Scoring) Total() int {
	return s.total
}

// Lines returns the total lines cleared.
func (s *Scoring) Lines() int {
	return s.lines
}

// AddSoftDrop adds points for a soft drop.
func (s *Scoring) AddSoftDrop(lines int) {
	s.total += lines
}

// AddHardDrop adds points for a hard drop.
func (s *Scoring) AddHardDrop(lines int) {
	s.total += lines * 2
}

// ProcessAction processes an action and updates the score, lines cleared, level, etc.
// The returned boolean indicates if the game should end.
func (s *Scoring) ProcessAction(a Action) (bool, error) {
	if a == Actions.None {
		return false, nil
	}

	points := float64(a.GetPoints())

	var err error
	var result bool
	backToBack := 0.0
	if result, err = a.EndsBackToBack(); result {
		s.backToBack = false
	} else if result, err = a.StartsBackToBack(); result {
		if s.backToBack {
			backToBack = points * 0.5
		}
		s.backToBack = true
	}
	if err != nil {
		return false, err
	}

	s.total += int(points+backToBack) * s.level
	s.lines += int((points + backToBack) / 100)

	// if max lines enabled, and max lines reached
	if s.maxLines > 0 && s.lines >= s.maxLines {
		s.lines = s.maxLines

		if s.endOnMaxLines {
			return true, nil
		}
	}

	// while increase level enabled, and the next level was reached
	for s.increaseLevel && s.lines >= s.level*5 {
		s.level++

		// if no max level, or max level not reached
		if s.maxLevel <= 0 || s.level < s.maxLevel {
			continue
		}

		// if max level reached
		s.level = s.maxLevel
		if s.endOnMaxLevel {
			return true, nil
		}
		break
	}

	return false, nil
}
