package tetris

import (
	"fmt"
)

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

func (s *Scoring) Level() int {
	return s.level
}

func (s *Scoring) Total() int {
	return s.total
}

func (s *Scoring) Lines() int {
	return s.lines
}

func (s *Scoring) AddSoftDrop(lines int) {
	s.total += lines
}

func (s *Scoring) AddHardDrop(lines int) {
	s.total += lines * 2
}

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

	// increase lines
	s.lines += int((points + backToBack) / 100)

	// if max lines enabled, and max lines reached
	if s.maxLines > 0 && s.lines >= s.maxLines {
		s.lines = s.maxLines

		if s.endOnMaxLines {
			return true, nil
		}
	}

	// increase level
	for s.increaseLevel && s.lines >= s.level*5 {
		s.level++

		// if no max level, or max level not reached
		if s.maxLevel <= 0 || s.level < s.maxLevel {
			continue
		}

		// if max level reached
		s.level = s.maxLevel
		return s.endOnMaxLevel, nil
	}

	return false, nil
}
