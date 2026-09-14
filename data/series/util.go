package series

func inRange(i, start, end int) bool {
	return i >= start && i <= end
}

// If in range, fit i to [0, len(range)-1].
func fitRange(i *int, start, end int) bool {
	if !inRange(*i, start, end) {
		return false
	}

	*i -= start
	return true
}

// Return true if i < x, else shift i by x.
func shiftRange(i *int, x int) bool {
	if *i < x {
		return true
	}

	*i -= x
	return false
}
