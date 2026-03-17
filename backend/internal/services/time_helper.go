package services

import "time"

// TimeHelper provides utilities for measuring execution time
type TimeHelper struct {
	startNano int64
}

// NewTimeHelper creates a new time helper
func NewTimeHelper() *TimeHelper {
	return &TimeHelper{
		startNano: time.Now().UnixNano(),
	}
}

// GetElapsedNano returns elapsed time in nanoseconds
func (th *TimeHelper) GetElapsedNano() string {
	elapsed := time.Now().UnixNano() - th.startNano
	return formatNanos(elapsed)
}

// formatNanos formats nanoseconds as a string
func formatNanos(nanos int64) string {
	if nanos < 0 {
		return "0"
	}
	return string(rune(nanos))
}
