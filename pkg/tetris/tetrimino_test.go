package tetris

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTetrimino_MoveDown(t *testing.T) {
	tt := map[string]struct {
		tet    *Tetrimino
		matrix Matrix
		want   bool
	}{
		"true": {
			tet: &Tetrimino{
				Cells:    [][]bool{{true}},
				Position: Coordinate{X: 0, Y: 0},
			},
			matrix: Matrix{
				{0},
				{0},
			},
			want: true,
		},
		"false - blocking cell": {
			tet: &Tetrimino{
				Cells:    [][]bool{{true}},
				Position: Coordinate{X: 0, Y: 0},
			},
			matrix: Matrix{
				{0},
				{'#'},
			},
			want: false,
		},
		"false - out of bounds": {
			tet: &Tetrimino{
				Cells:    [][]bool{{true}},
				Position: Coordinate{X: 0, Y: 0},
			},
			matrix: Matrix{
				{0},
			},
			want: false,
		},
		"false - blocking mino; uneven edge": {
			tet: &Tetrimino{
				Cells: [][]bool{
					{true, true},
					{false, true},
				},
				Position: Coordinate{X: 0, Y: 0},
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
			originalY := tc.tet.Position.Y

			moved := tc.tet.MoveDown(tc.matrix)

			assert.Equal(t, tc.want, moved)
			assert.Equal(t, originalMatrix, tc.matrix)

			if tc.want {
				assert.Equal(t, originalY+1, tc.tet.Position.Y)
			} else {
				assert.Equal(t, originalY, tc.tet.Position.Y)
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
				Cells: [][]bool{
					{true, true},
					{false, true},
				},
				Position: Coordinate{X: 1, Y: 0},
			},
			want: true,
		},
		"true - perfect fit": {
			matrix: Matrix{
				{'#', 0, 0},
			},
			tet: Tetrimino{
				Cells: [][]bool{
					{true},
				},
				Position: Coordinate{X: 2, Y: 0},
			},
			want: true,
		},
		"true - ghost cell": {
			matrix: Matrix{
				{'G', 0},
			},
			tet: Tetrimino{
				Cells: [][]bool{
					{true},
				},
				Position: Coordinate{X: 1, Y: 0},
			},
			want: true,
		},
		"false - blocking cell": {
			matrix: Matrix{
				{'#', 0},
			},
			tet: Tetrimino{
				Cells: [][]bool{
					{true},
				},
				Position: Coordinate{X: 1, Y: 0},
			},
			want: false,
		},
		"false - out of bounds": {
			matrix: Matrix{
				{0},
			},
			tet: Tetrimino{
				Cells: [][]bool{
					{true},
				},
				Position: Coordinate{X: 0, Y: 0},
			},
			want: false,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			originalMatrix := *tc.matrix.DeepCopy()
			originalX := tc.tet.Position.X

			moved := tc.tet.MoveLeft(tc.matrix)

			assert.Equal(t, tc.want, moved)
			assert.Equal(t, originalMatrix, tc.matrix)

			if tc.want {
				assert.Equal(t, originalX-1, tc.tet.Position.X)
			} else {
				assert.Equal(t, originalX, tc.tet.Position.X)
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
				Cells: [][]bool{
					{true, true},
					{true, false},
				},
				Position: Coordinate{X: 0, Y: 0},
			},
			want: true,
		},
		"true - perfect fit": {
			matrix: Matrix{
				{0, 0, '#'},
			},
			tet: Tetrimino{
				Cells: [][]bool{
					{true},
				},
				Position: Coordinate{X: 0, Y: 0},
			},
			want: true,
		},
		"true - ghost cell": {
			matrix: Matrix{
				{0, 'G'},
			},
			tet: Tetrimino{
				Cells: [][]bool{
					{true},
				},
				Position: Coordinate{X: 0, Y: 0},
			},
			want: true,
		},
		"false - blocking cell": {
			matrix: Matrix{
				{0, '#'},
			},
			tet: Tetrimino{
				Cells: [][]bool{
					{true},
				},
				Position: Coordinate{X: 0, Y: 0},
			},
			want: false,
		},
		"false - out of bounds": {
			matrix: Matrix{
				{0},
			},
			tet: Tetrimino{
				Cells: [][]bool{
					{true},
				},
				Position: Coordinate{X: 0, Y: 0},
			},
			want: false,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			originalMatrix := *tc.matrix.DeepCopy()
			originalX := tc.tet.Position.X

			moved := tc.tet.MoveRight(tc.matrix)

			assert.Equal(t, tc.want, moved)
			assert.Equal(t, originalMatrix, tc.matrix)

			if tc.want {
				assert.Equal(t, originalX+1, tc.tet.Position.X)
			} else {
				assert.Equal(t, originalX, tc.tet.Position.X)
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
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 1},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Position:         Coordinate{X: 2, Y: 0},
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
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 1},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Position:         Coordinate{X: 1, Y: 0},
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
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 1},
				CompassDirection: 0,
				RotationCompass:  RotationCompass{},
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 1},
				CompassDirection: 0,
				RotationCompass:  RotationCompass{},
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
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 1},
				CompassDirection: 0,
				RotationCompass:  RotationCompass{},
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 1},
				CompassDirection: 0,
				RotationCompass:  RotationCompass{},
			},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			err := tc.tet.Rotate(tc.matrix, tc.clockwise)

			require.NoError(t, err)
			assert.Equal(t, tc.wantTet, tc.tet)
		})
	}
}

