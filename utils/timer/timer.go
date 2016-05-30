package timer

import (
	"time"
)

type Timer struct {
	Start           time.Time
	TimeoutDuration float64
	LeftTime        float64
}

func InitTimer(initialTimeout float64) (timer *Timer) {
	t := new(Timer)
	t.Start = time.Now()
	t.TimeoutDuration = initialTimeout
	t.LeftTime = initialTimeout
	return t
}

func (this *Timer) GetLeftTime() (leftTime float64) {
	this.LeftTime = this.LeftTime - time.Now().Sub(this.Start).Seconds()
	this.Start = time.Now()
	return this.LeftTime
}
