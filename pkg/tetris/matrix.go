package tetris

import "fmt"

// Matrix represents the board of cells on which the game is played.
type Matrix [][]byte

// NewMatrix creates a new Matrix with the given height and width.
func NewMatrix(height, width uint) (Matrix, error) {
	if height <= 20 {
		return nil, ErrBufferZoneTooSmall
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

// GetStartRow returns the line at which new Tetriminos are created.
// This accounts for the buffer zone at the top of the Matrix.
func (m *Matrix) GetStartRow() int {
	return len(*m) - 20
}

// GetVisible returns the Matrix without the buffer zone at the top (ie. the visible portion of the Matrix).
func (m *Matrix) GetVisible() [][]byte {
	return (*m)[20:]
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
func (m *Matrix) RemoveTetrimino(tetrimino *Tetrimino) error {
	for row := range tetrimino.Cells {
		for col := range tetrimino.Cells[row] {
			if tetrimino.Cells[row][col] {
				// The following bounds checks should never be reached as long as tetriminos are added with the dedicated function
				if row+tetrimino.Pos.Y >= len(*m) || row+tetrimino.Pos.Y < 0 {
					return fmt.Errorf("row %d is out of bounds", row+tetrimino.Pos.Y)
				}
				if col+tetrimino.Pos.X >= len((*m)[row]) || col+tetrimino.Pos.X < 0 {
					return fmt.Errorf("col %d is out of bounds", col+tetrimino.Pos.X)
				}

				cellValue := (*m)[row+tetrimino.Pos.Y][col+tetrimino.Pos.X]
				if cellValue != tetrimino.Value {
					return fmt.Errorf("cell at row %d, col %d is '%s' (byte value %v) not the expected '%s' (byte value %v)", row+tetrimino.Pos.Y, col+tetrimino.Pos.X, string(cellValue), cellValue, string(tetrimino.Value), tetrimino.Value)
				}
				(*m)[row+tetrimino.Pos.Y][col+tetrimino.Pos.X] = 0
			}
		}
	}
	return nil
}

// AddTetrimino adds the given Tetrimino to the Matrix.
// It returns an error if the Tetrimino is out of bounds or if the Tetrimino overlaps with an occupied cell.
func (m *Matrix) AddTetrimino(tetrimino *Tetrimino) error {
	for row := range tetrimino.Cells {
		for col := range tetrimino.Cells[row] {
			if tetrimino.Cells[row][col] {
				if row+tetrimino.Pos.Y >= len(*m) || row+tetrimino.Pos.Y < 0 {
					return fmt.Errorf("row %d is out of bounds", row+tetrimino.Pos.Y)
				}
				if col+tetrimino.Pos.X >= len((*m)[row]) || col+tetrimino.Pos.X < 0 {
					return fmt.Errorf("col %d is out of bounds", col+tetrimino.Pos.X)
				}

				cellValue := (*m)[row+tetrimino.Pos.Y][col+tetrimino.Pos.X]
				if !isCellEmpty(cellValue) {
					return fmt.Errorf("cell at row %d, col %d is '%s' (byte value %v) not empty", row+tetrimino.Pos.Y, col+tetrimino.Pos.X, string(cellValue), cellValue)
				}
				(*m)[row+tetrimino.Pos.Y][col+tetrimino.Pos.X] = tetrimino.Value
			}
		}
	}
	return nil
}

// RemoveCompletedLines checks each row that the given Tetrimino occupies and removes any completed lines from the Matrix.
// It returns an action to be used for calculating the score.
func (m *Matrix) RemoveCompletedLines(tet *Tetrimino) action {
	lines := 0
	for row := range tet.Cells {
		if m.isLineComplete(tet.Pos.Y + row) {
			m.removeLine(tet.Pos.Y + row)
			lines++
		}
	}

	switch lines {
	case 0:
		return actionNone
	case 1:
		return actionSingle
	case 2:
		return actionDouble
	case 3:
		return actionTriple
	case 4:
		return actionTetris
	}
	return actionNone
}