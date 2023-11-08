package main

import (
	"fmt"
	"math/rand"
)

var tetriminos = []Tetrimino{
	{
		Value: 'I',
		Cells: [][]bool{
			{true, true, true, true},
		},
		X: 3, Y: -1,
	},
	{
		Value: 'O',
		Cells: [][]bool{
			{true, true},
			{true, true},
		},
		X: 4, Y: -2,
	},
	{
		Value: 'T',
		Cells: [][]bool{
			{false, true, false},
			{true, true, true},
		},
		X: 3, Y: -2,
	},
	{
		Value: 'S',
		Cells: [][]bool{
			{false, true, true},
			{true, true, false},
		},
		X: 3, Y: -2,
	},
	{
		Value: 'Z',
		Cells: [][]bool{
			{true, true, false},
			{false, true, true},
		},
		X: 3, Y: -2,
	},
	{
		Value: 'J',
		Cells: [][]bool{
			{true, false, false},
			{true, true, true},
		},
		X: 3, Y: -2,
	},
	{
		Value: 'L',
		Cells: [][]bool{
			{false, false, true},
			{true, true, true},
		},
		X: 3, Y: -2,
	},
}

type Tetrimino struct {
	Value byte
	Cells [][]bool
	X, Y  int // the top left cell of the tetrimino
}

// RandomTetrimino returns a random tetrimino.
// The tetrimino will be positioned just above the playfield.
func RandomTetrimino(playfieldHeight int) *Tetrimino {
	t := tetriminos[rand.Intn(len(tetriminos))]
	t.Y += playfieldHeight - 20
	return &t
}

// MoveDown moves the tetrimino down one row.
// If the tetrimino cannot move down, it will be added to the playfield and a new tetrimino will be returned.
func (t *Tetrimino) MoveDown(playfield *Playfield) (*Tetrimino, error) {
	if !t.canMoveDown(*playfield) {
		return RandomTetrimino(len(playfield)), nil
	}
	err := playfield.removeCells(t)
	if err != nil {
		return nil, err
	}
	t.Y++
	err = playfield.addCells(t)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// MoveLeft moves the tetrimino left one column.
// If the tetrimino cannot move left, it will not move.
func (t *Tetrimino) MoveLeft(playfield *Playfield) error {
	if !t.canMoveLeft(*playfield) {
		return nil
	}
	err := playfield.removeCells(t)
	if err != nil {
		return err
	}
	t.X--
	err = playfield.addCells(t)
	if err != nil {
		return err
	}
	return nil
}

// MoveRight moves the tetrimino right one column.
// If the tetrimino cannot move right, it will not move.
func (t *Tetrimino) MoveRight(playfield *Playfield) error {
	if !t.canMoveRight(*playfield) {
		return nil
	}
	err := playfield.removeCells(t)
	if err != nil {
		return err
	}
	t.X++
	err = playfield.addCells(t)
	if err != nil {
		return err
	}
	return nil
}

func (t *Tetrimino) canMoveDown(playfield Playfield) bool {
	bottomRow := len(t.Cells) - 1
	for col := range t.Cells[bottomRow] {
		if t.Cells[bottomRow][col] {
			if bottomRow+t.Y+1 >= len(playfield) {
				return false
			}
			if cell := playfield[bottomRow+t.Y+1][col+t.X]; cell != 0 && cell != 'G' {
				return false
			}
		}
	}
	return true
}

func (t *Tetrimino) canMoveLeft(playfield Playfield) bool {
	leftCol := 0
	for row := range t.Cells {
		if t.Cells[row][leftCol] {
			if leftCol+t.X-1 < 0 {
				return false
			}
			if cell := playfield[row+t.Y][leftCol+t.X-1]; cell != 0 && cell != 'G' {
				return false
			}
		}
	}
	return true
}

func (t *Tetrimino) canMoveRight(playfield Playfield) bool {
	rightCol := len(t.Cells[0]) - 1
	for row := range t.Cells {
		if t.Cells[row][rightCol] {
			if rightCol+t.X+1 >= len(playfield[0]) {
				return false
			}
			if cell := playfield[row+t.Y][rightCol+t.X+1]; cell != 0 && cell != 'G' {
				return false
			}
		}
	}
	return true
}

func (p *Playfield) removeCells(tetrimino *Tetrimino) error {
	for row := range tetrimino.Cells {
		for col := range tetrimino.Cells[row] {
			if tetrimino.Cells[row][col] {
				if p[row+tetrimino.Y][col+tetrimino.X] != tetrimino.Value {
					return fmt.Errorf("cell at row %d, col %d is not the expected value", row+tetrimino.Y, col+tetrimino.X)
				}
				p[row+tetrimino.Y][col+tetrimino.X] = 0
			}
		}
	}
	return nil
}

func (p *Playfield) addCells(tetrimino *Tetrimino) error {
	for row := range tetrimino.Cells {
		for col := range tetrimino.Cells[row] {
			if tetrimino.Cells[row][col] {
				if cell := p[row+tetrimino.Y][col+tetrimino.X]; cell != 0 && cell != 'G' {
					return fmt.Errorf("cell at row %d, col %d is not empty or a ghost", row+tetrimino.Y, col+tetrimino.X)
				}
				p[row+tetrimino.Y][col+tetrimino.X] = tetrimino.Value
			}
		}
	}
	return nil
}
