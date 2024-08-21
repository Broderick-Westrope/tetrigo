package tetris

type Scoring struct {
	level         uint
	maxLevel      uint
	increaseLevel bool
	endOnMaxLevel bool

	lines         uint
	maxLines      uint
	endOnMaxLines bool

	total      uint
	backToBack bool
}

func NewScoring(level, maxLevel uint, increaseLevel, endOnMaxLevel bool, maxLines uint, endOnMaxLines bool) *Scoring {
	return &Scoring{
		level:         level,
		maxLevel:      maxLevel,
		increaseLevel: increaseLevel,
		endOnMaxLevel: endOnMaxLevel,

		maxLines:      maxLines,
		endOnMaxLines: endOnMaxLines,
	}
}

func (s *Scoring) Level() uint {
	return s.level
}

func (s *Scoring) Total() uint {
	return s.total
}

func (s *Scoring) Lines() uint {
	return s.lines
}

func (s *Scoring) AddSoftDrop(lines uint) {
	s.total += lines
}

func (s *Scoring) AddHardDrop(lines uint) {
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

	s.total += uint(points+backToBack) * s.level

	// increase lines
	s.lines += uint((points + backToBack) / 100)

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
