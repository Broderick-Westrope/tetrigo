package main

import (
	"testing"
)

func TestRandomTetrimino(t *testing.T) {
	tt := []struct {
		name            string
		playfieldHeight int
	}{
		{
			name:            "height 20",
			playfieldHeight: 21,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tet := RandomTetrimino(tc.playfieldHeight)
			if tet == nil {
				t.Errorf("expected tetrimino, got nil")
			}

			var expectedY int
			switch tet.Value {
			case 'I':
				expectedY = (tc.playfieldHeight - 20) - 1
			case 'O':
				expectedY = (tc.playfieldHeight - 20) - 2
			case 'T':
				expectedY = (tc.playfieldHeight - 20) - 2
			case 'S':
				expectedY = (tc.playfieldHeight - 20) - 2
			case 'Z':
				expectedY = (tc.playfieldHeight - 20) - 2
			case 'J':
				expectedY = (tc.playfieldHeight - 20) - 2
			case 'L':
				expectedY = (tc.playfieldHeight - 20) - 2
			default:
				t.Errorf("tetrimino value %c not valid", tet.Value)
			}
			if tet.Pos.Y != expectedY {
				t.Errorf("expected tetrimino to be at row %d, got %d", expectedY, tet.Pos.Y)
			}
		})
	}
}

// func TestTetrimino_MoveDown(t *testing.T) {
// 	tt := []struct {
// 		name              string
// 		startingPlayfield Playfield
// 		startingTet       Tetrimino
// 		expectedPlayfield Playfield
// 		expectsTet        bool
// 		expectsErr        bool
// 	}{
// 		{
// 			name: "can, empty playfield",
// 			startingPlayfield: Playfield{
// 				{'T', 'T', 'T'},
// 				{0, 'T', 0},
// 			},
// 			startingTet: Tetrimino{
// 				Value: 'T',
// 				Cells: [][]bool{
// 					{true, true, true},
// 					{false, true, false},
// 				},
// 				Pos: Coordinate{X: 0, Y: 0},
// 			},
// 			expectedPlayfield: Playfield{
// 				{},
// 				{'T', 'T', 'T'},
// 				{0, 'T', 0},
// 			},
// 			expectsTet: false,
// 			expectsErr: false,
// 		},
// 		{
// 			name: "can, perfect fit",
// 			startingPlayfield: Playfield{
// 				{'T', 'T', 'T'},
// 				{0, 'T', 0},
// 				{'#', 0, '#'},
// 				{'#', '#', '#'},
// 			},
// 			startingTet: Tetrimino{
// 				Value: 'T',
// 				Cells: [][]bool{
// 					{true, true, true},
// 					{false, true, false},
// 				},
// 				Pos: Coordinate{X: 0, Y: 0},
// 			},
// 			expectedPlayfield: Playfield{
// 				{},
// 				{'T', 'T', 'T'},
// 				{'#', 'T', '#'},
// 				{'#', '#', '#'},
// 			},
// 			expectsTet: false,
// 			expectsErr: false,
// 		},
// 		{
// 			name: "can, ghost cells",
// 			startingPlayfield: Playfield{
// 				{'T', 'T', 'T'},
// 				{'G', 'T', 'G'},
// 				{'#', 'G', '#'},
// 				{'#', '#', '#'},
// 			},
// 			startingTet: Tetrimino{
// 				Value: 'T',
// 				Cells: [][]bool{
// 					{true, true, true},
// 					{false, true, false},
// 				},
// 				Pos: Coordinate{X: 0, Y: 0},
// 			},
// 			expectedPlayfield: Playfield{
// 				{},
// 				{'T', 'T', 'T'},
// 				{'#', 'T', '#'},
// 				{'#', '#', '#'},
// 			},
// 			expectsTet: false,
// 			expectsErr: false,
// 		},
// 		{
// 			name: "cannot, blocking tetrimino",
// 			startingPlayfield: Playfield{
// 				{'T', 'T', 'T'},
// 				{0, 'T', 0},
// 				{'#', '#', '#'},
// 			},
// 			startingTet: Tetrimino{
// 				Value: 'T',
// 				Cells: [][]bool{
// 					{true, true, true},
// 					{false, true, false},
// 				},
// 				Pos: Coordinate{X: 0, Y: 0},
// 			},
// 			expectedPlayfield: Playfield{
// 				{'T', 'T', 'T'},
// 				{0, 'T', 0},
// 				{'#', '#', '#'},
// 			},
// 			expectsTet: true,
// 			expectsErr: false,
// 		},
// 		{
// 			name: "cannot, end of playfield",
// 			startingPlayfield: Playfield{
// 				{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
// 				{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
// 				{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
// 				{}, {}, {}, {}, {}, {}, {}, {},
// 				{'T', 'T', 'T'},
// 				{0, 'T', 0},
// 			},
// 			startingTet: Tetrimino{
// 				Value: 'T',
// 				Cells: [][]bool{
// 					{true, true, true},
// 					{false, true, false},
// 				},
// 				Pos: Coordinate{X: 0, Y: 38},
// 			},
// 			expectedPlayfield: Playfield{
// 				{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
// 				{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
// 				{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
// 				{}, {}, {}, {}, {}, {}, {}, {},
// 				{'T', 'T', 'T'},
// 				{0, 'T', 0},
// 			},
// 			expectsTet: true,
// 			expectsErr: false,
// 		},
// 		{
// 			name: "error, wrong value",
// 			startingPlayfield: Playfield{
// 				{'X', 'X', 'X'},
// 				{0, 'X', 0},
// 			},
// 			startingTet: Tetrimino{
// 				Value: 'T',
// 				Cells: [][]bool{
// 					{true, true, true},
// 					{false, true, false},
// 				},
// 				Pos: Coordinate{X: 0, Y: 0},
// 			},
// 			expectedPlayfield: Playfield{},
// 			expectsTet:        false,
// 			expectsErr:        true,
// 		},
// 	}

