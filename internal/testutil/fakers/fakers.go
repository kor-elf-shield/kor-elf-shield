package fakers

import (
	"math/rand"
)

func RandBool() bool {
	return rand.Intn(2) == 0
}

func RandInt(min, max int) int {
	return rand.Intn(max-min) + min
}

func RandItem[T any](items []T) T {
	return items[RandInt(0, len(items)-1)]
}
