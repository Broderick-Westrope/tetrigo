package tetris

import (
	"fmt"
)

type Coordinate struct {
	X, Y int
}

var RotationCoords = map[byte][]Coordinate{
	'I': {
		{X: -1, Y: -1},
		{X: 2, Y: 1},
		{X: -2, Y: -2},
		{X: 1, Y: 2},
	},
	'O': {
		{X: 0, Y: 0},
		{X: 0, Y: 0},
		{X: 0, Y: 0},
		{X: 0, Y: 0},
	},
	'6': { // All tetriminos with 6 cells (T, S, Z, J, L)
		{X: 0, Y: 0},
		{X: 1, Y: 0},
		{X: -1, Y: -1},
		{X: 0, Y: 1},
	},
}

var startingPositions = map[byte]Coordinate{
	'I': {X: 3, Y: -1},
	'O': {X: 4, Y: -2},
	'6': {X: 3, Y: -2},
}

var Tetriminos = []Tetrimino{
	{
		Value: 'I',
		Cells: [][]bool{
			{true, true, true, true},
		},
		Pos:            startingPositions['I'],
		RotationCoords: RotationCoords['I'],
	},
	{
		Value: 'O',
		Cells: [][]bool{
			{true, true},
			{true, true},
		},
		Pos:            startingPositions['O'],
		RotationCoords: RotationCoords['O'],
	},
	{
		Value: 'T',
		Cells: [][]bool{
			{false, true, false},
			{true, true, true},
		},
		Pos:            startingPositions['6'],
		RotationCoords: RotationCoords['6'],
	},
	{
		Value: 'S',
		Cells: [][]bool{
			{false, true, true},
			{true, true, false},
		},
		Pos:            startingPositions['6'],
		RotationCoords: RotationCoords['6'],
	},
	{
		Value: 'Z',
		Cells: [][]bool{
			{true, true, false},
			{false, true, true},
		},
		Pos:            startingPositions['6'],
		RotationCoords: RotationCoords['6'],
	},
	{
		Value: 'J',
		Cells: [][]bool{
			{true, false, false},
			{true, true, true},
		},
		Pos:            startingPositions['6'],
		RotationCoords: RotationCoords['6'],
	},
	{
		Value: 'L',
		Cells: [][]bool{
			{false, false, true},
			{true, true, true},
		},
		Pos:            startingPositions['6'],
		RotationCoords: RotationCoords['6'],
	},
}

type Tetrimino struct {
	Value           byte
	Cells           [][]bool
	Pos             Coordinate // the top left cell of the tetrimino
	CurrentRotation int
	RotationCoords  []Coordinate
}

// MoveDown moves the tetrimino down one row.
// If the tetrimino cannot move down, it will be added to the matrix and a new tetrimino will be returned.
func (t *Tetrimino) MoveDown(matrix *Matrix) error {
	err := matrix.RemoveTetrimino(t)
	if err != nil {
		return fmt.Errorf("failed to remove cells: %w", err)
	}
	t.Pos.Y++
	err = matrix.AddTetrimino(t)
	if err != nil {
		return fmt.Errorf("failed to add cells: %w", err)
	}
	return nil
}

// MoveLeft moves the tetrimino left one column.
// If the tetrimino cannot move left, it will not move.
func (t *Tetrimino) MoveLeft(matrix *Matrix) error {
	if !t.canMoveLeft(*matrix) {
		return nil
	}
	err := matrix.RemoveTetrimino(t)
	if err != nil {
		return fmt.Errorf("failed to remove cells: %w", err)
	}
	t.Pos.X--
	err = matrix.AddTetrimino(t)
	if err != nil {
		return fmt.Errorf("failed to add cells: %w", err)
	}
	return nil
}

// MoveRight moves the tetrimino right one column.
// If the tetrimino cannot move right, it will not move.
func (t *Tetrimino) MoveRight(matrix *Matrix) error {
	if !t.canMoveRight(*matrix) {
		return nil
	}
	err := matrix.RemoveTetrimino(t)
	if err != nil {
		return fmt.Errorf("failed to remove cells: %w", err)
	}
	t.Pos.X++
	err = matrix.AddTetrimino(t)
	if err != nil {
		return fmt.Errorf("failed to add cells: %w", err)
	}
	return nil
}

func (t *Tetrimino) CanMoveDown(matrix Matrix) bool {
	bottomRow := len(t.Cells) - 1
	for col := range t.Cells[bottomRow] {
		if t.Cells[bottomRow][col] {
			if bottomRow+t.Pos.Y+1 >= len(matrix) {
				return false
			}
			if !isCellEmpty(matrix[bottomRow+t.Pos.Y+1][col+t.Pos.X]) {
				return false
			}
		}
	}
	return true
}

func (t *Tetrimino) canMoveLeft(matrix Matrix) bool {
	leftCol := 0
	for row := range t.Cells {
		if t.Cells[row][leftCol] {
			if leftCol+t.Pos.X-1 < 0 {
				return false
			}
			if !isCellEmpty(matrix[row+t.Pos.Y][leftCol+t.Pos.X-1]) {
				return false
			}
		}
	}
	return true
}

