package fakers

import (
	"math/rand"
)

func RandBool() bool {
	return rand.Intn(2) == 0
}

func RandInt(min, max int) int {
	if max <= min {
		panic("RandInt: max must be greater than min")
	}
	return rand.Intn(max-min+1) + min
}

func RandItem[T any](items []T) T {
	if len(items) == 0 {
		panic("RandItem: empty items")
	}
	return items[RandInt(0, len(items)-1)]
}

func RandProtocol() string {
	protocols := []string{"tcp", "udp"}
	return RandItem(protocols)
}
