package tetris

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: tests to add:
// - TestTetrimino_MoveDown failure cases

func TestTetrimino_MoveDown(t *testing.T) {
	tt := map[string]struct {
		tet    *Tetrimino
		matrix Matrix
		want   bool
	}{
		"true": {
			tet: &Tetrimino{
				Minos: [][]bool{{true}},
				Pos:   Coordinate{X: 0, Y: 0},
			},
			matrix: Matrix{
				{0},
				{0},
			},
			want: true,
		},
		"false - blocking cell": {
			tet: &Tetrimino{
				Minos: [][]bool{{true}},
				Pos:   Coordinate{X: 0, Y: 0},
			},
			matrix: Matrix{
				{0},
				{'#'},
			},
			want: false,
		},
		"false - out of bounds": {
			tet: &Tetrimino{
				Minos: [][]bool{{true}},
				Pos:   Coordinate{X: 0, Y: 0},
			},
			matrix: Matrix{
				{0},
			},
			want: false,
		},
		"false - blocking mino; uneven edge": {
			tet: &Tetrimino{
				Minos: [][]bool{
					{true, true},
					{false, true},
				},
				Pos: Coordinate{X: 0, Y: 0},
			},
			matrix: Matrix{
				{0, 0},
				{'#', 0},
				{0, 0},
			},
			want: false,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			originalMatrix := *tc.matrix.DeepCopy()
			originalY := tc.tet.Pos.Y

			moved := tc.tet.MoveDown(tc.matrix)

			assert.Equal(t, tc.want, moved)
			assert.EqualValues(t, originalMatrix, tc.matrix)

			if tc.want {
				assert.Equal(t, originalY+1, tc.tet.Pos.Y)
			} else {
				assert.Equal(t, originalY, tc.tet.Pos.Y)
			}
		})
	}
}

func TestTetrimino_MoveLeft(t *testing.T) {
	tt := map[string]struct {
		matrix Matrix
		tet    Tetrimino
		want   bool
	}{
		"true - empty matrix": {
			matrix: Matrix{
				{0, 0, 0},
				{0, 0, 0},
			},
			tet: Tetrimino{
				Minos: [][]bool{
					{true, true},
					{false, true},
				},
				Pos: Coordinate{X: 1, Y: 0},
			},
			want: true,
		},
		"true - perfect fit": {
			matrix: Matrix{
				{'#', 0, 0},
			},
			tet: Tetrimino{
				Minos: [][]bool{
					{true},
				},
				Pos: Coordinate{X: 2, Y: 0},
			},
			want: true,
		},
		"true - ghost cell": {
			matrix: Matrix{
				{'G', 0},
			},
			tet: Tetrimino{
				Minos: [][]bool{
					{true},
				},
				Pos: Coordinate{X: 1, Y: 0},
			},
			want: true,
		},
		"false - blocking cell": {
			matrix: Matrix{
				{'#', 0},
			},
			tet: Tetrimino{
				Minos: [][]bool{
					{true},
				},
				Pos: Coordinate{X: 1, Y: 0},
			},
			want: false,
		},
		"false - out of bounds": {
			matrix: Matrix{
				{0},
			},
			tet: Tetrimino{
				Minos: [][]bool{
					{true},
				},
				Pos: Coordinate{X: 0, Y: 0},
			},
			want: false,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			originalMatrix := *tc.matrix.DeepCopy()
			originalX := tc.tet.Pos.X

			moved := tc.tet.MoveLeft(tc.matrix)

			assert.Equal(t, tc.want, moved)
			assert.EqualValues(t, originalMatrix, tc.matrix)

			if tc.want {
				assert.Equal(t, originalX-1, tc.tet.Pos.X)
			} else {
				assert.Equal(t, originalX, tc.tet.Pos.X)
			}
		})
	}
}

