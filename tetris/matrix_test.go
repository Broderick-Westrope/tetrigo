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

func TestMatrix_AddTetrimino(t *testing.T) {
	errorCases := []bool{true, false}

	for _, expectsErr := range errorCases {
		for _, tet := range Tetriminos {
			t.Run(fmt.Sprintf("Tetrimino %s, error %t", string(tet.Value), expectsErr), func(t *testing.T) {
				tet.Pos = Coordinate{0, 0}

				m := Matrix{}
				if expectsErr {
					m[0] = [10]byte{'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X'}
				}

				err := m.AddTetrimino(&tet)
				if expectsErr && err == nil {
					t.Error("expected an error, got nil")
				} else if !expectsErr && err != nil {
					t.Errorf("expected no error, got %v", err)
				} else if !expectsErr {
					for row := range tet.Cells {
						for col := range tet.Cells[row] {
							if tet.Cells[row][col] {
								if m[row][col] != tet.Value {
									t.Errorf("tet.Cells[%d][%d]: expected %v, got %v", row, col, tet.Value, m[row][col])
								}
							}
						}
					}
				}
			})
		}
	}
}
