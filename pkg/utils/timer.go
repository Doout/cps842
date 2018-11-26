package utils

import "time"

type Timer struct {
	start time.Time
	end   time.Time
}

func (timer *Timer) Start() {
	timer.start = time.Now()
}

func (timer *Timer) Stop() {
	timer.end = time.Now()
}

func (timer *Timer) Duration() time.Duration {
	return timer.end.Sub(timer.start)
}
