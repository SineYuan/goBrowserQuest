package bqs

import (
	"os"
	"sync"
	"time"
	"math/rand"
	"log"
	"strconv"
)

type World struct {
	playersCount int
	lock         sync.Mutex
	Map          *Map
	IdCh         chan int
	ZonesMap     map[string]*Zone
	BroadcastCh  chan *Event
	stopCh       chan struct{}
	zonesStopCh  chan struct{}

	itemCount    int

	Logger       *log.Logger
}

func NewWrold(cfg *Config) *World {
	m, err := LoadMap(cfg.MapFilePath)
	if err != nil {
		panic(err)
	}
	w := &World{
		Map: m,
		stopCh: make(chan struct{}),
		zonesStopCh: make(chan struct{}),
		Logger: log.New(os.Stdout, "", 0),
	}
	w.startIDService()
	w.startBroadcastService()
	w.initZones()
	w.spawnStaticEntities()
	return w
}

func (w *World) initZones() {
	w.ZonesMap = make(map[string]*Zone)
	for x := 0; x < w.Map.GroupWidth; x++ {
		for y := 0; y < w.Map.GroupHeight; y++ {
			zone := NewZone(x, y, w.BroadcastCh, w.zonesStopCh, w.Map, w.Logger)
			go zone.Run()
			w.ZonesMap[zone.ID] = zone
		}
	}
}

func (w *World) spawnStaticEntities() {
	var count int
	for tidStr, kindName := range w.Map.StaticEntities {
		kind := GetEntityTypeByString(kindName)
		tid, err := strconv.Atoi(tidStr)
		if err != nil {
			w.Logger.Println("[ERROR] error in map json format")
		}
		x, y := w.Map.TileIndexToGridPosition(tid)
		cls := kind.Class
		switch cls {
		case "npc":
			npc := NewNPC(kind, x, y)
			w.addNPC(npc)
		case "mob":
			count++
			mob := NewMob(count, kind, x, y)
			w.addMob(mob)
		default:
			if IsItem(kind) {
				item := NewItem(90000 + w.itemCount, kind, x, y)
				item.IsStatic = true
				w.addItem(item)
			}
		}
	}
}

func (w *World) addNPC(npc *NPC) {
	zch := w.GetZoneChannelFromPosition(npc.X, npc.Y)
	e := AquireEvent(EVENT_ADD_ENTITY, "npc", npc)
	select {
	case zch <- e:
	}
}

func (w *World) addMob(mob *Mob) {
	zch := w.GetZoneChannelFromPosition(mob.X, mob.Y)
	e := AquireEvent(EVENT_ADD_ENTITY, "mob", mob)
	select {
	case zch <- e:
	}
}

func (w *World) addItem(item *Item) {
	zch := w.GetZoneChannelFromPosition(item.X, item.Y)
	e := AquireEvent(EVENT_ADD_ENTITY, "item", item)
	select {
	case zch <- e:
	}
	w.itemCount++
}

func (w *World) startBroadcastService() {
	w.BroadcastCh = make(chan *Event, 500)
	go func() {
		for {
			select {
			case <-w.stopCh:
				close(w.zonesStopCh)
				return
			case event := <-w.BroadcastCh:
				w.Logger.Println("recive boardcast event:", event)
				switch event.Type {
				case EVENT_POPULATION:
					delta := event.Data[0].(int)
					w.playersCount += delta
					event.Data[0] = w.playersCount
					w.broadcast(event)
				}
			}
		}
	}()
}

func (w *World) broadcast(e *Event) {
	for _, z := range w.ZonesMap {
		z.EventCh <- e
	}
}

func (w *World) Stop() {
	close(w.stopCh)
}

func (w *World) AddPlayer() {
	go func() {
		e := AquireEvent(EVENT_POPULATION, 1)
		select {
		case w.BroadcastCh <- e:
		}
	}()
}

func (w *World) RemovePlayer() {
	go func() {
		e := AquireEvent(EVENT_POPULATION, -1)
		select {
		case w.BroadcastCh <- e:
		}
	}()
}

func (w *World) IsValidPosition(x, y int) bool {
	return w.Map.IsValidPosition(x, y)
}

func (w *World) GetZoneChannelFromPosition(x, y int) chan *Event {
	zid := w.Map.GetGroupIdFromPosition(x, y)
	zone := w.ZonesMap[zid]
	return zone.EventCh
}

func (w *World) GetZoneIdFromPosition(x, y int) string {
	zid := w.Map.GetGroupIdFromPosition(x, y)
	return zid
}

func (w *World) GetStartPosition(p *Player) (int, int) {
	return w.Map.GetRandomStartingPosition()
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func (w *World) startIDService() {
	idCh := make(chan int)
	go func() {
		count := 0
		id := 50000 + 100 * rand.Intn(99) + count
		for {
			select {
			case <-w.stopCh:
				return
			case idCh <- id:
				count++
				id = 50000 + 100 * rand.Intn(99) + count
			}
		}
	}()
	w.IdCh = idCh
}
