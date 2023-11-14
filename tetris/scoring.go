package tetris

type Scoring struct {
	level      uint
	total      uint
	backToBack bool
}

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
		level:      level,
		total:      0,
		backToBack: false,
	}
}

func (s *Scoring) Level() uint {
	return s.level
}

func (s *Scoring) Total() uint {
	return s.total
}

func (s *Scoring) ProcessAction(a action) {
	if a == actionNone {
		return
	}

	switch a {
	case actionSingle:
		s.total += 100 * s.level
		s.backToBack = false
	case actionDouble:
		s.total += 300 * s.level
		s.backToBack = false
	case actionTriple:
		s.total += 500 * s.level
		s.backToBack = false
	case actionTetris:
		s.total += 800 * s.level
		if s.backToBack {
			s.total += 400 * s.level
		}
		s.backToBack = true
	case actionMiniTSpin:
		s.total += 100 * s.level
	case actionMiniTSpinSingle:
		s.total += 200 * s.level
		if s.backToBack {
			s.total += 100 * s.level
		}
		s.backToBack = true
	case actionTSpin:
		s.total += 400 * s.level
	case actionTSpinSingle:
		s.total += 800 * s.level
		if s.backToBack {
			s.total += 400 * s.level
		}
		s.backToBack = true
	case actionTSpinDouble:
		s.total += 1200 * s.level
		if s.backToBack {
			s.total += 600 * s.level
		}
		s.backToBack = true
	case actionTSpinTriple:
		s.total += 1600 * s.level
		if s.backToBack {
			s.total += 800 * s.level
		}
		s.backToBack = true
	}
}
