package tetris

import "fmt"

type Matrix [40][10]byte

func (p *Matrix) isLineComplete(row int) bool {
	for _, cell := range p[row] {
		if isCellEmpty(cell) {
			return false
		}
	}
	return true
}

func (p *Matrix) removeLine(row int) {
	p[0] = [10]byte{}
	for i := row; i > 0; i-- {
		p[i] = p[i-1]
	}
}

func (p *Matrix) RemoveTetrimino(tetrimino *Tetrimino) error {
	for row := range tetrimino.Cells {
		for col := range tetrimino.Cells[row] {
			if tetrimino.Cells[row][col] {
				cellValue := p[row+tetrimino.Pos.Y][col+tetrimino.Pos.X]
				if cellValue != tetrimino.Value {
					return fmt.Errorf("cell at row %d, col %d is '%v' not the expected value", row+tetrimino.Pos.Y, col+tetrimino.Pos.X, cellValue)
				}
				p[row+tetrimino.Pos.Y][col+tetrimino.Pos.X] = 0
			}
		}
	}
	return nil
}

func (p *Matrix) AddTetrimino(tetrimino *Tetrimino) error {
	for row := range tetrimino.Cells {
		for col := range tetrimino.Cells[row] {
			if tetrimino.Cells[row][col] {
				cellValue := p[row+tetrimino.Pos.Y][col+tetrimino.Pos.X]
				if !isCellEmpty(cellValue) {
					return fmt.Errorf("cell at row %d, col %d is '%v' not empty or a ghost", row+tetrimino.Pos.Y, col+tetrimino.Pos.X, cellValue)
				}
				p[row+tetrimino.Pos.Y][col+tetrimino.Pos.X] = tetrimino.Value
			}
		}
	}
	return nil
}

func (p *Matrix) RemoveCompletedLines(tet *Tetrimino) action {
	lines := 0
	for row := range tet.Cells {
		if p.isLineComplete(tet.Pos.Y + row) {
			p.removeLine(tet.Pos.Y + row)
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
