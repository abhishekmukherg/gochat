package goscs

import (
	"time"
)

type timer interface {
	Now() time.Time
}

type defaultTimer struct{}

func (defaultTimer) Now() time.Time {
	return time.Now()
}
