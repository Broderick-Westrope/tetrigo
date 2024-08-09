package tetris

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewScoring(t *testing.T) {
	tt := map[string]struct {
		level    uint
		maxLevel uint
	}{
		"level 1": {
			level:    1,
			maxLevel: 15,
		},
		"level 15": {
			level:    15,
			maxLevel: 15,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			s := NewScoring(tc.level, tc.maxLevel)

			assert.Equal(t, tc.level, s.level)
			assert.Equal(t, tc.maxLevel, s.maxLevel)
			assert.Equal(t, uint(0), s.total)
			assert.Equal(t, uint(0), s.lines)
			assert.False(t, s.backToBack)
		})
	}
}

func TestScoring_Level(t *testing.T) {
	tt := map[string]struct {
		level uint
	}{
		"level 1": {
			level: 1,
		},
		"level 15": {
			level: 15,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			s := &Scoring{
				level: tc.level,
			}

			assert.Equal(t, tc.level, s.Level())
		})
	}
}

func TestScoring_Total(t *testing.T) {
	tt := map[string]struct {
		total uint
	}{
		"total 0": {
			total: 0,
		},
		"total 100": {
			total: 100,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			s := &Scoring{
				total: tc.total,
			}

			assert.Equal(t, tc.total, s.Total())
		})
	}
}

func TestScoring_Lines(t *testing.T) {
	tt := map[string]struct {
		lines uint
	}{
		"lines 0": {
			lines: 0,
		},
		"lines 100": {
			lines: 100,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			s := &Scoring{
				lines: tc.lines,
			}

			assert.Equal(t, tc.lines, s.Lines())
		})
	}
}

func TestScoring_AddSoftDrop(t *testing.T) {
	tt := map[string]struct {
		lines uint
	}{
		"0 lines": {
			lines: 0,
		},
		"10 lines": {
			lines: 10,
		},
		"123 lines": {
			lines: 123,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			s := &Scoring{
				total: 0,
			}

			s.AddSoftDrop(tc.lines)

			assert.Equal(t, tc.lines, s.total)
		})
	}
}

func TestScoring_AddHardDrop(t *testing.T) {
	tt := map[string]struct {
		lines uint
	}{
		"0 lines": {
			lines: 0,
		},
		"10 lines": {
			lines: 10,
		},
		"123 lines": {
			lines: 123,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			s := &Scoring{
				total: 0,
			}

			s.AddHardDrop(tc.lines)

			assert.Equal(t, tc.lines*2, s.total)
		})
	}
}

