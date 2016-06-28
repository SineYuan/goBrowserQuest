package bqs

type NPC struct {
	Entity
}

func NewNPC(kind EntityType, x, y int) *NPC {
	npc := new(NPC)
	npc.Id = 80000 + x * 100 + y
	npc.X, npc.Y = x, y
	npc.Kind = kind
	return npc
}

func (npc *NPC) GetState() []interface{} {
	return npc.Entity.GetState()
}