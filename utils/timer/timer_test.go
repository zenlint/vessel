package timer

import (
	"testing"
	"time"
)

// TestLeftTime test left time
func TestLeftTime(t *testing.T) {
	initialTimeout := 10
	timeNanoSeconds := time.Second * time.Duration(initialTimeout)
	tim := InitHourglass(timeNanoSeconds)

	<-time.After(time.Second * time.Duration(1))
	t.Errorf("Testing timer.GetLeftTime for 1 sec err, TimeLeft: %d", tim.GetLeftNanoseconds())

	<-time.After(time.Second * time.Duration(2))
	t.Errorf("Testing timer.GetLeftTime for 3 sec err, TimeLeft: %d", tim.GetLeftNanoseconds())
}

// TestOverflow test over flow
func TestOverflow(t *testing.T) {
	initialTimeout := 10
	tim := InitHourglass(time.Second * time.Duration(initialTimeout))

	<-time.After(time.Second * time.Duration(11))
	t.Errorf("Testing timer.GetLeftTime overflow err, TimeLeft: %d", tim.GetLeftNanoseconds())

}
