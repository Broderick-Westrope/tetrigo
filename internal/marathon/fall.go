package marathon

import (
	"math"
	"time"

	"github.com/charmbracelet/bubbles/stopwatch"
)

type Fall struct {
	stopwatch    stopwatch.Model
	defaultTime  time.Duration
	softDropTime time.Duration
	isSoftDrop   bool
}

func (f *Fall) calculateFallSpeeds(level uint) {
	speed := math.Pow((0.8-float64(level-1)*0.007), float64(level-1)) * 1000000

	f.defaultTime = time.Microsecond * time.Duration(speed)
	f.softDropTime = time.Microsecond * time.Duration(speed/10)
}

func (f *Fall) toggleSoftDrop() {
	f.isSoftDrop = !f.isSoftDrop
	if f.isSoftDrop {
		f.stopwatch.Interval = f.softDropTime
		return
	}
	f.stopwatch.Interval = f.defaultTime
}

func defaultFall(level uint) *Fall {
	f := Fall{}
	f.calculateFallSpeeds(level)
	f.stopwatch = stopwatch.NewWithInterval(f.defaultTime)
	return &f
}
