package utils

import (
	"math/rand"
	"time"
)

// RandomSample returns n random elements from the input slice
// The type parameter T can be any type
func RandomSample[T any](arr []T, n int) []T {
	if n > len(arr) {
		n = len(arr)
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	shuffled := make([]T, len(arr))
	copy(shuffled, arr)
	r.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	return shuffled[:n]
}
