package clock

import "time"

// Clock keeps time-dependent workflows testable without leaking time.Now calls.
type Clock interface {
	Now() time.Time
}

type System struct{}

func (System) Now() time.Time {
	return time.Now().UTC()
}
