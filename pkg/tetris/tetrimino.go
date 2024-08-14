package tetris

import (
	"errors"
	"fmt"
	"math"
)

// A Tetrimino is a geometric Tetris shape formed by four Minos connected along their sides.
type Tetrimino struct {
	Value           byte         // The value of the Tetrimino. This is the character that will be used to represent the Tetrimino in the matrix.
	Minos           [][]bool     // A 2D slice of cells that make up the Tetrimino. True means the mino is occupied by a Mino.
	Pos             Coordinate   // The top left mino of the Tetrimino. Used as a reference point for movement and rotation.
	CurrentRotation int          // The index of the current rotation in the RotationCoords slice.
	RotationCoords  []Coordinate // The coordinates used during rotation to control the axis.
}

// Coordinate represents a point on a 2D plane.
type Coordinate struct {
	X, Y int
}

// RotationCoords is a map of Tetrimino values to the coordinates used for rotation.
// Each slice should contain a coordinate for north, east, south, and west in that order.
// These are added to (clockwise) or subtracted from (counter-clockwise) the Tetrimino's position when rotating to ensure it rotates around the correct axis.
var RotationCoords = map[byte][]Coordinate{
	'I': {
		{X: -1, Y: 1},
		{X: 2, Y: -1},
		{X: -2, Y: 2},
		{X: 1, Y: -2},
	},
	'O': {
		{X: 0, Y: 0},
		{X: 0, Y: 0},
		{X: 0, Y: 0},
		{X: 0, Y: 0},
	},
	// All tetriminos with 6 cells (T, S, Z, J, L) have the same rotation coordinates:
	'6': {
		{X: 0, Y: 0},
		{X: 1, Y: 0},
		{X: -1, Y: 1},
		{X: 0, Y: -1},
	},
}

// A map of Tetrimino values to the starting position of the Tetrimino.
// The starting position is the top left mino of the Tetrimino when it spawns.
// This does not take into account the buffer zone at the top of the Matrix.
var startingPositions = map[byte]Coordinate{
	'I': {X: 3, Y: -1},
	'O': {X: 4, Y: -2},
	'6': {X: 3, Y: -2},
}

// Tetriminos contains all seven valid Tetrimino values that can be made using four Minos.
var Tetriminos = []Tetrimino{
	{
		Value: 'I',
		Minos: [][]bool{
			{true, true, true, true},
		},
		Pos:            startingPositions['I'],
		RotationCoords: RotationCoords['I'],
	},
	{
		Value: 'O',
		Minos: [][]bool{
			{true, true},
			{true, true},
		},
		Pos:            startingPositions['O'],
		RotationCoords: RotationCoords['O'],
	},
	{
		Value: 'T',
		Minos: [][]bool{
			{false, true, false},
			{true, true, true},
		},
		Pos:            startingPositions['6'],
		RotationCoords: RotationCoords['6'],
	},
	{
		Value: 'S',
		Minos: [][]bool{
			{false, true, true},
			{true, true, false},
		},
		Pos:            startingPositions['6'],
		RotationCoords: RotationCoords['6'],
	},
	{
		Value: 'Z',
		Minos: [][]bool{
			{true, true, false},
			{false, true, true},
		},
		Pos:            startingPositions['6'],
		RotationCoords: RotationCoords['6'],
	},
	{
		Value: 'J',
		Minos: [][]bool{
			{true, false, false},
			{true, true, true},
		},
		Pos:            startingPositions['6'],
		RotationCoords: RotationCoords['6'],
	},
	{
		Value: 'L',
		Minos: [][]bool{
			{false, false, true},
			{true, true, true},
		},
		Pos:            startingPositions['6'],
		RotationCoords: RotationCoords['6'],
	},
}

