package tetris

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBag(t *testing.T) {
	tt := map[string]struct {
		matrixHeight int
	}{
		"matrix height 20": {
			matrixHeight: 20,
		},
		"matrix height 40": {
			matrixHeight: 40,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			b := NewBag(tc.matrixHeight)

			if len(b.elements) != 14 {
				t.Errorf("Length: want 14, got %d", len(b.elements))
			}

			for _, e := range b.elements {
				for _, tet := range Tetriminos {
					if tet.Value != e.Value {
						continue
					}

					assert.Equal(t, tet.Pos.X, e.Pos.X)
				}
			}
		})
	}
}

// Checks:
//   - that the tetrimino returned is the first element of the bag.
//   - that the first element of the bag is removed.
//   - that the bag is filled with 7 tetriminos if it has less than 7 elements.
func TestBag_Next(t *testing.T) {
	tt := map[string]struct {
		elements []Tetrimino
	}{
		"2 elements": {
			elements: make([]Tetrimino, 2, 14),
		},
		"7 elements": {
			elements: make([]Tetrimino, 7, 14),
		},
		"8 elements": {
			elements: make([]Tetrimino, 8, 14),
		},
		"14 elements": {
			elements: make([]Tetrimino, 14),
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			b := Bag{
				elements:  tc.elements,
				startLine: 40,
			}
			expected := tc.elements[0].DeepCopy()
			expected.Pos.Y += b.startLine

			var expectedElements []Tetrimino
			for _, e := range tc.elements[1:] {
				temp := e.DeepCopy()
				temp.Pos.Y += b.startLine
				expectedElements = append(expectedElements, *temp)
			}

			result := b.Next()
			assert.EqualValues(t, *expected, *result)

			for i := range b.elements {
				b.elements[i].Pos.Y += b.startLine
			}

			v := b.elements[:len(expectedElements)]
			assert.EqualValues(t, expectedElements, v)

			expectedLength := len(expectedElements)
			if expectedLength < 7 && len(b.elements) != expectedLength+7 {
				t.Errorf("Length: want %d, got %d", expectedLength+7, len(b.elements))
			}
		})
	}
}

// Checks:
//   - that the bag is filled with 7 tetriminos if it has less than 7 elements.
//   - that the tetriminos added to the bag are not duplicates.
func TestBag_Fill(t *testing.T) {
	tt := map[string]struct {
		elements    []Tetrimino
		timesToFill int
	}{
		"0 elements, 1 fill": {
			elements:    make([]Tetrimino, 0, 14),
			timesToFill: 1,
		},
		"4 elements, 1 fill": {
			elements:    make([]Tetrimino, 4, 14),
			timesToFill: 1,
		},
		"7 elements, 1 fill": {
			elements:    make([]Tetrimino, 7, 14),
			timesToFill: 1,
		},
		"0 elements, 2 fills": {
			elements:    make([]Tetrimino, 0, 14),
			timesToFill: 2,
		},
		"4 elements, 2 fills": {
			elements:    make([]Tetrimino, 4, 14),
			timesToFill: 2,
		},
		"7 elements, 2 fills": {
			elements:    make([]Tetrimino, 7, 14),
			timesToFill: 2,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			b := Bag{
				elements:  tc.elements,
				startLine: 40,
			}

			for i := 0; i < tc.timesToFill; i++ {
				b.fill()
			}

			expectedLength := len(tc.elements) + (7 * tc.timesToFill)
			for expectedLength > 14 {
				expectedLength -= 7
			}
			assert.Len(t, b.elements, expectedLength)

			tetCount := make(map[byte]int)
			for i := len(tc.elements); i < len(b.elements); i++ {
				tetCount[b.elements[i].Value]++
			}
			for value, count := range tetCount {
				assert.LessOrEqualf(t, count, tc.timesToFill, "Duplicate tetrimino '%v' in bag: %v", value, b.elements)
			}
		})
	}
}
