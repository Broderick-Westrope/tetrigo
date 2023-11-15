package tetris

import "testing"

func TestMatrix_IsLineComplete(t *testing.T) {
	m := &Matrix{
		[10]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		[10]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[10]byte{0, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		[10]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 0},
		[10]byte{1, 1, 1, 1, 0, 1, 1, 0, 0, 1},
		[10]byte{1, 1, 0, 1, 1, 1, 1, 1, 1, 1},
	}

	// these test cases correspond to the rows in the matrix above
	tt := []struct {
		name     string
		expected bool
	}{
		{"row 0", true},
		{"row 1", false},
		{"row 2", false},
		{"row 3", false},
		{"row 4", false},
		{"row 5", false},
	}

	for tcIndex, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if actual := m.isLineComplete(tcIndex); actual != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}
