package utils

import (
	"math/rand"
	"time"
)

type Random struct{}

func NewString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := range result {
		result[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(result)
}
