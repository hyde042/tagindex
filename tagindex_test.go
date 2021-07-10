package tagindex_test

import "testing"

func assert(t *testing.T, tag string, result, expected bool) {
	if result != expected {
		t.Fatalf("%s: result mismatch", tag)
	}
}
