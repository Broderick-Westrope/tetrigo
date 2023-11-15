package tetris

import "testing"

func TestNewScoring(t *testing.T) {
	tt := []struct {
		name  string
		level uint
	}{
		{
			"level 1",
			1,
		},
		{
			"level 15",
			15,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s := NewScoring(tc.level)

			if s.level != tc.level {
				t.Errorf("Level: expected %d, got %d", tc.level, s.level)
			}
			if s.total != 0 {
				t.Errorf("Total: expected 0, got %d", s.total)
			}
			if s.lines != 0 {
				t.Errorf("Lines: expected 0, got %d", s.lines)
			}
			if s.backToBack {
				t.Error("BackToBack: expected false, got true")
			}
		})
	}
}

func TestScoring_Level(t *testing.T) {
	tt := []struct {
		name  string
		level uint
	}{
		{
			"level 1",
			1,
		},
		{
			"level 15",
			15,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s := &Scoring{
				level: tc.level,
			}

			if s.Level() != tc.level {
				t.Errorf("Level: expected %d, got %d", tc.level, s.Level())
			}
		})
	}
}

func TestScoring_Total(t *testing.T) {
	tt := []struct {
		name  string
		total uint
	}{
		{
			"total 0",
			0,
		},
		{
			"total 100",
			100,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s := &Scoring{
				total: tc.total,
			}

			if s.Total() != tc.total {
				t.Errorf("Total: expected %d, got %d", tc.total, s.Total())
			}
		})
	}
}

func TestScoring_Lines(t *testing.T) {
	tt := []struct {
		name  string
		lines uint
	}{
		{
			"lines 0",
			0,
		},
		{
			"lines 100",
			100,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s := &Scoring{
				lines: tc.lines,
			}

			if s.Lines() != tc.lines {
				t.Errorf("Lines: expected %d, got %d", tc.lines, s.Lines())
			}
		})
	}
}

func TestScoring_ProcessAction(t *testing.T) {
	// Each action should have 2 test cases: one with back-to-back enabled, and one without.
	tt := []struct {
		name               string
		a                  action
		isBackToBack       bool
		expectedTotal      uint
		expectedBackToBack bool
	}{
		// Back-to-back disabled
		{
			name:               "no action, no back to back",
			a:                  actionNone,
			isBackToBack:       false,
			expectedTotal:      0,
			expectedBackToBack: false,
		},
		{
			name:               "single, no back to back",
			a:                  actionSingle,
			isBackToBack:       false,
			expectedTotal:      100,
			expectedBackToBack: false,
		},
		{
			name:               "double, no back to back",
			a:                  actionDouble,
			isBackToBack:       false,
			expectedTotal:      300,
			expectedBackToBack: false,
		},
		{
			name:               "triple, no back to back",
			a:                  actionTriple,
			isBackToBack:       false,
			expectedTotal:      500,
			expectedBackToBack: false,
		},
		{
			name:               "tetris, no back to back",
			a:                  actionTetris,
			isBackToBack:       false,
			expectedTotal:      800,
			expectedBackToBack: true,
		},
		{
			name:               "mini T-spin, no back to back",
			a:                  actionMiniTSpin,
			isBackToBack:       false,
			expectedTotal:      100,
			expectedBackToBack: false,
		},
		{
			name:               "mini T-spin single, no back to back",
			a:                  actionMiniTSpinSingle,
			isBackToBack:       false,
			expectedTotal:      200,
			expectedBackToBack: true,
		},
		{
			name:               "T-spin, no back to back",
			a:                  actionTSpin,
			isBackToBack:       false,
			expectedTotal:      400,
			expectedBackToBack: false,
		},
		{
			name:               "T-spin single, no back to back",
			a:                  actionTSpinSingle,
			isBackToBack:       false,
			expectedTotal:      800,
			expectedBackToBack: true,
		},
		{
			name:               "T-spin double, no back to back",
			a:                  actionTSpinDouble,
			isBackToBack:       false,
			expectedTotal:      1200,
			expectedBackToBack: true,
		},
		{
			name:               "T-spin triple, no back to back",
			a:                  actionTSpinTriple,
			isBackToBack:       false,
			expectedTotal:      1600,
			expectedBackToBack: true,
		},
		// Back-to-back enabled
		{
			name:               "no action, back to back",
			a:                  actionNone,
			isBackToBack:       true,
			expectedTotal:      0,
			expectedBackToBack: true,
		},
		{
			name:               "single, back to back",
			a:                  actionSingle,
			isBackToBack:       true,
			expectedTotal:      100,
			expectedBackToBack: false,
		},
		{
			name:               "double, back to back",
			a:                  actionDouble,
			isBackToBack:       true,
			expectedTotal:      300,
			expectedBackToBack: false,
		},
		{
			name:               "triple, back to back",
			a:                  actionTriple,
			isBackToBack:       true,
			expectedTotal:      500,
			expectedBackToBack: false,
		},
		{
			name:               "tetris, back to back",
			a:                  actionTetris,
			isBackToBack:       true,
			expectedTotal:      800 * 1.5,
			expectedBackToBack: true,
		},
		{
			name:               "mini T-spin, back to back",
			a:                  actionMiniTSpin,
			isBackToBack:       true,
			expectedTotal:      100,
			expectedBackToBack: true,
		},
		{
			name:               "mini T-spin single, back to back",
			a:                  actionMiniTSpinSingle,
			isBackToBack:       true,
			expectedTotal:      200 * 1.5,
			expectedBackToBack: true,
		},
		{
			name:               "T-spin, back to back",
			a:                  actionTSpin,
			isBackToBack:       true,
			expectedTotal:      400,
			expectedBackToBack: true,
		},
		{
			name:               "T-spin single, back to back",
			a:                  actionTSpinSingle,
			isBackToBack:       true,
			expectedTotal:      800 * 1.5,
			expectedBackToBack: true,
		},
		{
			name:               "T-spin double, back to back",
			a:                  actionTSpinDouble,
			isBackToBack:       true,
			expectedTotal:      1200 * 1.5,
			expectedBackToBack: true,
		},
		{
			name:               "T-spin triple, back to back",
			a:                  actionTSpinTriple,
			isBackToBack:       true,
			expectedTotal:      1600 * 1.5,
			expectedBackToBack: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s := &Scoring{
				backToBack: tc.isBackToBack,
				total:      0,
				lines:      0,
				level:      1,
			}

			s.ProcessAction(tc.a)

			if s.total != tc.expectedTotal {
				t.Errorf("Total: expected %d, got %d", tc.expectedTotal, s.total)
			}

			if s.backToBack != tc.expectedBackToBack {
				t.Errorf("BackToBack: expected %t, got %t", tc.expectedBackToBack, s.backToBack)
			}

			expectedLines := (tc.expectedTotal / 100)
			if s.lines != expectedLines {
				t.Errorf("Lines: expected %d, got %d", expectedLines, s.lines)
			}

			expectedLevel := 1 + uint(expectedLines/5)
			if s.level != expectedLevel {
				t.Errorf("Level: expected %d, got %d", expectedLevel, s.level)
			}
		})
	}
}
