package tetris

import (
	"errors"
	"fmt"
)

// Matrix represents the board of cells on which the game is played.
type Matrix [][]byte

func DefaultMatrix() Matrix {
	m, err := NewMatrix(40, 10)
	if err != nil {
		panic(fmt.Errorf("failed to create default matrix: %w", err))
	}
	return m
}

// NewMatrix creates a new Matrix with the given height and width.
func NewMatrix(height, width uint) (Matrix, error) {
	if height <= 20 {
		return nil, errors.New("matrix height must be greater than 20 to allow for a buffer zone of 20 lines")
	}

	matrix := make(Matrix, height)
	for i := range matrix {
		matrix[i] = make([]byte, width)
	}
	return matrix, nil
}

// GetHeight returns the height of the Matrix.
func (m *Matrix) GetHeight() int {
	return len(*m)
}

// GetSkyline returns the skyline; the highest row that the player can see.
func (m *Matrix) GetSkyline() int {
	return len(*m) - 20
}

// GetVisible returns the Matrix without the buffer zone at the top (ie. the visible portion of the Matrix).
func (m *Matrix) GetVisible() Matrix {
	return (*m)[20:]
}

func (m *Matrix) DeepCopy() *Matrix {
	duplicate := make(Matrix, len(*m))
	for i := range *m {
		duplicate[i] = make([]byte, len((*m)[i]))
		copy(duplicate[i], (*m)[i])
	}
	return &duplicate
}

func (m *Matrix) isLineComplete(row int) bool {
	for _, cell := range (*m)[row] {
		if isCellEmpty(cell) {
			return false
		}
	}
	return true
}

func (m *Matrix) removeLine(row int) {
	(*m)[0] = make([]byte, len((*m)[0]))
	for i := row; i > 0; i-- {
		(*m)[i] = (*m)[i-1]
	}
}

// RemoveTetrimino removes the given Tetrimino from the Matrix.
// It returns an error if the Tetrimino is out of bounds or if the Tetrimino is not found in the Matrix.
func (m *Matrix) RemoveTetrimino(tet *Tetrimino) error {
	isExpectedValue := func(cell byte) bool {
		return cell == tet.Value
	}

	return m.modifyCell(tet.Minos, tet.Pos, 0, isExpectedValue)
}

// AddTetrimino adds the given Tetrimino to the Matrix.
// It returns an error if the Tetrimino is out of bounds or if the Tetrimino overlaps with an occupied mino.
func (m *Matrix) AddTetrimino(tet *Tetrimino) error {
	return m.modifyCell(tet.Minos, tet.Pos, tet.Value, isCellEmpty)
}

func (m *Matrix) modifyCell(minos [][]bool, pos Coordinate, newValue byte, isExpectedValue func(byte) bool) error {
	for row := range minos {
		for col := range minos[row] {
			if minos[row][col] {
				minoAbsRow := row + pos.Y
				minoAbsCol := col + pos.X

				if minoAbsRow >= len(*m) || minoAbsRow < 0 {
					return fmt.Errorf("row %d is out of bounds", minoAbsRow)
				}
				if minoAbsCol >= len((*m)[row]) || minoAbsCol < 0 {
					return fmt.Errorf("col %d is out of bounds", minoAbsCol)
				}

				minoValue := (*m)[minoAbsRow][minoAbsCol]
				if !isExpectedValue(minoValue) {
					err := fmt.Errorf("mino at row %d, col %d is '%s' (byte value %v) not the expected value",
						minoAbsRow, minoAbsCol, string(minoValue), minoValue)
					return errors.Join(err, ErrUnexpectedMatrixCellValue)
				}
				(*m)[minoAbsRow][minoAbsCol] = newValue
			}
		}
	}
	return nil
}

// RemoveCompletedLines checks each row that the given Tetrimino occupies and removes any completed lines from the Matrix.
// It returns an Action to be used for calculating the score.
func (m *Matrix) RemoveCompletedLines(tet *Tetrimino) Action {
	lines := 0
	for row := range tet.Minos {
		if m.isLineComplete(tet.Pos.Y + row) {
			m.removeLine(tet.Pos.Y + row)
			lines++
		}
	}

	switch lines {
	case 0:
		return Actions.NONE
	case 1:
		return Actions.SINGLE
	case 2:
		return Actions.DOUBLE
	case 3:
		return Actions.TRIPLE
	case 4:
		return Actions.TETRIS
	}
	return Actions.UNKNOWN
}

func (m *Matrix) isOutOfBoundsHorizontally(col int) bool {
	return col < 0 || col >= len((*m)[0])
}

func (m *Matrix) isOutOfBoundsVertically(row int) bool {
	return row < 0 || row >= len(*m)
}

func isCellEmpty(cell byte) bool {
	return cell == 0 || cell == 'G'
}

func (m *Matrix) canPlaceInCell(row, col int) bool {
	if m.isOutOfBoundsHorizontally(col) {
		return false
	}
	if m.isOutOfBoundsVertically(row) {
		return false
	}
	if !isCellEmpty((*m)[row][col]) {
		return false
	}
	return true
}
