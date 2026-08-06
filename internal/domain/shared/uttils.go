package shared

import "time"

// CoalesceTime retorna t se não for zero, caso contrário fallback
func CoalesceTime(t, fallback time.Time) time.Time {
	if t.IsZero() {
		return fallback
	}
	return t
}
