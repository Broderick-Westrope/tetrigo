package tetris

import (
	"fmt"
	"reflect"
	"testing"
)

func TestTetrimino_MoveDown(t *testing.T) {
	rows := []int{0, 1, 2, 3, 4, 5, 6, 7, 8}
	cols := []int{0, 1, 2, 3, 4, 5, 6, 7}
	for _, row := range rows {
		for _, col := range cols {
			coord := Coordinate{X: col, Y: row}
			for _, tet := range Tetriminos {
				t.Run(fmt.Sprintf("value: %s, coord: %v", string(tet.Value), coord), func(t *testing.T) {
					m := &Matrix{}
					tet.Pos = coord
					m.AddTetrimino(&tet)

					tet.MoveDown(m)

					if tet.Pos.Y == coord.Y+1 {
						return
					}
					if tet.Pos.Y+len(tet.Minos)-1 == coord.Y {
						return
					}
					t.Errorf("got %v, want %v", tet.Pos.Y, coord.Y+1)
				})
			}
		}
	}
}

func TestTetrimino_MoveLeft(t *testing.T) {
	tt := []struct {
		name              string
		startingPlayfield Matrix
		startingTet       Tetrimino
		expectedPlayfield Matrix
		expectsErr        bool
	}{
		{
			name: "can, empty matrix",
			startingPlayfield: Matrix{
				{0, 'T', 'T', 'T'},
				{0, 0, 'T', 0},
			},
			startingTet: Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Pos: Coordinate{X: 1, Y: 0},
			},
			expectedPlayfield: Matrix{
				{'T', 'T', 'T'},
				{0, 'T', 0},
			},
			expectsErr: false,
		},
		{
			name: "can, perfect fit",
			startingPlayfield: Matrix{
				{'#', 0, 'T', 'T', 'T'},
				{'#', '#', 0, 'T', 0},
			},
			startingTet: Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Pos: Coordinate{X: 2, Y: 0},
			},
			expectedPlayfield: Matrix{
				{'#', 'T', 'T', 'T'},
				{'#', '#', 'T', 0},
			},
			expectsErr: false,
		},
		{
			name: "can, ghost cells",
			startingPlayfield: Matrix{
				{'#', 'G', 'T', 'T', 'T'},
				{'#', '#', 'G', 'T', 0},
			},
			startingTet: Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Pos: Coordinate{X: 2, Y: 0},
			},
			expectedPlayfield: Matrix{
				{'#', 'T', 'T', 'T'},
				{'#', '#', 'T', 0},
			},
			expectsErr: false,
		},
		{
			name: "cannot, blocking tetrimino",
			startingPlayfield: Matrix{
				{'#', 'T', 'T', 'T'},
				{'#', 0, 'T', 0},
			},
			startingTet: Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Pos: Coordinate{X: 1, Y: 0},
			},
			expectedPlayfield: Matrix{
				{'#', 'T', 'T', 'T'},
				{'#', 0, 'T', 0},
			},
			expectsErr: false,
		},
		{
			name: "cannot, end of matrix",
			startingPlayfield: Matrix{
				{'T', 'T', 'T'},
				{0, 'T', 0},
			},
			startingTet: Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Pos: Coordinate{X: 0, Y: 0},
			},
			expectedPlayfield: Matrix{
				{'T', 'T', 'T'},
				{0, 'T', 0},
			},
			expectsErr: false,
		},
		{
			name: "error, wrong value",
			startingPlayfield: Matrix{
				{0, 'X', 'X', 'X'},
				{0, 0, 'X', 0},
			},
			startingTet: Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Pos: Coordinate{X: 1, Y: 0},
			},
			expectedPlayfield: Matrix{},
			expectsErr:        true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.startingTet.MoveLeft(&tc.startingPlayfield)

			if tc.expectsErr && err == nil {
				t.Errorf("expected error, got nil")
			} else if !tc.expectsErr && err != nil {
				t.Errorf("expected nil, got error")
			}

			if err == nil && tc.startingPlayfield != tc.expectedPlayfield {
				t.Errorf("expected matrix %v, got %v", tc.expectedPlayfield, tc.startingPlayfield)
			}
		})
	}
}

