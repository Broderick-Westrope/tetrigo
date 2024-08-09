package tetris

import (
	"fmt"
	"reflect"
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

func TestMatrix_RemoveTetrimino(t *testing.T) {
	matrixRows := []int{0, 10, 20, 30, 39}
	matrixCols := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	for _, tet := range Tetriminos {
		for _, mRow := range matrixRows {
			for _, mCol := range matrixCols {
				t.Run(fmt.Sprintf("Tetrimino %s, row %d, col %d", string(tet.Value), mRow, mCol), func(t *testing.T) {
					m := Matrix{}
					tet.Pos = Coordinate{mCol, mRow}

					err := m.AddTetrimino(&tet)
					if err != nil {
						return
					}

					err = m.RemoveTetrimino(&tet)
					if err != nil {
						t.Errorf("expected no error, got %v", err)
					}

					for row := range tet.Minos {
						for col := range tet.Minos[row] {
							if m[row+tet.Pos.Y][col+tet.Pos.X] != 0 {
								t.Errorf("expected 0, got %v", m[row+tet.Pos.Y][col+tet.Pos.X])
							}
						}
					}
				})
			}
		}
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
					for row := range tet.Minos {
						for col := range tet.Minos[row] {
							if tet.Minos[row][col] {
								if m[row][col] != tet.Value {
									t.Errorf("tet.Minos[%d][%d]: expected %v, got %v", row, col, tet.Value, m[row][col])
								}
							}
						}
					}
				}
			})
		}
	}
}

func TestMatrix_RemoveCompletedLines(t *testing.T) {
	tt := []struct {
		name           string
		matrix         *Matrix
		posY           int
		cells          [][]bool
		expectedAction action
		expectedMatrix *Matrix
	}{
		{
			"0 lines, height 1",
			&Matrix{
				{},
				{'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X'},
			},
			0,
			[][]bool{{}},
			actionNone,
			&Matrix{
				{},
				{'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X'},
			},
		},
		{
			"1 line, height 1",
			&Matrix{
				{},
				{'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X'},
			},
			1,
			[][]bool{{}},
			actionSingle,
			&Matrix{},
		},
		{
			"1 line, height 1",
			&Matrix{
				{},
				{'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X'},
			},
			0,
			[][]bool{{}, {}},
			actionSingle,
			&Matrix{},
		},
		{
			"2 lines, height 2",
			&Matrix{
				{},
				{'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X'},
				{'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X'},
			},
			1,
			[][]bool{{}, {}},
			actionDouble,
			&Matrix{},
		},
		{
			"2 lines, height 4",
			&Matrix{
				{},
				{},
				{'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X'},
				{'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X'},
			},
			0,
			[][]bool{{}, {}, {}, {}},
			actionDouble,
			&Matrix{},
		},
		{
			"3 lines",
			&Matrix{
				{'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X'},
				{'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X'},
				{'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X'},
			},
			0,
			[][]bool{{}, {}, {}},
			actionTriple,
			&Matrix{},
		},
		{
			"4 lines",
			&Matrix{
				{'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X'},
				{'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X'},
				{'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X'},
				{'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X'},
			},
			0,
			[][]bool{{}, {}, {}, {}},
			actionTetris,
			&Matrix{},
		},
		{
			"5 lines (unknown action)",
			&Matrix{
				{'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X'},
				{'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X'},
				{'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X'},
				{'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X'},
				{'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X', 'X'},
			},
			0,
			[][]bool{{}, {}, {}, {}, {}},
			actionNone,
			&Matrix{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tet := Tetrimino{
				Pos:   Coordinate{0, tc.posY},
				Minos: tc.cells,
			}

			actualAction := tc.matrix.RemoveCompletedLines(&tet)

			if actualAction != tc.expectedAction {
				t.Errorf("expected action %v, got %v", tc.expectedAction, actualAction)
			}

			if !reflect.DeepEqual(tc.matrix, tc.expectedMatrix) {
				t.Errorf("expected matrix %v, got %v", tc.expectedMatrix, tc.matrix)
			}
		})
	}
}
