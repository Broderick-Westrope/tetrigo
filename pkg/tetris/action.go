package tetris

import "fmt"

// An Action performed by the user that corresponds to points.
type Action struct {
	action
}

type action int

const (
	actionUnknown action = iota
	actionNone
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

var (
	actionToStrMap = map[action]string{
		actionNone:            "NONE",
		actionSingle:          "SINGLE",
		actionDouble:          "DOUBLE",
		actionTriple:          "TRIPLE",
		actionTetris:          "TETRIS",
		actionMiniTSpin:       "MINI_T_SPIN",
		actionMiniTSpinSingle: "MINI_T_SPIN_SINGLE",
		actionTSpin:           "T_SPIN",
		actionTSpinSingle:     "T_SPIN_SINGLE",
		actionTSpinDouble:     "T_SPIN_DOUBLE",
		actionTSpinTriple:     "T_SPIN_TRIPLE",
	}

	strToActionMap = map[string]action{
		"NONE":               actionNone,
		"SINGLE":             actionSingle,
		"DOUBLE":             actionDouble,
		"TRIPLE":             actionTriple,
		"TETRIS":             actionTetris,
		"MINI_T_SPIN":        actionMiniTSpin,
		"MINI_T_SPIN_SINGLE": actionMiniTSpinSingle,
		"T_SPIN":             actionTSpin,
		"T_SPIN_SINGLE":      actionTSpinSingle,
		"T_SPIN_DOUBLE":      actionTSpinDouble,
		"T_SPIN_TRIPLE":      actionTSpinTriple,
	}

	actionToPointsMap = map[action]int{
		actionNone:            0,
		actionSingle:          100,
		actionDouble:          300,
		actionTriple:          500,
		actionTetris:          800,
		actionMiniTSpin:       100,
		actionMiniTSpinSingle: 200,
		actionTSpin:           400,
		actionTSpinSingle:     800,
		actionTSpinDouble:     1200,
		actionTSpinTriple:     1600,
	}
)

// String returns the string representation of the Action
func (a action) String() string {
	return actionToStrMap[a]
}

// ParseAction attempts to parse the given value into Action.
// It supports string, fmt.Stringer, int, int64, and int32.
// If the value is not a valid Action or the value is not a supported type, it will return the enums unknown value (unknownAction).
func ParseAction(a any) Action {
	switch v := a.(type) {
	case Action:
		return v
	case string:
		return Action{stringToAction(v)}
	case fmt.Stringer:
		return Action{stringToAction(v.String())}
	case int:
		return Action{action(v)}
	case int64:
		return Action{action(int(v))}
	case int32:
		return Action{action(int(v))}
	}
	return Action{actionUnknown}
}

func stringToAction(s string) action {
	if v, ok := strToActionMap[s]; ok {
		return v
	}
	return actionUnknown
}

func (a action) IsValid() bool {
	return a >= action(1) && a <= action(len(actionToStrMap))
}

func (a action) GetPoints() int {
	return actionToPointsMap[a]
}

func (a action) EndsBackToBack() bool {
	switch a {
	case actionSingle, actionDouble, actionTriple:
		return true
	default:
		return false
	}
}

func (a action) StartsBackToBack() bool {
	switch a {
	case actionTetris, actionMiniTSpinSingle, actionTSpinSingle, actionTSpinDouble, actionTSpinTriple:
		return true
	default:
		return false
	}
}

type ActionContainer struct {
	UNKNOWN            Action
	NONE               Action
	SINGLE             Action
	DOUBLE             Action
	TRIPLE             Action
	TETRIS             Action
	MINI_T_SPIN        Action
	MINI_T_SPIN_SINGLE Action
	T_SPIN             Action
	T_SPIN_SINGLE      Action
	T_SPIN_DOUBLE      Action
	T_SPIN_TRIPLE      Action
}

// Actions is a global instance of ActionContainer that contains all possible values of type Action.
var Actions = ActionContainer{
	UNKNOWN:            Action{actionUnknown},
	NONE:               Action{actionNone},
	SINGLE:             Action{actionSingle},
	DOUBLE:             Action{actionDouble},
	TRIPLE:             Action{actionTriple},
	TETRIS:             Action{actionTetris},
	MINI_T_SPIN:        Action{actionMiniTSpin},
	MINI_T_SPIN_SINGLE: Action{actionMiniTSpinSingle},
	T_SPIN:             Action{actionTSpin},
	T_SPIN_SINGLE:      Action{actionTSpinSingle},
	T_SPIN_DOUBLE:      Action{actionTSpinDouble},
	T_SPIN_TRIPLE:      Action{actionTSpinTriple},
}
