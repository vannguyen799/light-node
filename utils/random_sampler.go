package utils

import (
	"math/rand"
	"time"
)

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

func RandomElement[T any](arr []T) T {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomIndex := r.Intn(len(arr))
	return arr[randomIndex]
}
