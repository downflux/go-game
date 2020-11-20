package visitor

import (
	"github.com/downflux/game/curve/curve"
)

type Visitor interface {
	// Schedule adds a Visitor-specific command to the Visitor. This
	// function will be called concurrently.
	Schedule(tick float64, args interface{}) error

	// Execute will run appropriate commands for the current tick. If
	// a timeout occurs, the function will return early. This function
	// will be called serially.
	Execute(tick float64) ([]curve.Curve, error)
}
