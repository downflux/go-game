package entity

import (
	"sync"

	"github.com/downflux/game/curve/curve"

	gcpb "github.com/downflux/game/api/constants_go_proto"
)

type NoCurveEntity struct{}

func (e *NoCurveEntity) Curve(c gcpb.CurveCategory) curve.Curve { return nil }
func (e *NoCurveEntity) CurveCategories() []gcpb.CurveCategory  { return nil }

type BaseEntity struct {
	lifetimeMux sync.RWMutex
	start       float64
	end         float64
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
