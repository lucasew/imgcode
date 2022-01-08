package utils

import (
	"math/rand"
	"testing"
)

func TestSizeWorks(t *testing.T) {
	// Preparation
	const numTests = 10
	for k := 0; k < numTests; k++ {
		in := rand.Int()
		side := Size(in)
		out := side * side
		if in > out {
			t.Errorf("%d should be >= %d (x: %d, y: %d)", out, in, side, side)
		}
	}
}
