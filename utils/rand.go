package utils

import (
	"math/rand"
	"time"
)

func RangeRand(max int) int {
	if max <= 1 {
		return max
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max)
}
