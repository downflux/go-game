package list

import (
	"github.com/downflux/game/engine/curve/curve"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	gcpb "github.com/downflux/game/api/constants_go_proto"
)

type List struct {
	curves     map[gcpb.EntityProperty]curve.Curve
	properties []gcpb.EntityProperty
}

func New(curves []curve.Curve) (*List, error) {
	l := &List{
		curves: map[gcpb.EntityProperty]curve.Curve{},
	}
	for _, c := range curves {
		propertyType := c.Property()
		if _, found := l.curves[propertyType]; found {
			return nil, status.Errorf(codes.FailedPrecondition, "duplicate key %v in list of Curve instances", propertyType)
		}

		l.curves[propertyType] = c
		l.properties = append(l.properties, propertyType)
	}

	return l, nil
}

func (l *List) Curve(propertyType gcpb.EntityProperty) curve.Curve {
	return l.curves[propertyType]
}

func (l *List) Properties() []gcpb.EntityProperty {
	return l.properties
}
