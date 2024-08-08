package tetris

type Scoring struct {
	level      uint
	total      uint
	lines      uint
	backToBack bool
}

// TODO: make action exported

// Actions that score points. Defined in chapter 8 of the 2009 Guideline
type action int8

const (
	actionNone = iota
	actionSingle
	actionDouble
	actionTriple
	actionTetris
	actionMiniTSpin
	actionMiniTSpinSingle
	actionTSpin
	actionTSpinSingle
	actionTSpinDouble
	actionTSpinTriple
)

func NewScoring(level uint) *Scoring {
	return &Scoring{
		level: level,
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

func (s *Scoring) ProcessAction(a action, maxLevel uint) {
	if a == actionNone {
		return
	}

	points := 0.0
	switch a {
	case actionSingle:
		points = 100
	case actionDouble:
		points = 300
	case actionTriple:
		points = 500
	case actionTetris:
		points = 800
	case actionMiniTSpin:
		points = 100
	case actionMiniTSpinSingle:
		points = 200
	case actionTSpin:
		points = 400
	case actionTSpinSingle:
		points = 800
	case actionTSpinDouble:
		points = 1200
	case actionTSpinTriple:
		points = 1600
	}

	backToBack := 0.0
	switch a {
	case actionSingle:
		s.backToBack = false
	case actionDouble:
		s.backToBack = false
	case actionTriple:
		s.backToBack = false
	case actionTetris:
		if s.backToBack {
			backToBack = points * 0.5
		}
		s.backToBack = true
	case actionMiniTSpinSingle:
		if s.backToBack {
			backToBack = points * 0.5
		}
		s.backToBack = true
	case actionTSpinSingle:
		if s.backToBack {
			backToBack = points * 0.5
		}
		s.backToBack = true
	case actionTSpinDouble:
		if s.backToBack {
			backToBack = points * 0.5
		}
		s.backToBack = true
	case actionTSpinTriple:
		if s.backToBack {
			backToBack = points * 0.5
		}
		s.backToBack = true
	}

	s.total += uint(points+backToBack) * s.level
	s.lines += uint((points + backToBack) / 100)

	for s.lines >= s.level*5 {
		s.level++
		if maxLevel > 0 && s.level >= maxLevel {
			s.level = maxLevel
			return
		}
	}
}
