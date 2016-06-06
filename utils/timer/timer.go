package timer

import (
	"time"
)

// Hourglass calculate left time
type Hourglass struct {
	endTime time.Time
}

// InitHourglass create hourglass
func InitHourglass(nanoseconds time.Duration) *Hourglass {
	t := new(Hourglass)
	t.endTime = time.Now().Add(nanoseconds)
	return t
}

// GetLeftNanoseconds get left time for nanoseconds
func (this *Hourglass) GetLeftNanoseconds() time.Duration {
	return time.Duration(this.endTime.Sub(time.Now()).Nanoseconds())
}
