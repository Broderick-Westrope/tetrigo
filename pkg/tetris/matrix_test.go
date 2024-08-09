package tetris

import (
	"fmt"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatrix_IsLineComplete(t *testing.T) {
	matrix := &Matrix{
		[]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]byte{0, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		[]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 0},
		[]byte{1, 1, 1, 1, 0, 1, 1, 0, 0, 1},
		[]byte{1, 1, 0, 1, 1, 1, 1, 1, 1, 1},
	}

	// these test cases correspond to the rows in the matrix above
	tt := []struct {
		expected bool
	}{
		{
			expected: true,
		},
		{
			expected: false,
		},
		{
			expected: false,
		},
		{
			expected: false,
		},
		{
			expected: false,
		},
		{
			expected: false,
		},
	}

	for row, tc := range tt {
		name := fmt.Sprintf("row %d", row)
		t.Run(name, func(t *testing.T) {
			if actual := matrix.isLineComplete(row); actual != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, actual)
			}
		})
	}
}

func TestMatrix_RemoveLine(t *testing.T) {
	matrix := Matrix{
		[]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]byte{0, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		[]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 0},
		[]byte{1, 1, 1, 1, 0, 1, 1, 0, 0, 1},
		[]byte{1, 1, 0, 1, 1, 1, 1, 1, 1, 1},
	}

	for row := range matrix {
		t.Run(fmt.Sprint("row ", row), func(t *testing.T) {
			m := DefaultMatrix()
			m[row] = slices.Clone(matrix[row])

			assert.ElementsMatch(t, matrix[row], m[row])

			m.removeLine(row)

			for colIndex := range m[row] {
				assert.Equal(t, byte(0), m[row][colIndex])
			}
		})
	}
}

func TestMatrix_RemoveTetrimino(t *testing.T) {
	testRows := []int{0, 10, 20, 30, 38}
	testCols := []int{0, 1, 2, 3, 4, 5, 6}

	for _, tet := range Tetriminos {
		for _, mRow := range testRows {
			for _, mCol := range testCols {
				name := fmt.Sprintf("Tetrimino %s, row %d, col %d", string(tet.Value), mRow, mCol)
				t.Run(name, func(t *testing.T) {
					m := DefaultMatrix()
					tet.Pos = Coordinate{X: mCol, Y: mRow}

					err := m.AddTetrimino(&tet)
					if !assert.NoError(t, err) {
						return
					}

					err = m.RemoveTetrimino(&tet)
					if !assert.NoError(t, err) {
						return
					}

					for row := range tet.Minos {
						for col := range tet.Minos[row] {
							assert.Equal(t, byte(0), m[row+tet.Pos.Y][col+tet.Pos.X])
						}
					}
				})
			}
		}
	}
}

func TestMatrix_AddTetrimino(t *testing.T) {
	errorCases := []bool{true, false}

	for _, wantErr := range errorCases {
		for _, tet := range Tetriminos {
			t.Run(fmt.Sprintf("Tetrimino %s, error %t", string(tet.Value), wantErr), func(t *testing.T) {
				tet.Pos = Coordinate{0, 0}

				m := DefaultMatrix()
				if wantErr {
					m[0] = []byte{'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X'}
				}

				err := m.AddTetrimino(&tet)
				if wantErr {
					assert.Error(t, err)
					return
				}
				assert.NoError(t, err)
				for row := range tet.Minos {
					for col := range tet.Minos[row] {
						if tet.Minos[row][col] {
							assert.Equal(t, tet.Value, m[row][col])
						}
					}
				}
			})
		}
	}
}

func TestMatrix_RemoveCompletedLines(t *testing.T) {
	tt := map[string]struct {
		matrix         Matrix
		posY           int
		cells          [][]bool
		expectedAction Action
		wantMatrix     Matrix
	}{
		"none - tet height 1": {
			matrix:         Matrix{{0}, {'X'}},
			posY:           0,
			cells:          [][]bool{{}},
			expectedAction: actionNone,
			wantMatrix:     Matrix{{0}, {'X'}},
		},
		"single - tet height 1": {
			matrix:         Matrix{{0}, {'X'}},
			posY:           1,
			cells:          [][]bool{{}},
			expectedAction: actionSingle,
			wantMatrix:     Matrix{{0}, {0}},
		},
		"single - tet height 2": {
			matrix:         Matrix{{0}, {'X'}},
			posY:           0,
			cells:          [][]bool{{}, {}},
			expectedAction: actionSingle,
			wantMatrix:     Matrix{{0}, {0}},
		},
		"double - tet height 2": {
			matrix:         Matrix{{0}, {'X'}, {'X'}},
			posY:           1,
			cells:          [][]bool{{}, {}},
			expectedAction: actionDouble,
			wantMatrix:     Matrix{{0}, {0}, {0}},
		},
		"double - tet height 4": {
			matrix:         Matrix{{0}, {0}, {'X'}, {'X'}},
			posY:           0,
			cells:          [][]bool{{}, {}, {}, {}},
			expectedAction: actionDouble,
			wantMatrix:     Matrix{{0}, {0}, {0}, {0}},
		},
		"triple - tet height 3": {
			matrix:         Matrix{{'X'}, {'X'}, {'X'}},
			posY:           0,
			cells:          [][]bool{{}, {}, {}},
			expectedAction: actionTriple,
			wantMatrix:     Matrix{{0}, {0}, {0}},
		},
		"tetris - tet height 4": {
			matrix:         Matrix{{'X'}, {'X'}, {'X'}, {'X'}},
			posY:           0,
			cells:          [][]bool{{}, {}, {}, {}},
			expectedAction: actionTetris,
			wantMatrix:     Matrix{{0}, {0}, {0}, {0}},
		},
		"unknown - tet height 5": {
			matrix:         Matrix{{'X'}, {'X'}, {'X'}, {'X'}, {'X'}},
			posY:           0,
			cells:          [][]bool{{}, {}, {}, {}, {}},
			expectedAction: actionNone,
			wantMatrix:     Matrix{{0}, {0}, {0}, {0}, {0}},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			tet := &Tetrimino{
				Pos:   Coordinate{0, tc.posY},
				Minos: tc.cells,
			}

			act := tc.matrix.RemoveCompletedLines(tet)

			assert.Equal(t, tc.expectedAction, act)
			assert.EqualValues(t, tc.wantMatrix, tc.matrix)
		})
	}
}
