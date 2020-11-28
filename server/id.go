// Package id generates random ID strings for server-side objects. This should
// not be called by client-side logic.
package id

import (
	"math/rand"
)

type ID string

// TODO(minkezhang): Migrate into separate directory.
type EntityID ID

func (id EntityID) Value() string { return string(id) }

type ClientID ID

func (id ClientID) Value() string { return string(id) }

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

// RandomString returns a random string of the specified length.
// TODO(minkezhang): Rename to Generate.
func RandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
