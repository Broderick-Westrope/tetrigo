package tetris

import (
	"math"
	"time"
)

type Fall struct {
	DefaultTime  time.Duration
	SoftDropTime time.Duration
	IsSoftDrop   bool
}

func NewFall(level uint) *Fall {
	f := Fall{}
	f.CalculateFallSpeeds(level)
	return &f
}

func (f *Fall) CalculateFallSpeeds(level uint) {
	speed := math.Pow(0.8-float64(level-1)*0.007, float64(level-1)) * 1000000

	f.DefaultTime = time.Microsecond * time.Duration(speed)
	f.SoftDropTime = time.Microsecond * time.Duration(speed/10)
}

func (f *Fall) ToggleSoftDrop() {
	f.IsSoftDrop = !f.IsSoftDrop
}
