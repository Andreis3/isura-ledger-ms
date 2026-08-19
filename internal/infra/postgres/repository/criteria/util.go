package criteria

import "strconv"

// Ultra-fast and zero-allocation helper function to convert argument indexes up to 9 to string
func argNumToString(n int) string {
	return strconv.Itoa(n)
}