func TestTetrimino_MoveRight(t *testing.T) {
	tt := []struct {
		name              string
		startingPlayfield Matrix
		startingTet       Tetrimino
		expectedPlayfield Matrix
		expectsErr        bool
	}{
		{
			name: "can, empty matrix",
			startingPlayfield: Matrix{
				{'T', 'T', 'T'},
				{0, 'T', 0},
			},
			startingTet: Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Pos: Coordinate{X: 0, Y: 0},
			},
			expectedPlayfield: Matrix{
				{0, 'T', 'T', 'T'},
				{0, 0, 'T', 0},
			},
			expectsErr: false,
		},
		{
			name: "can, perfect fit",
			startingPlayfield: Matrix{
				{'T', 'T', 'T', 0, '#'},
				{0, 'T', 0, '#'},
			},
			startingTet: Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Pos: Coordinate{X: 0, Y: 0},
			},
			expectedPlayfield: Matrix{
				{0, 'T', 'T', 'T', '#'},
				{0, 0, 'T', '#'},
			},
			expectsErr: false,
		},
		{
			name: "can, ghost cells",
			startingPlayfield: Matrix{
				{'T', 'T', 'T', 'G', '#'},
				{0, 'T', 'G', '#'},
			},
			startingTet: Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Pos: Coordinate{X: 0, Y: 0},
			},
			expectedPlayfield: Matrix{
				{0, 'T', 'T', 'T', '#'},
				{0, 0, 'T', '#'},
			},
			expectsErr: false,
		},
		{
			name: "cannot, blocking tetrimino",
			startingPlayfield: Matrix{
				{'T', 'T', 'T', '#'},
				{0, 'T', '#'},
			},
			startingTet: Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Pos: Coordinate{X: 0, Y: 0},
			},
			expectedPlayfield: Matrix{
				{'T', 'T', 'T', '#'},
				{0, 'T', '#'},
			},
			expectsErr: false,
		},
		{
			name: "cannot, end of matrix",
			startingPlayfield: Matrix{
				{0, 0, 0, 0, 0, 0, 0, 'T', 'T', 'T'},
				{0, 0, 0, 0, 0, 0, 0, 0, 'T', 0},
			},
			startingTet: Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Pos: Coordinate{X: 7, Y: 0},
			},
			expectedPlayfield: Matrix{
				{0, 0, 0, 0, 0, 0, 0, 'T', 'T', 'T'},
				{0, 0, 0, 0, 0, 0, 0, 0, 'T', 0},
			},
			expectsErr: false,
		},
		{
			name: "error, wrong value",
			startingPlayfield: Matrix{
				{'X', 'X', 'X'},
				{0, 'X', 0},
			},
			startingTet: Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Pos: Coordinate{X: 0, Y: 0},
			},
			expectedPlayfield: Matrix{},
			expectsErr:        true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.startingTet.MoveRight(&tc.startingPlayfield)

			if tc.expectsErr && err == nil {
				t.Errorf("expected error, got nil")
			} else if !tc.expectsErr && err != nil {
				t.Errorf("expected nil, got error")
			}

			if err == nil && tc.startingPlayfield != tc.expectedPlayfield {
				t.Errorf("expected matrix %v, got %v", tc.expectedPlayfield, tc.startingPlayfield)
			}
		})
	}
}

func TestTetrimino_CanMoveDown(t *testing.T) {
	tt := []struct {
		name     string
		tet      *Tetrimino
		matrix   *Matrix
		expected bool
	}{
		{
			name: "I, can, top of matrix",
			tet: &Tetrimino{
				Value:           'I',
				Minos:           [][]bool{{true}, {true}, {true}, {true}},
				Pos:             Coordinate{X: 0, Y: 0},
				CurrentRotation: 1,
				RotationCoords:  RotationCoords['I'],
			},
			matrix:   &Matrix{},
			expected: true,
		},
		{
			name: "Z, can, starting position",
			tet: &Tetrimino{
				Value: 'Z',
				Minos: [][]bool{
					{true, true, false},
					{false, true, true},
				},
				Pos:             Coordinate{X: 3, Y: 18},
				CurrentRotation: 0,
				RotationCoords:  RotationCoords['6'],
			},
			matrix:   &Matrix{},
			expected: true,
		},
		{
			name: "S, cannot, blocking tetrimino",
			tet: &Tetrimino{
				Value: 'S',
				Minos: [][]bool{
					{false, true, true},
					{true, true, false},
				},
				Pos:             Coordinate{X: 0, Y: 0},
				CurrentRotation: 0,
				RotationCoords:  RotationCoords['6'],
			},
			matrix: &Matrix{
				{}, {}, {'X'},
			},
			expected: false,
		},
		{
			name: "O, cannot, bottom of matrix",
			tet: &Tetrimino{
				Value: 'O',
				Minos: [][]bool{
					{true, true},
					{true, true},
				},
				Pos:             Coordinate{X: 0, Y: 38},
				CurrentRotation: 0,
				RotationCoords:  RotationCoords['O'],
			},
			matrix:   &Matrix{},
			expected: false,
		},
		{
			name: "GitHub Issue #1",
			tet: &Tetrimino{
				Value: 'J',
				Minos: [][]bool{
					{true, true},
					{false, true},
					{false, true},
				},
				Pos:             Coordinate{X: 0, Y: 0},
				CurrentRotation: 3,
				RotationCoords:  RotationCoords['J'],
			},
			matrix: &Matrix{
				{}, {'X'},
			},
			expected: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.matrix.AddTetrimino(tc.tet)
			if err != nil {
				t.Errorf("Failed to add tetrimino: got error, want nil")
			}

			result := tc.tet.CanMoveDown(*tc.matrix)

			if result != tc.expected {
				t.Errorf("got %v, want %v", result, tc.expected)
			}
		})
	}
}

