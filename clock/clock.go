package clock

import "time"

type Clocker interface {
	Now() time.Time
}

type RealClock struct{}

func (r RealClock) Now() time.Time {
	return time.Now()
}