func TestTetrimino_MoveRight(t *testing.T) {
	tt := map[string]struct {
		matrix Matrix
		tet    Tetrimino
		want   bool
	}{
		"true - empty matrix": {
			matrix: Matrix{
				{0, 0, 0},
				{0, 0, 0},
			},
			tet: Tetrimino{
				Minos: [][]bool{
					{true, true},
					{true, false},
				},
				Pos: Coordinate{X: 0, Y: 0},
			},
			want: true,
		},
		"true - perfect fit": {
			matrix: Matrix{
				{0, 0, '#'},
			},
			tet: Tetrimino{
				Minos: [][]bool{
					{true},
				},
				Pos: Coordinate{X: 0, Y: 0},
			},
			want: true,
		},
		"true - ghost cell": {
			matrix: Matrix{
				{0, 'G'},
			},
			tet: Tetrimino{
				Minos: [][]bool{
					{true},
				},
				Pos: Coordinate{X: 0, Y: 0},
			},
			want: true,
		},
		"false - blocking cell": {
			matrix: Matrix{
				{0, '#'},
			},
			tet: Tetrimino{
				Minos: [][]bool{
					{true},
				},
				Pos: Coordinate{X: 0, Y: 0},
			},
			want: false,
		},
		"false - out of bounds": {
			matrix: Matrix{
				{0},
			},
			tet: Tetrimino{
				Minos: [][]bool{
					{true},
				},
				Pos: Coordinate{X: 0, Y: 0},
			},
			want: false,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			originalMatrix := *tc.matrix.DeepCopy()
			originalX := tc.tet.Pos.X

			moved := tc.tet.MoveRight(tc.matrix)

			assert.Equal(t, tc.want, moved)
			assert.EqualValues(t, originalMatrix, tc.matrix)

			if tc.want {
				assert.Equal(t, originalX+1, tc.tet.Pos.X)
			} else {
				assert.Equal(t, originalX, tc.tet.Pos.X)
			}
		})
	}
}

