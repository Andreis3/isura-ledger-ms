package dto

import (
	"strings"
)

// MaskMiddleVisible optimized with direct rune manipulation without excessive strings.Builder
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

// MaskTotal optimized with direct rune manipulation without excessive strings.Builder
func MaskTotal(s string) string {
	return strings.Repeat("*", len([]rune(s)))
}
