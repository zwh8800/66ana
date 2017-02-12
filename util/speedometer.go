package util

import (
	"sync/atomic"
	"time"
)

type Speedometer struct {
	counter   uint64
	startTime time.Time
}

func NewSpeedometer() *Speedometer {
	return &Speedometer{startTime: time.Now()}
}

func (s *Speedometer) Add() {
	atomic.AddUint64(&s.counter, 1)
}

func (s *Speedometer) GetSpeed() float64 {
	diffTime := float64(time.Now().Sub(s.startTime)) / float64(time.Second)
	if diffTime == 0 {
		return 0
	}
	speed := float64(s.counter) / diffTime
	if diffTime > 1 {
		s.startTime = time.Now()
		atomic.StoreUint64(&s.counter, 0)
	}
	return speed
}
