package tetris

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewScoring(t *testing.T) {
	tt := map[string]struct {
		level         uint
		maxLevel      uint
		increaseLevel bool
		endOnMaxLevel bool
		maxLines      uint
		endOnMaxLines bool
	}{
		"0; false": {
			level:         0,
			maxLevel:      0,
			increaseLevel: false,
			endOnMaxLevel: false,
			maxLines:      0,
			endOnMaxLines: false,
		},
		"10; true": {
			level:         10,
			maxLevel:      10,
			increaseLevel: true,
			endOnMaxLevel: true,
			maxLines:      10,
			endOnMaxLines: true,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			s := NewScoring(tc.level, tc.maxLevel, tc.increaseLevel, tc.endOnMaxLevel, tc.maxLines, tc.endOnMaxLines)

			assert.Equal(t, tc.level, s.level)
			assert.Equal(t, tc.maxLevel, s.maxLevel)
			assert.Equal(t, tc.increaseLevel, s.increaseLevel)
			assert.Equal(t, tc.endOnMaxLevel, s.endOnMaxLevel)
			assert.Equal(t, tc.maxLines, s.maxLines)
			assert.Equal(t, tc.endOnMaxLines, s.endOnMaxLines)

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
			a:                  Actions.None,
			isBackToBack:       false,
			maxLevel:           0,
			expectedTotal:      0,
			expectedBackToBack: false,
		},
		"single, no back to back": {
			a:                  Actions.Single,
			isBackToBack:       false,
			maxLevel:           0,
			expectedTotal:      100,
			expectedBackToBack: false,
		},
		"double, no back to back": {
			a:                  Actions.Double,
			isBackToBack:       false,
			maxLevel:           0,
			expectedTotal:      300,
			expectedBackToBack: false,
		},
		"triple, no back to back": {
			a:                  Actions.Triple,
			isBackToBack:       false,
			maxLevel:           0,
			expectedTotal:      500,
			expectedBackToBack: false,
		},
		"tetris, no back to back": {
			a:                  Actions.Tetris,
			isBackToBack:       false,
			maxLevel:           0,
			expectedTotal:      800,
			expectedBackToBack: true,
		},
		"mini T-spin, no back to back": {
			a:                  Actions.MiniTSpin,
			isBackToBack:       false,
			maxLevel:           0,
			expectedTotal:      100,
			expectedBackToBack: false,
		},
		"mini T-spin single, no back to back": {
			a:                  Actions.MiniTSpinSingle,
			isBackToBack:       false,
			maxLevel:           0,
			expectedTotal:      200,
			expectedBackToBack: true,
		},
		"T-spin, no back to back": {
			a:                  Actions.TSpin,
			isBackToBack:       false,
			maxLevel:           0,
			expectedTotal:      400,
			expectedBackToBack: false,
		},
		"T-spin single, no back to back": {
			a:                  Actions.TSpinSingle,
			isBackToBack:       false,
			maxLevel:           0,
			expectedTotal:      800,
			expectedBackToBack: true,
		},
		"T-spin double, no back to back": {
			a:                  Actions.TSpinDouble,
			isBackToBack:       false,
			maxLevel:           0,
			expectedTotal:      1200,
			expectedBackToBack: true,
		},
		"T-spin triple, no back to back": {
			a:                  Actions.TSpinTriple,
			isBackToBack:       false,
			maxLevel:           0,
			expectedTotal:      1600,
			expectedBackToBack: true,
		},
		// Back-to-back enabled
		"none, back to back": {
			a:                  Actions.None,
			isBackToBack:       true,
			maxLevel:           0,
			expectedTotal:      0,
			expectedBackToBack: true,
		},
		"single, back to back": {
			a:                  Actions.Single,
			isBackToBack:       true,
			maxLevel:           0,
			expectedTotal:      100,
			expectedBackToBack: false,
		},
		"double, back to back": {
			a:                  Actions.Double,
			isBackToBack:       true,
			maxLevel:           0,
			expectedTotal:      300,
			expectedBackToBack: false,
		},
		"triple, back to back": {
			a:                  Actions.Triple,
			isBackToBack:       true,
			maxLevel:           0,
			expectedTotal:      500,
			expectedBackToBack: false,
		},
		"tetris, back to back": {
			a:                  Actions.Tetris,
			isBackToBack:       true,
			maxLevel:           0,
			expectedTotal:      800 * 1.5,
			expectedBackToBack: true,
		},
		"mini T-spin, back to back": {
			a:                  Actions.MiniTSpin,
			isBackToBack:       true,
			maxLevel:           0,
			expectedTotal:      100,
			expectedBackToBack: true,
		},
		"mini T-spin single, back to back": {
			a:                  Actions.MiniTSpinSingle,
			isBackToBack:       true,
			maxLevel:           0,
			expectedTotal:      200 * 1.5,
			expectedBackToBack: true,
		},
		"T-spin, back to back": {
			a:                  Actions.TSpin,
			isBackToBack:       true,
			maxLevel:           0,
			expectedTotal:      400,
			expectedBackToBack: true,
		},
		"T-spin single, back to back": {
			a:                  Actions.TSpinSingle,
			isBackToBack:       true,
			maxLevel:           0,
			expectedTotal:      800 * 1.5,
			expectedBackToBack: true,
		},
		"T-spin double, back to back": {
			a:                  Actions.TSpinDouble,
			isBackToBack:       true,
			maxLevel:           0,
			expectedTotal:      1200 * 1.5,
			expectedBackToBack: true,
		},
		"T-spin triple, back to back": {
			a:                  Actions.TSpinTriple,
			isBackToBack:       true,
			maxLevel:           0,
			expectedTotal:      1600 * 1.5,
			expectedBackToBack: true,
		},
		"max level 1": {
			a:                  Actions.TSpinTriple,
			isBackToBack:       true,
			maxLevel:           1,
			expectedTotal:      1600 * 1.5,
			expectedBackToBack: true,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			s := &Scoring{
				level:         1,
				maxLevel:      tc.maxLevel,
				increaseLevel: true,
				endOnMaxLevel: false,

				lines:         0,
				maxLines:      0,
				endOnMaxLines: false,

				total:      0,
				backToBack: tc.isBackToBack,
			}

			// TODO: check gameOver (from endsOnMaxLevel)
			_, err := s.ProcessAction(tc.a)

			require.NoError(t, err)

			assert.Equal(t, tc.expectedTotal, s.total)
			assert.Equal(t, tc.expectedBackToBack, s.backToBack)

			expectedLines := tc.expectedTotal / 100
			assert.Equal(t, expectedLines, s.lines)

			var expectedLevel uint
			if tc.maxLevel == 0 {
				expectedLevel = 1 + (expectedLines / 5)
			} else {
				expectedLevel = tc.maxLevel
			}

			assert.Equal(t, expectedLevel, s.level)
		})
	}
}
