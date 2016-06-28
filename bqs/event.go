package bqs

import (
	"encoding/json"
)

type Event struct {
	Type     int
	PlayerID int
	Data     []interface{}
}

const (
	EVENT_HELLO = 0
	EVENT_WELCOME = 1
	EVENT_SWAPN = 2
	EVENT_DESPAWN = 3
	EVENT_MOVE = 4
	EVENT_LOOTMOVE = 5
	EVENT_ATTACk = 7
	EVENT_HIT = 8
	EVENT_HURT = 9
	EVENT_HEALTH = 10
	EVENT_CHAT = 11
	EVENT_LOOT = 12
	EVENT_EQUIP = 13
	EVENT_DROP = 14
	EVENT_TELEPORT = 15
	EVENT_DAMAGE = 16
	EVENT_POPULATION = 17
	EVENT_KILL = 18
	EVENT_LIST = 19
	EVENT_WHO = 20
	EVENT_ZONE = 21
	EVENT_DESTROY = 22
	EVENT_HP = 23
	EVENT_BLINK = 24
	EVENT_CHECK = 26

	EVENT_ADD_ENTITY = 100
	EVENT_LOGOUT = 101
	EVENT_ITEM_RESWAPN = 102
	EVENT_ITEM_DESWAPN = 103
	EVENT_MOB_RESWAPN = 104
	EVENT_MOB_CALM = 105
)

func AquireEvent(t int, data... interface{}) *Event {
	return &Event{
		Type: t,
		Data: data,
	}
}

func ReleaseEvent(e *Event) {

}

func (event *Event) ToBytes() ([]byte, error) {
	objs := make([]interface{}, len(event.Data) + 1)
	objs = append([]interface{}{event.Type}, event.Data...)
	return json.Marshal(objs)
}