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
