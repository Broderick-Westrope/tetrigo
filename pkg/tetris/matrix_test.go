package tetris

import (
	"errors"
	"fmt"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMatrix(t *testing.T) {
	tt := map[string]struct {
		width   uint
		height  uint
		want    Matrix
		wantErr error
	}{
		"success": {
			width:  1,
			height: 21,
			want: Matrix{
				{0}, {0}, {0}, {0}, {0}, {0}, {0}, {0}, {0}, {0}, {0}, {0}, {0}, {0}, {0}, {0}, {0}, {0}, {0}, {0}, {0},
			},
		},
		"failure": {
			width:   1,
			height:  20,
			wantErr: ErrBufferZoneTooSmall,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			got, err := NewMatrix(tc.height, tc.width)
			if tc.wantErr != nil {
				assert.EqualError(t, tc.wantErr, err.Error())
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

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

func TestMatrix_modifyCell(t *testing.T) {
	tt := map[string]struct {
		matrix  Matrix
		pos     Coordinate
		wantErr error
	}{
		"success - top left": {
			matrix: Matrix{{0, 0}, {0, 0}},
			pos:    Coordinate{0, 0},
		},
		"success - bottom left": {
			matrix: Matrix{{0, 0}, {0, 0}},
			pos:    Coordinate{0, 1},
		},
		"success - top right": {
			matrix: Matrix{{0, 0}, {0, 0}},
			pos:    Coordinate{1, 0},
		},
		"success - bottom right": {
			matrix: Matrix{{0, 0}, {0, 0}},
			pos:    Coordinate{1, 1},
		},
		"failure - row out of bounds": {
			matrix:  Matrix{{0}},
			pos:     Coordinate{0, 1},
			wantErr: errors.New("row 1 is out of bounds"),
		},
		"failure - col out of bounds": {
			matrix:  Matrix{{0}},
			pos:     Coordinate{1, 0},
			wantErr: errors.New("col 1 is out of bounds"),
		},
		"failure - mino not expected value": {
			matrix:  Matrix{{'X'}},
			pos:     Coordinate{0, 0},
			wantErr: errors.New("mino at row 0, col 0 is 'X' (byte value 88) not the expected value"),
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			minos := [][]bool{{true}}
			newValue := byte('X')

			err := tc.matrix.modifyCell(minos, tc.pos, newValue, isCellEmpty)

			if tc.wantErr != nil {
				assert.EqualError(t, err, tc.wantErr.Error())
				return
			}

			assert.NoError(t, err)

			for row := range minos {
				for col := range minos[row] {
					if minos[row][col] {
						assert.Equal(t, newValue, tc.matrix[row+tc.pos.Y][col+tc.pos.X])
					}
				}
			}
		})
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
			expectedAction: Actions.NONE,
			wantMatrix:     Matrix{{0}, {'X'}},
		},
		"single - tet height 1": {
			matrix:         Matrix{{0}, {'X'}},
			posY:           1,
			cells:          [][]bool{{}},
			expectedAction: Actions.SINGLE,
			wantMatrix:     Matrix{{0}, {0}},
		},
		"single - tet height 2": {
			matrix:         Matrix{{0}, {'X'}},
			posY:           0,
			cells:          [][]bool{{}, {}},
			expectedAction: Actions.SINGLE,
			wantMatrix:     Matrix{{0}, {0}},
		},
		"double - tet height 2": {
			matrix:         Matrix{{0}, {'X'}, {'X'}},
			posY:           1,
			cells:          [][]bool{{}, {}},
			expectedAction: Actions.DOUBLE,
			wantMatrix:     Matrix{{0}, {0}, {0}},
		},
		"double - tet height 4": {
			matrix:         Matrix{{0}, {0}, {'X'}, {'X'}},
			posY:           0,
			cells:          [][]bool{{}, {}, {}, {}},
			expectedAction: Actions.DOUBLE,
			wantMatrix:     Matrix{{0}, {0}, {0}, {0}},
		},
		"triple - tet height 3": {
			matrix:         Matrix{{'X'}, {'X'}, {'X'}},
			posY:           0,
			cells:          [][]bool{{}, {}, {}},
			expectedAction: Actions.TRIPLE,
			wantMatrix:     Matrix{{0}, {0}, {0}},
		},
		"tetris - tet height 4": {
			matrix:         Matrix{{'X'}, {'X'}, {'X'}, {'X'}},
			posY:           0,
			cells:          [][]bool{{}, {}, {}, {}},
			expectedAction: Actions.TETRIS,
			wantMatrix:     Matrix{{0}, {0}, {0}, {0}},
		},
		"unknown - tet height 5": {
			matrix:         Matrix{{'X'}, {'X'}, {'X'}, {'X'}, {'X'}},
			posY:           0,
			cells:          [][]bool{{}, {}, {}, {}, {}},
			expectedAction: Actions.UNKNOWN,
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

func TestMatrix_IsOutOfBoundsHorizontally(t *testing.T) {
	tt := map[string]struct {
		col     int
		expects bool
	}{
		"out left": {
			-1, true,
		},
		"in left": {
			0, false,
		},
		"in right": {
			9, false,
		},
		"out right": {
			10, true,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			m := DefaultMatrix()
			result := m.isOutOfBoundsHorizontally(tc.col)

			assert.EqualValues(t, result, tc.expects)
		})
	}
}

func TestMatrix_IsOutOfBoundsVertically(t *testing.T) {
	tt := map[string]struct {
		row     int
		expects bool
	}{
		"out up": {
			-1, true,
		},
		"in up": {
			0, false,
		},
		"in down": {
			39, false,
		},
		"out down": {
			40, true,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			m := DefaultMatrix()
			result := m.isOutOfBoundsVertically(tc.row)

			assert.Equal(t, tc.expects, result)
		})
	}
}
