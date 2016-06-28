package bqs

import (
	"log"
	"math/rand"
	"time"
	"strings"
	"strconv"

	"github.com/bitly/go-simplejson"
	"github.com/iris-contrib/websocket"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

const DEFAULT_CONNECT_TIME = 60 * 5

type Player struct {
	Character
	Armor, Weapon           int
	ArmorLevel, WeaponLevel int
	Name                    string
	previousZoneId          string
	zoneId                  string
	Orientation             int
	hasEnteredGame          bool
	conn                    *websocket.Conn
	EventCh                 chan *Event
	PacketChan              chan *Packet
	exit                    chan bool
	world                   *World
	timer                   *time.Timer

	logger                  *log.Logger
}

func NewPlayer(conn *websocket.Conn, world *World) *Player {
	p := &Player{
		conn: conn,
		world: world,
		PacketChan: make(chan *Packet, 100),
		exit: make(chan bool),
		timer: time.NewTimer(time.Second * DEFAULT_CONNECT_TIME),
	}
	p.Kind = TYPE_WARRIOR
	p.Id = <-world.IdCh
	return p
}

func (p *Player) HandleLoop() {
	for {
		select {
		case msg := <-p.PacketChan:
			if msg.Err != nil {
				log.Printf("Get from %v : %s\n", p.Id, string(msg.Data))
				if _, ok := msg.Err.(*websocket.CloseError); ok {
					log.Println("get a websocket.CloseError")
				}
				p.onDisconnect()
				return
			}
			cmd, args, err := ParseCmd(msg.Data)
			if err != nil {
				log.Println("Wrong data format:", err)
				continue
			}
			p.HandleCmd(cmd, args)
			log.Printf("Get from %v : %s\n", p.Id, string(msg.Data))
			p.timer.Reset(time.Second * DEFAULT_CONNECT_TIME)
		case <-p.timer.C:
			log.Println(p.Id, "Timeout!!")
			p.onDisconnect()
			p.conn.Close()
			return
		}
	}
}

func ParseCmd(data []byte) (int, *simplejson.Json, error) {
	js, err := simplejson.NewJson(data)
	if err != nil {
		log.Println(err)
		return -1, nil, err
	}
	cmd, err := js.GetIndex(0).Int()
	if err != nil {
		log.Println(err)
		return -1, nil, err
	}
	return cmd, js, nil
}

func (p *Player) HandleCmd(cmd int, args *simplejson.Json) {
	switch cmd {
	case EVENT_HELLO:
		log.Println("Hello")
		name := args.GetIndex(1).MustString()
		armorId := args.GetIndex(2).MustInt(-1)
		weaponId := args.GetIndex(3).MustInt(-1)
		if name != "" && armorId > 0 && weaponId > 0 {
			p.onHello(name, armorId, weaponId)
		}
	case EVENT_MOVE:
		x := args.GetIndex(1).MustInt(-1)
		y := args.GetIndex(2).MustInt(-1)
		if x > 0 && y > 0 {
			log.Println("Move", x, y)
			p.onMove(x, y)
		}
	case EVENT_LOOTMOVE:
		x := args.GetIndex(1).MustInt(-1)
		y := args.GetIndex(2).MustInt(-1)
		itemId := args.GetIndex(3).MustInt(-1)
		if x > 0 && y > 0 && itemId > 0 {
			log.Println("LootMove", x, y, itemId)
			p.onMove(x, y)
		}
	case EVENT_ATTACk:
		targetId := args.GetIndex(1).MustInt(-1)
		if targetId > 0 {
			p.onAttack(targetId)
		}
	case EVENT_HIT:
		targetId := args.GetIndex(1).MustInt(-1)
		if targetId > 0 {
			p.onHit(targetId)
		}
	case EVENT_HURT:
		targetId := args.GetIndex(1).MustInt(-1)
		if targetId > 0 {
			p.onHurt(targetId)
		}
	case EVENT_WHO:
		n := len(args.MustArray())
		if n > 1 {
			ids := make([]int, 0)
			for i := 1; i < n; i++ {
				if id, err := args.GetIndex(i).Int(); err == nil {
					ids = append(ids, id)
				}
			}
			if len(ids) > 0 {
				p.onWho(ids)
			}
		}
	case EVENT_CHAT:
		msg := args.GetIndex(1).MustString("")
		msg = strings.TrimSpace(msg)
		if msg != "" {
			p.onChat(msg)
		}
	case EVENT_LOOT:
		itemId := args.GetIndex(1).MustInt(-1)
		if itemId > 0 {
			p.onLoot(itemId)
		}
	case EVENT_TELEPORT:
		x := args.GetIndex(1).MustInt(-1)
		y := args.GetIndex(2).MustInt(-1)
		if x > 0 && y > 0 {
			log.Println("onTeleport", x, y)
			p.onTeleport(x, y)
		}
	case EVENT_ZONE:
		p.onZone()
	default:
		log.Println("No Such Event Type", cmd)
	}
}

func (p *Player) onHello(name string, armorId, weaponId int) {
	p.Name = name
	log.Println("on your mark")
	world := p.world
	p.X, p.Y = world.GetStartPosition(p)
	p.equipArmor(armorId)
	p.equipWeapon(weaponId)
	p.updateHP()
	p.EventCh = world.GetZoneChannelFromPosition(p.X, p.Y)
	p.zoneId = p.world.GetZoneIdFromPosition(p.X, p.Y)

	e := AquireEvent(EVENT_HELLO, p)
	e.PlayerID = p.Id

	p.EventCh <- e
	world.AddPlayer()
}

func (p *Player) onMove(x, y int) {
	world := p.world
	if world.IsValidPosition(x, y) {
		p.X, p.Y = x, y
		e := AquireEvent(EVENT_MOVE, p.Id, x, y)
		e.PlayerID = p.Id

		p.EventCh <- e
	} else {
		world.Logger.Println("invalide Pos")
	}
}

func (p *Player) onLootMove(x, y int) {
	world := p.world
	if world.IsValidPosition(x, y) {
		p.X, p.Y = x, y
		e := AquireEvent(EVENT_LOOTMOVE, p.Id, x, y)
		e.PlayerID = p.Id

		p.EventCh <- e
	} else {
		world.Logger.Println("invalide Pos")
	}
}

func (p *Player) onTeleport(x, y int) {
	world := p.world
	if world.IsValidPosition(x, y) {
		p.X, p.Y = x, y
		e := AquireEvent(EVENT_TELEPORT, p.Id, x, y)
		e.PlayerID = p.Id

		p.EventCh <- e
	} else {
		world.Logger.Println("invalide Pos")
	}
}

func (p *Player) onWho(ids []int) {
	e := AquireEvent(EVENT_WHO, ids)
	e.PlayerID = p.Id
	p.EventCh <- e
}

func (p *Player) onZone() {
	oldZoneCh := p.EventCh
	newZoneCh := p.world.GetZoneChannelFromPosition(p.X, p.Y)
	p.previousZoneId = p.zoneId
	p.zoneId = p.world.GetZoneIdFromPosition(p.X, p.Y)
	p.updateOrientationAfterZone()

	p.EventCh = newZoneCh
	e1 := AquireEvent(EVENT_DESTROY, p.Id)
	e1.PlayerID = p.Id
	e2 := AquireEvent(EVENT_ZONE, p)
	e2.PlayerID = p.Id
	oldZoneCh <- e1
	newZoneCh <- e2
}

func (p *Player) updateOrientationAfterZone() {
	if p.previousZoneId == "" {
		p.Orientation = ORIENTATION_DOWN
		return
	}
	xy1 := strings.Split(p.previousZoneId, "-")
	xy2 := strings.Split(p.zoneId, "-")

	zx1, _ := strconv.Atoi(xy1[0])
	zy1, _ := strconv.Atoi(xy1[1])
	zx2, _ := strconv.Atoi(xy2[0])
	zy2, _ := strconv.Atoi(xy2[1])

	if zx1 > zx2 {
		p.Orientation = ORIENTATION_LEFT
	} else if zx1 < zx2 {
		p.Orientation = ORIENTATION_RIGTH
	} else if zy1 > zy2 {
		p.Orientation = ORIENTATION_UP
	} else {
		p.Orientation = ORIENTATION_DOWN
	}
}

func (p *Player) onChat(msg string) {
	e := AquireEvent(EVENT_CHAT, p.Id, msg)

	p.send(e)

	e.PlayerID = p.Id
	p.EventCh <- e
}

func (p *Player) onLoot(itemId int) {
	e := AquireEvent(EVENT_LOOT, itemId)
	e.PlayerID = p.Id

	p.EventCh <- e
}

func (p *Player) onAttack(targetId int) {
	e := AquireEvent(EVENT_ATTACk, p.Id, targetId)
	e.PlayerID = p.Id

	p.EventCh <- e
}

func (p *Player) onHit(targetId int) {
	e := AquireEvent(EVENT_HIT, targetId)
	e.PlayerID = p.Id

	p.EventCh <- e
}

func (p *Player) onHurt(targetId int) {
	e := AquireEvent(EVENT_HURT, targetId)
	e.PlayerID = p.Id

	p.EventCh <- e
}

func (p *Player) send(event *Event) error {
	bytes, _ := event.ToBytes()

	log.Printf("send to %v: %v\n", p.Id, string(bytes))
	err := p.conn.WriteMessage(websocket.TextMessage, bytes)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (p *Player) sendBatch(events []*Event) error {
	bytes := []byte("[")
	for _, e := range events {
		if e != nil {
			eventBytes, _ := e.ToBytes()
			bytes = append(bytes, eventBytes...)
			bytes = append(bytes, ',')
		}
	}
	if len(bytes) > 1 {
		bytes[len(bytes) - 1] = ']'

		log.Printf("send to %v: %v\n", p.Id, string(bytes))
		err := p.conn.WriteMessage(websocket.TextMessage, bytes)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func (p *Player) onDisconnect() {
	close(p.exit)
	log.Println("disconnect")
	responseCh := make(chan struct{})
	e := AquireEvent(EVENT_LOGOUT, responseCh)
	e.PlayerID = p.Id
	p.EventCh <- e

	<-responseCh

	p.world.RemovePlayer()
}

func (p *Player) GetState() []interface{} {
	state := p.Entity.GetState()
	state = append(state, p.Name, p.Orientation, p.Armor, p.Weapon)
	return state
}

func (p *Player) equipArmor(armorId int) {
	p.Armor = armorId
	p.ArmorLevel = GetArmorRank(armorId)
}

func (p *Player) equipWeapon(weaponId int) {
	p.Weapon = weaponId
	p.WeaponLevel = GetWeaponRank(weaponId)
}

func (p *Player) updateHP() {
	p.ResetMaxHP(HpFormula(p.ArmorLevel))
}

func (p *Player) updatePosition() {
	p.X, p.Y = p.world.GetStartPosition(p)
}
