package single

import (
	"math/rand/v2"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToggleSoftDrop(t *testing.T) {
	tests := map[string]struct {
		initialSoftDrop bool
		tetriminoStartY int
		tetriminoEndY   int
		skylineY        int
		wantInterval    time.Duration
		wantScoreAdded  bool
	}{
		"Enable soft drop": {
			initialSoftDrop: false,
			tetriminoStartY: 5,
			tetriminoEndY:   5,
			skylineY:        0,
			wantInterval:    ((100 / 3) * 2) * time.Millisecond, // Based on level 1 soft drop speed
			wantScoreAdded:  false,
		},
		"Disable soft drop with lines cleared": {
			initialSoftDrop: true,
			tetriminoStartY: 5,
			tetriminoEndY:   10,
			skylineY:        0,
			wantInterval:    time.Second, // Based on level 1 default speed
			wantScoreAdded:  true,
		},
		"Disable soft drop without lines cleared": {
			initialSoftDrop: true,
			tetriminoStartY: 5,
			tetriminoEndY:   5,
			skylineY:        0,
			wantInterval:    time.Second,
			wantScoreAdded:  false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup game state
			game, err := NewGame(&Input{
				Level:         1,
				MaxLevel:      0,
				IncreaseLevel: false,
				EndOnMaxLevel: false,
				MaxLines:      0,
				EndOnMaxLines: false,
				GhostEnabled:  false,
				Rand:          rand.New(rand.NewPCG(0, 0)),
			})
			require.NoError(t, err)

			game.fall.IsSoftDrop = tt.initialSoftDrop
			game.tetInPlay.Position.Y = tt.tetriminoEndY

			game.ToggleSoftDrop()
			interval := game.GetFallInterval()

			assert.InDelta(t, tt.wantInterval, interval, float64(time.Millisecond))
			assert.Equal(t, game.fall.IsSoftDrop, !tt.initialSoftDrop)
		})
	}
}
