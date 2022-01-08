package utils

import (
	"math"
)

// Size Gets the best size for x and y given length
func Size(length int) (side int) {
	squared := math.Sqrt(float64(length))
	if squared == float64(int(squared)) {
		return int(squared)
	}
	return int(squared) + 1
}
