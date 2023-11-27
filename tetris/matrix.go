package tetris

import "fmt"

type Matrix [40][10]byte

func (m *Matrix) isLineComplete(row int) bool {
	for _, cell := range m[row] {
		if isCellEmpty(cell) {
			return false
		}
	}
	return true
}

func (m *Matrix) removeLine(row int) {
	m[0] = [10]byte{}
	for i := row; i > 0; i-- {
		m[i] = m[i-1]
	}
}

func (m *Matrix) RemoveTetrimino(tetrimino *Tetrimino) error {
	for row := range tetrimino.Cells {
		for col := range tetrimino.Cells[row] {
			if tetrimino.Cells[row][col] {
				// The following bounds checks should never be reached as long as tetriminos are added with the dedicated function
				if row+tetrimino.Pos.Y >= len(m) || row+tetrimino.Pos.Y < 0 {
					return fmt.Errorf("row %d is out of bounds", row+tetrimino.Pos.Y)
				}
				if col+tetrimino.Pos.X >= len(m[row]) || col+tetrimino.Pos.X < 0 {
					return fmt.Errorf("col %d is out of bounds", col+tetrimino.Pos.X)
				}

				cellValue := m[row+tetrimino.Pos.Y][col+tetrimino.Pos.X]
				if cellValue != tetrimino.Value {
					return fmt.Errorf("cell at row %d, col %d is '%s' (byte value %v) not the expected '%s' (byte value %v)", row+tetrimino.Pos.Y, col+tetrimino.Pos.X, string(cellValue), cellValue, string(tetrimino.Value), tetrimino.Value)
				}
				m[row+tetrimino.Pos.Y][col+tetrimino.Pos.X] = 0
			}
		}
	}
	return nil
}

func (m *Matrix) AddTetrimino(tetrimino *Tetrimino) error {
	for row := range tetrimino.Cells {
		for col := range tetrimino.Cells[row] {
			if tetrimino.Cells[row][col] {
				if row+tetrimino.Pos.Y >= len(m) || row+tetrimino.Pos.Y < 0 {
					return fmt.Errorf("row %d is out of bounds", row+tetrimino.Pos.Y)
				}
				if col+tetrimino.Pos.X >= len(m[row]) || col+tetrimino.Pos.X < 0 {
					return fmt.Errorf("col %d is out of bounds", col+tetrimino.Pos.X)
				}

				cellValue := m[row+tetrimino.Pos.Y][col+tetrimino.Pos.X]
				if !isCellEmpty(cellValue) {
					return fmt.Errorf("cell at row %d, col %d is '%s' (byte value %v) not empty", row+tetrimino.Pos.Y, col+tetrimino.Pos.X, string(cellValue), cellValue)
				}
				m[row+tetrimino.Pos.Y][col+tetrimino.Pos.X] = tetrimino.Value
			}
		}
	}
	return nil
}

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