func TestTetrimino_CanMoveLeft(t *testing.T) {
	tt := []struct {
		name     string
		tet      *Tetrimino
		matrix   *Matrix
		expected bool
	}{
		{
			name: "Z, can, starting position",
			tet: &Tetrimino{
				Value: 'Z',
				Minos: [][]bool{
					{true, true, false},
					{false, true, true},
				},
				Pos:             Coordinate{X: 3, Y: 18},
				CurrentRotation: 0,
				RotationCoords:  RotationCoords['6'],
			},
			matrix:   &Matrix{},
			expected: true,
		},
		{
			name: "O, cannot, blocking tetrimino",
			tet: &Tetrimino{
				Value: 'O',
				Minos: [][]bool{
					{true, true},
					{true, true},
				},
				Pos:             Coordinate{X: 1, Y: 0},
				CurrentRotation: 0,
				RotationCoords:  RotationCoords['O'],
			},
			matrix: &Matrix{
				{'X'},
			},
			expected: false,
		},
		{
			name: "O, cannot, end of matrix",
			tet: &Tetrimino{
				Value: 'O',
				Minos: [][]bool{
					{true, true},
					{true, true},
				},
				Pos:             Coordinate{X: 0, Y: 38},
				CurrentRotation: 0,
				RotationCoords:  RotationCoords['O'],
			},
			matrix:   &Matrix{},
			expected: false,
		},
		{
			name: "GitHub Issue #1",
			tet: &Tetrimino{
				Value: 'L',
				Minos: [][]bool{
					{false, false, true},
					{true, true, true},
				},
				Pos:             Coordinate{X: 1, Y: 0},
				CurrentRotation: 3,
				RotationCoords:  RotationCoords['J'],
			},
			matrix: &Matrix{
				{'X', 'X', 'X'},
			},
			expected: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.matrix.AddTetrimino(tc.tet)
			if err != nil {
				t.Errorf("Failed to add tetrimino: got error, want nil")
			}

			result := tc.tet.canMoveLeft(*tc.matrix)

			if result != tc.expected {
				t.Errorf("got %v, want %v", result, tc.expected)
			}
		})
	}
}

func TestTetrimino_CanMoveRight(t *testing.T) {
	tt := []struct {
		name     string
		tet      *Tetrimino
		matrix   *Matrix
		expected bool
	}{
		{
			name: "Z, can, starting position",
			tet: &Tetrimino{
				Value: 'Z',
				Minos: [][]bool{
					{true, true, false},
					{false, true, true},
				},
				Pos:             Coordinate{X: 3, Y: 18},
				CurrentRotation: 0,
				RotationCoords:  RotationCoords['6'],
			},
			matrix:   &Matrix{},
			expected: true,
		},
		{
			name: "O, cannot, blocking tetrimino",
			tet: &Tetrimino{
				Value: 'O',
				Minos: [][]bool{
					{true, true},
					{true, true},
				},
				Pos:             Coordinate{X: 0, Y: 0},
				CurrentRotation: 0,
				RotationCoords:  RotationCoords['O'],
			},
			matrix: &Matrix{
				{0, 0, 'X'},
			},
			expected: false,
		},
		{
			name: "O, cannot, end of matrix",
			tet: &Tetrimino{
				Value: 'O',
				Minos: [][]bool{
					{true, true},
					{true, true},
				},
				Pos:             Coordinate{X: 8, Y: 38},
				CurrentRotation: 0,
				RotationCoords:  RotationCoords['O'],
			},
			matrix:   &Matrix{},
			expected: false,
		},
		{
			name: "GitHub Issue #1",
			tet: &Tetrimino{
				Value: 'J',
				Minos: [][]bool{
					{true, false, false},
					{true, true, true},
				},
				Pos:             Coordinate{X: 0, Y: 0},
				CurrentRotation: 0,
				RotationCoords:  RotationCoords['J'],
			},
			matrix: &Matrix{
				{0, 'X', 'X', 'X'},
			},
			expected: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.matrix.AddTetrimino(tc.tet)
			if err != nil {
				t.Errorf("Failed to add tetrimino: got error, want nil")
			}

			result := tc.tet.canMoveRight(*tc.matrix)

			if result != tc.expected {
				t.Errorf("got %v, want %v", result, tc.expected)
			}
		})
	}
}

