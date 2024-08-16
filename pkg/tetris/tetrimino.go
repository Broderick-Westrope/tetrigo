package tetris

import (
	"errors"
	"fmt"
	"math"
)

// A Tetrimino is a geometric Tetris shape formed by four Minos connected along their sides.
type Tetrimino struct {
	Value byte       // The value of the Tetrimino. This is the character that will be used to represent the Tetrimino in the matrix.
	Minos [][]bool   // A 2D slice of cells that make up the Tetrimino. True means the mino is occupied by a Mino.
	Pos   Coordinate // The top left mino of the Tetrimino. Used as a reference point for movement and rotation.

	CurrentRotation int // The index of the current rotation in the RotationCompass.
	RotationCompass rotationCompass
}

// Coordinate represents a point on a 2D plane.
type Coordinate struct {
	X, Y int
}

// A rotatationCompass contains a rotationSet corresponding to each of the four compass directions in the order N, E, S, W.
// These compass directions represent the four rotations of a Tetrimino.
type rotationCompass [4]rotationSet

// A rotationSet contains coordinates to be used for a single rotation/compass direction.
// If the first coordinate cannot be used the next will be attempted.
// This continues until there are no more coordinates to fall back on (in which case rotation is not possible).
// This is part of the Super Rotation System (SRS).
type rotationSet []Coordinate

