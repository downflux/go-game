package visitorlist

import (
	"github.com/downflux/game/server/service/visitor/visitor"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	vcpb "github.com/downflux/game/server/service/visitor/api/constants_go_proto"
)

type List struct {
	visitors map[vcpb.VisitorType]visitor.Visitor
	order    []vcpb.VisitorType
}

func New(visitors []visitor.Visitor) (*List, error) {
	l := &List{}

	for _, v := range visitors {
		if err := l.appendVisitor(v); err != nil {
			return nil, err
		}
	}

	return l, nil
}

func (l *List) Get(visitorType vcpb.VisitorType) visitor.Visitor {
	return l.visitors[visitorType]
}

func (l *List) Iter() []visitor.Visitor {
	var visitors []visitor.Visitor
	for _, visitorType := range l.order {
		visitors = append(visitors, l.Get(visitorType))
	}

	return visitors
}

func (l *List) appendVisitor(v visitor.Visitor) error {
	if l.visitors == nil {
		l.visitors = map[vcpb.VisitorType]visitor.Visitor{}
	}

	if _, found := l.visitors[v.Type()]; found {
		return status.Errorf(codes.InvalidArgument, "a Visitor with type %v has already been specified", v.Type())
	}

	l.visitors[v.Type()] = v
	l.order = append(l.order, v.Type())

	return nil
}