func TestTetrimino_RotateClockwise(t *testing.T) {
	tt := []struct {
		name        string
		tet         *Tetrimino
		expectedTet *Tetrimino
	}{
		{
			name: "I, rotation 0",
			tet: &Tetrimino{
				Value: 'I',
				Minos: [][]bool{
					{true, true, true, true},
				},
				Pos:             Coordinate{X: 0, Y: 0},
				CurrentRotation: 0,
				RotationCoords:  RotationCoords['I'],
			},
			expectedTet: &Tetrimino{
				Value: 'I',
				Minos: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Pos:             Coordinate{X: 2, Y: -1},
				CurrentRotation: 1,
				RotationCoords:  RotationCoords['I'],
			},
		},
		{
			name: "I, rotation 1",
			tet: &Tetrimino{
				Value: 'I',
				Minos: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Pos:             Coordinate{X: 0, Y: 0},
				CurrentRotation: 1,
				RotationCoords:  RotationCoords['I'],
			},
			expectedTet: &Tetrimino{
				Value: 'I',
				Minos: [][]bool{
					{true, true, true, true},
				},
				Pos:             Coordinate{X: -2, Y: 2},
				CurrentRotation: 2,
				RotationCoords:  RotationCoords['I'],
			},
		},
		{
			name: "I, rotation 2",
			tet: &Tetrimino{
				Value: 'I',
				Minos: [][]bool{
					{true, true, true, true},
				},
				Pos:             Coordinate{X: 0, Y: 0},
				CurrentRotation: 2,
				RotationCoords:  RotationCoords['I'],
			},
			expectedTet: &Tetrimino{
				Value: 'I',
				Minos: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Pos:             Coordinate{X: 1, Y: -2},
				CurrentRotation: 3,
				RotationCoords:  RotationCoords['I'],
			},
		},
		{
			name: "I, rotation 3",
			tet: &Tetrimino{
				Value: 'I',
				Minos: [][]bool{
					{true, true, true, true},
				},
				Pos:             Coordinate{X: 0, Y: 0},
				CurrentRotation: 3,
				RotationCoords:  RotationCoords['I'],
			},
			expectedTet: &Tetrimino{
				Value: 'I',
				Minos: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Pos:             Coordinate{X: -1, Y: 1},
				CurrentRotation: 0,
				RotationCoords:  RotationCoords['I'],
			},
		},
		{
			name: "O, rotation 0",
			tet: &Tetrimino{
				Value: 'O',
				Minos: [][]bool{
					{true, true},
					{true, true},
				},
				Pos:             Coordinate{X: 0, Y: 0},
				CurrentRotation: 0,
				RotationCoords:  RotationCoords['O'],
			},
			expectedTet: &Tetrimino{
				Value: 'O',
				Minos: [][]bool{
					{true, true},
					{true, true},
				},
				Pos:             Coordinate{X: 0, Y: 0},
				CurrentRotation: 1,
				RotationCoords:  RotationCoords['O'],
			},
		},
		{
			name: "O, rotation 1",
			tet: &Tetrimino{
				Value: 'O',
				Minos: [][]bool{
					{true, true},
					{true, true},
				},
				Pos:             Coordinate{X: 0, Y: 0},
				CurrentRotation: 1,
				RotationCoords:  RotationCoords['O'],
			},
			expectedTet: &Tetrimino{
				Value: 'O',
				Minos: [][]bool{
					{true, true},
					{true, true},
				},
				Pos:             Coordinate{X: 0, Y: 0},
				CurrentRotation: 2,
				RotationCoords:  RotationCoords['O'],
			},
		},
		{
			name: "O, rotation 2",
			tet: &Tetrimino{
				Value: 'O',
				Minos: [][]bool{
					{true, true},
					{true, true},
				},
				Pos:             Coordinate{X: 0, Y: 0},
				CurrentRotation: 2,
				RotationCoords:  RotationCoords['O'],
			},
			expectedTet: &Tetrimino{
				Value: 'O',
				Minos: [][]bool{
					{true, true},
					{true, true},
				},
				Pos:             Coordinate{X: 0, Y: 0},
				CurrentRotation: 3,
				RotationCoords:  RotationCoords['O'],
			},
		},
		{
			name: "O, rotation 3",
			tet: &Tetrimino{
				Value: 'O',
				Minos: [][]bool{
					{true, true},
					{true, true},
				},
				Pos:             Coordinate{X: 0, Y: 0},
				CurrentRotation: 3,
				RotationCoords:  RotationCoords['O'],
			},
			expectedTet: &Tetrimino{
				Value: 'O',
				Minos: [][]bool{
					{true, true},
					{true, true},
				},
				Pos:             Coordinate{X: 0, Y: 0},
				CurrentRotation: 0,
				RotationCoords:  RotationCoords['O'],
			},
		},
		{
			name: "T, rotation 0",
			tet: &Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{false, true, false},
					{true, true, true},
				},
				Pos:             Coordinate{X: 0, Y: 0},
				CurrentRotation: 0,
				RotationCoords:  RotationCoords['6'],
			},
			expectedTet: &Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{true, false},
					{true, true},
					{true, false},
				},
				Pos:             Coordinate{X: 1, Y: 0},
				CurrentRotation: 1,
				RotationCoords:  RotationCoords['6'],
			},
		},
		{
			name: "T, rotation 1",
			tet: &Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{true, false},
					{true, true},
					{true, false},
				},
				Pos:             Coordinate{X: 0, Y: 0},
				CurrentRotation: 1,
				RotationCoords:  RotationCoords['6'],
			},
			expectedTet: &Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Pos:             Coordinate{X: -1, Y: 1},
				CurrentRotation: 2,
				RotationCoords:  RotationCoords['6'],
			},
		},
		{
			name: "T, rotation 2",
			tet: &Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Pos:             Coordinate{X: 0, Y: 0},
				CurrentRotation: 2,
				RotationCoords:  RotationCoords['6'],
			},
			expectedTet: &Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{false, true},
					{true, true},
					{false, true},
				},
				Pos:             Coordinate{X: 0, Y: -1},
				CurrentRotation: 3,
				RotationCoords:  RotationCoords['6'],
			},
		},
		{
			name: "T, rotation 3",
			tet: &Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{false, true},
					{true, true},
					{false, true},
				},
				Pos:             Coordinate{X: 0, Y: 0},
				CurrentRotation: 3,
				RotationCoords:  RotationCoords['6'],
			},
			expectedTet: &Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{false, true, false},
					{true, true, true},
				},
				Pos:             Coordinate{X: 0, Y: 0},
				CurrentRotation: 0,
				RotationCoords:  RotationCoords['6'],
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			newTet, err := tc.tet.rotateClockwise()

			if err != nil {
				t.Errorf("got error, want nil")
			}

			if !reflect.DeepEqual(newTet, tc.expectedTet) {
				t.Errorf("Tetrimino: got %v, want %v", newTet, tc.expectedTet)
			}
		})
	}
}

