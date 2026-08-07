package shared

import "time"

// CoalesceTime returns t if it is not zero, otherwise fallback
func CoalesceTime(t, fallback time.Time) time.Time {
	if t.IsZero() {
		return fallback
	}
	return t
}
