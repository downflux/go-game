package targetable

import (
	"github.com/downflux/game/server/entity/component/positionable"
)

type Component interface {
	positionable.Component
}

type Base struct{}

func New() *Base {
	return &Base{}
}
