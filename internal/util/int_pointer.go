package util

func Int64Pointer(i int64) *int64 {
	return &i
}

func Int64(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}
