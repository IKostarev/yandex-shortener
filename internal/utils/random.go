package utils

import (
	"math/rand"
	"time"
)

const letters = 8

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func RandomString() string {
	alphabet := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

	b := make([]rune, letters)
	for i := range b {
		b[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(b)
}
