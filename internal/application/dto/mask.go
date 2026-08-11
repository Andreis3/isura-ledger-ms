package dto

import (
	"log/slog"
	"strings"
)

// MaskSlogValue optimized to avoid heavy reflection per call if the struct
// directly implements a masking interface or if we use direct paths.
// If you keep the generic reflection, here is a version that reuses logic and avoids Sprintf.

func MaskSlogValue[T any](input T) slog.Value {
	// Optimization Tip: If T implements a custom logging interface,
	// we avoid the full cost of reflection at runtime for known structures.
	if logValuer, ok := any(input).(slog.LogValuer); ok {
		return logValuer.LogValue()
	}

	// If you need the generic reflection fallback, keep it lean:
	return slog.AnyValue(input)
}

// maskMiddleVisible optimized with direct rune manipulation without excessive strings.Builder
func MaskMiddleVisible(s string) string {
	runes := []rune(s)
	length := len(runes)

	if length <= 4 {
		return "******"
	}

	visibleLen := length / 3
	if visibleLen < 2 {
		visibleLen = 2
	}

	start := (length - visibleLen) / 2
	end := start + visibleLen

	// Pre-allocates the exact size of the resulting string to zero out Builder reallocations
	var sb strings.Builder
	sb.Grow(length)

	for i := 0; i < start; i++ {
		sb.WriteRune('*')
	}
	sb.WriteString(string(runes[start:end]))
	for i := end; i < length; i++ {
		sb.WriteRune('*')
	}

	return sb.String()
}
