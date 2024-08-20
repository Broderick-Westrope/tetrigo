package tetris

type Scoring struct {
	level         uint
	maxLevel      uint
	endOnMaxLevel bool

	lines         uint
	maxLines      uint
	endOnMaxLines bool

	total      uint
	backToBack bool
}

func NewScoring(level, maxLevel uint, endOnMaxLevel bool, maxLines uint, endOnMaxLines bool) *Scoring {
	return &Scoring{
		level:         level,
		maxLevel:      maxLevel,
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
		return s.isGameOver(), nil
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
		return s.isGameOver(), err
	}

	s.total += uint(points+backToBack) * s.level
	s.lines += uint((points + backToBack) / 100)

	if s.maxLines > 0 && s.lines > s.maxLines {
		s.lines = s.maxLines
	}

	for s.lines >= s.level*5 {
		s.level++
		if s.maxLevel > 0 && s.level >= s.maxLevel {
			s.level = s.maxLevel
			return s.isGameOver(), nil
		}
	}

	return s.isGameOver(), nil
}

func (s *Scoring) isGameOver() bool {
	return s.level >= s.maxLevel && s.endOnMaxLevel ||
		s.lines >= s.maxLines && s.endOnMaxLines
}
