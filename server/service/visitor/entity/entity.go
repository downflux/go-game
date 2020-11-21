package entity

import (
	"sync"

	"github.com/downflux/game/curve/curve"

	gcpb "github.com/downflux/game/api/constants_go_proto"
)

type Entity interface {
	ID() string
	Curve(t gcpb.CurveCategory) curve.Curve

	// CurveCategories returns list of curve categories defined in a specific
	// entity. This list is created at init time and is immutable.
	CurveCategories() []gcpb.CurveCategory

	Start() float64
	End() float64

	Delete(tick float64)
}

type BaseEntity struct {
	lifetimeMux sync.RWMutex
	start float64
	end float64
}

func (e *BaseEntity) Start() float64 {
	e.lifetimeMux.RLock()
	defer e.lifetimeMux.RUnlock()

	return e.start
}

func (e *BaseEntity) End() float64 {
	e.lifetimeMux.RLock()
	defer e.lifetimeMux.RUnlock()

	return e.end
}

func (e *BaseEntity) Delete(tick float64) {
	e.lifetimeMux.Lock()
	defer e.lifetimeMux.Unlock()

	e.end = tick
}
