package bqs

type Area struct {
	Id            int
	X, Y          int
	Width, Height int
	Entities      []*Entity
}

type ChestArea struct {
	Id     int
	X, Y   int
	TX, TY int
	W, H   int
	Items  []int `json:"i"`
}

