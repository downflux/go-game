package curve

import (
	"github.com/downflux/game/engine/curve/list"
)

type Component struct {
	curves *list.List
}

func New(cs *list.List) *Component {
	return &Component{
		curves: cs,
	}
}

func (c Component) Curves() *list.List { return c.curves }