// 	for _, tc := range tt {
// 		t.Run(tc.name, func(t *testing.T) {
// 			tet, err := tc.startingTet.MoveDown(&tc.startingPlayfield)

// 			if tc.expectsErr && err == nil {
// 				t.Errorf("expected error, got nil")
// 			} else if !tc.expectsErr && err != nil {
// 				t.Errorf("expected nil, got error")
// 			}

// 			// TODO: look into a better way to avoid comparing new tetriminos
// 			if err == nil && (!reflect.DeepEqual(tc.startingPlayfield[:10], tc.expectedPlayfield[:10]) || !reflect.DeepEqual(tc.startingPlayfield[30:], tc.expectedPlayfield[30:])) {
// 				t.Errorf("expected playfield %v, got %v", tc.expectedPlayfield, tc.startingPlayfield)
// 			}

// 			if tc.expectsTet && tet == nil {
// 				t.Errorf("expected tetrimino, got nil")
// 			} else if !tc.expectsTet && tet != nil {
// 				t.Errorf("expected nil, got tetrimino")
// 			}
// 		})
// 	}
// }

func TestTetrimino_MoveLeft(t *testing.T) {
	tt := []struct {
		name              string
		startingPlayfield Playfield
		startingTet       Tetrimino
		expectedPlayfield Playfield
		expectsErr        bool
	}{
		{
			name: "can, empty playfield",
			startingPlayfield: Playfield{
				{0, 'T', 'T', 'T'},
				{0, 0, 'T', 0},
			},
			startingTet: Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Pos: Coordinate{X: 1, Y: 0},
			},
			expectedPlayfield: Playfield{
				{'T', 'T', 'T'},
				{0, 'T', 0},
			},
			expectsErr: false,
		},
		{
			name: "can, perfect fit",
			startingPlayfield: Playfield{
				{'#', 0, 'T', 'T', 'T'},
				{'#', '#', 0, 'T', 0},
			},
			startingTet: Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Pos: Coordinate{X: 2, Y: 0},
			},
			expectedPlayfield: Playfield{
				{'#', 'T', 'T', 'T'},
				{'#', '#', 'T', 0},
			},
			expectsErr: false,
		},
		{
			name: "can, ghost cells",
			startingPlayfield: Playfield{
				{'#', 'G', 'T', 'T', 'T'},
				{'#', '#', 'G', 'T', 0},
			},
			startingTet: Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Pos: Coordinate{X: 2, Y: 0},
			},
			expectedPlayfield: Playfield{
				{'#', 'T', 'T', 'T'},
				{'#', '#', 'T', 0},
			},
			expectsErr: false,
		},
		{
			name: "cannot, blocking tetrimino",
			startingPlayfield: Playfield{
				{'#', 'T', 'T', 'T'},
				{'#', 0, 'T', 0},
			},
			startingTet: Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Pos: Coordinate{X: 1, Y: 0},
			},
			expectedPlayfield: Playfield{
				{'#', 'T', 'T', 'T'},
				{'#', 0, 'T', 0},
			},
			expectsErr: false,
		},
		{
			name: "cannot, end of playfield",
			startingPlayfield: Playfield{
				{'T', 'T', 'T'},
				{0, 'T', 0},
			},
			startingTet: Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Pos: Coordinate{X: 0, Y: 0},
			},
			expectedPlayfield: Playfield{
				{'T', 'T', 'T'},
				{0, 'T', 0},
			},
			expectsErr: false,
		},
		{
			name: "error, wrong value",
			startingPlayfield: Playfield{
				{0, 'X', 'X', 'X'},
				{0, 0, 'X', 0},
			},
			startingTet: Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Pos: Coordinate{X: 1, Y: 0},
			},
			expectedPlayfield: Playfield{},
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
				t.Errorf("expected playfield %v, got %v", tc.expectedPlayfield, tc.startingPlayfield)
			}
		})
	}
}

