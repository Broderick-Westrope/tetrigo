package tetris

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewNextQueue(t *testing.T) {
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
			b := NewNextQueue(tc.matrixHeight)

			assert.GreaterOrEqual(t, len(b.elements), 7, "Length: want at least 7, got %d", len(b.elements))

			for _, e := range b.elements {
				tet, err := GetTetrimino(e.Value)
				require.NoError(t, err)

				assert.Equal(t, tet.Position.X, e.Position.X)
			}
		})
	}
}

// Checks:
//   - that the tetrimino returned is the first element of the queue.
//   - that the first element of the queue is removed.
//   - that the queue is filled with 7 tetriminos if it has less than 7 elements.
func TestNextQueue_Next(t *testing.T) {
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
			nq := NextQueue{
				elements: tc.elements,
				skyline:  40,
			}
			expected := tc.elements[0].DeepCopy()
			expected.Position.Y += nq.skyline

			var expectedElements []Tetrimino
			for _, e := range tc.elements[1:] {
				temp := e.DeepCopy()
				temp.Position.Y += nq.skyline
				expectedElements = append(expectedElements, *temp)
			}

			result := nq.Next()
			assert.EqualValues(t, *expected, *result)

			for i := range nq.elements {
				nq.elements[i].Position.Y += nq.skyline
			}

			v := nq.elements[:len(expectedElements)]
			assert.EqualValues(t, expectedElements, v)

			expectedLength := len(expectedElements)
			if expectedLength < 7 && len(nq.elements) != expectedLength+7 {
				t.Errorf("Length: want %d, got %d", expectedLength+7, len(nq.elements))
			}
		})
	}
}

// Checks:
//   - that the queue is filled with 7 tetriminos if it has less than 7 elements.
//   - that the tetriminos added to the queue are not duplicates.
func TestNextQueue_Fill(t *testing.T) {
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
			nq := NextQueue{
				elements: tc.elements,
				skyline:  40,
			}

			for range tc.timesToFill {
				nq.fill()
			}

			expectedLength := len(tc.elements) + (7 * tc.timesToFill)
			for expectedLength > 14 {
				expectedLength -= 7
			}
			assert.Len(t, nq.elements, expectedLength)

			tetCount := make(map[byte]int)
			for i := len(tc.elements); i < len(nq.elements); i++ {
				tetCount[nq.elements[i].Value]++
			}
			for value, count := range tetCount {
				assert.LessOrEqualf(t, count, tc.timesToFill, "Duplicate tetrimino '%v' in next queue: %v", value, nq.elements)
			}
		})
	}
}
