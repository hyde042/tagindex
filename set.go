package tagindex

// TODO: move to internal package

// SetContains returns true if a has all elements of b.
// input slice values must be sorted from lowest to highest.
func SetContains(a, b []uint32) bool {
	var (
		al = len(a)
		bl = len(b)
	)
	if al == 0 {
		return bl == 0
	}
	if bl == 0 {
		return true
	}
	if b[bl-1] > a[al-1] {
		return false
	}
	var i, j int
	for i < bl {
		if b[i] == a[j] {
			i++
		} else if b[i] < a[j] {
			return false
		}
		j++
	}
	return true
}
