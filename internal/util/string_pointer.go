package util

func StringPointer(s string) *string {
	return &s
}

func String(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