// EmptyTetrimino is a tetrimino with no cells or value. To be used for the starting (empty) hold.
var EmptyTetrimino = &Tetrimino{
	Minos: [][]bool{},
	Value: 0,
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
// If the tetrimino cannot move left, nothing will happen.
func (t *Tetrimino) MoveLeft(matrix Matrix) error {
	if !t.canMoveLeft(matrix) {
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
// If the tetrimino cannot move right, nothing will happen.
func (t *Tetrimino) MoveRight(matrix Matrix) error {
	if !t.canMoveRight(matrix) {
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

// Returns true if the tetrimino can move down one row.
// This gets the lowest mino in each column of the tetrimino, and checks if it is at the bottom of the matrix or if the mino below is occupied.
func (t *Tetrimino) CanMoveDown(matrix Matrix) bool {
	for col := range t.Minos[0] {
		for row := len(t.Minos) - 1; row >= 0; row-- {
			if t.Minos[row][col] {
				if row+t.Pos.Y+1 >= len(matrix) {
					return false
				}
				if !isMinoEmpty(matrix[row+t.Pos.Y+1][col+t.Pos.X]) {
					return false
				}
				break
			}
		}
	}
	return true
}

func (t *Tetrimino) canMoveLeft(matrix Matrix) bool {
	for row := range t.Minos {
		for col := range t.Minos[row] {
			if t.Minos[row][col] {
				if col+t.Pos.X-1 < 0 {
					return false
				}
				if !isMinoEmpty(matrix[row+t.Pos.Y][col+t.Pos.X-1]) {
					return false
				}
				break
			}
		}
	}
	return true
}

func (t *Tetrimino) canMoveRight(matrix Matrix) bool {
	for row := range t.Minos {
		for col := len(t.Minos[row]) - 1; col >= 0; col-- {
			if t.Minos[row][col] {
				if col+t.Pos.X+1 >= len(matrix[0]) {
					return false
				}
				if !isMinoEmpty(matrix[row+t.Pos.Y][col+t.Pos.X+1]) {
					return false
				}
				break
			}
		}
	}
	return true
}

func (t *Tetrimino) Rotate(matrix Matrix, clockwise bool) error {
	if t.Value == 'O' {
		return nil
	}

	rotated := t.DeepCopy()
	var err error
	if clockwise {
		err = rotated.rotateClockwise()
	} else {
		err = rotated.rotateCounterClockwise()
	}
	if err != nil {
		return fmt.Errorf("failed to rotate tetrimino: %w", err)
	}

	err = matrix.RemoveTetrimino(t)
	if err != nil {
		return fmt.Errorf("failed to remove cells: %w", err)
	}

	if rotated.canBePlaced(matrix) {
		t.Pos = rotated.Pos
		t.Minos = rotated.Minos
		t.CurrentRotation = rotated.CurrentRotation
	}

	err = matrix.AddTetrimino(t)
	if err != nil {
		return fmt.Errorf("failed to add cells: %w", err)
	}
	return nil
}

func (t *Tetrimino) rotateClockwise() error {
	// Reverse the order of the rows
	for i, j := 0, len(t.Minos)-1; i < j; i, j = i+1, j-1 {
		t.Minos[i], t.Minos[j] = t.Minos[j], t.Minos[i]
	}

	t.transpose()

	var err error
	t.CurrentRotation, err = positiveMod(t.CurrentRotation+1, len(t.RotationCoords))
	if err != nil {
		return fmt.Errorf("failed to get positive mod: %w", err)
	}

	t.Pos.X += t.RotationCoords[t.CurrentRotation].X
	t.Pos.Y += t.RotationCoords[t.CurrentRotation].Y

	return nil
}

func (t *Tetrimino) rotateCounterClockwise() error {
	// Reverse the order of the columns
	for _, row := range t.Minos {
		for i, j := 0, len(row)-1; i < j; i, j = i+1, j-1 {
			row[i], row[j] = row[j], row[i]
		}
	}

	t.transpose()

	t.Pos.X -= t.RotationCoords[t.CurrentRotation].X
	t.Pos.Y -= t.RotationCoords[t.CurrentRotation].Y

	var err error
	t.CurrentRotation, err = positiveMod(t.CurrentRotation-1, len(t.RotationCoords))
	if err != nil {
		return fmt.Errorf("failed to get positive mod: %w", err)
	}

	return nil
}

func (t *Tetrimino) transpose() {
	xl := len(t.Minos[0])
	yl := len(t.Minos)
	result := make([][]bool, xl)
	for i := range result {
		result[i] = make([]bool, yl)
	}
	for i := 0; i < xl; i++ {
		for j := 0; j < yl; j++ {
			result[i][j] = t.Minos[j][i]
		}
	}
	t.Minos = result
}

// canBePlaced returns true if the given Tetrimino is within the bounds of the matrix and does not overlap with any occupied cells.
// The Tetrimino being checked should not be in the Matrix yet.
func (t *Tetrimino) canBePlaced(matrix Matrix) bool {
	for cellRow := range t.Minos {
		for cellCol := range t.Minos[cellRow] {
			if t.Minos[cellRow][cellCol] {
				if isOutOfBoundsHorizontally(t.Pos.X, cellCol, matrix) {
					return false
				}
				if isOutOfBoundsVertically(t.Pos.Y, cellRow, matrix) {
					return false
				}
				if !isMinoEmpty(matrix[t.Pos.Y+cellRow][t.Pos.X+cellCol]) {
					return false
				}
			}
		}
	}
	return true
}

func isOutOfBoundsHorizontally(tetPosX, cellCol int, matrix Matrix) bool {
	tetPosX += cellCol
	return tetPosX < 0 || tetPosX >= len(matrix[0])
}

func isOutOfBoundsVertically(tetPosY, cellRow int, matrix Matrix) bool {
	tetPosY += cellRow
	return tetPosY < 0 || tetPosY >= len(matrix)
}

func positiveMod(dividend, divisor int) (int, error) {
	result := math.Mod(float64(dividend), float64(divisor))
	if math.IsNaN(result) {
		return 0, errors.New("result is NaN")
	}
	if result < 0 && divisor > 0 {
		return int(result) + divisor, nil
	}
	return int(result), nil
}

func deepCopyMinos(cells [][]bool) [][]bool {
	cellsCopy := make([][]bool, len(cells))
	for i := range cells {
		cellsCopy[i] = make([]bool, len(cells[i]))
		copy(cellsCopy[i], cells[i])
	}
	return cellsCopy
}

func isMinoEmpty(cell byte) bool {
	return cell == 0 || cell == 'G'
}

func (t *Tetrimino) DeepCopy() *Tetrimino {
	var cells [][]bool
	if t.Minos == nil {
		cells = nil
	} else {
		cells = deepCopyMinos(t.Minos)
	}

	var rotationCoords []Coordinate
	if t.RotationCoords == nil {
		rotationCoords = nil
	} else {
		rotationCoords = make([]Coordinate, len(t.RotationCoords))
		copy(rotationCoords, t.RotationCoords)
	}

	return &Tetrimino{
		Value:           t.Value,
		Minos:           cells,
		Pos:             t.Pos,
		CurrentRotation: t.CurrentRotation,
		RotationCoords:  rotationCoords,
	}
}

// IsAboveSkyline returns true if the entire Tetrimino is above the skyline.
// This can be helpful when checking for Lock Out.
func (t *Tetrimino) IsAboveSkyline(skyline int) bool {
	for row := len(t.Minos) - 1; row >= 0; row-- {
		for col := range t.Minos[row] {
			if t.Minos[row][col] {
				if t.Pos.Y+row >= skyline {
					return false
				}
				break
			}
		}
	}
	return true
}

// IsOverlapping checks whether the Tetrimino would be overlapping with an occupied Mino if it were on the Matrix.
// The Tetrmino should not yet be added to the Matrix, otherwise this will always return true as a Tetrimino is always overlapping with itself.
func (t *Tetrimino) IsOverlapping(matrix Matrix) bool {
	for col := range t.Minos[0] {
		for row := range t.Minos {
			if t.Minos[row][col] {
				if !isMinoEmpty(matrix[row+t.Pos.Y][col+t.Pos.X]) {
					return true
				}
				break
			}
		}
	}
	return false
}
