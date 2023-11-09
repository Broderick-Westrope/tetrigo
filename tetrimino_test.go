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
