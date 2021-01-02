package targetable

type Component interface{}

type Base struct{}

func New() *Base {
	return &Base{}
}
