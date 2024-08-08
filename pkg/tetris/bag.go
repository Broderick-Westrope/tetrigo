package tetris

import (
	"math/rand"
)

// Bag is a collection of up to 14 Tetriminos that are drawn from randomly.
// The bag is refilled when it has less than 7 Tetriminos.
type Bag struct {
	elements  []Tetrimino
	startLine int
}

// NewBag creates a new Bag of Tetriminos.
func NewBag(startLine int) *Bag {
	bag := Bag{
		elements:  make([]Tetrimino, 0, 14),
		startLine: startLine,
	}
	bag.fill()
	bag.fill()
	return &bag
}

func (b *Bag) GetElements() []Tetrimino {
	return b.elements
}

// Next returns the next Tetrimino, removing it from the bag and fill the bag if necessary.
func (b *Bag) Next() *Tetrimino {
	tet := b.elements[0]
	b.elements = b.elements[1:]

	if len(b.elements) <= 7 {
		b.fill()
	}

	tet.Pos.Y += b.startLine
	return &tet
}

func (b *Bag) fill() {
	if len(b.elements) > 7 {
		return
	}

	perm := rand.Perm(len(Tetriminos))
	for _, i := range perm {
		if len(b.elements) == 14 {
			// This is impossible whilst there are only 7 tetriminos and we check that there is space for 7 in the bag
			return
		}
		b.elements = append(b.elements, Tetriminos[i])
	}
}
