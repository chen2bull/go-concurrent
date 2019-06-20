package atomic

func Max(x, y int64) int64 {
	if x < y {
		return y
	}
	return x
}

// Min returns the smaller of x or y.
func Min(x, y int64) int64 {
	if x > y {
		return y
	}
	return x
}
