package implement

import (
	icpb "github.com/downflux/game/entity/implement/api/constants_go_proto"
)

type List struct {
	implements map[icpb.Implement]bool
}

func New(implements []icpb.Implement) *List {
	l := &List{
		implements: map[icpb.Implement]bool{},
	}

	for _, i := range implements {
		l.implements[i] = true
	}

	return l
}

func (l *List) Implements(i icpb.Implement) bool {
	return l.implements[i]
}
