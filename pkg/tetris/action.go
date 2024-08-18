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
		actionUnknown:         "Unknown",
		actionNone:            "None",
		actionSingle:          "Single",
		actionDouble:          "Double",
		actionTriple:          "Triple",
		actionTetris:          "Tetris",
		actionMiniTSpin:       "MiniTSpin",
		actionMiniTSpinSingle: "MiniTSpinSingle",
		actionTSpin:           "TSpin",
		actionTSpinSingle:     "TSpinSingle",
		actionTSpinDouble:     "TSpinDouble",
		actionTSpinTriple:     "TSpinTriple",
	}

	strToActionMap = map[string]action{
		"Unknown":         actionUnknown,
		"None":            actionNone,
		"Single":          actionSingle,
		"Double":          actionDouble,
		"Triple":          actionTriple,
		"Tetris":          actionTetris,
		"MiniTSpin":       actionMiniTSpin,
		"MiniTSpinSingle": actionMiniTSpinSingle,
		"TSpin":           actionTSpin,
		"TSpinSingle":     actionTSpinSingle,
		"TSpinDouble":     actionTSpinDouble,
		"TSpinTriple":     actionTSpinTriple,
	}

	actionToPointsMap = map[action]int{
		actionUnknown:         0,
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

// String returns the string representation of the Action.
func (a action) String() string {
	return actionToStrMap[a]
}

// ParseAction attempts to parse the given value into Action.
// It supports string, fmt.Stringer, int, int64, and int32.
// If the value is not a valid Action or the value is not a supported type, it will
// return the enums unknown value (unknownAction).
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

func (a action) EndsBackToBack() (bool, error) {
	switch a {
	case actionSingle, actionDouble, actionTriple:
		return true, nil
	case actionUnknown, actionNone, actionTetris, actionMiniTSpin, actionMiniTSpinSingle,
		actionTSpin, actionTSpinSingle, actionTSpinDouble, actionTSpinTriple:
		return false, nil
	default:
		return false, fmt.Errorf("unknown action: %v", a)
	}
}

func (a action) StartsBackToBack() (bool, error) {
	switch a {
	case actionTetris, actionMiniTSpinSingle, actionTSpinSingle, actionTSpinDouble, actionTSpinTriple:
		return true, nil
	case actionUnknown, actionNone, actionSingle, actionDouble, actionTriple, actionMiniTSpin, actionTSpin:
		return false, nil
	default:
		return false, fmt.Errorf("unknown action: %v", a)
	}
}

type ActionContainer struct {
	Unknown         Action
	None            Action
	Single          Action
	Double          Action
	Triple          Action
	Tetris          Action
	MiniTSpin       Action
	MiniTSpinSingle Action
	TSpin           Action
	TSpinSingle     Action
	TSpinDouble     Action
	TSpinTriple     Action
}

// Actions is a global instance of ActionContainer that contains all possible values of type Action.
var Actions = ActionContainer{
	Unknown:         Action{actionUnknown},
	None:            Action{actionNone},
	Single:          Action{actionSingle},
	Double:          Action{actionDouble},
	Triple:          Action{actionTriple},
	Tetris:          Action{actionTetris},
	MiniTSpin:       Action{actionMiniTSpin},
	MiniTSpinSingle: Action{actionMiniTSpinSingle},
	TSpin:           Action{actionTSpin},
	TSpinSingle:     Action{actionTSpinSingle},
	TSpinDouble:     Action{actionTSpinDouble},
	TSpinTriple:     Action{actionTSpinTriple},
}
