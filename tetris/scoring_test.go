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
