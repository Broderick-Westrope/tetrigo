package main

import "fmt"

type Playfield [40][10]byte

func (p *Playfield) isLineComplete(row int) bool {
	for _, cell := range p[row] {
		if isCellEmpty(cell) {
			return false
		}
	}
	return true
}

func (p *Playfield) removeLine(row int) {
	p[0] = [10]byte{}
	for i := row; i > 0; i-- {
		p[i] = p[i-1]
	}
}

func (p *Playfield) removeCells(tetrimino *Tetrimino) error {
	for row := range tetrimino.Cells {
		for col := range tetrimino.Cells[row] {
			if tetrimino.Cells[row][col] {
				if v := p[row+tetrimino.Pos.Y][col+tetrimino.Pos.X]; v != tetrimino.Value {
					return fmt.Errorf("cell at row %d, col %d is not the expected value", row+tetrimino.Pos.Y, col+tetrimino.Pos.X)
				}
				p[row+tetrimino.Pos.Y][col+tetrimino.Pos.X] = 0
			}
		}
	}
	return nil
}

func (p *Playfield) addCells(tetrimino *Tetrimino) error {
	for row := range tetrimino.Cells {
		for col := range tetrimino.Cells[row] {
			if tetrimino.Cells[row][col] {
				if !isCellEmpty(p[row+tetrimino.Pos.Y][col+tetrimino.Pos.X]) {
					return fmt.Errorf("cell at row %d, col %d is not empty or a ghost", row+tetrimino.Pos.Y, col+tetrimino.Pos.X)
				}
				p[row+tetrimino.Pos.Y][col+tetrimino.Pos.X] = tetrimino.Value
			}
		}
	}
	return nil
}

func (p *Playfield) NewTetrimino() *Tetrimino {
	tet := RandomTetrimino(len(p))
	p.addCells(tet)
	return tet
}

func (p *Playfield) removeCompletedLines(tet *Tetrimino) action {
	for row := range tet.Cells {
		if p.isLineComplete(tet.Pos.Y + row) {
			p.removeLine(tet.Pos.Y + row)
		}
	}
	return actionNone
}
