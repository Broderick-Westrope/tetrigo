package tetris

import (
	"reflect"
	"testing"
)

// Checks:
//   - that the tetrimino returned is the first element of the bag.
//   - that the first element of the bag is removed.
//   - that the bag is filled with 7 tetriminos if it has less than 7 elements.
func TestBag_Next(t *testing.T) {
	tt := []struct {
		name     string
		elements []Tetrimino
	}{
		{
			name:     "2 elements",
			elements: make([]Tetrimino, 2, 14),
		},
		{
			name:     "7 elements",
			elements: make([]Tetrimino, 7, 14),
		},
		{
			name:     "8 elements",
			elements: make([]Tetrimino, 8, 14),
		},
		{
			name:     "14 elements",
			elements: make([]Tetrimino, 14),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			b := Bag{
				Elements:     tc.elements,
				matrixHeight: 40,
			}
			expected := tc.elements[0].Copy()
			expected.Pos.Y += b.matrixHeight - 20

			var expectedElements []Tetrimino
			for _, e := range tc.elements[1:] {
				temp := e.Copy()
				temp.Pos.Y += b.matrixHeight - 20
				expectedElements = append(expectedElements, *temp)
			}

			result := b.Next()
			if !reflect.DeepEqual(*result, *expected) {
				t.Errorf("Tetrimino: want %v, got %v", *expected, *result)
			}

			for i := range b.Elements {
				b.Elements[i].Pos.Y += b.matrixHeight - 20
			}

			if v := b.Elements[:len(expectedElements)]; !reflect.DeepEqual(v, expectedElements) {
				t.Errorf("Elements: want %v, got %v", expectedElements, v)
			}

			expectedLength := len(expectedElements)
			if expectedLength < 7 && len(b.Elements) != expectedLength+7 {
				t.Errorf("Length: want %d, got %d", expectedLength+7, len(b.Elements))
			}
		})
	}
}

// Checks:
//   - that the bag is filled with 7 tetriminos if it has less than 7 elements.
//   - that the tetriminos added to the bag are not duplicates.
func TestBag_Fill(t *testing.T) {
	tt := []struct {
		name        string
		elements    []Tetrimino
		timesToFill int
	}{
		{
			name:        "0 elements, 1 fill",
			elements:    make([]Tetrimino, 0, 14),
			timesToFill: 1,
		},
		{
			name:        "4 elements, 1 fill",
			elements:    make([]Tetrimino, 4, 14),
			timesToFill: 1,
		},
		{
			name:        "7 elements, 1 fill",
			elements:    make([]Tetrimino, 7, 14),
			timesToFill: 1,
		},
		{
			name:        "0 elements, 2 fills",
			elements:    make([]Tetrimino, 0, 14),
			timesToFill: 2,
		},
		{
			name:        "4 elements, 2 fills",
			elements:    make([]Tetrimino, 4, 14),
			timesToFill: 2,
		},
		{
			name:        "7 elements, 2 fills",
			elements:    make([]Tetrimino, 7, 14),
			timesToFill: 2,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			b := Bag{
				Elements:     tc.elements,
				matrixHeight: 40,
			}

			for i := 0; i < tc.timesToFill; i++ {
				b.fill()
			}

			length := len(b.Elements)
			expectedLength := len(tc.elements) + (7 * tc.timesToFill)
			for expectedLength > 14 {
				expectedLength -= 7
			}
			if length != expectedLength {
				t.Errorf("Length: want %d, got %d", expectedLength, length)
			}

			tetCount := make(map[byte]int)
			for i := len(tc.elements); i < length; i++ {
				tetCount[b.Elements[i].Value]++
			}
			for value, count := range tetCount {
				if count > tc.timesToFill {
					t.Errorf("Duplicate tetrimino '%v' in bag: %v", value, b.Elements)
				}
			}
		})
	}
}
