package tetris

import (
	"math/rand"
)

// NextQueue is a collection of up to 14 Tetriminos that are drawn from randomly.
// The queue is refilled when it has less than 7 Tetriminos.
type NextQueue struct {
	elements  []Tetrimino
	startLine int
}

// NewNextQueue creates a new NextQueue of Tetriminos.
func NewNextQueue(startLine int) *NextQueue {
	nq := NextQueue{
		elements:  make([]Tetrimino, 0, 14),
		startLine: startLine,
	}
	nq.fill()
	nq.fill()
	return &nq
}

func (nq *NextQueue) GetElements() []Tetrimino {
	return nq.elements
}

// Next returns the next Tetrimino, removing it from the queue and refilling if necessary.
// This applies the skyline value (provided in NewNextQueue) to the Tetriminos Y axis.
func (nq *NextQueue) Next() *Tetrimino {
	tet := nq.elements[0]
	nq.elements = nq.elements[1:]

	if len(nq.elements) <= 7 {
		nq.fill()
	}

	tet.Pos.Y += nq.startLine
	return &tet
}

func (nq *NextQueue) fill() {
	if len(nq.elements) > 7 {
		return
	}

	tetriminos := GetValidTetriminos()
	perm := rand.Perm(len(tetriminos))
	for _, i := range perm {
		if len(nq.elements) == 14 {
			// This should be impossible whilst there are only 7 Tetriminos and we check that there is space for 7 in the queue
			return
		}
		nq.elements = append(nq.elements, tetriminos[i])
	}
}