func TestScoring_ProcessAction(t *testing.T) {
	// Each Action should have 2 test cases: one with back-to-back enabled, and one without.
	tt := map[string]struct {
		a                  Action
		isBackToBack       bool
		maxLevel           uint
		expectedTotal      uint
		expectedBackToBack bool
	}{
		// Back-to-back disabled
		"none, no back to back": {
			a:                  actionNone,
			isBackToBack:       false,
			maxLevel:           0,
			expectedTotal:      0,
			expectedBackToBack: false,
		},
		"single, no back to back": {
			a:                  actionSingle,
			isBackToBack:       false,
			maxLevel:           0,
			expectedTotal:      100,
			expectedBackToBack: false,
		},
		"double, no back to back": {
			a:                  actionDouble,
			isBackToBack:       false,
			maxLevel:           0,
			expectedTotal:      300,
			expectedBackToBack: false,
		},
		"triple, no back to back": {
			a:                  actionTriple,
			isBackToBack:       false,
			maxLevel:           0,
			expectedTotal:      500,
			expectedBackToBack: false,
		},
		"tetris, no back to back": {
			a:                  actionTetris,
			isBackToBack:       false,
			maxLevel:           0,
			expectedTotal:      800,
			expectedBackToBack: true,
		},
		"mini T-spin, no back to back": {
			a:                  actionMiniTSpin,
			isBackToBack:       false,
			maxLevel:           0,
			expectedTotal:      100,
			expectedBackToBack: false,
		},
		"mini T-spin single, no back to back": {
			a:                  actionMiniTSpinSingle,
			isBackToBack:       false,
			maxLevel:           0,
			expectedTotal:      200,
			expectedBackToBack: true,
		},
		"T-spin, no back to back": {
			a:                  actionTSpin,
			isBackToBack:       false,
			maxLevel:           0,
			expectedTotal:      400,
			expectedBackToBack: false,
		},
		"T-spin single, no back to back": {
			a:                  actionTSpinSingle,
			isBackToBack:       false,
			maxLevel:           0,
			expectedTotal:      800,
			expectedBackToBack: true,
		},
		"T-spin double, no back to back": {
			a:                  actionTSpinDouble,
			isBackToBack:       false,
			maxLevel:           0,
			expectedTotal:      1200,
			expectedBackToBack: true,
		},
		"T-spin triple, no back to back": {
			a:                  actionTSpinTriple,
			isBackToBack:       false,
			maxLevel:           0,
			expectedTotal:      1600,
			expectedBackToBack: true,
		},
		// Back-to-back enabled
		"none, back to back": {
			a:                  actionNone,
			isBackToBack:       true,
			maxLevel:           0,
			expectedTotal:      0,
			expectedBackToBack: true,
		},
		"single, back to back": {
			a:                  actionSingle,
			isBackToBack:       true,
			maxLevel:           0,
			expectedTotal:      100,
			expectedBackToBack: false,
		},
		"double, back to back": {
			a:                  actionDouble,
			isBackToBack:       true,
			maxLevel:           0,
			expectedTotal:      300,
			expectedBackToBack: false,
		},
		"triple, back to back": {
			a:                  actionTriple,
			isBackToBack:       true,
			maxLevel:           0,
			expectedTotal:      500,
			expectedBackToBack: false,
		},
		"tetris, back to back": {
			a:                  actionTetris,
			isBackToBack:       true,
			maxLevel:           0,
			expectedTotal:      800 * 1.5,
			expectedBackToBack: true,
		},
		"mini T-spin, back to back": {
			a:                  actionMiniTSpin,
			isBackToBack:       true,
			maxLevel:           0,
			expectedTotal:      100,
			expectedBackToBack: true,
		},
		"mini T-spin single, back to back": {
			a:                  actionMiniTSpinSingle,
			isBackToBack:       true,
			maxLevel:           0,
			expectedTotal:      200 * 1.5,
			expectedBackToBack: true,
		},
		"T-spin, back to back": {
			a:                  actionTSpin,
			isBackToBack:       true,
			maxLevel:           0,
			expectedTotal:      400,
			expectedBackToBack: true,
		},
		"T-spin single, back to back": {
			a:                  actionTSpinSingle,
			isBackToBack:       true,
			maxLevel:           0,
			expectedTotal:      800 * 1.5,
			expectedBackToBack: true,
		},
		"T-spin double, back to back": {
			a:                  actionTSpinDouble,
			isBackToBack:       true,
			maxLevel:           0,
			expectedTotal:      1200 * 1.5,
			expectedBackToBack: true,
		},
		"T-spin triple, back to back": {
			a:                  actionTSpinTriple,
			isBackToBack:       true,
			maxLevel:           0,
			expectedTotal:      1600 * 1.5,
			expectedBackToBack: true,
		},
		"max level 1": {
			a:                  actionTSpinTriple,
			isBackToBack:       true,
			maxLevel:           1,
			expectedTotal:      1600 * 1.5,
			expectedBackToBack: true,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			s := &Scoring{
				backToBack: tc.isBackToBack,
				total:      0,
				lines:      0,
				level:      1,
				maxLevel:   tc.maxLevel,
			}

			s.ProcessAction(tc.a)

			if s.total != tc.expectedTotal {
				t.Errorf("Total: expected %d, got %d", tc.expectedTotal, s.total)
			}

			if s.backToBack != tc.expectedBackToBack {
				t.Errorf("BackToBack: expected %t, got %t", tc.expectedBackToBack, s.backToBack)
			}

			expectedLines := tc.expectedTotal / 100
			if s.lines != expectedLines {
				t.Errorf("Lines: expected %d, got %d", expectedLines, s.lines)
			}

			var expectedLevel uint
			if tc.maxLevel == 0 {
				expectedLevel = 1 + uint(expectedLines/5)
			} else {
				expectedLevel = tc.maxLevel
			}

			if s.level != expectedLevel {
				t.Errorf("Level: expected %d, got %d", expectedLevel, s.level)
			}
		})
	}
}
