package bqs


type MobArea struct {
	Id int
	X, Y int
	Width, Height int
	Type string
	NB int
	Entities []*Entity
}