func TestTetrimino_MoveRight(t *testing.T) {
	tt := []struct {
		name              string
		startingPlayfield Playfield
		startingTet       Tetrimino
		expectedPlayfield Playfield
		expectsErr        bool
	}{
		{
			name: "can, empty playfield",
			startingPlayfield: Playfield{
				{'T', 'T', 'T'},
				{0, 'T', 0},
			},
			startingTet: Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Pos: Coordinate{X: 0, Y: 0},
			},
			expectedPlayfield: Playfield{
				{0, 'T', 'T', 'T'},
				{0, 0, 'T', 0},
			},
			expectsErr: false,
		},
		{
			name: "can, perfect fit",
			startingPlayfield: Playfield{
				{'T', 'T', 'T', 0, '#'},
				{0, 'T', 0, '#'},
			},
			startingTet: Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Pos: Coordinate{X: 0, Y: 0},
			},
			expectedPlayfield: Playfield{
				{0, 'T', 'T', 'T', '#'},
				{0, 0, 'T', '#'},
			},
			expectsErr: false,
		},
		{
			name: "can, ghost cells",
			startingPlayfield: Playfield{
				{'T', 'T', 'T', 'G', '#'},
				{0, 'T', 'G', '#'},
			},
			startingTet: Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Pos: Coordinate{X: 0, Y: 0},
			},
			expectedPlayfield: Playfield{
				{0, 'T', 'T', 'T', '#'},
				{0, 0, 'T', '#'},
			},
			expectsErr: false,
		},
		{
			name: "cannot, blocking tetrimino",
			startingPlayfield: Playfield{
				{'T', 'T', 'T', '#'},
				{0, 'T', '#'},
			},
			startingTet: Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Pos: Coordinate{X: 0, Y: 0},
			},
			expectedPlayfield: Playfield{
				{'T', 'T', 'T', '#'},
				{0, 'T', '#'},
			},
			expectsErr: false,
		},
		{
			name: "cannot, end of playfield",
			startingPlayfield: Playfield{
				{0, 0, 0, 0, 0, 0, 0, 'T', 'T', 'T'},
				{0, 0, 0, 0, 0, 0, 0, 0, 'T', 0},
			},
			startingTet: Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Pos: Coordinate{X: 7, Y: 0},
			},
			expectedPlayfield: Playfield{
				{0, 0, 0, 0, 0, 0, 0, 'T', 'T', 'T'},
				{0, 0, 0, 0, 0, 0, 0, 0, 'T', 0},
			},
			expectsErr: false,
		},
		{
			name: "error, wrong value",
			startingPlayfield: Playfield{
				{'X', 'X', 'X'},
				{0, 'X', 0},
			},
			startingTet: Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Pos: Coordinate{X: 0, Y: 0},
			},
			expectedPlayfield: Playfield{},
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
				t.Errorf("expected playfield %v, got %v", tc.expectedPlayfield, tc.startingPlayfield)
			}
		})
	}
}

func TestIsOutOfBoundsHorizontally(t *testing.T) {
	playfield := &Playfield{}

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
			result := isOutOfBoundsHorizontally(tc.tetPosX, tc.cellCol, playfield)

			if result != tc.expects {
				t.Errorf("got %v, want %v", result, tc.expects)
			}
		})
	}
}

func TestCanRotate(t *testing.T) {
	playfield := &Playfield{}

	tt := []struct {
		name      string
		playfield *Playfield
		rotated   *Tetrimino
		expects   bool
	}{
		{
			"can rotate, empty board & starting position",
			playfield,
			&Tetrimino{
				Value: 'T',
				Cells: [][]bool{
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
			&Playfield{
				{0, 0, 0, 'X'},
				{'X', 0, 'X'},
				{'X', 'X', 'X'},
			},
			&Tetrimino{
				Value: 'T',
				Cells: [][]bool{
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
			&Playfield{
				{'X'},
			},
			&Tetrimino{
				Value: 'T',
				Cells: [][]bool{
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
			playfield,
			&Tetrimino{
				Value: 'T',
				Cells: [][]bool{
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
			playfield,
			&Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Pos: Coordinate{
					X: len(playfield[0]) - 2,
					Y: 0,
				},
				CurrentRotation: 2,
				RotationCoords:  RotationCoords['6'],
			},
			false,
		},
		{
			"cannot rotate, out of bounds up",
			playfield,
			&Tetrimino{
				Value: 'T',
				Cells: [][]bool{
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
			playfield,
			&Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				Pos: Coordinate{
					X: 0,
					Y: len(playfield) - 1,
				},
				CurrentRotation: 2,
				RotationCoords:  RotationCoords['6'],
			},
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.rotated.canRotate(tc.playfield)

			if result != tc.expects {
				t.Errorf("got %v, want %v", result, tc.expects)
			}
		})
	}
}

func TestIsOutOfBoundsVertically(t *testing.T) {
	playfield := &Playfield{}

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
			result := isOutOfBoundsVertically(tc.tetPosY, tc.cellRow, playfield)

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
