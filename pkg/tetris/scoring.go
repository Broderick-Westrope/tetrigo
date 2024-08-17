package tetris

type Scoring struct {
	level         uint
	maxLevel      uint
	endOnMaxLevel bool
	total         uint
	lines         uint
	backToBack    bool
}

func NewScoring(level, maxLevel uint, endOnMaxLevel bool) *Scoring {
	return &Scoring{
		level:         level,
		maxLevel:      maxLevel,
		endOnMaxLevel: endOnMaxLevel,
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

func (s *Scoring) ProcessAction(a Action) bool {
	if a == Actions.NONE {
		return s.isGameOver()
	}

	points := float64(a.GetPoints())

	backToBack := 0.0
	if a.EndsBackToBack() {
		s.backToBack = false
	} else if a.StartsBackToBack() {
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
			return s.isGameOver()
		}
	}

	return s.isGameOver()
}

func (s *Scoring) isGameOver() bool {
	return s.level >= s.maxLevel && s.endOnMaxLevel
}
