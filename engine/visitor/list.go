// Package list encapsulates logic for a managed list of Visitor
// instances.
package list

import (
	"github.com/downflux/game/engine/visitor/visitor"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	fcpb "github.com/downflux/game/engine/fsm/api/constants_go_proto"
)

// List implements a managed list of Visitor instances.
type List struct {
	visitors map[fcpb.FSMType]visitor.Visitor
	order    []fcpb.FSMType
}

// New creates a new instance of the List object.
func New(visitors []visitor.Visitor) (*List, error) {
	l := &List{}

	for _, v := range visitors {
		if err := l.appendVisitor(v); err != nil {
			return nil, err
		}
	}

	return l, nil
}

// Get returns a concrete Visitor instance given the registered VisitorType.
func (l *List) Visitor(visitorType fcpb.FSMType) visitor.Visitor {
	return l.visitors[visitorType]
}

// Iter returns a list of Visitor instances. This is used in for loops when
// ranging over all registered instances.
func (l *List) Iter() []visitor.Visitor {
	var visitors []visitor.Visitor
	for _, visitorType := range l.order {
		visitors = append(visitors, l.Visitor(visitorType))
	}

	return visitors
}

// appendVisitor registers a new Visitor instance into the managed List.
func (l *List) appendVisitor(v visitor.Visitor) error {
	if l.visitors == nil {
		l.visitors = map[fcpb.FSMType]visitor.Visitor{}
	}

	if _, found := l.visitors[v.Type()]; found {
		return status.Errorf(codes.InvalidArgument, "a Visitor with type %v has already been specified", v.Type())
	}

	l.visitors[v.Type()] = v
	l.order = append(l.order, v.Type())

	return nil
}
