package ports

import "time"

// Clock allows tests to inject a fixed time.
type Clock interface {
	Now() time.Time
}

// SystemClock returns time.Now().
type SystemClock struct{}

func (SystemClock) Now() time.Time { return time.Now().UTC() }
