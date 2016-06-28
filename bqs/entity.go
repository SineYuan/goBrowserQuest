package bqs

type Entity struct {
	Id         int
	X, Y int
	Kind EntityType
}

const (
	ORIENTATION_UP = 1
	ORIENTATION_DOWN = 2
	ORIENTATION_LEFT = 3
	ORIENTATION_RIGTH = 4
)

type EntityInterface interface {
	GetState() []interface{}
}

func (e *Entity) GetState() []interface{} {
	return []interface{} {e.Id, e.Kind.Id, e.X, e.Y}
}