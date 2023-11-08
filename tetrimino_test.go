package main

import "testing"

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
			if tet.Y != expectedY {
				t.Errorf("expected tetrimino to be at row %d, got %d", expectedY, tet.Y)
			}
		})
	}
}

func TestTetrimino_MoveDown(t *testing.T) {
	tt := []struct {
		name              string
		startingPlayfield Playfield
		startingTet       Tetrimino
		expectedPlayfield Playfield
		expectsTet        bool
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
				X: 0, Y: 0,
			},
			expectedPlayfield: Playfield{
				{},
				{'T', 'T', 'T'},
				{0, 'T', 0},
			},
			expectsTet: false,
			expectsErr: false,
		},
		{
			name: "can, perfect fit",
			startingPlayfield: Playfield{
				{'T', 'T', 'T'},
				{0, 'T', 0},
				{'#', 0, '#'},
				{'#', '#', '#'},
			},
			startingTet: Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				X: 0, Y: 0,
			},
			expectedPlayfield: Playfield{
				{},
				{'T', 'T', 'T'},
				{'#', 'T', '#'},
				{'#', '#', '#'},
			},
			expectsTet: false,
			expectsErr: false,
		},
		{
			name: "can, ghost cells",
			startingPlayfield: Playfield{
				{'T', 'T', 'T'},
				{'G', 'T', 'G'},
				{'#', 'G', '#'},
				{'#', '#', '#'},
			},
			startingTet: Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				X: 0, Y: 0,
			},
			expectedPlayfield: Playfield{
				{},
				{'T', 'T', 'T'},
				{'#', 'T', '#'},
				{'#', '#', '#'},
			},
			expectsTet: false,
			expectsErr: false,
		},
		{
			name: "cannot, blocking tetrimino",
			startingPlayfield: Playfield{
				{'T', 'T', 'T'},
				{0, 'T', 0},
				{'#', '#', '#'},
			},
			startingTet: Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				X: 0, Y: 0,
			},
			expectedPlayfield: Playfield{
				{'T', 'T', 'T'},
				{0, 'T', 0},
				{'#', '#', '#'},
			},
			expectsTet: true,
			expectsErr: false,
		},
		{
			name: "cannot, end of playfield",
			startingPlayfield: Playfield{
				{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
				{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
				{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
				{}, {}, {}, {}, {}, {}, {}, {},
				{'T', 'T', 'T'},
				{0, 'T', 0},
			},
			startingTet: Tetrimino{
				Value: 'T',
				Cells: [][]bool{
					{true, true, true},
					{false, true, false},
				},
				X: 0, Y: 38,
			},
			expectedPlayfield: Playfield{
				{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
				{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
				{}, {}, {}, {}, {}, {}, {}, {}, {}, {},
				{}, {}, {}, {}, {}, {}, {}, {},
				{'T', 'T', 'T'},
				{0, 'T', 0},
			},
			expectsTet: true,
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
				X: 0, Y: 0,
			},
			expectedPlayfield: Playfield{},
			expectsTet:        false,
			expectsErr:        true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tet, err := tc.startingTet.MoveDown(&tc.startingPlayfield)

			if tc.expectsErr && err == nil {
				t.Errorf("expected error, got nil")
			} else if !tc.expectsErr && err != nil {
				t.Errorf("expected nil, got error")
			}

			if err == nil && tc.startingPlayfield != tc.expectedPlayfield {
				t.Errorf("expected playfield %v, got %v", tc.expectedPlayfield, tc.startingPlayfield)
			}

			if tc.expectsTet && tet == nil {
				t.Errorf("expected tetrimino, got nil")
			} else if !tc.expectsTet && tet != nil {
				t.Errorf("expected nil, got tetrimino")
			}
		})
	}
}
