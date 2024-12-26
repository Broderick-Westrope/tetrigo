package tetris

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// calculateExpectedSpeed helper function to match the formula exactly.
func calculateExpectedSpeed(level int) time.Duration {
	decrementedLevel := float64(level - 1)
	speed := math.Pow(0.8-(decrementedLevel*0.007), decrementedLevel)
	return time.Duration(speed * float64(time.Second))
}

func TestNewFall(t *testing.T) {
	// Test initialization with different levels.
	tt := map[string]struct {
		level int
	}{
		"level 1":  {level: 1},
		"level 15": {level: 15},
		"level 30": {level: 30},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			f := NewFall(tc.level)

			expectedDefault := calculateExpectedSpeed(tc.level)
			expectedSoftDrop := time.Duration(float64(expectedDefault) / 15)

			assert.Equal(t, expectedDefault, f.DefaultInterval,
				"default fall speed should be %v, got %v", expectedDefault, f.DefaultInterval)

			assert.Equal(t, expectedSoftDrop, f.SoftDropInterval,
				"soft drop speed should be %v, got %v", expectedSoftDrop, f.SoftDropInterval)

			assert.False(t, f.IsSoftDrop, "IsSoftDrop should be false by default")
		})
	}
}

func TestFall_CalculateFallSpeeds(t *testing.T) {
	f := &Fall{}

	// Test that recalculating speeds works.
	t.Run("recalculate speeds", func(t *testing.T) {
		// Start with level 1
		f.CalculateFallSpeeds(1)
		initialDefault := f.DefaultInterval
		initialSoftDrop := f.SoftDropInterval

		// Recalculate for level 10.
		f.CalculateFallSpeeds(10)

		assert.Greater(t, initialDefault, f.DefaultInterval,
			"level 10 should be faster (smaller interval) than level 1")
		assert.Greater(t, initialSoftDrop, f.SoftDropInterval,
			"soft drop at level 10 should be faster than level 1")
	})

	// Test that soft drop is always faster.
	t.Run("soft drop faster than default", func(t *testing.T) {
		levels := []int{1, 5, 10, 15, 20}
		for _, level := range levels {
			f.CalculateFallSpeeds(level)
			assert.Greater(t, f.DefaultInterval, f.SoftDropInterval,
				"soft drop should be faster (smaller interval) than default at level %d", level)
		}
	})
}

func TestFall_ToggleSoftDrop(t *testing.T) {
	f := NewFall(1)

	assert.False(t, f.IsSoftDrop, "should start with soft drop disabled")

	f.ToggleSoftDrop()
	assert.True(t, f.IsSoftDrop, "should enable soft drop after first toggle")

	f.ToggleSoftDrop()
	assert.False(t, f.IsSoftDrop, "should disable soft drop after second toggle")
}

func TestFall_SpeedProgression(t *testing.T) {
	f := &Fall{}
	var prevInterval time.Duration

	// Test first 20 levels
	for level := 1; level <= 20; level++ {
		f.CalculateFallSpeeds(level)

		if level > 1 {
			assert.Greater(t, prevInterval, f.DefaultInterval,
				"fall speed should increase (interval should decrease) from level %d to %d", level-1, level)
		}

		prevInterval = f.DefaultInterval
	}
}
