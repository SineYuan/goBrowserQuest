package bqs

import (
	"time"
	"fmt"
)

type Item struct {
	Entity
	IsStatic bool
	Blink bool
	IsDestroy bool
}

const DEFAULT_ITEM_RESWAPN_TIME = 30 * time.Second
const DEFAULT_DROP_ITEM_TIME = 5 * time.Second

var itemIDCh chan int

func init() {
	itemIDCh := make(chan int)
	count := 0
	go func() {
		for {
			select {
			case itemIDCh <- 110000 + count:
				count++
			}

		}
	}()
}

func NewItem(id int, kind EntityType, x, y int) *Item {
	item := new(Item)
	item.Id = id
	item.X, item.Y = x, y
	item.Kind = kind

	return item
}

func NewDropItem(kind EntityType, x, y int) *Item {
	fmt.Println("drop item", kind)
	item := new(Item)
	item.Id = <- itemIDCh
	item.X, item.Y = x, y
	item.Kind = kind
	fmt.Println("return drop item", kind)

	return item
}

func IsItem(kind EntityType) bool {
	cls := kind.Class
	return cls == "armor" || cls == "weapon" || cls == "object"
}

func IsHealingItem(kind EntityType) bool {
	kindId := kind.Id
	return kindId == TYPE_FLASK.Id || kindId == TYPE_BURGER.Id
}

func IsArmor(kind EntityType) bool {
	cls := kind.Class
	return cls == "armor"
}

func IsWeapon(kind EntityType) bool {
	cls := kind.Class
	return cls == "weapon"
}

func (i *Item) ReswapnLater(eventCh chan *Event) {
	// TODO optimizate with priorityqueue
	go func() {
		timer := time.NewTimer(DEFAULT_ITEM_RESWAPN_TIME)
		<- timer.C
		e := AquireEvent(EVENT_ITEM_RESWAPN, i.Id)
		select {
		case eventCh <- e:
		}
	}()
}

func (i *Item) DeswapnLater(eventCh chan *Event) {
	// TODO optimizate with priorityqueue
	go func() {
		timer := time.NewTimer(DEFAULT_DROP_ITEM_TIME)
		<- timer.C
		e := AquireEvent(EVENT_ITEM_DESWAPN, i.Id)
		select {
		case eventCh <- e:
		}
	}()
}