func TestTetrimino_RotateCounterClockwise(t *testing.T) {
	tt := []struct {
		name        string
		tet         *Tetrimino
		expectedTet *Tetrimino
	}{
		{
			name: "I, rotation 0",
			tet: &Tetrimino{
				Value: 'I',
				Minos: [][]bool{
					{true, true, true, true},
				},
				Pos:             Coordinate{X: 0, Y: 0},
				CurrentRotation: 0,
				RotationCoords:  RotationCoords['I'],
			},
			expectedTet: &Tetrimino{
				Value: 'I',
				Minos: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Pos:             Coordinate{X: 1, Y: -1},
				CurrentRotation: 3,
				RotationCoords:  RotationCoords['I'],
			},
		},
		{
			name: "I, rotation 1",
			tet: &Tetrimino{
				Value: 'I',
				Minos: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Pos:             Coordinate{X: 0, Y: 0},
				CurrentRotation: 1,
				RotationCoords:  RotationCoords['I'],
			},
			expectedTet: &Tetrimino{
				Value: 'I',
				Minos: [][]bool{
					{true, true, true, true},
				},
				Pos:             Coordinate{X: -2, Y: 1},
				CurrentRotation: 0,
				RotationCoords:  RotationCoords['I'],
			},
		},
		{
			name: "I, rotation 2",
			tet: &Tetrimino{
				Value: 'I',
				Minos: [][]bool{
					{true, true, true, true},
				},
				Pos:             Coordinate{X: 0, Y: 0},
				CurrentRotation: 2,
				RotationCoords:  RotationCoords['I'],
			},
			expectedTet: &Tetrimino{
				Value: 'I',
				Minos: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Pos:             Coordinate{X: 2, Y: -2},
				CurrentRotation: 1,
				RotationCoords:  RotationCoords['I'],
			},
		},
		{
			name: "I, rotation 3",
			tet: &Tetrimino{
				Value: 'I',
				Minos: [][]bool{
					{true},
					{true},
					{true},
					{true},
				},
				Pos:             Coordinate{X: 0, Y: 0},
				CurrentRotation: 3,
				RotationCoords:  RotationCoords['I'],
			},
			expectedTet: &Tetrimino{
				Value: 'I',
				Minos: [][]bool{
					{true, true, true, true},
				},
				Pos:             Coordinate{X: -1, Y: 2},
				CurrentRotation: 2,
				RotationCoords:  RotationCoords['I'],
			},
		},
		{
			name: "O, rotation 0",
			tet: &Tetrimino{
				Value: 'O',
				Minos: [][]bool{
					{true, true},
					{true, true},
				},
				Pos:             Coordinate{X: 0, Y: 0},
				CurrentRotation: 0,
				RotationCoords:  RotationCoords['O'],
			},
			expectedTet: &Tetrimino{
				Value: 'O',
				Minos: [][]bool{
					{true, true},
					{true, true},
				},
				Pos:             Coordinate{X: 0, Y: 0},
				CurrentRotation: 3,
				RotationCoords:  RotationCoords['O'],
			},
		},
		{
			name: "O, rotation 1",
			tet: &Tetrimino{
				Value: 'O',
				Minos: [][]bool{
					{true, true},
					{true, true},
				},
				Pos:             Coordinate{X: 0, Y: 0},
				CurrentRotation: 1,
				RotationCoords:  RotationCoords['O'],
			},
			expectedTet: &Tetrimino{
				Value: 'O',
				Minos: [][]bool{
					{true, true},
					{true, true},
				},
				Pos:             Coordinate{X: 0, Y: 0},
				CurrentRotation: 0,
				RotationCoords:  RotationCoords['O'],
			},
		},
		{
			name: "O, rotation 2",
			tet: &Tetrimino{
				Value: 'O',
				Minos: [][]bool{
					{true, true},
					{true, true},
				},
				Pos:             Coordinate{X: 0, Y: 0},
				CurrentRotation: 2,
				RotationCoords:  RotationCoords['O'],
			},
			expectedTet: &Tetrimino{
				Value: 'O',
				Minos: [][]bool{
					{true, true},
					{true, true},
				},
				Pos:             Coordinate{X: 0, Y: 0},
				CurrentRotation: 1,
				RotationCoords:  RotationCoords['O'],
			},
		},
		{
			name: "O, rotation 3",
			tet: &Tetrimino{
				Value: 'O',
				Minos: [][]bool{
					{true, true},
					{true, true},
				},
				Pos:             Coordinate{X: 0, Y: 0},
				CurrentRotation: 3,
				RotationCoords:  RotationCoords['O'],
			},
			expectedTet: &Tetrimino{
				Value: 'O',
				Minos: [][]bool{
					{true, true},
					{true, true},
				},
				Pos:             Coordinate{X: 0, Y: 0},
				CurrentRotation: 2,
				RotationCoords:  RotationCoords['O'],
			},
		},
		{
			name: "T, rotation 0",
			tet: &Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{false, true, false},
					{true, true, true},
				},
				Pos:             Coordinate{X: 0, Y: 0},
				CurrentRotation: 0,
				RotationCoords:  RotationCoords['6'],
			},
			expectedTet: &Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{false, true},
					{true, true},
					{false, true},
				},
				Pos:             Coordinate{X: 0, Y: 0},
				CurrentRotation: 3,
				RotationCoords:  RotationCoords['6'],
			},
		},
		{
			name: "T, rotation 1",
			tet: &Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{true, false},
					{true, true},
					{true, false},
				},
				Pos:             Coordinate{X: 0, Y: 0},
				CurrentRotation: 1,
				RotationCoords:  RotationCoords['6'],
			},
			expectedTet: &Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{false, true, false},
					{true, true, true},
				},
				Pos:             Coordinate{X: -1, Y: 0},
				CurrentRotation: 0,
				RotationCoords:  RotationCoords['6'],
			},
		},
		{
			name: "T, rotation 2",
			tet: &Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Pos:             Coordinate{X: 0, Y: 0},
				CurrentRotation: 2,
				RotationCoords:  RotationCoords['6'],
			},
			expectedTet: &Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{true, false},
					{true, true},
					{true, false},
				},
				Pos:             Coordinate{X: 1, Y: -1},
				CurrentRotation: 1,
				RotationCoords:  RotationCoords['6'],
			},
		},
		{
			name: "T, rotation 3",
			tet: &Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{false, true},
					{true, true},
					{false, true},
				},
				Pos:             Coordinate{X: 0, Y: 0},
				CurrentRotation: 3,
				RotationCoords:  RotationCoords['6'],
			},
			expectedTet: &Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Pos:             Coordinate{X: 0, Y: 1},
				CurrentRotation: 2,
				RotationCoords:  RotationCoords['6'],
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			newTet, err := tc.tet.rotateCounterClockwise()

			if err != nil {
				t.Errorf("got error, want nil")
			}

			if !reflect.DeepEqual(newTet, tc.expectedTet) {
				t.Errorf("Tetrimino: got %v, want %v", newTet, tc.expectedTet)

				if newTet.Value != tc.expectedTet.Value {
					t.Errorf("Tetrimino.Value: got %v, want %v", newTet.Value, tc.expectedTet.Value)
				}
				if !reflect.DeepEqual(newTet.Cells, tc.expectedTet.Minos) {
					t.Errorf("Tetrimino.Minos: got %v, want %v", newTet.Cells, tc.expectedTet.Minos)
				}
				if newTet.Pos != tc.expectedTet.Pos {
					t.Errorf("Tetrimino.Pos: got %v, want %v", newTet.Pos, tc.expectedTet.Pos)
				}
				if newTet.CurrentRotation != tc.expectedTet.CurrentRotation {
					t.Errorf("Tetrimino.CurrentRotation: got %v, want %v", newTet.CurrentRotation, tc.expectedTet.CurrentRotation)
				}
				if !reflect.DeepEqual(newTet.RotationCoords, tc.expectedTet.RotationCoords) {
					t.Errorf("Tetrimino.RotationCoords: got %v, want %v", newTet.RotationCoords, tc.expectedTet.RotationCoords)
				}
			}
		})
	}
}

