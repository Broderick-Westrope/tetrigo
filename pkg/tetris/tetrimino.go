package tetris

import (
	"errors"
	"fmt"
	"math"
)

const invalidRotationPoint = -1

// A Tetrimino is a geometric Tetris piece formed by connecting four square blocks (Minos) along their edges.
// Each Tetrimino has a unique shape, position, rotation state, and rotation behavior defined by the Super Rotation System (SRS).
type Tetrimino struct {
	// Value is the character identifier for the Tetrimino (I, O, T, S, Z, J, L).
	// This is used internally and may differ from the display representation.
	Value byte

	// Cells represents the shape of the Tetrimino as a 2D grid.
	// True indicates a Mino is present in the cell, false indicates empty space.
	Cells [][]bool

	// Position tracks the coordinate of the Tetrimino's top-left corner in the game matrix.
	// This serves as the reference point for all movement and rotation operations.
	Position Coordinate

	// CompassDirection tracks the current rotation state (0-3, representing North, East, South, West).
	CompassDirection int

	// RotationCompass defines the piece's rotation behavior according to the Super Rotation System (SRS).
	// It contains offset data for each rotation state to handle various kicks.
	RotationCompass RotationCompass
}

// Coordinate represents a position in the game matrix using X (horizontal) and Y (vertical) coordinates.
// The origin (0,0) is the top-left of a matrix (within any buffer zone that exists),
// with positive X going right and positive Y going down.
type Coordinate struct {
	X, Y int
}

// RotationCompass defines rotation offset data for all four possible orientations
// of a Tetrimino (North, East, South, West). This implements the Super Rotation System (SRS)
// which allows pieces to "kick" off walls and other blocks when rotating.
type RotationCompass [4]RotationOffsets

// RotationOffsets contains a sequence of offset coordinates to try when rotating a Tetrimino.
// When rotating, these offsets are tried in order until a valid position is found.
// If no valid position is found using any offset, the rotation fails.
// This is a key component of the Super Rotation System (SRS).
type RotationOffsets []*Coordinate

