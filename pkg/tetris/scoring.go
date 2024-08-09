package tetris

type Scoring struct {
	level      uint
	maxLevel   uint
	total      uint
	lines      uint
	backToBack bool
}

func NewScoring(level, maxLevel uint) *Scoring {
	return &Scoring{
		level:    level,
		maxLevel: maxLevel,
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

func (s *Scoring) ProcessAction(a Action) {
	if a == Actions.NONE {
		return
	}

	points := float64(a.GetPoints())

	backToBack := 0.0
	switch a {
	case Actions.SINGLE:
		s.backToBack = false
	case Actions.DOUBLE:
		s.backToBack = false
	case Actions.TRIPLE:
		s.backToBack = false
	case Actions.TETRIS:
		if s.backToBack {
			backToBack = points * 0.5
		}
		s.backToBack = true
	case Actions.MINI_T_SPIN_SINGLE:
		if s.backToBack {
			backToBack = points * 0.5
		}
		s.backToBack = true
	case Actions.T_SPIN_SINGLE:
		if s.backToBack {
			backToBack = points * 0.5
		}
		s.backToBack = true
	case Actions.T_SPIN_DOUBLE:
		if s.backToBack {
			backToBack = points * 0.5
		}
		s.backToBack = true
	case Actions.T_SPIN_TRIPLE:
		if s.backToBack {
			backToBack = points * 0.5
		}
		s.backToBack = true
	}

	s.total += uint(points+backToBack) * s.level
	s.lines += uint((points + backToBack) / 100)

	for s.lines >= s.level*5 {
		s.level++
		if s.maxLevel > 0 && s.level >= s.maxLevel {
			s.level = s.maxLevel
			return
		}
	}
}