func TestTranspose(t *testing.T) {
	tt := []struct {
		name          string
		tet           *Tetrimino
		expectedCells [][]bool
	}{
		{
			"1x2",
			&Tetrimino{
				Minos: [][]bool{
					{true, false},
				},
			},
			[][]bool{
				{true},
				{false},
			},
		},
		{
			"2x2",
			&Tetrimino{
				Minos: [][]bool{
					{true, false},
					{true, false},
				},
			},
			[][]bool{
				{true, true},
				{false, false},
			},
		},
		{
			"3x3",
			&Tetrimino{
				Minos: [][]bool{
					{true, false, true},
					{true, false, false},
					{true, false, true},
				},
			},
			[][]bool{
				{true, true, true},
				{false, false, false},
				{true, false, true},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.tet.transpose()

			if !reflect.DeepEqual(tc.tet.Minos, tc.expectedCells) {
				t.Errorf("got %v,\nwant %v", tc.tet.Minos, tc.expectedCells)
			}
		})
	}
}

func TestCanRotate(t *testing.T) {
	matrix := &Matrix{}

	tt := []struct {
		name    string
		matrix  *Matrix
		rotated *Tetrimino
		expects bool
	}{
		{
			"can rotate, empty board & starting position",
			matrix,
			&Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{false, true, false},
					{true, true, true},
				},
				Pos: Coordinate{
					X: startingPositions['6'].X,
					Y: startingPositions['6'].Y + 20,
				},
				CurrentRotation: 0,
				RotationCoords:  RotationCoords['6'],
			},
			true,
		},
		{
			"can rotate, perfect fit",
			&Matrix{
				{0, 0, 0, 'X'},
				{'X', 0, 'X'},
				{'X', 'X', 'X'},
			},
			&Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Pos: Coordinate{
					X: 0,
					Y: 0,
				},
				CurrentRotation: 2,
				RotationCoords:  RotationCoords['6'],
			},
			true,
		},
		{
			"cannot rotate, blocking cell",
			&Matrix{
				{'X'},
			},
			&Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Pos: Coordinate{
					X: 0,
					Y: 0,
				},
				CurrentRotation: 2,
				RotationCoords:  RotationCoords['6'],
			},
			false,
		},
		{
			"cannot rotate, out of bounds left",
			matrix,
			&Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Pos: Coordinate{
					X: -1,
					Y: 0,
				},
				CurrentRotation: 2,
				RotationCoords:  RotationCoords['6'],
			},
			false,
		},
		{
			"cannot rotate, out of bounds right",
			matrix,
			&Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Pos: Coordinate{
					X: len(matrix[0]) - 2,
					Y: 0,
				},
				CurrentRotation: 2,
				RotationCoords:  RotationCoords['6'],
			},
			false,
		},
		{
			"cannot rotate, out of bounds up",
			matrix,
			&Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Pos: Coordinate{
					X: 0,
					Y: -1,
				},
				CurrentRotation: 2,
				RotationCoords:  RotationCoords['6'],
			},
			false,
		},
		{
			"cannot rotate, out of bounds down",
			matrix,
			&Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Pos: Coordinate{
					X: 0,
					Y: len(matrix) - 1,
				},
				CurrentRotation: 2,
				RotationCoords:  RotationCoords['6'],
			},
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			result := tc.rotated.canBePlaced(tc.matrix)

			if result != tc.expects {
				t.Errorf("got %v, want %v", result, tc.expects)
			}
		})
	}
}

