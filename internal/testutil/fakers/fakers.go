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

func RandString(length int) string {
	if length <= 0 {
		panic("RandString: length must be greater than 0")
	}

	b := make([]byte, length)
	const firstSymbolCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b[0] = firstSymbolCharset[RandInt(0, len(firstSymbolCharset)-1)]
	for i := 1; i < length; i++ {
		b[i] = charset[RandInt(0, len(charset)-1)]
	}
	return string(b)
}