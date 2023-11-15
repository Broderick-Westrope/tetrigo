package tetris

import (
	"fmt"
	"testing"
)

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

func TestMatrix_RemoveLine(t *testing.T) {
	testRows := Matrix{
		[10]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		[10]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[10]byte{0, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		[10]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 0},
		[10]byte{1, 1, 1, 1, 0, 1, 1, 0, 0, 1},
		[10]byte{1, 1, 0, 1, 1, 1, 1, 1, 1, 1},
	}

	for rowIndex := range testRows {
		t.Run(fmt.Sprint("row ", rowIndex), func(t *testing.T) {
			m := Matrix{}
			m[rowIndex] = testRows[rowIndex]

			if m[rowIndex] != testRows[rowIndex] {
				t.Errorf("Before: expected %v, got %v", testRows[rowIndex], m[rowIndex])
			}

			m.removeLine(rowIndex)

			for colIndex := range m[rowIndex] {
				if m[rowIndex][colIndex] != 0 {
					t.Errorf("After: expected 0, got %v", m[rowIndex][colIndex])
				}
			}
		})
	}
}
