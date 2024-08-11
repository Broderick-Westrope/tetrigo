package tetris

import (
	"math"
	"time"
)

type Fall struct {
	DefaultInterval  time.Duration
	SoftDropInterval time.Duration
	IsSoftDrop       bool
}

func NewFall(level uint) *Fall {
	f := Fall{}
	f.CalculateFallSpeeds(level)
	return &f
}

func (f *Fall) CalculateFallSpeeds(level uint) {
	decrementedLevel := float64(level - 1)
	speed := math.Pow(0.8-(decrementedLevel*0.007), decrementedLevel)
	speed *= float64(time.Second)

	f.DefaultInterval = time.Duration(speed)
	f.SoftDropInterval = time.Duration(speed / 15)
}

func (f *Fall) ToggleSoftDrop() {
	f.IsSoftDrop = !f.IsSoftDrop
}