func TestTetrimino_rotateClockwise(t *testing.T) {
	tt := map[string]struct {
		matrix            Matrix
		tet               *Tetrimino
		wantTet           *Tetrimino
		wantRotationPoint int
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
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 1},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Position:         Coordinate{X: 2, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['I'],
			},
			wantRotationPoint: 1,
		},
		"I; starting rotation 0 (north); rotation point 2": {
			matrix: Matrix{
				{0, 0, '#', 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 1},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['I'],
			},
			wantRotationPoint: 2,
		},
		"I; starting rotation 0 (north); rotation point 3": {
			matrix: Matrix{
				{'#', 0, '#', 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 1},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Position:         Coordinate{X: 3, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['I'],
			},
			wantRotationPoint: 3,
		},
		"I; starting rotation 0 (north); rotation point 4": {
			matrix: Matrix{
				{'#', 0, '#', 0},
				{0, 0, 0, 0},
				{0, 0, 0, '#'},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 1},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Position:         Coordinate{X: 0, Y: 1},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['I'],
			},
			wantRotationPoint: 4,
		},
		"I; starting rotation 0 (north); rotation point 5": {
			matrix: Matrix{
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{'#', 0, '#', 0},
				{0, 0, 0, 0},
				{0, 0, 0, '#'},
				{0, 0, 0, 0},
				{'#', 0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 3},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Position:         Coordinate{X: 3, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['I'],
			},
			wantRotationPoint: 5,
		},
		"I; starting rotation 0 (north); no valid rotation": {
			matrix: Matrix{
				{0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 3},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: nil,
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
				Cells: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Position:         Coordinate{X: 2, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 2},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['I'],
			},
			wantRotationPoint: 1,
		},
		"I; starting rotation 1 (east); rotation point 2": {
			matrix: Matrix{
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Position:         Coordinate{X: 3, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 2},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['I'],
			},
			wantRotationPoint: 2,
		},
		"I; starting rotation 1 (east); rotation point 3": {
			matrix: Matrix{
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 2},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['I'],
			},
			wantRotationPoint: 3,
		},
		"I; starting rotation 1 (east); rotation point 4": {
			matrix: Matrix{
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{'#', 0, 0, 0},
				{0, 0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Position:         Coordinate{X: 3, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['I'],
			},
			wantRotationPoint: 4,
		},
		"I; starting rotation 1 (east); rotation point 5": {
			matrix: Matrix{
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, '#'},
				{0, 0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 3},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['I'],
			},
			wantRotationPoint: 5,
		},
		"I; starting rotation 1 (east); no valid rotation": {
			matrix: Matrix{
				{0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: nil,
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
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 2},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Position:         Coordinate{X: 1, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['I'],
			},
			wantRotationPoint: 1,
		},
		"I; starting rotation 2 (south); rotation point 2": {
			matrix: Matrix{
				{0, '#', 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 2},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Position:         Coordinate{X: 3, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['I'],
			},
			wantRotationPoint: 2,
		},
		"I; starting rotation 2 (south); rotation point 3": {
			matrix: Matrix{
				{0, '#', 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, '#'},
			},
			tet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 2},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['I'],
			},
			wantRotationPoint: 3,
		},
		"I; starting rotation 2 (south); rotation point 4": {
			matrix: Matrix{
				{0, 0, 0, 0},
				{'#', '#', 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, '#'},
			},
			tet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 3},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Position:         Coordinate{X: 3, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['I'],
			},
			wantRotationPoint: 4,
		},
		"I; starting rotation 2 (south); rotation point 5": {
			matrix: Matrix{
				{'#', '#', 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, '#'},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 2},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Position:         Coordinate{X: 0, Y: 2},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['I'],
			},
			wantRotationPoint: 5,
		},
		"I; starting rotation 2 (south); no valid rotation": {
			matrix: Matrix{
				{0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: nil,
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
				Cells: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Position:         Coordinate{X: 1, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 1},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['I'],
			},
			wantRotationPoint: 1,
		},
		"I; starting rotation 3 (west); rotation point 2": {
			matrix: Matrix{
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 1},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['I'],
			},
			wantRotationPoint: 2,
		},
		"I; starting rotation 3 (west); rotation point 3": {
			matrix: Matrix{
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Position:         Coordinate{X: 3, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 1},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['I'],
			},
			wantRotationPoint: 3,
		},
		"I; starting rotation 3 (west); rotation point 4": {
			matrix: Matrix{
				{0, 0, 0, 0},
				{0, '#', 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 3},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['I'],
			},
			wantRotationPoint: 4,
		},
		"I; starting rotation 3 (west); rotation point 5": {
			matrix: Matrix{
				{0, 0, 0, 0},
				{'#', 0, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Position:         Coordinate{X: 3, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['I'],
			},
			wantRotationPoint: 5,
		},
		"I; starting rotation 3 (west); no valid rotation": {
			matrix: Matrix{
				{0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Position:         Coordinate{X: 3, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: nil,
		},

		"O; starting rotation 0 (north); rotation point 1": {
			matrix: Matrix{
				{0, 0},
				{0, 0},
			},
			tet: &Tetrimino{
				Value: 'O',
				Cells: [][]bool{
					{true, true},
					{true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['O'],
			},
			wantTet: &Tetrimino{
				Value: 'O',
				Cells: [][]bool{
					{true, true},
					{true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['O'],
			},
			wantRotationPoint: 1,
		},
		"O; starting rotation 0 (north); no valid rotation": {
			matrix: Matrix{
				{0},
			},
			tet: &Tetrimino{
				Value: 'O',
				Cells: [][]bool{
					{true, true},
					{true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['O'],
			},
			wantTet: nil,
		},

		"O; starting rotation 1 (east); rotation point 1": {
			matrix: Matrix{
				{0, 0},
				{0, 0},
			},
			tet: &Tetrimino{
				Value: 'O',
				Cells: [][]bool{
					{true, true},
					{true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['O'],
			},
			wantTet: &Tetrimino{
				Value: 'O',
				Cells: [][]bool{
					{true, true},
					{true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['O'],
			},
			wantRotationPoint: 1,
		},
		"O; starting rotation 1 (east); no valid rotation": {
			matrix: Matrix{
				{0},
			},
			tet: &Tetrimino{
				Value: 'O',
				Cells: [][]bool{
					{true, true},
					{true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['O'],
			},
			wantTet: nil,
		},

		"O; starting rotation 2 (south); rotation point 1": {
			matrix: Matrix{
				{0, 0},
				{0, 0},
			},
			tet: &Tetrimino{
				Value: 'O',
				Cells: [][]bool{
					{true, true},
					{true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['O'],
			},
			wantTet: &Tetrimino{
				Value: 'O',
				Cells: [][]bool{
					{true, true},
					{true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['O'],
			},
			wantRotationPoint: 1,
		},
		"O; starting rotation 2 (south); no valid rotation": {
			matrix: Matrix{
				{0},
			},
			tet: &Tetrimino{
				Value: 'O',
				Cells: [][]bool{
					{true, true},
					{true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['O'],
			},
			wantTet: nil,
		},

		"O; starting rotation 3 (west); rotation point 1": {
			matrix: Matrix{
				{0, 0},
				{0, 0},
			},
			tet: &Tetrimino{
				Value: 'O',
				Cells: [][]bool{
					{true, true},
					{true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['O'],
			},
			wantTet: &Tetrimino{
				Value: 'O',
				Cells: [][]bool{
					{true, true},
					{true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['O'],
			},
			wantRotationPoint: 1,
		},
		"O; starting rotation 3 (west); no valid rotation": {
			matrix: Matrix{
				{0},
			},
			tet: &Tetrimino{
				Value: 'O',
				Cells: [][]bool{
					{true, true},
					{true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['O'],
			},
			wantTet: nil,
		},

		"6; starting rotation 0 (north); rotation point 1": {
			matrix: Matrix{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{false, true, false},
					{true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, false},
					{true, true},
					{true, false},
				},
				Position:         Coordinate{X: 1, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['6'],
			},
			wantRotationPoint: 1,
		},
		"6; starting rotation 0 (north); rotation point 2": {
			matrix: Matrix{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{false, true, false},
					{true, true, true},
				},
				Position:         Coordinate{X: 1, Y: 0},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, false},
					{true, true},
					{true, false},
				},
				Position:         Coordinate{X: 1, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['6'],
			},
			wantRotationPoint: 2,
		},
		"6; starting rotation 0 (north); rotation point 3": {
			matrix: Matrix{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{false, true, false},
					{true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 1},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, false},
					{true, true},
					{true, false},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['6'],
			},
			wantRotationPoint: 3,
		},
		"6; starting rotation 0 (north); rotation point 4": {
			matrix: Matrix{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{false, true, false},
					{true, true, true},
				},
				Position:         Coordinate{X: -1, Y: -2},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, false},
					{true, true},
					{true, false},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['6'],
			},
			wantRotationPoint: 4,
		},
		"6; starting rotation 0 (north); rotation point 5": {
			matrix: Matrix{
				{0, 0},
				{0, 0},
				{0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{false, true, false},
					{true, true, true},
				},
				Position:         Coordinate{X: 0, Y: -2},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, false},
					{true, true},
					{true, false},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['6'],
			},
			wantRotationPoint: 5,
		},
		"6; starting rotation 0 (north); no valid rotation": {
			matrix: Matrix{
				{0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{false, true, false},
					{true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: nil,
		},

		"6; starting rotation 1 (east); rotation point 1": {
			matrix: Matrix{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, false},
					{true, true},
					{true, false},
				},
				Position:         Coordinate{X: 1, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Position:         Coordinate{X: 0, Y: 1},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['6'],
			},
			wantRotationPoint: 1,
		},
		"6; starting rotation 1 (east); rotation point 2": {
			matrix: Matrix{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, false},
					{true, true},
					{true, false},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Position:         Coordinate{X: 0, Y: 1},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['6'],
			},
			wantRotationPoint: 2,
		},
		"6; starting rotation 1 (east); rotation point 3": {
			matrix: Matrix{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, false},
					{true, true},
					{true, false},
				},
				Position:         Coordinate{X: 0, Y: -2},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['6'],
			},
			wantRotationPoint: 3,
		},
		"6; starting rotation 1 (east); rotation point 4": {
			matrix: Matrix{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, false},
					{true, true},
					{true, false},
				},
				Position:         Coordinate{X: 1, Y: 1},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['6'],
			},
			wantRotationPoint: 4,
		},
		"6; starting rotation 1 (east); rotation point 5": {
			matrix: Matrix{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, false},
					{true, true},
					{true, false},
				},
				Position:         Coordinate{X: 0, Y: 1},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['6'],
			},
			wantRotationPoint: 5,
		},
		"6; starting rotation 1 (east); no valid rotation": {
			matrix: Matrix{
				{0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, false},
					{true, true},
					{true, false},
				},
				Position:         Coordinate{X: 0, Y: 1},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: nil,
		},

		"6; starting rotation 2 (south); rotation point 1": {
			matrix: Matrix{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Position:         Coordinate{X: 0, Y: 1},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{false, true},
					{true, true},
					{false, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['6'],
			},
			wantRotationPoint: 1,
		},
		"6; starting rotation 2 (south); rotation point 2": {
			matrix: Matrix{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Position:         Coordinate{X: -1, Y: 1},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{false, true},
					{true, true},
					{false, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['6'],
			},
			wantRotationPoint: 2,
		},
		"6; starting rotation 2 (south); rotation point 3": {
			matrix: Matrix{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Position:         Coordinate{X: -1, Y: 2},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{false, true},
					{true, true},
					{false, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['6'],
			},
			wantRotationPoint: 3,
		},
		"6; starting rotation 2 (south); rotation point 4": {
			matrix: Matrix{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Position:         Coordinate{X: 0, Y: -1},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{false, true},
					{true, true},
					{false, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['6'],
			},
			wantRotationPoint: 4,
		},
		"6; starting rotation 2 (south); rotation point 5": {
			matrix: Matrix{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Position:         Coordinate{X: -1, Y: -1},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{false, true},
					{true, true},
					{false, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['6'],
			},
			wantRotationPoint: 5,
		},
		"6; starting rotation 2 (south); no valid rotation": {
			matrix: Matrix{
				{0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: nil,
		},

		"6; starting rotation 3 (west); rotation point 1": {
			matrix: Matrix{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{false, true},
					{true, true},
					{false, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{false, true, false},
					{true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['6'],
			},
			wantRotationPoint: 1,
		},
		"6; starting rotation 3 (west); rotation point 2": {
			matrix: Matrix{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{false, true},
					{true, true},
					{false, true},
				},
				Position:         Coordinate{X: 1, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{false, true, false},
					{true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['6'],
			},
			wantRotationPoint: 2,
		},
		"6; starting rotation 3 (west); rotation point 3": {
			matrix: Matrix{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{false, true},
					{true, true},
					{false, true},
				},
				Position:         Coordinate{X: 1, Y: -1},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{false, true, false},
					{true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['6'],
			},
			wantRotationPoint: 3,
		},
		"6; starting rotation 3 (west); rotation point 4": {
			matrix: Matrix{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{false, true},
					{true, true},
					{false, true},
				},
				Position:         Coordinate{X: 0, Y: 2},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{false, true, false},
					{true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['6'],
			},
			wantRotationPoint: 4,
		},
		"6; starting rotation 3 (west); rotation point 5": {
			matrix: Matrix{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{false, true},
					{true, true},
					{false, true},
				},
				Position:         Coordinate{X: 1, Y: 2},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{false, true, false},
					{true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['6'],
			},
			wantRotationPoint: 5,
		},
		"6; starting rotation 3 (west); no valid rotation": {
			matrix: Matrix{
				{0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{false, true},
					{true, true},
					{false, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: nil,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			rotationPoint, err := tc.tet.rotateClockwise(tc.matrix)

			require.NoError(t, err)
			if tc.wantTet != nil {
				assert.Equal(t, tc.wantTet, tc.tet)
				assert.Equal(t, tc.wantRotationPoint, rotationPoint)
			} else {
				assert.Equal(t, invalidRotationPoint, rotationPoint)
			}
		})
	}
}

func TestTetrimino_rotateCounterClockwise(t *testing.T) {
	tt := map[string]struct {
		matrix            Matrix
		tet               *Tetrimino
		wantTet           *Tetrimino
		wantRotationPoint int
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
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 1},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Position:         Coordinate{X: 1, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['I'],
			},
			wantRotationPoint: 1,
		},
		"I; starting rotation 0 (north); rotation point 2": {
			matrix: Matrix{
				{0},
				{0},
				{0},
				{0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 1},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['I'],
			},
			wantRotationPoint: 2,
		},
		"I; starting rotation 0 (north); rotation point 3": {
			matrix: Matrix{
				{0},
				{0},
				{0},
				{0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: -3, Y: 1},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['I'],
			},
			wantRotationPoint: 3,
		},
		"I; starting rotation 0 (north); rotation point 4": {
			matrix: Matrix{
				{0},
				{0},
				{0},
				{0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 3},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['I'],
			},
			wantRotationPoint: 4,
		},
		"I; starting rotation 0 (north); rotation point 5": {
			matrix: Matrix{
				{0},
				{0},
				{0},
				{0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: -3, Y: 0},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['I'],
			},
			wantRotationPoint: 5,
		},
		"I; starting rotation 0 (north); no valid rotation": {
			matrix: Matrix{
				{0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: nil,
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
				Cells: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Position:         Coordinate{X: 2, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 1},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['I'],
			},
			wantRotationPoint: 1,
		},
		"I; starting rotation 1 (east); rotation point 2": {
			matrix: Matrix{
				{0, 0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Position:         Coordinate{X: 0, Y: -1},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['I'],
			},
			wantRotationPoint: 2,
		},
		"I; starting rotation 1 (east); rotation point 3": {
			matrix: Matrix{
				{0, 0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Position:         Coordinate{X: 3, Y: -1},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['I'],
			},
			wantRotationPoint: 3,
		},
		"I; starting rotation 1 (east); rotation point 4": {
			matrix: Matrix{
				{0, 0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['I'],
			},
			wantRotationPoint: 4,
		},
		"I; starting rotation 1 (east); rotation point 5": {
			matrix: Matrix{
				{0, 0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Position:         Coordinate{X: 3, Y: -3},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['I'],
			},
			wantRotationPoint: 5,
		},
		"I; starting rotation 1 (east); no valid rotation": {
			matrix: Matrix{
				{0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Position:         Coordinate{X: 2, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: nil,
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
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 2},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Position:         Coordinate{X: 2, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['I'],
			},
			wantRotationPoint: 1,
		},
		"I; starting rotation 2 (south); rotation point 2": {
			matrix: Matrix{
				{0},
				{0},
				{0},
				{0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: -3, Y: 2},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['I'],
			},
			wantRotationPoint: 2,
		},
		"I; starting rotation 2 (south); rotation point 3": {
			matrix: Matrix{
				{0},
				{0},
				{0},
				{0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 2},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['I'],
			},
			wantRotationPoint: 3,
		},
		"I; starting rotation 2 (south); rotation point 4": {
			matrix: Matrix{
				{0},
				{0},
				{0},
				{0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: -3, Y: 0},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['I'],
			},
			wantRotationPoint: 4,
		},
		"I; starting rotation 2 (south); rotation point 5": {
			matrix: Matrix{
				{0},
				{0},
				{0},
				{0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 3},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['I'],
			},
			wantRotationPoint: 5,
		},
		"I; starting rotation 2 (south); no valid rotation": {
			matrix: Matrix{
				{0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 2},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: nil,
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
				Cells: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Position:         Coordinate{X: 1, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 2},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['I'],
			},
			wantRotationPoint: 1,
		},
		"I; starting rotation 3 (west); rotation point 2": {
			matrix: Matrix{
				{0, 0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Position:         Coordinate{X: 3, Y: -2},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['I'],
			},
			wantRotationPoint: 2,
		},
		"I; starting rotation 3 (west); rotation point 3": {
			matrix: Matrix{
				{0, 0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Position:         Coordinate{X: 0, Y: -2},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['I'],
			},
			wantRotationPoint: 3,
		},
		"I; starting rotation 3 (west); rotation point 4": {
			matrix: Matrix{
				{0, 0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Position:         Coordinate{X: 3, Y: -3},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['I'],
			},
			wantRotationPoint: 4,
		},
		"I; starting rotation 3 (west); rotation point 5": {
			matrix: Matrix{
				{0, 0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true, true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['I'],
			},
			wantRotationPoint: 5,
		},
		"I; starting rotation 3 (west); no valid rotation": {
			matrix: Matrix{
				{0},
			},
			tet: &Tetrimino{
				Value: 'I',
				Cells: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Position:         Coordinate{X: 1, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['I'],
			},
			wantTet: nil,
		},

		"O; starting rotation 0 (north); rotation point 1": {
			matrix: Matrix{
				{0, 0},
				{0, 0},
			},
			tet: &Tetrimino{
				Value: 'O',
				Cells: [][]bool{
					{true, true},
					{true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['O'],
			},
			wantTet: &Tetrimino{
				Value: 'O',
				Cells: [][]bool{
					{true, true},
					{true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['O'],
			},
			wantRotationPoint: 1,
		},
		"O; starting rotation 0 (north); no valid rotation": {
			matrix: Matrix{
				{0},
			},
			tet: &Tetrimino{
				Value: 'O',
				Cells: [][]bool{
					{true, true},
					{true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['O'],
			},
			wantTet: nil,
		},

		"O; starting rotation 1 (east); rotation point 1": {
			matrix: Matrix{
				{0, 0},
				{0, 0},
			},
			tet: &Tetrimino{
				Value: 'O',
				Cells: [][]bool{
					{true, true},
					{true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['O'],
			},
			wantTet: &Tetrimino{
				Value: 'O',
				Cells: [][]bool{
					{true, true},
					{true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['O'],
			},
			wantRotationPoint: 1,
		},
		"O; starting rotation 1 (east); no valid rotation": {
			matrix: Matrix{
				{0},
			},
			tet: &Tetrimino{
				Value: 'O',
				Cells: [][]bool{
					{true, true},
					{true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['O'],
			},
			wantTet: nil,
		},

		"O; starting rotation 2 (south); rotation point 1": {
			matrix: Matrix{
				{0, 0},
				{0, 0},
			},
			tet: &Tetrimino{
				Value: 'O',
				Cells: [][]bool{
					{true, true},
					{true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['O'],
			},
			wantTet: &Tetrimino{
				Value: 'O',
				Cells: [][]bool{
					{true, true},
					{true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['O'],
			},
			wantRotationPoint: 1,
		},
		"O; starting rotation 2 (south); no valid rotation": {
			matrix: Matrix{
				{0},
			},
			tet: &Tetrimino{
				Value: 'O',
				Cells: [][]bool{
					{true, true},
					{true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['O'],
			},
			wantTet: nil,
		},

		"O; starting rotation 3 (west); rotation point 1": {
			matrix: Matrix{
				{0, 0},
				{0, 0},
			},
			tet: &Tetrimino{
				Value: 'O',
				Cells: [][]bool{
					{true, true},
					{true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['O'],
			},
			wantTet: &Tetrimino{
				Value: 'O',
				Cells: [][]bool{
					{true, true},
					{true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['O'],
			},
			wantRotationPoint: 1,
		},
		"O; starting rotation 3 (west); no valid rotation": {
			matrix: Matrix{
				{0},
			},
			tet: &Tetrimino{
				Value: 'O',
				Cells: [][]bool{
					{true, true},
					{true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['O'],
			},
			wantTet: nil,
		},

		"6; starting rotation 0 (north); rotation point 1": {
			matrix: Matrix{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{false, true, false},
					{true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{false, true},
					{true, true},
					{false, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['6'],
			},
			wantRotationPoint: 1,
		},
		"6; starting rotation 0 (north); rotation point 2": {
			matrix: Matrix{
				{0, 0},
				{0, 0},
				{0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{false, true, false},
					{true, true, true},
				},
				Position:         Coordinate{X: -1, Y: 0},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{false, true},
					{true, true},
					{false, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['6'],
			},
			wantRotationPoint: 2,
		},
		"6; starting rotation 0 (north); rotation point 3": {
			matrix: Matrix{
				{0, 0},
				{0, 0},
				{0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{false, true, false},
					{true, true, true},
				},
				Position:         Coordinate{X: -1, Y: 1},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{false, true},
					{true, true},
					{false, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['6'],
			},
			wantRotationPoint: 3,
		},
		"6; starting rotation 0 (north); rotation point 4": {
			matrix: Matrix{
				{0, 0},
				{0, 0},
				{0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{false, true, false},
					{true, true, true},
				},
				Position:         Coordinate{X: 0, Y: -2},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{false, true},
					{true, true},
					{false, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['6'],
			},
			wantRotationPoint: 4,
		},
		"6; starting rotation 0 (north); rotation point 5": {
			matrix: Matrix{
				{0, 0},
				{0, 0},
				{0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{false, true, false},
					{true, true, true},
				},
				Position:         Coordinate{X: -1, Y: -2},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{false, true},
					{true, true},
					{false, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['6'],
			},
			wantRotationPoint: 5,
		},
		"6; starting rotation 0 (north); no valid rotation": {
			matrix: Matrix{
				{0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{false, true, false},
					{true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: nil,
		},

		"6; starting rotation 1 (east); rotation point 1": {
			matrix: Matrix{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, false},
					{true, true},
					{true, false},
				},
				Position:         Coordinate{X: 1, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{false, true, false},
					{true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['6'],
			},
			wantRotationPoint: 1,
		},
		"6; starting rotation 1 (east); rotation point 2": {
			matrix: Matrix{
				{0, 0, 0},
				{0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, false},
					{true, true},
					{true, false},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{false, true, false},
					{true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['6'],
			},
			wantRotationPoint: 2,
		},
		"6; starting rotation 1 (east); rotation point 3": {
			matrix: Matrix{
				{0, 0, 0},
				{0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, false},
					{true, true},
					{true, false},
				},
				Position:         Coordinate{X: 0, Y: -1},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{false, true, false},
					{true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['6'],
			},
			wantRotationPoint: 3,
		},
		"6; starting rotation 1 (east); rotation point 4": {
			matrix: Matrix{
				{0, 0, 0},
				{0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, false},
					{true, true},
					{true, false},
				},
				Position:         Coordinate{X: 1, Y: 2},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{false, true, false},
					{true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['6'],
			},
			wantRotationPoint: 4,
		},
		"6; starting rotation 1 (east); rotation point 5": {
			matrix: Matrix{
				{0, 0, 0},
				{0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, false},
					{true, true},
					{true, false},
				},
				Position:         Coordinate{X: 0, Y: 2},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{false, true, false},
					{true, true, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 0,
				RotationCompass:  RotationCompasses['6'],
			},
			wantRotationPoint: 5,
		},
		"6; starting rotation 1 (east); no valid rotation": {
			matrix: Matrix{
				{0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, false},
					{true, true},
					{true, false},
				},
				Position:         Coordinate{X: 1, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: nil,
		},

		"6; starting rotation 2 (south); rotation point 1": {
			matrix: Matrix{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Position:         Coordinate{X: 0, Y: 1},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, false},
					{true, true},
					{true, false},
				},
				Position:         Coordinate{X: 1, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['6'],
			},
			wantRotationPoint: 1,
		},
		"6; starting rotation 2 (south); rotation point 2": {
			matrix: Matrix{
				{0, 0},
				{0, 0},
				{0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Position:         Coordinate{X: 0, Y: 1},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, false},
					{true, true},
					{true, false},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['6'],
			},
			wantRotationPoint: 2,
		},
		"6; starting rotation 2 (south); rotation point 3": {
			matrix: Matrix{
				{0, 0},
				{0, 0},
				{0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Position:         Coordinate{X: 0, Y: 2},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, false},
					{true, true},
					{true, false},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['6'],
			},
			wantRotationPoint: 3,
		},
		"6; starting rotation 2 (south); rotation point 4": {
			matrix: Matrix{
				{0, 0},
				{0, 0},
				{0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Position:         Coordinate{X: -1, Y: -1},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, false},
					{true, true},
					{true, false},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['6'],
			},
			wantRotationPoint: 4,
		},
		"6; starting rotation 2 (south); rotation point 5": {
			matrix: Matrix{
				{0, 0},
				{0, 0},
				{0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Position:         Coordinate{X: 0, Y: -1},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, false},
					{true, true},
					{true, false},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 1,
				RotationCompass:  RotationCompasses['6'],
			},
			wantRotationPoint: 5,
		},
		"6; starting rotation 2 (south); no valid rotation": {
			matrix: Matrix{
				{0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Position:         Coordinate{X: 0, Y: 1},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: nil,
		},

		"6; starting rotation 3 (west); rotation point 1": {
			matrix: Matrix{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{false, true},
					{true, true},
					{false, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Position:         Coordinate{X: 0, Y: 1},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['6'],
			},
			wantRotationPoint: 1,
		},
		"6; starting rotation 3 (west); rotation point 2": {
			matrix: Matrix{
				{0, 0, 0},
				{0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{false, true},
					{true, true},
					{false, true},
				},
				Position:         Coordinate{X: 1, Y: -1},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['6'],
			},
			wantRotationPoint: 2,
		},
		"6; starting rotation 3 (west); rotation point 3": {
			matrix: Matrix{
				{0, 0, 0},
				{0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{false, true},
					{true, true},
					{false, true},
				},
				Position:         Coordinate{X: 1, Y: -2},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['6'],
			},
			wantRotationPoint: 3,
		},
		"6; starting rotation 3 (west); rotation point 4": {
			matrix: Matrix{
				{0, 0, 0},
				{0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{false, true},
					{true, true},
					{false, true},
				},
				Position:         Coordinate{X: 0, Y: 1},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['6'],
			},
			wantRotationPoint: 4,
		},
		"6; starting rotation 3 (west); rotation point 5": {
			matrix: Matrix{
				{0, 0, 0},
				{0, 0, 0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{false, true},
					{true, true},
					{false, true},
				},
				Position:         Coordinate{X: 1, Y: 1},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 2,
				RotationCompass:  RotationCompasses['6'],
			},
			wantRotationPoint: 5,
		},
		"6; starting rotation 3 (west); no valid rotation": {
			matrix: Matrix{
				{0},
			},
			tet: &Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{false, true},
					{true, true},
					{false, true},
				},
				Position:         Coordinate{X: 0, Y: 0},
				CompassDirection: 3,
				RotationCompass:  RotationCompasses['6'],
			},
			wantTet: nil,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			rotationPoint, err := tc.tet.rotateCounterClockwise(tc.matrix)

			require.NoError(t, err)
			if tc.wantTet != nil {
				assert.Equal(t, tc.wantTet, tc.tet)
				assert.Equal(t, tc.wantRotationPoint, rotationPoint)
			} else {
				assert.Equal(t, invalidRotationPoint, rotationPoint)
			}
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
				Cells: [][]bool{
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
				Cells: [][]bool{
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
				Cells: [][]bool{
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

			assert.ElementsMatch(t, tc.tet.Cells, tc.wantMinos)
		})
	}
}

func TestTetrimino_IsValid(t *testing.T) {
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
				Cells:    [][]bool{{true}},
				Position: Coordinate{X: 0, Y: 0},
			},
			want: true,
		},
		"true; perfect fit": {
			matrix: Matrix{
				{0, 0, 0},
				{'X', 0, 'X'},
			},
			tet: &Tetrimino{
				Cells: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Position: Coordinate{X: 0, Y: 0},
			},
			want: true,
		},
		"false; blocking mino": {
			matrix: Matrix{
				{'#'},
			},
			tet: &Tetrimino{
				Cells:    [][]bool{{true}},
				Position: Coordinate{X: 0, Y: 0},
			},
			want: false,
		},
		"false; out of bounds left": {
			Matrix{{0}},
			&Tetrimino{
				Cells:    [][]bool{{true}},
				Position: Coordinate{X: -1, Y: 0},
			},
			false,
		},
		"false; out of bounds right": {
			Matrix{{0}},
			&Tetrimino{
				Cells:    [][]bool{{true}},
				Position: Coordinate{X: 1, Y: 0},
			},
			false,
		},
		"false; out of bounds up": {
			Matrix{{0}},
			&Tetrimino{
				Cells:    [][]bool{{true}},
				Position: Coordinate{X: 0, Y: -1},
			},
			false,
		},
		"false; out of bounds down": {
			Matrix{{0}},
			&Tetrimino{
				Cells:    [][]bool{{true}},
				Position: Coordinate{X: 0, Y: 1},
			},
			false,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			result := tc.tet.IsValid(tc.matrix, true)

			assert.Equal(t, tc.want, result)
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
				require.EqualError(t, err, tc.wantErr.Error())
				assert.Equal(t, tc.want, result)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_deepCopyCells(t *testing.T) {
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
			cellsCopy := deepCopyCells(tc.cells)
			tc.modify(&cellsCopy)

			assert.False(t, tc.cells[0][0])
		})
	}
}

func TestTetrimino_DeepCopy(t *testing.T) {
	tet := Tetrimino{
		Value:            byte(0),
		Cells:            [][]bool{{false}},
		Position:         Coordinate{0, 0},
		CompassDirection: 0,
		RotationCompass: RotationCompass{
			{&Coordinate{0, 0}},
		},
	}

	// Create a copy manually
	manualCopy := Tetrimino{
		Value:            byte(0),
		Cells:            [][]bool{{false}},
		Position:         Coordinate{0, 0},
		CompassDirection: 0,
		RotationCompass:  RotationCompass{{&Coordinate{0, 0}}},
	}

	// Create a (dereferences) copy with the helper function
	easyCopy := *tet.DeepCopy()

	// Assert that all are equal
	assert.Equal(t, tet, manualCopy)
	assert.Equal(t, tet, easyCopy)
	assert.Equal(t, manualCopy, easyCopy)

	// Modify the original
	tet.Value = 1
	tet.Cells = [][]bool{{true}}
	tet.Position = Coordinate{1, 1}
	tet.CompassDirection = 1
	tet.RotationCompass[0][0].X = 1

	// Assert that the original changed but both copies stayed the same
	assert.NotEqual(t, tet, manualCopy)
	assert.NotEqual(t, tet, easyCopy)
	assert.Equal(t, manualCopy, easyCopy)
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
				Position: Coordinate{X: 0, Y: 9},
				Cells:    [][]bool{{true}},
			},
			want: true,
		},
		"true; skyline 20": {
			skyline: 20,
			tet: &Tetrimino{
				Position: Coordinate{X: 0, Y: 19},
				Cells:    [][]bool{{true}},
			},
			want: true,
		},
		"false; skyline 10": {
			skyline: 10,
			tet: &Tetrimino{
				Position: Coordinate{X: 0, Y: 10},
				Cells:    [][]bool{{true}},
			},
			want: false,
		},
		"false; skyline 20": {
			skyline: 20,
			tet: &Tetrimino{
				Position: Coordinate{X: 0, Y: 20},
				Cells:    [][]bool{{true}},
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