func (t *Tetrimino) canMoveRight(matrix Matrix) bool {
	rightCol := len(t.Cells[0]) - 1
	for row := range t.Cells {
		if t.Cells[row][rightCol] {
			if rightCol+t.Pos.X+1 >= len(matrix[0]) {
				return false
			}
			if !isCellEmpty(matrix[row+t.Pos.Y][rightCol+t.Pos.X+1]) {
				return false
			}
		}
	}
	return true
}

func (t *Tetrimino) Rotate(matrix *Matrix, clockwise bool) error {
	if t.Value == 'O' {
		return nil
	}

	var rotated *Tetrimino
	var err error
	if clockwise {
		rotated, err = t.rotateClockwise()
	} else {
		rotated, err = t.rotateCounterClockwise()
	}
	if err != nil {
		return fmt.Errorf("failed to rotate tetrimino: %w", err)
	}

	err = matrix.RemoveTetrimino(t)
	if err != nil {
		return fmt.Errorf("failed to remove cells: %w", err)
	}

	if rotated.canRotate(matrix) {
		t.Cells = rotated.Cells
	}

	err = matrix.AddTetrimino(t)
	if err != nil {
		return fmt.Errorf("failed to add cells: %w", err)
	}
	return nil
}

func (t Tetrimino) rotateClockwise() (*Tetrimino, error) {
	t.Cells = deepCopyCells(t.Cells)

	// Reverse the order of the rows
	for i, j := 0, len(t.Cells)-1; i < j; i, j = i+1, j-1 {
		t.Cells[i], t.Cells[j] = t.Cells[j], t.Cells[i]
	}

	t.transpose()

	var err error
	t.CurrentRotation, err = positiveMod(t.CurrentRotation+1, len(t.RotationCoords))
	if err != nil {
		return nil, fmt.Errorf("failed to get positive mod: %w", err)
	}

	t.Pos.X += t.RotationCoords[t.CurrentRotation].X
	t.Pos.Y += t.RotationCoords[t.CurrentRotation].Y

	return &t, nil
}

func (t Tetrimino) rotateCounterClockwise() (*Tetrimino, error) {
	t.Cells = deepCopyCells(t.Cells)

	// Reverse the order of the columns
	for _, row := range t.Cells {
		for i, j := 0, len(row)-1; i < j; i, j = i+1, j-1 {
			row[i], row[j] = row[j], row[i]
		}
	}

	t.transpose()

	var err error
	t.CurrentRotation, err = positiveMod(t.CurrentRotation-1, len(t.RotationCoords))
	if err != nil {
		return nil, fmt.Errorf("failed to get positive mod: %w", err)
	}

	t.Pos.X -= t.RotationCoords[t.CurrentRotation].X
	t.Pos.Y -= t.RotationCoords[t.CurrentRotation].Y

	return &t, nil
}

func (t *Tetrimino) transpose() {
	xl := len(t.Cells[0])
	yl := len(t.Cells)
	result := make([][]bool, xl)
	for i := range result {
		result[i] = make([]bool, yl)
	}
	for i := 0; i < xl; i++ {
		for j := 0; j < yl; j++ {
			result[i][j] = t.Cells[j][i]
		}
	}
	t.Cells = result
}

func (t *Tetrimino) canRotate(matrix *Matrix) bool {
	for cellRow := range t.Cells {
		for cellCol := range t.Cells[cellRow] {
			if t.Cells[cellRow][cellCol] {
				if isOutOfBoundsHorizontally(t.Pos.X, cellCol, matrix) {
					return false
				}
				if isOutOfBoundsVertically(t.Pos.Y, cellRow, matrix) {
					return false
				}
				if !isCellEmpty(matrix[t.Pos.Y+cellRow][t.Pos.X+cellCol]) {
					return false
				}
			}
		}
	}
	return true
}

func isOutOfBoundsHorizontally(tetPosX, cellCol int, matrix *Matrix) bool {
	tetPosX += cellCol
	return tetPosX < 0 || tetPosX >= len(matrix[0])
}
func isOutOfBoundsVertically(tetPosY, cellRow int, matrix *Matrix) bool {
	tetPosY += cellRow
	return tetPosY < 0 || tetPosY >= len(matrix)
}

func positiveMod(dividend, divisor int) (int, error) {
	if divisor == 0 {
		return 0, fmt.Errorf("cannot %v divide by %v", dividend, divisor)
	}
	result := dividend % divisor
	if result < 0 && divisor > 0 {
		result += divisor
	}
	return result, nil
}

func deepCopyCells(cells [][]bool) [][]bool {
	cellsCopy := make([][]bool, len(cells))
	for i := range cells {
		cellsCopy[i] = make([]bool, len(cells[i]))
		copy(cellsCopy[i], cells[i])
	}
	return cellsCopy
}

func isCellEmpty(cell byte) bool {
	return cell == 0 || cell == 'G'
}