func TestIsOutOfBoundsHorizontally(t *testing.T) {
	matrix := &Matrix{}

	tt := []struct {
		name    string
		tetPosX int
		cellCol int
		expects bool
	}{
		{
			"out left", 0, -1, true,
		},
		{
			"in left", -1, 1, false,
		},
		{
			"in right", 10, -1, false,
		},
		{
			"out right", 10, 0, true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := isOutOfBoundsHorizontally(tc.tetPosX, tc.cellCol, matrix)

			if result != tc.expects {
				t.Errorf("got %v, want %v", result, tc.expects)
			}
		})
	}
}

func TestIsOutOfBoundsVertically(t *testing.T) {
	matrix := &Matrix{}

	tt := []struct {
		name    string
		tetPosY int
		cellRow int
		expects bool
	}{
		{
			"out up", 0, -1, true,
		},
		{
			"in up", -1, 1, false,
		},
		{
			"in down", 40, -1, false,
		},
		{
			"out down", 40, 0, true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := isOutOfBoundsVertically(tc.tetPosY, tc.cellRow, matrix)

			if result != tc.expects {
				t.Errorf("got %v, want %v", result, tc.expects)
			}
		})
	}
}

func TestPositiveMod(t *testing.T) {
	tt := []struct {
		name       string
		dividend   int
		divisor    int
		expected   int
		expectsErr bool
	}{
		{
			"0 mod 0",
			0, 0, 0, true,
		},
		{
			"1 mod 0",
			0, 0, 0, true,
		},
		{
			"3 mod 5",
			3, 5, 3, false,
		},
		{
			"5 mod 5",
			5, 5, 0, false,
		},
		{
			"3 mod -5",
			3, -5, 3, false,
		},
		{
			"5 mod -5",
			5, -5, 0, false,
		},
		{
			"-4 mod -5",
			-4, -5, -4, false,
		},
		{
			"-4 mod 5",
			-4, 5, 1, false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result, err := positiveMod(tc.dividend, tc.divisor)
			if result != tc.expected {
				t.Errorf("result: got %v, want %v", result, tc.expected)
			}
			if tc.expectsErr != (err != nil) {
				t.Error("err: got error, want nil")
			}
		})
	}
}

func TestDeepCopyCells(t *testing.T) {
	var cells, cellsCopy [][]bool

	// Case: modify byte
	cells = [][]bool{
		{false, false},
	}

	cellsCopy = deepCopyCells(cells)
	cellsCopy[0][0] = true

	if cells[0][0] {
		t.Errorf("byte was modified")
	}

	// Case: modify inner array
	cells = [][]bool{
		{false, false},
	}

	cellsCopy = deepCopyCells(cells)
	cellsCopy[0] = []bool{true, false}

	if cells[0][0] {
		t.Errorf("byte was modified")
	}

}

