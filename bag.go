package main

import (
	"fmt"
	"math/rand"
)

var RotationCoords = map[byte][]Coordinate{
	'I': {
		{X: -1, Y: -1},
		{X: 2, Y: 1},
		{X: -2, Y: -2},
		{X: 1, Y: 2},
	},
	'O': {
		{X: 0, Y: 0},
		{X: 0, Y: 0},
		{X: 0, Y: 0},
		{X: 0, Y: 0},
	},
	'6': { // All tetriminos with 6 cells (T, S, Z, J, L)
		{X: 0, Y: 0},
		{X: 1, Y: 0},
		{X: -1, Y: -1},
		{X: 0, Y: 1},
	},
}

var startingPositions = map[byte]Coordinate{
	'I': {X: 3, Y: -1},
	'O': {X: 4, Y: -2},
	'6': {X: 3, Y: -2},
}

var tetriminos = []Tetrimino{
	{
		Value: 'I',
		Cells: [][]bool{
			{true, true, true, true},
		},
		Pos:            startingPositions['I'],
		RotationCoords: RotationCoords['I'],
	},
	{
		Value: 'O',
		Cells: [][]bool{
			{true, true},
			{true, true},
		},
		Pos:            startingPositions['O'],
		RotationCoords: RotationCoords['O'],
	},
	{
		Value: 'T',
		Cells: [][]bool{
			{false, true, false},
			{true, true, true},
		},
		Pos:            startingPositions['6'],
		RotationCoords: RotationCoords['6'],
	},
	{
		Value: 'S',
		Cells: [][]bool{
			{false, true, true},
			{true, true, false},
		},
		Pos:            startingPositions['6'],
		RotationCoords: RotationCoords['6'],
	},
	{
		Value: 'Z',
		Cells: [][]bool{
			{true, true, false},
			{false, true, true},
		},
		Pos:            startingPositions['6'],
		RotationCoords: RotationCoords['6'],
	},
	{
		Value: 'J',
		Cells: [][]bool{
			{true, false, false},
			{true, true, true},
		},
		Pos:            startingPositions['6'],
		RotationCoords: RotationCoords['6'],
	},
	{
		Value: 'L',
		Cells: [][]bool{
			{false, false, true},
			{true, true, true},
		},
		Pos:            startingPositions['6'],
		RotationCoords: RotationCoords['6'],
	},
}

type bag struct {
	elements        []*Tetrimino
	playfieldHeight int
}

func defaultBag(playfieldHeight int) *bag {
	b := bag{
		elements:        make([]*Tetrimino, 0, 14),
		playfieldHeight: playfieldHeight,
	}
	b.fillBag()
	b.fillBag()
	return &b
}

func (b *bag) next() *Tetrimino {
	tet := b.elements[0]
	b.elements = b.elements[1:]

	if len(b.elements) <= 7 {
		b.fillBag()
	}

	tet.Pos.Y += b.playfieldHeight - 20
	return tet
}

func (b *bag) fillBag() {
	perm := rand.Perm(len(tetriminos))
	for _, i := range perm {
		if len(b.elements) == 14 {
			return
		}
		b.enqueue(&tetriminos[i])
	}
}

// Adds an item to the end of the queue
func (b *bag) enqueue(value *Tetrimino) {
	b.elements = append(b.elements, value)
}

// Removes an item from the front of the queue and returns it
// If the queue is empty, returns 0 and an error
func (b *bag) dequeue() (*Tetrimino, error) {
	if len(b.elements) == 0 {
		return nil, fmt.Errorf("queue is empty")
	}
	element := b.elements[0]    // Get the first element
	b.elements = b.elements[1:] // Remove it from the queue
	return element, nil
}
