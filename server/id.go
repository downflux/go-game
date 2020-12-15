// Package id generates random ID strings for server-side objects. This should
// not be called by client-side logic.
//
// TODO(minkezhang): Migrate EntityID, ClientID into separate directory(s).
package id

import (
	"math/rand"
)

// Tick represents the internal game time. This is incremented once per game
// loop. Non-integer values of the Tick are used to mark special moments in
// the game for Curves (e.g. inflection points, zero-values, etc.).
type Tick float64

// Value returns the basic value of the Tick.
func (t Tick) Value() float64 { return float64(t) }

type ID string

func (id ID) Value() string { return string(id) }

type InstanceID ID
type EntityID ID
type ClientID ID

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
