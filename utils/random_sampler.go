package utils

import (
	"math/rand"
	"time"
)

func RandomSample[T any](arr []T, n int) []T {
	if len(arr) == 0 {
		return []T{}
	}

	if n > len(arr) {
		n = len(arr)
	}
	
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

func RandomElement[T any](arr []T) (result T) {
	if len(arr) == 0 {
		return result
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomIndex := r.Intn(len(arr))
	return arr[randomIndex]
}