// RotationCompasses maps each Tetrimino type to its rotation behavior definition.
// The offsets are used to adjust piece position during rotation to implement kicks
// according to the Super Rotation System (SRS).
//
// When rotating clockwise (right), the offsets are applied in order.
// When rotating counter-clockwise (left), the offsets are applied in reverse order.
// As an example, when rotating the I Tetrimino clockwise from North to East
// (clockwise) the first offset added to its position is (2, -1). If this position is not valid,
// the next offset (0, -1) is tried, and so on.
var RotationCompasses = map[byte]RotationCompass{
	'I': {
		{ // North
			{X: -1, Y: 1}, {X: 0, Y: 1}, {X: -3, Y: 1}, {X: 0, Y: 3}, {X: -3, Y: 0},
		},
		{ // East
			{X: 2, Y: -1}, {X: 0, Y: -1}, {X: 3, Y: -1}, {X: 0, Y: 0}, {X: 3, Y: -3},
		},
		{ // South
			{X: -2, Y: 2}, {X: -3, Y: 2}, {X: 0, Y: 2}, {X: -3, Y: 0}, {X: 0, Y: 3},
		},
		{ // West
			{X: 1, Y: -2}, {X: 3, Y: -2}, {X: 0, Y: -2}, {X: 3, Y: -3}, {X: 0, Y: 0},
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
	// All tetriminos with a 2x3 grid (6 total) of minos (T, S, Z, J, L) have the same rotation compass:
	'6': {
		{ // North
			{X: 0, Y: 0}, {X: -1, Y: 0}, {X: -1, Y: 1}, {X: 0, Y: -2}, {X: -1, Y: -2},
		},
		{ // East
			{X: 1, Y: 0}, {X: 0, Y: 0}, {X: 0, Y: -1}, {X: 1, Y: 2}, {X: 0, Y: 2},
		},
		{ // South
			{X: -1, Y: 1}, {X: 0, Y: 1}, {X: 0, Y: 2}, {X: -1, Y: -1}, {X: 0, Y: -1},
		},
		{ // West
			{X: 0, Y: -1}, {X: 1, Y: -1}, {X: 1, Y: -2}, {X: 0, Y: 1}, {X: 1, Y: 1},
		},
	},
}

// startingPositions defines the initial spawn position for each Tetrimino type when it enters the game matrix.
// This position is relative to the matrix buffer zone, meaning the negative Y coordinates account for
// the pieces spawning partially above the visible playfield.
var startingPositions = map[byte]Coordinate{
	'I': {X: 3, Y: -1},
	'O': {X: 4, Y: -2},
	'6': {X: 3, Y: -2},
}

func getMapOfValidTetriminos() map[byte]Tetrimino {
	return map[byte]Tetrimino{
		'I': {
			Value: 'I',
			Cells: [][]bool{
				{true, true, true, true},
			},
			Position:        startingPositions['I'],
			RotationCompass: RotationCompasses['I'],
		},
		'O': {
			Value: 'O',
			Cells: [][]bool{
				{true, true},
				{true, true},
			},
			Position:        startingPositions['O'],
			RotationCompass: RotationCompasses['O'],
		},
		'T': {
			Value: 'T',
			Cells: [][]bool{
				{false, true, false},
				{true, true, true},
			},
			Position:        startingPositions['6'],
			RotationCompass: RotationCompasses['6'],
		},
		'S': {
			Value: 'S',
			Cells: [][]bool{
				{false, true, true},
				{true, true, false},
			},
			Position:        startingPositions['6'],
			RotationCompass: RotationCompasses['6'],
		},
		'Z': {
			Value: 'Z',
			Cells: [][]bool{
				{true, true, false},
				{false, true, true},
			},
			Position:        startingPositions['6'],
			RotationCompass: RotationCompasses['6'],
		},
		'J': {
			Value: 'J',
			Cells: [][]bool{
				{true, false, false},
				{true, true, true},
			},
			Position:        startingPositions['6'],
			RotationCompass: RotationCompasses['6'],
		},
		'L': {
			Value: 'L',
			Cells: [][]bool{
				{false, false, true},
				{true, true, true},
			},
			Position:        startingPositions['6'],
			RotationCompass: RotationCompasses['6'],
		},
	}
}

// GetValidTetriminos returns a slice containing all seven valid Tetriminos (I, O, T, S, Z, J, L).
func GetValidTetriminos() []Tetrimino {
	validTetriminos := getMapOfValidTetriminos()
	result := make([]Tetrimino, 0, len(validTetriminos))
	for _, t := range validTetriminos {
		result = append(result, t)
	}
	return result
}

// GetTetrimino returns the Tetrmino with the given value.
// Valid values are: I, O, T, S, Z, J, L.
func GetTetrimino(value byte) (*Tetrimino, error) {
	result, ok := getMapOfValidTetriminos()[value]
	if !ok {
		return nil, errors.New("invalid value")
	}
	return result.DeepCopy(), nil
}

// GetEmptyTetrimino returns a tetrimino to be used for the starting (empty) hold.
func GetEmptyTetrimino() *Tetrimino {
	return &Tetrimino{
		Cells: [][]bool{},
		Value: 0,
	}
}

// MoveDown moves the tetrimino down one row.
// This does not modify the matrix.
// If the tetrimino cannot move down, it will not be modified and false will be returned.
func (t *Tetrimino) MoveDown(matrix Matrix) bool {
	// Loop through the cells of each column from bottom to top to find the lowest Mino in each column.
	for col := range t.Cells[0] {
		for row := len(t.Cells) - 1; row >= 0; row-- {
			if !t.Cells[row][col] {
				continue
			}
			if !matrix.canPlaceInCell(row+t.Position.Y+1, col+t.Position.X) {
				return false
			}
			break
		}
	}

	t.Position.Y++
	return true
}

// MoveLeft moves the tetrimino left one column.
// This does not modify the matrix.
// If the tetrimino cannot move left false will be returned.
func (t *Tetrimino) MoveLeft(matrix Matrix) bool {
	// Loop through the cells of each row from left to right to find the rightmost Mino in each row.
	for row := range t.Cells {
		for col := range t.Cells[row] {
			if !t.Cells[row][col] {
				continue
			}
			if !matrix.canPlaceInCell(row+t.Position.Y, col+t.Position.X-1) {
				return false
			}
			break
		}
	}

	t.Position.X--
	return true
}

// MoveRight moves the tetrimino right one column.
// This does not modify the matrix.
// If the tetrimino cannot move right false will be returned.
func (t *Tetrimino) MoveRight(matrix Matrix) bool {
	// Loop through the cells of each row from right to left to find the leftmost Mino in each row.
	for row := range t.Cells {
		for col := len(t.Cells[row]) - 1; col >= 0; col-- {
			if !t.Cells[row][col] {
				continue
			}
			if !matrix.canPlaceInCell(row+t.Position.Y, col+t.Position.X+1) {
				return false
			}
			break
		}
	}

	t.Position.X++
	return true
}

// Rotate rotates the Tetrimino clockwise or counter-clockwise.
// This does not modify the matrix.
// This will automatically use Super Rotation System (SRS).
// If no valid rotation is found, the Tetrimino will not be modified and an error will be returned.
func (t *Tetrimino) Rotate(matrix Matrix, clockwise bool) error {
	if t.Value == 'O' {
		// O Tetrimino does not rotate.
		return nil
	}

	rotated := t.DeepCopy()
	var err error
	var rotationPoint int
	if clockwise {
		rotationPoint, err = rotated.rotateClockwise(matrix)
	} else {
		rotationPoint, err = rotated.rotateCounterClockwise(matrix)
	}
	if err != nil {
		return fmt.Errorf("failed to rotate tetrimino: %w", err)
	}

	foundValid := rotationPoint != invalidRotationPoint
	if !foundValid {
		return nil
	}

	t.Position = rotated.Position
	t.Cells = rotated.Cells
	t.CompassDirection = rotated.CompassDirection
	return nil
}

// rotateClockwise rotates the Tetrimino clockwise.
// This does not modify the matrix.
// If a valid rotation is found, the rotation point is returned.
// If no valid rotation is found, invalidRotationPoint is returned.
// This rotation is done by reversing the order of the rows then transposing the matrix.
func (t *Tetrimino) rotateClockwise(matrix Matrix) (int, error) {
	// Reverse the order of the rows
	for i, j := 0, len(t.Cells)-1; i < j; i, j = i+1, j-1 {
		t.Cells[i], t.Cells[j] = t.Cells[j], t.Cells[i]
	}

	t.transpose()

	var err error
	t.CompassDirection, err = positiveMod(t.CompassDirection+1, len(t.RotationCompass))
	if err != nil {
		return invalidRotationPoint, fmt.Errorf("getting positive mod: %w", err)
	}

	originalX, originalY := t.Position.X, t.Position.Y
	for i, coord := range t.RotationCompass[t.CompassDirection] {
		t.Position.X = originalX + coord.X
		t.Position.Y = originalY + coord.Y

		if t.IsValid(matrix, true) {
			return i + 1, nil
		}
	}

	return invalidRotationPoint, nil
}

// rotateCounterClockwise rotates the Tetrimino counter-clockwise.
// This does not modify the matrix.
// If a valid rotation is found, the rotation point is returned.
// If no valid rotation is found, invalidRotationPoint is returned.
// This rotation is done by reversing the order of the columns then transposing the matrix.
func (t *Tetrimino) rotateCounterClockwise(matrix Matrix) (int, error) {
	// Reverse the order of the columns
	for _, row := range t.Cells {
		for i, j := 0, len(row)-1; i < j; i, j = i+1, j-1 {
			row[i], row[j] = row[j], row[i]
		}
	}

	t.transpose()

	rotationPoint := invalidRotationPoint
	originalX, originalY := t.Position.X, t.Position.Y
	for i, coord := range t.RotationCompass[t.CompassDirection] {
		t.Position.X = originalX - coord.X
		t.Position.Y = originalY - coord.Y

		if t.IsValid(matrix, true) {
			rotationPoint = i + 1
			break
		}
	}
	if rotationPoint == invalidRotationPoint {
		return rotationPoint, nil
	}

	var err error
	t.CompassDirection, err = positiveMod(t.CompassDirection-1, len(t.RotationCompass))
	if err != nil {
		return invalidRotationPoint, fmt.Errorf("getting positive mod: %w", err)
	}

	return rotationPoint, nil
}

// transpose transposes the Tetrimino's cells.
// This is used during rotation to rotate the Tetrimino.
func (t *Tetrimino) transpose() {
	xl := len(t.Cells[0])
	yl := len(t.Cells)
	result := make([][]bool, xl)
	for i := range result {
		result[i] = make([]bool, yl)
	}
	for i := range xl {
		for j := range yl {
			result[i][j] = t.Cells[j][i]
		}
	}
	t.Cells = result
}

// positiveMod calculates the positive modulus of the dividend and divisor.
// This is used to handle negative numbers in a way that is consistent with the Super Rotation System (SRS).
// If the result is NaN, an error is returned.
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

// deepCopyCells creates a deep copy of the given 2D grid of cells.
func deepCopyCells(cells [][]bool) [][]bool {
	cellsCopy := make([][]bool, len(cells))
	for i := range cells {
		cellsCopy[i] = make([]bool, len(cells[i]))
		copy(cellsCopy[i], cells[i])
	}
	return cellsCopy
}

// DeepCopy creates a deep copy of the Tetrimino.
// This is useful when you need to modify a Tetrimino without affecting the original.
func (t *Tetrimino) DeepCopy() *Tetrimino {
	var cells [][]bool
	if t.Cells == nil {
		cells = nil
	} else {
		cells = deepCopyCells(t.Cells)
	}

	var compass RotationCompass
	for i := range t.RotationCompass {
		if t.RotationCompass[i] == nil {
			compass[i] = nil
			continue
		}

		compass[i] = make(RotationOffsets, len(t.RotationCompass[i]))
		for j := range t.RotationCompass[i] {
			compass[i][j] = &Coordinate{
				X: t.RotationCompass[i][j].X,
				Y: t.RotationCompass[i][j].Y,
			}
		}
	}

	return &Tetrimino{
		Value:            t.Value,
		Cells:            cells,
		Position:         t.Position,
		CompassDirection: t.CompassDirection,
		RotationCompass:  compass,
	}
}

// IsAboveSkyline returns true if the entire Tetrimino is above the skyline.
// This can be helpful when checking for Lock Out.
func (t *Tetrimino) IsAboveSkyline(skyline int) bool {
	for row := len(t.Cells) - 1; row >= 0; row-- {
		for col := range t.Cells[row] {
			if t.Cells[row][col] {
				if t.Position.Y+row >= skyline {
					return false
				}
				break
			}
		}
	}
	return true
}

// IsValid returns true if the given Tetrimino is within the bounds of the matrix and
// does not overlap with any occupied cells.
// If checkBounds is false, only overlap is checked.
// The Tetrimino being checked should not be in the Matrix yet.
func (t *Tetrimino) IsValid(matrix Matrix, checkBounds bool) bool {
	for row := range t.Cells {
		for col := range t.Cells[row] {
			if !t.Cells[row][col] {
				continue
			}

			if checkBounds {
				// Check for out of bounds and overlap
				if !matrix.canPlaceInCell(row+t.Position.Y, col+t.Position.X) {
					return false
				}
			} else {
				// Only check for overlap
				if !isCellEmpty(matrix[row+t.Position.Y][col+t.Position.X]) {
					return false
				}
			}
		}
	}
	return true
}
