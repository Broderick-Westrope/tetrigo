package tetris

import (
	"math/rand"
)

type Bag struct {
	Elements     []Tetrimino
	matrixHeight int
}

func NewBag(matrixHeight int) *Bag {
	b := Bag{
		Elements:     make([]Tetrimino, 0, 14),
		matrixHeight: matrixHeight,
	}
	b.fill()
	b.fill()
	return &b
}

func (b *Bag) Next() *Tetrimino {
	tet := b.Elements[0]
	b.Elements = b.Elements[1:]

	if len(b.Elements) <= 7 {
		b.fill()
	}

	tet.Pos.Y += b.matrixHeight - 20
	return &tet
}

func (b *Bag) fill() {
	if cap(b.Elements)-len(b.Elements) < 7 {
		return
	}

	perm := rand.Perm(len(Tetriminos))
	for _, i := range perm {
		if len(b.Elements) == 14 {
			return
		}
		b.Elements = append(b.Elements, Tetriminos[i])
	}
}
