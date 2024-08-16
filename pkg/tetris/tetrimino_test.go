package tetris

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: tests to add:
// - TestTetrimino_MoveDown failure cases

func TestTetrimino_MoveDown_success(t *testing.T) {
	for _, tet := range Tetriminos {
		maxRow := 40 - len(tet.Minos)
		maxCol := 10 - len(tet.Minos[0])
		for row := 0; row < maxRow; row++ {
			for col := 0; col < maxCol; col++ {
				name := fmt.Sprintf("value %s, row %d, col %d", string(tet.Value), row, col)
				t.Run(name, func(t *testing.T) {
					matrix := DefaultMatrix()
					startingCoord := Coordinate{X: col, Y: row}
					tet.Pos = startingCoord
					err := matrix.AddTetrimino(&tet)
					if !assert.NoError(t, err) {
						return
					}

					err = tet.MoveDown(&matrix)

					assert.NoError(t, err)
					assert.Equal(t, tet.Pos.Y, startingCoord.Y+1)
				})
			}
		}
	}
}

func TestTetrimino_MoveLeft(t *testing.T) {
	tt := map[string]struct {
		matrix      Matrix
		startingTet Tetrimino
		wantMatrix  Matrix
		wantErr     bool
	}{
		"success - can, empty matrix": {
			matrix: Matrix{
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
			wantMatrix: Matrix{
				{'T', 'T', 'T', 0},
				{0, 'T', 0, 0},
			},
			wantErr: false,
		},
		"success - can, perfect fit": {
			matrix: Matrix{
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
			wantMatrix: Matrix{
				{'#', 'T', 'T', 'T', 0},
				{'#', '#', 'T', 0, 0},
			},
			wantErr: false,
		},
		"success - can, ghost cells": {
			matrix: Matrix{
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
			wantMatrix: Matrix{
				{'#', 'T', 'T', 'T', 0},
				{'#', '#', 'T', 0, 0},
			},
			wantErr: false,
		},
		"success - cannot, blocking tetrimino": {
			matrix: Matrix{
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
			wantMatrix: Matrix{
				{'#', 'T', 'T', 'T'},
				{'#', 0, 'T', 0},
			},
			wantErr: false,
		},
		"success - cannot, end of matrix": {
			matrix: Matrix{
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
			wantMatrix: Matrix{
				{'T', 'T', 'T'},
				{0, 'T', 0},
			},
			wantErr: false,
		},
		"failure - wrong value": {
			matrix: Matrix{
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
			wantMatrix: Matrix{},
			wantErr:    true,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			err := tc.startingTet.MoveLeft(tc.matrix)

			if tc.wantErr {
				assert.NotNil(t, err)
				return
			}
			assert.NoError(t, err)
			assert.ElementsMatch(t, tc.matrix, tc.wantMatrix)
		})
	}
}

func TestTetrimino_MoveRight(t *testing.T) {
	tt := map[string]struct {
		matrix      Matrix
		startingTet Tetrimino
		wantMatrix  Matrix
		wantErr     bool
	}{
		"success - can, empty matrix": {
			matrix: Matrix{
				{'T', 'T', 'T', 0},
				{0, 'T', 0, 0},
			},
			startingTet: Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Pos: Coordinate{X: 0, Y: 0},
			},
			wantMatrix: Matrix{
				{0, 'T', 'T', 'T'},
				{0, 0, 'T', 0},
			},
			wantErr: false,
		},
		"success - can, perfect fit": {
			matrix: Matrix{
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
			wantMatrix: Matrix{
				{0, 'T', 'T', 'T', '#'},
				{0, 0, 'T', '#'},
			},
			wantErr: false,
		},
		"success - can, ghost cells": {
			matrix: Matrix{
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
			wantMatrix: Matrix{
				{0, 'T', 'T', 'T', '#'},
				{0, 0, 'T', '#'},
			},
			wantErr: false,
		},
		"success - cannot, blocking tetrimino": {
			matrix: Matrix{
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
			wantMatrix: Matrix{
				{'T', 'T', 'T', '#'},
				{0, 'T', '#'},
			},
			wantErr: false,
		},
		"success - cannot, end of matrix": {
			matrix: Matrix{
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
			wantMatrix: Matrix{
				{0, 0, 0, 0, 0, 0, 0, 'T', 'T', 'T'},
				{0, 0, 0, 0, 0, 0, 0, 0, 'T', 0},
			},
			wantErr: false,
		},
		"failure - wrong value": {
			matrix: Matrix{
				{'X', 'X', 'X', 0},
				{0, 'X', 0, 0},
			},
			startingTet: Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Pos: Coordinate{X: 0, Y: 0},
			},
			wantMatrix: Matrix{},
			wantErr:    true,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			err := tc.startingTet.MoveRight(tc.matrix)

			if tc.wantErr {
				assert.NotNil(t, err)
				return
			}
			assert.NoError(t, err)
			assert.ElementsMatch(t, tc.matrix, tc.wantMatrix)
		})
	}
}

func TestTetrimino_CanMoveDown(t *testing.T) {
	tt := map[string]struct {
		tet      *Tetrimino
		matrix   Matrix
		expected bool
	}{
		"I, can, top of matrix": {
			tet: &Tetrimino{
				Value:           'I',
				Minos:           [][]bool{{true}, {true}, {true}, {true}},
				Pos:             Coordinate{X: 0, Y: 0},
				CurrentRotation: 1,
				RotationCoords:  RotationCoords['I'],
			},
			matrix:   DefaultMatrix(),
			expected: true,
		},
		"Z, can, starting position": {
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
			matrix:   DefaultMatrix(),
			expected: true,
		},
		"S, cannot, blocking tetrimino": {
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
			matrix: Matrix{
				{0, 0, 0},
				{0, 0, 0},
				{'X', 0, 0},
			},
			expected: false,
		},
		"O, cannot, bottom of matrix": {
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
			matrix:   DefaultMatrix(),
			expected: false,
		},
		"GitHub Issue #1": {
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
			matrix: Matrix{
				{0, 0},
				{'X', 0},
				{0, 0},
			},
			expected: false,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			err := tc.matrix.AddTetrimino(tc.tet)
			if !assert.NoError(t, err) {
				return
			}

			result := tc.tet.canMoveDown(tc.matrix)

			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestTetrimino_CanMoveLeft(t *testing.T) {
	tt := map[string]struct {
		tet      *Tetrimino
		matrix   Matrix
		expected bool
	}{
		"Z, can, starting position": {
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
			matrix:   DefaultMatrix(),
			expected: true,
		},
		"O, cannot, blocking tetrimino": {
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
			matrix: Matrix{
				{'X', 0, 0},
				{0, 0, 0},
			},
			expected: false,
		},
		"O, cannot, end of matrix": {
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
			matrix:   DefaultMatrix(),
			expected: false,
		},
		"GitHub Issue #1": {
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
			matrix: Matrix{
				{'X', 'X', 'X', 0},
				{0, 0, 0, 0},
			},
			expected: false,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			err := tc.matrix.AddTetrimino(tc.tet)
			if !assert.NoError(t, err) {
				return
			}

			result := tc.tet.canMoveLeft(tc.matrix)

			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestTetrimino_CanMoveRight(t *testing.T) {
	tt := map[string]struct {
		tet      *Tetrimino
		matrix   Matrix
		expected bool
	}{
		"Z, can, starting position": {
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
			matrix:   DefaultMatrix(),
			expected: true,
		},
		"O, cannot, blocking tetrimino": {
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
			matrix: Matrix{
				{0, 0, 'X'},
				{0, 0, 0},
			},
			expected: false,
		},
		"O, cannot, end of matrix": {
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
			matrix:   DefaultMatrix(),
			expected: false,
		},
		"GitHub Issue #1": {
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
			matrix: Matrix{
				{0, 'X', 'X', 'X'},
				{0, 0, 0, 0},
			},
			expected: false,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			err := tc.matrix.AddTetrimino(tc.tet)
			if !assert.NoError(t, err) {
				return
			}

			result := tc.tet.canMoveRight(tc.matrix)

			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestTetrimino_RotateClockwise_success(t *testing.T) {
	tt := map[string]struct {
		tet         *Tetrimino
		expectedTet *Tetrimino
	}{
		"I, rotation 0": {
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
		"I, rotation 1": {
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
		"I, rotation 2": {
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
		"I, rotation 3": {
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
		"O, rotation 0": {
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
		"O, rotation 1": {
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
		"O, rotation 2": {
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
		"O, rotation 3": {
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
		"T, rotation 0": {
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
		"T, rotation 1": {
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
		"T, rotation 2": {
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
		"T, rotation 3": {
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

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			newTet := tc.tet.DeepCopy()

			err := newTet.rotateClockwise()

			assert.NoError(t, err)
			assert.EqualValues(t, tc.expectedTet, newTet)
		})
	}
}

func TestTetrimino_RotateCounterClockwise_success(t *testing.T) {
	tt := map[string]struct {
		tet         *Tetrimino
		expectedTet *Tetrimino
	}{
		"I, rotation 0": {
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
		"I, rotation 1": {
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
		"I, rotation 2": {
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
		"I, rotation 3": {
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
		"O, rotation 0": {
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
		"O, rotation 1": {
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
		"O, rotation 2": {
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
		"O, rotation 3": {
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
		"T, rotation 0": {
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
		"T, rotation 1": {
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
		"T, rotation 2": {
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
		"T, rotation 3": {
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

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			newTet := tc.tet.DeepCopy()

			err := newTet.rotateCounterClockwise()

			assert.NoError(t, err)
			assert.EqualValues(t, tc.expectedTet, newTet)
		})
	}
}

func TestTetrimino_Transpose(t *testing.T) {
	tt := map[string]struct {
		tet           *Tetrimino
		expectedMinos [][]bool
	}{
		"1x2": {
			tet: &Tetrimino{
				Minos: [][]bool{
					{true, false},
				},
			},
			expectedMinos: [][]bool{
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
			expectedMinos: [][]bool{
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
			expectedMinos: [][]bool{
				{true, true, true},
				{false, false, false},
				{true, false, true},
			},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			tc.tet.transpose()

			assert.ElementsMatch(t, tc.tet.Minos, tc.expectedMinos)
		})
	}
}

func TestTetrimino_CanBePlaced(t *testing.T) {
	tt := map[string]struct {
		matrix  Matrix
		rotated *Tetrimino
		expects bool
	}{
		"can rotate, empty board & starting position": {
			DefaultMatrix(),
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
		"can rotate, perfect fit": {
			Matrix{
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
		"cannot rotate, blocking mino": {
			Matrix{
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
		"cannot rotate, out of bounds left": {
			DefaultMatrix(),
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
		"cannot rotate, out of bounds right": {
			DefaultMatrix(),
			&Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Pos: Coordinate{
					X: len(DefaultMatrix()[0]) - 2,
					Y: 0,
				},
				CurrentRotation: 2,
				RotationCoords:  RotationCoords['6'],
			},
			false,
		},
		"cannot rotate, out of bounds up": {
			DefaultMatrix(),
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
		"cannot rotate, out of bounds down": {
			DefaultMatrix(),
			&Tetrimino{
				Value: 'T',
				Minos: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Pos: Coordinate{
					X: 0,
					Y: len(DefaultMatrix()) - 1,
				},
				CurrentRotation: 2,
				RotationCoords:  RotationCoords['6'],
			},
			false,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			result := tc.rotated.isValid(tc.matrix)

			assert.EqualValues(t, result, tc.expects)
		})
	}
}

func Test_IsOutOfBoundsHorizontally(t *testing.T) {
	tt := map[string]struct {
		tetPosX int
		cellCol int
		expects bool
	}{
		"out left": {
			0, -1, true,
		},
		"in left": {
			-1, 1, false,
		},
		"in right": {
			10, -1, false,
		},
		"out right": {
			10, 0, true,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			result := isOutOfBoundsHorizontally(tc.tetPosX, tc.cellCol, DefaultMatrix())

			assert.EqualValues(t, result, tc.expects)
		})
	}
}

func Test_IsOutOfBoundsVertically(t *testing.T) {
	tt := map[string]struct {
		tetPosY int
		cellRow int
		expects bool
	}{
		"out up": {
			0, -1, true,
		},
		"in up": {
			-1, 1, false,
		},
		"in down": {
			40, -1, false,
		},
		"out down": {
			40, 0, true,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			result := isOutOfBoundsVertically(tc.tetPosY, tc.cellRow, DefaultMatrix())

			assert.Equal(t, tc.expects, result)
		})
	}
}

func TestPositiveMod(t *testing.T) {
	errResultIsNan := errors.New("result is NaN")

	tt := map[string]struct {
		dividend int
		divisor  int
		expected int
		wantErr  error
	}{
		"0 mod 0": {
			0, 0, 0, errResultIsNan,
		},
		"1 mod 0": {
			0, 0, 0, errResultIsNan,
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
				assert.Equal(t, tc.expected, result)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeepCopyMinos(t *testing.T) {
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

func TestIsMinoEmpty(t *testing.T) {
	tt := []struct {
		mino     byte
		expected bool
	}{
		{
			mino:     0,
			expected: true,
		},
		{
			mino:     'G',
			expected: true,
		},
		{
			mino:     'I',
			expected: false,
		},
		{
			mino:     'O',
			expected: false,
		},
		{
			mino:     'T',
			expected: false,
		},
		{
			mino:     'S',
			expected: false,
		},
		{
			mino:     'Z',
			expected: false,
		},
		{
			mino:     'J',
			expected: false,
		},
		{
			mino:     'L',
			expected: false,
		},
	}

	for _, tc := range tt {
		t.Run(string(tc.mino), func(t *testing.T) {
			result := isCellEmpty(tc.mino)

			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestTetrimino_DeepCopy(t *testing.T) {
	tet := Tetrimino{
		Value: byte('T'),
		Minos: [][]bool{
			{false, true, false},
			{true, true, true},
		},
		Pos:             Coordinate{0, 0},
		CurrentRotation: 0,
		RotationCoords:  []Coordinate{{0, 0}},
	}

	// Create a copy manually
	manualCopy := Tetrimino{
		Value:           tet.Value,
		Minos:           deepCopyMinos(tet.Minos),
		Pos:             Coordinate{tet.Pos.X, tet.Pos.Y},
		CurrentRotation: tet.CurrentRotation,
		RotationCoords:  make([]Coordinate, len(tet.RotationCoords)),
	}
	copy(manualCopy.RotationCoords, tet.RotationCoords)

	// Create a (dereferences) copy with the helper function
	easyCopy := *tet.DeepCopy()

	assert.EqualValues(t, tet, manualCopy)
	assert.EqualValues(t, tet, easyCopy)
	assert.EqualValues(t, manualCopy, easyCopy)

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

	assert.NotEqualValues(t, tet, manualCopy)
	assert.NotEqualValues(t, tet, easyCopy)
	assert.EqualValues(t, manualCopy, easyCopy)
}

func TestTetrimino_IsAbovePlayfield(t *testing.T) {
	tt := map[string]struct {
		skyline  int
		tet      *Tetrimino
		expected bool
	}{
		"true, skyline 10": {
			skyline: 10,
			tet: &Tetrimino{
				Pos:   Coordinate{X: 0, Y: 9},
				Minos: [][]bool{{true}},
			},
			expected: true,
		},
		"true, skyline 20": {
			skyline: 20,
			tet: &Tetrimino{
				Pos:   Coordinate{X: 0, Y: 19},
				Minos: [][]bool{{true}},
			},
			expected: true,
		},
		"false, skyline 10": {
			skyline: 10,
			tet: &Tetrimino{
				Pos:   Coordinate{X: 0, Y: 10},
				Minos: [][]bool{{true}},
			},
			expected: false,
		},
		"false, skyline 20": {
			skyline: 20,
			tet: &Tetrimino{
				Pos:   Coordinate{X: 0, Y: 20},
				Minos: [][]bool{{true}},
			},
			expected: false,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			result := tc.tet.IsAboveSkyline(tc.skyline)

			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestTetrimino_IsOverlapping(t *testing.T) {
	tt := map[string]struct {
		tet      *Tetrimino
		matrix   Matrix
		expected bool
	}{
		"true": {
			tet: &Tetrimino{
				Pos:   Coordinate{X: 0, Y: 0},
				Minos: [][]bool{{true, true, true, true}},
			},
			matrix: Matrix{
				{0, 'X'},
			},
			expected: true,
		},
		"false": {
			tet: &Tetrimino{
				Pos:   Coordinate{X: 0, Y: 0},
				Minos: [][]bool{{true, true, true, true}},
			},
			matrix: Matrix{
				{0, 0, 0, 0, 'X'},
				{'X', 'X', 'X', 'X', 'X'},
			},
			expected: false,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			result := tc.tet.IsOverlapping(tc.matrix)

			assert.Equal(t, tc.expected, result)
		})
	}
}
