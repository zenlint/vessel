package timer

import (
	"testing"
	"time"
)

func TestLeftTime(t *testing.T) {
	initialTimeout := 10
	tim := InitTimer(float64(initialTimeout))

	<-time.After(time.Second * time.Duration(1))
	if (tim.GetLeftTime() > float64(initialTimeout-1)) || (tim.GetLeftTime() < float64(initialTimeout-2)) {
		t.Errorf("Testing timer.GetLeftTime for 1 sec err, Timer: %v", tim)
	}

	<-time.After(time.Second * time.Duration(2))
	if (tim.GetLeftTime() > float64(initialTimeout-3)) || (tim.GetLeftTime() < float64(initialTimeout-4)) {
		t.Errorf("Testing timer.GetLeftTime for 3 sec err, Timer: %v", tim)
	}
}

func TestOverflow(t *testing.T) {
	initialTimeout := 10
	tim := InitTimer(float64(initialTimeout))

	<-time.After(time.Second * time.Duration(11))
	leftTime := tim.GetLeftTime()
	if leftTime > float64(-1) || leftTime < float64(-2) {
		t.Errorf("Testing timer.GetLeftTime overflow err, Timer: %v", tim)
	}

}