func TestIsCellEmpty(t *testing.T) {
	tt := []struct {
		name     string
		cell     byte
		expected bool
	}{
		{
			name:     "0",
			cell:     0,
			expected: true,
		},
		{
			name:     "G",
			cell:     'G',
			expected: true,
		},
		{
			name:     "I",
			cell:     'I',
			expected: false,
		},
		{
			name:     "O",
			cell:     'O',
			expected: false,
		},
		{
			name:     "T",
			cell:     'T',
			expected: false,
		},
		{
			name:     "S",
			cell:     'S',
			expected: false,
		},
		{
			name:     "Z",
			cell:     'Z',
			expected: false,
		},
		{
			name:     "J",
			cell:     'J',
			expected: false,
		},
		{
			name:     "L",
			cell:     'L',
			expected: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := isCellEmpty(tc.cell)

			if result != tc.expected {
				t.Errorf("got %v, want %v", result, tc.expected)
			}
		})
	}
}

func TestTetrimino_Copy(t *testing.T) {
	// Create a tetrimino
	tetValue := byte('T')
	tetCells := [][]bool{
		{false, true, false},
		{true, true, true},
	}
	tetPosX := 0
	tetPosY := 0
	tetCurrentRotation := 0
	tetRotationCoords := []Coordinate{
		{0, 0},
	}

	tet := Tetrimino{
		Value:           tetValue,
		Minos:           tetCells,
		Pos:             Coordinate{tetPosX, tetPosY},
		CurrentRotation: tetCurrentRotation,
		RotationCoords:  tetRotationCoords,
	}

	// Create a copy manually
	manualCopy := Tetrimino{
		Value:           tetValue,
		Minos:           deepCopyCells(tetCells),
		Pos:             Coordinate{tetPosX, tetPosY},
		CurrentRotation: tetCurrentRotation,
	}
	manualCopy.RotationCoords = make([]Coordinate, len(tetRotationCoords))
	copy(manualCopy.RotationCoords, tet.RotationCoords)

	// Create a copy with the function
	easyCopy := tet.DeepCopy()

	// Check that the copy is the same as the original and the manual copy
	if !reflect.DeepEqual(easyCopy, &tet) {
		t.Errorf("easy copy doesn't match original: easy %v, original %v", easyCopy, &tet)
	} else if !reflect.DeepEqual(easyCopy, &manualCopy) {
		t.Errorf("easy copy doesn't match manual copy: easy %v, manual %v", easyCopy, &manualCopy)
	}

	// Modify the original
	tet.Value = 'I'
	tet.Minos = [][]bool{
		{true, true, true, true},
	}
	tet.Pos = Coordinate{1, 1}
	tet.CurrentRotation = 1
	tet.RotationCoords = []Coordinate{
		{1, 1},
	}

	// Check that the copy is still the same as the manual copy
	if !reflect.DeepEqual(easyCopy, &manualCopy) {
		t.Errorf("easy copy was modified: easy %v, manual %v", easyCopy, &manualCopy)
	}
}

func TestTetrimino_IsAbovePlayfield(t *testing.T) {
	tt := []struct {
		name         string
		matrixLength int
		tet          *Tetrimino
		expected     bool
	}{
		{
			name:         "true, length 30",
			matrixLength: 30,
			tet: &Tetrimino{
				Pos:   Coordinate{X: 0, Y: 9},
				Minos: [][]bool{{true, true, true, true}},
			},
			expected: true,
		},
		{
			name:         "true, length 40",
			matrixLength: 40,
			tet: &Tetrimino{
				Pos:   Coordinate{X: 0, Y: 19},
				Minos: [][]bool{{true, true, true, true}},
			},
			expected: true,
		},
		{
			name:         "false, length 30",
			matrixLength: 30,
			tet: &Tetrimino{
				Pos:   Coordinate{X: 0, Y: 10},
				Minos: [][]bool{{true, true, true, true}},
			},
			expected: false,
		},
		{
			name:         "false, length 40",
			matrixLength: 40,
			tet: &Tetrimino{
				Pos:   Coordinate{X: 0, Y: 20},
				Minos: [][]bool{{true, true, true, true}},
			},
			expected: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.tet.IsAbovePlayfield(tc.matrixLength)

			if result != tc.expected {
				t.Errorf("got %v, want %v", result, tc.expected)
			}
		})
	}
}

func TestTetrimino_IsOverlapping(t *testing.T) {
	tt := []struct {
		name     string
		tet      *Tetrimino
		matrix   *Matrix
		expected bool
	}{
		{
			name: "true",
			tet: &Tetrimino{
				Pos:   Coordinate{X: 0, Y: 0},
				Minos: [][]bool{{true, true, true, true}},
			},
			matrix: &Matrix{
				{0, 'X'},
			},
			expected: true,
		},
		{
			name: "false",
			tet: &Tetrimino{
				Pos:   Coordinate{X: 0, Y: 0},
				Minos: [][]bool{{true, true, true, true}},
			},
			matrix: &Matrix{
				{0, 0, 0, 0, 'X'},
				{'X', 'X', 'X', 'X', 'X'},
			},
			expected: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.tet.IsOverlapping(tc.matrix)

			if result != tc.expected {
				t.Errorf("got %v, want %v", result, tc.expected)
			}
		})
	}
}
