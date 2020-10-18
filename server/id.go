package id

import (
	"math/rand"
)

// TODO(minkezhang): Export to shared for server.
func RandomString(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