// RotationCompasses is a map of Tetrimino values to the coordinates used for rotation.
// Each slice should contain a coordinate for north, east, south, and west in that order.
// These are added to (clockwise) or subtracted from (counter-clockwise) the Tetrimino's position when rotating to ensure it rotates around the correct axis.
var RotationCompasses = map[byte]rotationCompass{
	'I': {
		{ // North
			{X: -1, Y: 1},
		},
		{ // East
			{X: 2, Y: -1},
		},
		{ // South
			{X: -2, Y: 2},
		},
		{ // West
			{X: 1, Y: -2},
		},
	},
	'O': {
		{ // North
			{X: 0, Y: 0},
		},
		{ // East
			{X: 0, Y: 0},
		},
		{ // South
			{X: 0, Y: 0},
		},
		{ // West
			{X: 0, Y: 0},
		},
	},
	// All tetriminos with 6 cells (T, S, Z, J, L) have the same rotation compass:
	'6': {
		{ // North
			{X: 0, Y: 0},
		},
		{ // East
			{X: 1, Y: 0},
		},
		{ // South
			{X: -1, Y: 1},
		},
		{ // West
			{X: 0, Y: -1},
		},
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

var validTetriminos = map[byte]Tetrimino{
	'I': {
		Value: 'I',
		Minos: [][]bool{
			{true, true, true, true},
		},
		Pos:             startingPositions['I'],
		RotationCompass: RotationCompasses['I'],
	},
	'O': {
		Value: 'O',
		Minos: [][]bool{
			{true, true},
			{true, true},
		},
		Pos:             startingPositions['O'],
		RotationCompass: RotationCompasses['O'],
	},
	'T': {
		Value: 'T',
		Minos: [][]bool{
			{false, true, false},
			{true, true, true},
		},
		Pos:             startingPositions['6'],
		RotationCompass: RotationCompasses['6'],
	},
	'S': {
		Value: 'S',
		Minos: [][]bool{
			{false, true, true},
			{true, true, false},
		},
		Pos:             startingPositions['6'],
		RotationCompass: RotationCompasses['6'],
	},
	'Z': {
		Value: 'Z',
		Minos: [][]bool{
			{true, true, false},
			{false, true, true},
		},
		Pos:             startingPositions['6'],
		RotationCompass: RotationCompasses['6'],
	},
	'J': {
		Value: 'J',
		Minos: [][]bool{
			{true, false, false},
			{true, true, true},
		},
		Pos:             startingPositions['6'],
		RotationCompass: RotationCompasses['6'],
	},
	'L': {
		Value: 'L',
		Minos: [][]bool{
			{false, false, true},
			{true, true, true},
		},
		Pos:             startingPositions['6'],
		RotationCompass: RotationCompasses['6'],
	},
}

// GetValidTetriminos returns a slice containing all seven valid Tetriminos (I, O, T, S, Z, J, L)
func GetValidTetriminos() []Tetrimino {
	result := make([]Tetrimino, 0, len(validTetriminos))
	for _, t := range validTetriminos {
		result = append(result, *t.DeepCopy())
	}
	return result
}

// GetTetrimino returns the Tetrmino with the given value.
// Valid values include: I, O, T, S, Z, J, L.
func GetTetrimino(value byte) (*Tetrimino, error) {
	result, ok := validTetriminos[value]
	if !ok {
		return nil, errors.New("invalid value")
	}
	return result.DeepCopy(), nil
}

// GetEmptyTetrimino returns a tetrimino with no minos or value. To be used for the starting (empty) hold.
func GetEmptyTetrimino() *Tetrimino {
	return &Tetrimino{
		Minos: [][]bool{},
		Value: 0,
	}
}

// MoveDown moves the tetrimino down one row.
// This does not modify the matrix.
// If the tetrimino cannot move down, it will not be modified and false will be returned.
func (t *Tetrimino) MoveDown(matrix Matrix) bool {
	if t.canMoveDown(matrix) {
		t.Pos.Y++
		return true
	}
	return false
}

// MoveLeft moves the tetrimino left one column.
// This does not modify the matrix.
// If the tetrimino cannot move left false will be returned.
func (t *Tetrimino) MoveLeft(matrix Matrix) bool {
	if t.canMoveLeft(matrix) {
		t.Pos.X--
		return true
	}
	return false
}

// MoveRight moves the tetrimino right one column.
// This does not modify the matrix.
// If the tetrimino cannot move right false will be returned.
func (t *Tetrimino) MoveRight(matrix Matrix) bool {
	if t.canMoveRight(matrix) {
		t.Pos.X++
		return true
	}
	return false
}

// Returns true if the tetrimino can move down one row.
// This gets the lowest mino in each column of the tetrimino, and checks if it is at the bottom of the matrix or if the mino below is occupied.
func (t *Tetrimino) canMoveDown(matrix Matrix) bool {
	for col := range t.Minos[0] {
		for row := len(t.Minos) - 1; row >= 0; row-- {
			if !t.Minos[row][col] {
				continue
			}
			if !matrix.canPlaceInCell(row+t.Pos.Y+1, col+t.Pos.X) {
				return false
			}
			break
		}
	}
	return true
}

func (t *Tetrimino) canMoveLeft(matrix Matrix) bool {
	for row := range t.Minos {
		for col := range t.Minos[row] {
			if !t.Minos[row][col] {
				continue
			}
			if !matrix.canPlaceInCell(row+t.Pos.Y, col+t.Pos.X-1) {
				return false
			}
			break
		}
	}
	return true
}

func (t *Tetrimino) canMoveRight(matrix Matrix) bool {
	for row := range t.Minos {
		for col := len(t.Minos[row]) - 1; col >= 0; col-- {
			if !t.Minos[row][col] {
				continue
			}
			if !matrix.canPlaceInCell(row+t.Pos.Y, col+t.Pos.X+1) {
				return false
			}
			break
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
	var foundValid bool
	if clockwise {
		foundValid, err = rotated.rotateClockwise(matrix)
	} else {
		foundValid, err = rotated.rotateCounterClockwise(matrix)
	}
	if err != nil {
		return fmt.Errorf("failed to rotate tetrimino: %w", err)
	}

	if foundValid {
		t.Pos = rotated.Pos
		t.Minos = rotated.Minos
		t.CurrentRotation = rotated.CurrentRotation
	}

	return nil
}

func (t *Tetrimino) rotateClockwise(matrix Matrix) (bool, error) {
	// Reverse the order of the rows
	for i, j := 0, len(t.Minos)-1; i < j; i, j = i+1, j-1 {
		t.Minos[i], t.Minos[j] = t.Minos[j], t.Minos[i]
	}

	t.transpose()

	var err error
	t.CurrentRotation, err = positiveMod(t.CurrentRotation+1, len(t.RotationCompass))
	if err != nil {
		return false, fmt.Errorf("failed to get positive mod: %w", err)
	}

	originalX, originalY := t.Pos.X, t.Pos.Y
	for _, coord := range t.RotationCompass[t.CurrentRotation] {
		t.Pos.X = originalX + coord.X
		t.Pos.Y = originalY + coord.Y

		if t.isValid(matrix) {
			return true, nil
		}
	}

	return false, nil
}

func (t *Tetrimino) rotateCounterClockwise(matrix Matrix) (bool, error) {
	// Reverse the order of the columns
	for _, row := range t.Minos {
		for i, j := 0, len(row)-1; i < j; i, j = i+1, j-1 {
			row[i], row[j] = row[j], row[i]
		}
	}

	t.transpose()

	var err error
	t.CurrentRotation, err = positiveMod(t.CurrentRotation+1, len(t.RotationCompass))
	if err != nil {
		return false, fmt.Errorf("failed to get positive mod: %w", err)
	}

	originalX, originalY := t.Pos.X, t.Pos.Y
	for _, coord := range t.RotationCompass[t.CurrentRotation] {
		t.Pos.X = originalX - coord.X
		t.Pos.Y = originalY - coord.Y

		if t.isValid(matrix) {
			return true, nil
		}
	}

	return false, nil
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

// isValid returns true if the given Tetrimino is within the bounds of the matrix and does not overlap with any occupied cells.
// The Tetrimino being checked should not be in the Matrix yet.
func (t *Tetrimino) isValid(matrix Matrix) bool {
	for row := range t.Minos {
		for col := range t.Minos[row] {
			if !t.Minos[row][col] {
				continue
			}
			if !matrix.canPlaceInCell(row+t.Pos.Y, col+t.Pos.X) {
				return false
			}
		}
	}
	return true
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

func (t *Tetrimino) DeepCopy() *Tetrimino {
	var cells [][]bool
	if t.Minos == nil {
		cells = nil
	} else {
		cells = deepCopyMinos(t.Minos)
	}

	var compass rotationCompass
	for i := range t.RotationCompass {
		compass[i] = make(rotationSet, len(t.RotationCompass[i]))
		copy(compass[i], t.RotationCompass[i])
	}

	return &Tetrimino{
		Value:           t.Value,
		Minos:           cells,
		Pos:             t.Pos,
		CurrentRotation: t.CurrentRotation,
		RotationCompass: compass,
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
			if !t.Minos[row][col] {
				continue
			}
			if !isCellEmpty(matrix[row+t.Pos.Y][col+t.Pos.X]) {
				return true
			}
		}
	}
	return false
}
