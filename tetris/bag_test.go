package tetris

import "testing"

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
