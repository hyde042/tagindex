package tagindex

// input slice values must be sorted from lowest to highest
func setContains(a, b []uint32) bool {
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
	var i, j int
	if b[bl-1] > a[al-1] {
		return false
	}
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
