package entity

import (
	"sync"

	"github.com/downflux/game/curve/curve"
	"github.com/downflux/game/server/id"

	gcpb "github.com/downflux/game/api/constants_go_proto"
)

type NoCurveEntity struct{}

func (e *NoCurveEntity) Curve(c gcpb.CurveCategory) curve.Curve { return nil }
func (e *NoCurveEntity) CurveCategories() []gcpb.CurveCategory  { return nil }

type BaseEntity struct {
	lifetimeMux sync.RWMutex
	start       id.Tick
	end         id.Tick
}

func (e *BaseEntity) Start() id.Tick {
	e.lifetimeMux.RLock()
	defer e.lifetimeMux.RUnlock()

	return e.start
}

func (e *BaseEntity) End() id.Tick {
	e.lifetimeMux.RLock()
	defer e.lifetimeMux.RUnlock()

	return e.end
}

func (e *BaseEntity) Delete(tick id.Tick) {
	e.lifetimeMux.Lock()
	defer e.lifetimeMux.Unlock()

	e.end = tick
}