func TestTetrimino_Rotate(t *testing.T) {
	tt := map[string]struct {
		clockwise bool
		matrix    Matrix
		tet       *Tetrimino
		wantTet   *Tetrimino
	}{
		"O unmodified": {
			tet:     &Tetrimino{Value: 'O'},
			wantTet: &Tetrimino{Value: 'O'},
		},
		"success; clockwise": {
			clockwise: true,
			matrix: Matrix{
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Minos: [][]bool{
					{true, true, true, true},
				},
				Pos:              Coordinate{X: 0, Y: 1},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Minos: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Pos:              Coordinate{X: 2, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['I'],
			},
		},
		"success; counter clockwise": {
			clockwise: false,
			matrix: Matrix{
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Minos: [][]bool{
					{true, true, true, true},
				},
				Pos:              Coordinate{X: 0, Y: 1},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Minos: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Pos:              Coordinate{X: 1, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['I'],
			},
		},
		"failure - no valid rotation; clockwise": {
			clockwise: true,
			matrix: Matrix{
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Minos: [][]bool{
					{true, true, true, true},
				},
				Pos:              Coordinate{X: 0, Y: 1},
				CompassDirection: 0,
				RotationCompass:  rotationCompass{},
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Minos: [][]bool{
					{true, true, true, true},
				},
				Pos:              Coordinate{X: 0, Y: 1},
				CompassDirection: 0,
				RotationCompass:  rotationCompass{},
			},
		},
		"failure - no valid rotation; counter clockwise": {
			clockwise: true,
			matrix: Matrix{
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Minos: [][]bool{
					{true, true, true, true},
				},
				Pos:              Coordinate{X: 0, Y: 1},
				CompassDirection: 0,
				RotationCompass:  rotationCompass{},
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Minos: [][]bool{
					{true, true, true, true},
				},
				Pos:              Coordinate{X: 0, Y: 1},
				CompassDirection: 0,
				RotationCompass:  rotationCompass{},
			},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			err := tc.tet.Rotate(tc.matrix, tc.clockwise)

			assert.NoError(t, err)
			assert.EqualValues(t, tc.wantTet, tc.tet)
		})
	}
}

func TestTetrimino_rotateClockwise(t *testing.T) {
	tt := map[string]struct {
		matrix  Matrix
		tet     *Tetrimino
		wantTet *Tetrimino
	}{
		"I; starting rotation 0 (north); rotation point 1": {
			matrix: Matrix{
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Minos: [][]bool{
					{true, true, true, true},
				},
				Pos:              Coordinate{X: 0, Y: 1},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Minos: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Pos:              Coordinate{X: 2, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['I'],
			},
		},
		"I; starting rotation 1 (east); rotation point 1": {
			matrix: Matrix{
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Minos: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Pos:              Coordinate{X: 2, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Minos: [][]bool{
					{true, true, true, true},
				},
				Pos:              Coordinate{X: 0, Y: 2},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['I'],
			},
		},
		"I; starting rotation 2 (south); rotation point 1": {
			matrix: Matrix{
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Minos: [][]bool{
					{true, true, true, true},
				},
				Pos:              Coordinate{X: 0, Y: 2},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Minos: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Pos:              Coordinate{X: 1, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['I'],
			},
		},
		"I; starting rotation 3 (west); rotation point 1": {
			matrix: Matrix{
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Minos: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Pos:              Coordinate{X: 1, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Minos: [][]bool{
					{true, true, true, true},
				},
				Pos:              Coordinate{X: 0, Y: 1},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['I'],
			},
		},

		"O; starting rotation 0 (north); rotation point 1": {
			matrix: Matrix{
				{0, 0},
				{0, 0},
			},
			tet: &Tetrimino{
				Value: 'O',
				Minos: [][]bool{
					{true, true},
					{true, true},
				},
				Pos:              Coordinate{X: 0, Y: 0},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['O'],
			},
			wantTet: &Tetrimino{
				Value: 'O',
				Minos: [][]bool{
					{true, true},
					{true, true},
				},
				Pos:              Coordinate{X: 0, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['O'],
			},
		},
		"O; starting rotation 1 (east); rotation point 1": {
			matrix: Matrix{
				{0, 0},
				{0, 0},
			},
			tet: &Tetrimino{
				Value: 'O',
				Minos: [][]bool{
					{true, true},
					{true, true},
				},
				Pos:              Coordinate{X: 0, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['O'],
			},
			wantTet: &Tetrimino{
				Value: 'O',
				Minos: [][]bool{
					{true, true},
					{true, true},
				},
				Pos:              Coordinate{X: 0, Y: 0},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['O'],
			},
		},
		"O; starting rotation 2 (south); rotation point 1": {
			matrix: Matrix{
				{0, 0},
				{0, 0},
			},
			tet: &Tetrimino{
				Value: 'O',
				Minos: [][]bool{
					{true, true},
					{true, true},
				},
				Pos:              Coordinate{X: 0, Y: 0},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['O'],
			},
			wantTet: &Tetrimino{
				Value: 'O',
				Minos: [][]bool{
					{true, true},
					{true, true},
				},
				Pos:              Coordinate{X: 0, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['O'],
			},
		},
		"O; starting rotation 3 (west); rotation point 1": {
			matrix: Matrix{
				{0, 0},
				{0, 0},
			},
			tet: &Tetrimino{
				Value: 'O',
				Minos: [][]bool{
					{true, true},
					{true, true},
				},
				Pos:              Coordinate{X: 0, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['O'],
			},
			wantTet: &Tetrimino{
				Value: 'O',
				Minos: [][]bool{
					{true, true},
					{true, true},
				},
				Pos:              Coordinate{X: 0, Y: 0},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['O'],
			},
		},

		"6; starting rotation 0 (north); rotation point 1": {
			matrix: Matrix{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{false, true, false},
					{true, true, true},
				},
				Pos:              Coordinate{X: 0, Y: 0},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{true, false},
					{true, true},
					{true, false},
				},
				Pos:              Coordinate{X: 1, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['6'],
			},
		},
		"6; starting rotation 1 (east); rotation point 1": {
			matrix: Matrix{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{true, false},
					{true, true},
					{true, false},
				},
				Pos:              Coordinate{X: 1, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Pos:              Coordinate{X: 0, Y: 1},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['6'],
			},
		},
		"6; starting rotation 2 (south); rotation point 1": {
			matrix: Matrix{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Pos:              Coordinate{X: 0, Y: 1},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{false, true},
					{true, true},
					{false, true},
				},
				Pos:              Coordinate{X: 0, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['6'],
			},
		},
		"6; starting rotation 3 (west); rotation point 1": {
			matrix: Matrix{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{false, true},
					{true, true},
					{false, true},
				},
				Pos:              Coordinate{X: 0, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{false, true, false},
					{true, true, true},
				},
				Pos:              Coordinate{X: 0, Y: 0},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['6'],
			},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			result, err := tc.tet.rotateClockwise(tc.matrix)

			assert.NoError(t, err)
			assert.True(t, result)
			assert.EqualValues(t, tc.wantTet, tc.tet)
		})
	}
}

func TestTetrimino_rotateCounterClockwise(t *testing.T) {
	tt := map[string]struct {
		matrix  Matrix
		tet     *Tetrimino
		wantTet *Tetrimino
	}{
		"I; starting rotation 0 (north); rotation point 1": {
			matrix: Matrix{
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Minos: [][]bool{
					{true, true, true, true},
				},
				Pos:              Coordinate{X: 0, Y: 1},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Minos: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Pos:              Coordinate{X: 1, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['I'],
			},
		},
		"I; starting rotation 1 (east); rotation point 1": {
			matrix: Matrix{
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Minos: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Pos:              Coordinate{X: 2, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Minos: [][]bool{
					{true, true, true, true},
				},
				Pos:              Coordinate{X: 0, Y: 1},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['I'],
			},
		},
		"I; starting rotation 2 (south); rotation point 1": {
			matrix: Matrix{
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Minos: [][]bool{
					{true, true, true, true},
				},
				Pos:              Coordinate{X: 0, Y: 2},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Minos: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Pos:              Coordinate{X: 2, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['I'],
			},
		},
		"I; starting rotation 3 (west); rotation point 1": {
			matrix: Matrix{
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Minos: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Pos:              Coordinate{X: 1, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Minos: [][]bool{
					{true, true, true, true},
				},
				Pos:              Coordinate{X: 0, Y: 2},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['I'],
			},
		},

		"O; starting rotation 0 (north); rotation point 1": {
			matrix: Matrix{
				{0, 0},
				{0, 0},
			},
			tet: &Tetrimino{
				Value: 'O',
				Minos: [][]bool{
					{true, true},
					{true, true},
				},
				Pos:              Coordinate{X: 0, Y: 0},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['O'],
			},
			wantTet: &Tetrimino{
				Value: 'O',
				Minos: [][]bool{
					{true, true},
					{true, true},
				},
				Pos:              Coordinate{X: 0, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['O'],
			},
		},
		"O; starting rotation 1 (east); rotation point 1": {
			matrix: Matrix{
				{0, 0},
				{0, 0},
			},
			tet: &Tetrimino{
				Value: 'O',
				Minos: [][]bool{
					{true, true},
					{true, true},
				},
				Pos:              Coordinate{X: 0, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['O'],
			},
			wantTet: &Tetrimino{
				Value: 'O',
				Minos: [][]bool{
					{true, true},
					{true, true},
				},
				Pos:              Coordinate{X: 0, Y: 0},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['O'],
			},
		},
		"O; starting rotation 2 (south); rotation point 1": {
			matrix: Matrix{
				{0, 0},
				{0, 0},
			},
			tet: &Tetrimino{
				Value: 'O',
				Minos: [][]bool{
					{true, true},
					{true, true},
				},
				Pos:              Coordinate{X: 0, Y: 0},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['O'],
			},
			wantTet: &Tetrimino{
				Value: 'O',
				Minos: [][]bool{
					{true, true},
					{true, true},
				},
				Pos:              Coordinate{X: 0, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['O'],
			},
		},
		"O; starting rotation 3 (west); rotation point 1": {
			matrix: Matrix{
				{0, 0},
				{0, 0},
			},
			tet: &Tetrimino{
				Value: 'O',
				Minos: [][]bool{
					{true, true},
					{true, true},
				},
				Pos:              Coordinate{X: 0, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['O'],
			},
			wantTet: &Tetrimino{
				Value: 'O',
				Minos: [][]bool{
					{true, true},
					{true, true},
				},
				Pos:              Coordinate{X: 0, Y: 0},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['O'],
			},
		},

		"6; starting rotation 0 (north); rotation point 1": {
			matrix: Matrix{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{false, true, false},
					{true, true, true},
				},
				Pos:              Coordinate{X: 0, Y: 0},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{false, true},
					{true, true},
					{false, true},
				},
				Pos:              Coordinate{X: 0, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['6'],
			},
		},
		"6; starting rotation 1 (east); rotation point 1": {
			matrix: Matrix{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{true, false},
					{true, true},
					{true, false},
				},
				Pos:              Coordinate{X: 1, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{false, true, false},
					{true, true, true},
				},
				Pos:              Coordinate{X: 0, Y: 0},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['6'],
			},
		},
		"6; starting rotation 2 (south); rotation point 1": {
			matrix: Matrix{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Pos:              Coordinate{X: 0, Y: 1},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{true, false},
					{true, true},
					{true, false},
				},
				Pos:              Coordinate{X: 1, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['6'],
			},
		},
		"6; starting rotation 3 (west); rotation point 1": {
			matrix: Matrix{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{false, true},
					{true, true},
					{false, true},
				},
				Pos:              Coordinate{X: 0, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Pos:              Coordinate{X: 0, Y: 1},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['6'],
			},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			result, err := tc.tet.rotateCounterClockwise(tc.matrix)

			assert.NoError(t, err)
			assert.True(t, result)
			assert.EqualValues(t, tc.wantTet, tc.tet)
		})
	}
}

func TestTetrimino_transpose(t *testing.T) {
	tt := map[string]struct {
		tet       *Tetrimino
		wantMinos [][]bool
	}{
		"1x2": {
			tet: &Tetrimino{
				Minos: [][]bool{
					{true, false},
				},
			},
			wantMinos: [][]bool{
				{true},
				{false},
			},
		},
		"2x2": {
			tet: &Tetrimino{
				Minos: [][]bool{
					{true, false},
					{true, false},
				},
			},
			wantMinos: [][]bool{
				{true, true},
				{false, false},
			},
		},
		"3x3": {
			tet: &Tetrimino{
				Minos: [][]bool{
					{true, false, true},
					{true, false, false},
					{true, false, true},
				},
			},
			wantMinos: [][]bool{
				{true, true, true},
				{false, false, false},
				{true, false, true},
			},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			tc.tet.transpose()

			assert.ElementsMatch(t, tc.tet.Minos, tc.wantMinos)
		})
	}
}

func TestTetrimino_isValid(t *testing.T) {
	tt := map[string]struct {
		matrix Matrix
		tet    *Tetrimino
		want   bool
	}{
		"true; empty board": {
			matrix: Matrix{
				{0},
			},
			tet: &Tetrimino{
				Minos: [][]bool{{true}},
				Pos:   Coordinate{X: 0, Y: 0},
			},
			want: true,
		},
		"true; perfect fit": {
			matrix: Matrix{
				{0, 0, 0},
				{'X', 0, 'X'},
			},
			tet: &Tetrimino{
				Minos: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Pos: Coordinate{X: 0, Y: 0},
			},
			want: true,
		},
		"false; blocking mino": {
			matrix: Matrix{
				{'#'},
			},
			tet: &Tetrimino{
				Minos: [][]bool{{true}},
				Pos:   Coordinate{X: 0, Y: 0},
			},
			want: false,
		},
		"false; out of bounds left": {
			Matrix{{0}},
			&Tetrimino{
				Minos: [][]bool{{true}},
				Pos:   Coordinate{X: -1, Y: 0},
			},
			false,
		},
		"false; out of bounds right": {
			Matrix{{0}},
			&Tetrimino{
				Minos: [][]bool{{true}},
				Pos:   Coordinate{X: 1, Y: 0},
			},
			false,
		},
		"false; out of bounds up": {
			Matrix{{0}},
			&Tetrimino{
				Minos: [][]bool{{true}},
				Pos:   Coordinate{X: 0, Y: -1},
			},
			false,
		},
		"false; out of bounds down": {
			Matrix{{0}},
			&Tetrimino{
				Minos: [][]bool{{true}},
				Pos:   Coordinate{X: 0, Y: 1},
			},
			false,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			result := tc.tet.isValid(tc.matrix)

			assert.EqualValues(t, result, tc.want)
		})
	}
}

func Test_positiveMod(t *testing.T) {
	tt := map[string]struct {
		dividend int
		divisor  int
		want     int
		wantErr  error
	}{
		"0 mod 0": {
			0, 0, 0, errors.New("result is NaN"),
		},
		"1 mod 0": {
			0, 0, 0, errors.New("result is NaN"),
		},
		"3 mod 5": {
			3, 5, 3, nil,
		},
		"5 mod 5": {
			5, 5, 0, nil,
		},
		"3 mod -5": {
			3, -5, 3, nil,
		},
		"5 mod -5": {
			5, -5, 0, nil,
		},
		"-4 mod -5": {
			-4, -5, -4, nil,
		},
		"-4 mod 5": {
			-4, 5, 1, nil,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			result, err := positiveMod(tc.dividend, tc.divisor)
			if tc.wantErr != nil {
				assert.EqualError(t, err, tc.wantErr.Error())
				assert.Equal(t, tc.want, result)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_deepCopyMinos(t *testing.T) {
	tt := map[string]struct {
		cells  [][]bool
		modify func(*[][]bool)
	}{
		"modify byte": {
			cells: [][]bool{
				{false, false},
			},
			modify: func(cells *[][]bool) {
				(*cells)[0][0] = true
			},
		},
		"modify inner array": {
			cells: [][]bool{
				{false, false},
			},
			modify: func(cells *[][]bool) {
				(*cells)[0] = []bool{true, false}
			},
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			cellsCopy := deepCopyMinos(tc.cells)
			tc.modify(&cellsCopy)

			assert.False(t, tc.cells[0][0])
		})
	}
}

func TestTetrimino_DeepCopy(t *testing.T) {
	tet := Tetrimino{
		Value:            byte(0),
		Minos:            [][]bool{{false}},
		Pos:              Coordinate{0, 0},
		CompassDirection: 0,
		RotationCompass: rotationCompass{
			{&Coordinate{0, 0}},
		},
	}

	// Create a copy manually
	manualCopy := Tetrimino{
		Value:            byte(0),
		Minos:            [][]bool{{false}},
		Pos:              Coordinate{0, 0},
		CompassDirection: 0,
		RotationCompass:  rotationCompass{{&Coordinate{0, 0}}},
	}

	// Create a (dereferences) copy with the helper function
	easyCopy := *tet.DeepCopy()

	// Assert that all are equal
	assert.EqualValues(t, tet, manualCopy)
	assert.EqualValues(t, tet, easyCopy)
	assert.EqualValues(t, manualCopy, easyCopy)

	// Modify the original
	tet.Value = 1
	tet.Minos = [][]bool{{true}}
	tet.Pos = Coordinate{1, 1}
	tet.CompassDirection = 1
	tet.RotationCompass[0][0].X = 1

	// Assert that the original changed but both copies stayed the same
	assert.NotEqualValues(t, tet, manualCopy)
	assert.NotEqualValues(t, tet, easyCopy)
	assert.EqualValues(t, manualCopy, easyCopy)
}

func TestTetrimino_IsAboveSkyline(t *testing.T) {
	tt := map[string]struct {
		skyline int
		tet     *Tetrimino
		want    bool
	}{
		"true; skyline 10": {
			skyline: 10,
			tet: &Tetrimino{
				Pos:   Coordinate{X: 0, Y: 9},
				Minos: [][]bool{{true}},
			},
			want: true,
		},
		"true; skyline 20": {
			skyline: 20,
			tet: &Tetrimino{
				Pos:   Coordinate{X: 0, Y: 19},
				Minos: [][]bool{{true}},
			},
			want: true,
		},
		"false; skyline 10": {
			skyline: 10,
			tet: &Tetrimino{
				Pos:   Coordinate{X: 0, Y: 10},
				Minos: [][]bool{{true}},
			},
			want: false,
		},
		"false; skyline 20": {
			skyline: 20,
			tet: &Tetrimino{
				Pos:   Coordinate{X: 0, Y: 20},
				Minos: [][]bool{{true}},
			},
			want: false,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			result := tc.tet.IsAboveSkyline(tc.skyline)

			assert.Equal(t, tc.want, result)
		})
	}
}

func TestTetrimino_IsOverlapping(t *testing.T) {
	tt := map[string]struct {
		matrix Matrix
		tet    *Tetrimino
		want   bool
	}{
		"true": {
			matrix: Matrix{
				{0, 'X'},
			},
			tet: &Tetrimino{
				Pos:   Coordinate{X: 0, Y: 0},
				Minos: [][]bool{{true, true}},
			},
			want: true,
		},
		"false": {
			matrix: Matrix{
				{0, 'X', 0, 'X'},
				{'X', 'X', 'X'},
			},
			tet: &Tetrimino{
				Pos:   Coordinate{X: 0, Y: 0},
				Minos: [][]bool{{true, false, true}},
			},
			want: false,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			result := tc.tet.IsOverlapping(tc.matrix)

			assert.Equal(t, tc.want, result)
		})
	}
}
