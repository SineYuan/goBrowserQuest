package bqs

import (
	"fmt"
	"log"
)

type Zone struct {
	ID            string
	ZX, ZY        int
	Map           *Map
	broadcastCh   chan *Event
	PlayersMap    map[int]*Player
	ItemsMap      map[int]*Item
	MobsMap       map[int]*Mob
	EntitiesIndex map[int]string
	NPCMap        map[int]*NPC
	EventCh       chan *Event
	Incoming      []int
	exitCh        chan struct{}

	Logger        *log.Logger
}

func NewZone(zx, zy int, bch chan *Event, ech chan struct{}, m *Map, logger *log.Logger) *Zone {
	return &Zone{
		ID: fmt.Sprintf("%v-%v", zx, zy),
		ZX: zx,
		ZY: zy,
		Map: m,
		PlayersMap: make(map[int]*Player),
		NPCMap: make(map[int]*NPC),
		ItemsMap: make(map[int]*Item),
		MobsMap: make(map[int]*Mob),
		EntitiesIndex: make(map[int]string),
		EventCh: make(chan *Event, 500),
		exitCh: ech,
		broadcastCh: bch,
		Logger: logger,
	}
}

func (z *Zone) Run() {
	for {
		select {
		case <-z.exitCh:
			z.Logger.Printf("Zone %v exit\n", z.ID)
			return
		case event := <-z.EventCh:
			if event.Type != EVENT_POPULATION {
				z.Logger.Println("Receive Event:", event)
			}
			z.HandleEvent(event)
		}
	}
}

func (z *Zone) HandleEvent(e *Event) {
	switch e.Type {
	case EVENT_ADD_ENTITY:
		z.onAddEntity(e)
	case EVENT_ITEM_RESWAPN:
		z.onItemReswapn(e)
	case EVENT_ITEM_DESWAPN:
		z.onItemDeswapn(e)
	case EVENT_MOB_RESWAPN:
		z.onMobReswapn(e)
	case EVENT_MOB_CALM:
		z.onMobCalm(e)
	case EVENT_HELLO:
		z.onHello(e)
	case EVENT_LOGOUT:
		z.onLogout(e)
	case EVENT_POPULATION:
		z.broadcastZone(e)
	case EVENT_MOVE:
		//z.broadcastZone(e)
		z.onPlayerMove(e)
	case EVENT_LOOTMOVE:
		z.broadcastZone(e)
	case EVENT_ATTACk:
		z.onAttack(e)
	case EVENT_HIT:
		z.onHit(e)
	case EVENT_HURT:
		z.onHurt(e)
	case EVENT_CHAT:
		z.broadcastZone(e)
	case EVENT_TELEPORT:
		z.broadcastZone(e)
	case EVENT_ZONE:
		z.onZone(e)
	case EVENT_DESTROY:
		z.onDestroy(e)
	case EVENT_WHO:
		z.onWho(e)
	case EVENT_LOOT:
		z.onLoot(e)
	default:
		z.Logger.Println("NO SUCH EVENT TYPE OF", e.Type)
	}

}

func (z *Zone) broadcastWorld(e *Event) {
	go func() {
		select {
		case z.broadcastCh <- e:
		}
	}()
}

func (z *Zone) broadcastZone(e *Event) {
	for _, p := range z.PlayersMap {
		if p.Id != e.PlayerID {
			if err := p.send(e); err != nil {
				z.Logger.Println(err)
			}
		}
	}
}

func (z *Zone) onHello(e *Event) {
	p := e.Data[0].(*Player)

	z.Logger.Println(z.ID, "add player", p.Id)
	event := AquireEvent(EVENT_WELCOME, p.Id, p.Name, p.X, p.Y, p.HP)
	p.send(event)

	z.onZone(e)
}

func (z *Zone) onPlayerMove(e *Event) {
	z.broadcastZone(e)

	pid := e.Data[0].(int)
	z.Logger.Println("[DEBUG]on Player", pid, "move")
	if player := z.PlayersMap[pid]; player != nil {
		for _, mob := range z.MobsMap {
			if mob.TargetId == pid {
				// chase player
				mob.ChaseTarget(z.ID, z.Map, player.X, player.Y)
				z.Logger.Println("[DEBUG]chase Player", pid, "to", mob.X, mob.Y)
			}
		}
	}
}

func (z *Zone) isOutOfBound(x, y int) bool {
	return z.ID == z.Map.GetGroupIdFromPosition(x, y)
}

func (z *Zone) onWho(e *Event) {
	ids, ok := e.Data[0].([]int)
	events := make([]*Event, len(ids))
	if !ok {
		return
	}
	p := z.PlayersMap[e.PlayerID]
	if p == nil {
		z.Logger.Println("[ERROR] NO SUCH USER", e.PlayerID, "IN ZONE", z.ID)
		return
	}
	for i, id := range ids {
		entity := z.getEntityById(id)
		if entity != nil {
			events[i] = AquireEvent(EVENT_SWAPN, entity.GetState()...)
		}
	}
	if len(events) > 0 {
		p.sendBatch(events)
	}
}

func (z *Zone) onZone(e *Event) {
	p := e.Data[0].(*Player)
	z.PlayersMap[p.Id] = p
	z.EntitiesIndex[p.Id] = "player"

	z.Logger.Println("player", p.Id, "into zone", z.ID)

	listEvent := z.listEntites(p)
	p.send(listEvent)

	spawnEvent := AquireEvent(EVENT_SWAPN, p.GetState()...)
	spawnEvent.PlayerID = p.Id
	z.broadcastZone(spawnEvent)
}

func (z *Zone) listEntites(p *Player) *Event {
	listEvent := AquireEvent(EVENT_LIST)
	for k, _ := range z.NPCMap {
		listEvent.Data = append(listEvent.Data, k)
	}
	for k, player := range z.PlayersMap {
		if player.Id != p.Id {
			listEvent.Data = append(listEvent.Data, k)
		}
	}
	for k, item := range z.ItemsMap {
		if !item.IsDestroy {
			listEvent.Data = append(listEvent.Data, k)
		}
	}
	for k, mob := range z.MobsMap {
		if !mob.IsDead {
			listEvent.Data = append(listEvent.Data, k)
		}
	}
	return listEvent
}

func (z *Zone) onDestroy(e *Event) {
	pid := e.PlayerID
	p := z.PlayersMap[e.PlayerID]
	if p == nil {
		z.Logger.Println("[BUG] Destory player", pid, "not in zone", z.ID)
	}
	z.clearPlayer(pid)

	z.broadcastZone(e)
	z.Logger.Println("[INFO] Player", pid, "left zone", z.ID)
}

func (z *Zone) onLoot(e *Event) {
	itemId := e.Data[0].(int)
	p := z.PlayersMap[e.PlayerID]
	if p == nil {
		return
	}
	if item := z.ItemsMap[itemId]; item != nil {
		deswapnEvent := AquireEvent(EVENT_DESPAWN, itemId)
		z.broadcastZone(deswapnEvent)
		item.IsDestroy = true
		if item.IsStatic {
			item.ReswapnLater(z.EventCh)
		}
		kind := item.Kind
		if kind.Id == TYPE_FIREPOTION.Id {
			// TODO
		} else if IsHealingItem(kind) {

			amount := 0
			switch kind.Id {
			case TYPE_FLASK.Id:
				amount = 40
			case TYPE_BURGER.Id:
				amount = 100
			}
			if amount > 0 && !p.HasFullHealth() {
				p.RegenHealthBy(amount)
				healthEvent := AquireEvent(EVENT_HEALTH, p.HP)
				p.send(healthEvent)
			}
		} else if IsArmor(kind) || IsWeapon(kind) {
			equipEvent := AquireEvent(EVENT_EQUIP, p.Id, kind.Id)
			z.broadcastZone(equipEvent)

			if IsArmor(kind) {
				p.equipArmor(kind.Id)
				p.updateHP()
				HPEvent := AquireEvent(EVENT_HP, p.MaxHP)
				p.send(HPEvent)
			} else {
				p.equipWeapon(kind.Id)
			}
		}
	}
}

func (z *Zone) onAttack(e *Event) {
	p := z.PlayersMap[e.PlayerID]
	targetId := e.Data[1].(int)
	if target := z.MobsMap[targetId]; target != nil && p != nil {
		p.SetTarget(targetId)
		z.broadcastZone(e)
	}
}

func (z *Zone) onHit(e *Event) {
	p := z.PlayersMap[e.PlayerID]
	targetId := e.Data[0].(int)
	if target := z.MobsMap[targetId]; target != nil && p != nil {
		z.broadcastZone(e)

		z.onMobAttacked(target, p)
	}
}

func (z *Zone) onMobAttacked(m *Mob, p *Player) {
	m.ResetHateLater(z.EventCh)

	dmg := DmgFormula(p.WeaponLevel, m.ArmorLevel)
	if dmg > 0 {
		m.HP -= dmg
		if m.HP > 0 {
			dmgEvent := AquireEvent(EVENT_DAMAGE, m.Id, dmg)
			p.send(dmgEvent)

			m.AddHate(p.Id, dmg)
			if maxHateTarget := m.ChooseMobTarget(); maxHateTarget > 0 {
				if maxHateTarget != m.TargetId {
					m.TargetId = maxHateTarget
				}
				attackEvent := AquireEvent(EVENT_ATTACk, m.Id, m.TargetId)
				z.broadcastZone(attackEvent)
			}
		} else {
			z.Logger.Println("[DEBUG] m", m.Id, "DEAD!")
			m.IsDead = true
			if dropItem := m.DropItem(); dropItem != nil {
				z.Logger.Println("[DEBUG] m", m.Id, "DROP!", dropItem)
				dropItem.DeswapnLater(z.EventCh)
				z.ItemsMap[dropItem.Id] = dropItem

				swanItemEvent := AquireEvent(EVENT_SWAPN, dropItem.GetState()...)
				z.broadcastZone(swanItemEvent)
			}
			z.Logger.Println("[DEBUG] m", m.Id, "DESPAWN LATER!")
			m.ReswapnLater(z.EventCh)

			deswapnEvent := AquireEvent(EVENT_DESPAWN, m.Id)
			z.broadcastZone(deswapnEvent)
			z.Logger.Println("[DEBUG] m", m.Id, "DESPAWN!")
		}
	}
}

func (z *Zone) onHurt(e *Event) {
	// hurt by mob
	p := z.PlayersMap[e.PlayerID]
	targetId := e.Data[0].(int)
	if target := z.MobsMap[targetId]; target != nil && p != nil {
		z.onPlayerAttacked(p, target)
	}
}

func (z *Zone) onPlayerAttacked(p *Player, m *Mob) {
	dmg := DmgFormula(m.WeaponLevel, p.ArmorLevel)
	if dmg > 0 {
		p.HP -= dmg
		if p.HP > 0 {
			healthEvent := AquireEvent(EVENT_HEALTH, p.HP)
			p.send(healthEvent)
		} else {
			deswapnEvent := AquireEvent(EVENT_DESTROY, p.Id)
			deswapnEvent.PlayerID = p.Id
			z.onDestroy(deswapnEvent)

			delete(z.PlayersMap, p.Id)
			e := AquireEvent(EVENT_POPULATION, -1)
			z.broadcastWorld(e)
		}
	}
}

func (z *Zone) clearPlayer(pid int) {
	if p := z.PlayersMap[pid]; p != nil {
		delete(z.PlayersMap, pid)
		delete(z.EntitiesIndex, pid)

		for _, mob := range z.MobsMap {
			if mob.TargetId == pid {
				if maxHateTarget := mob.ChooseMobTarget(); maxHateTarget > 0 {
					mob.TargetId = maxHateTarget
				} else {
					mob.TargetId = 0
				}
			}
			delete(mob.Haters, pid)
		}
	}
}


func (z *Zone) onLogout(e *Event) {
	despawnEvent := AquireEvent(EVENT_DESPAWN, e.PlayerID)
	z.broadcastZone(despawnEvent)

	z.clearPlayer(e.PlayerID)
	if respCh, ok := e.Data[0].(chan struct{}); ok {
		close(respCh)
	}
}

func (z *Zone) onAddEntity(e *Event) {
	class := e.Data[0].(string)
	switch class {
	case "npc":
		npc := e.Data[1].(*NPC)
		z.NPCMap[npc.Id] = npc
		z.EntitiesIndex[npc.Id] = "npc"
	case "item":
		item := e.Data[1].(*Item)
		z.ItemsMap[item.Id] = item
		z.EntitiesIndex[item.Id] = "item"
	case "mob":
		mob := e.Data[1].(*Mob)
		z.MobsMap[mob.Id] = mob
		z.EntitiesIndex[mob.Id] = "mob"
		mob.SetEventCh(z.EventCh)
	}
}

func (z *Zone) onItemReswapn(e *Event) {
	itemId := e.Data[0].(int)
	if item := z.ItemsMap[itemId]; item != nil {
		z.Logger.Println("[DEBUG] Item", item, "Reswapn")

		item.IsDestroy = false
		swapnEvent := AquireEvent(EVENT_SWAPN, item.GetState()...)
		z.broadcastZone(swapnEvent)
	}
}

func (z *Zone) onItemDeswapn(e *Event) {
	itemId := e.Data[0].(int)
	if item := z.ItemsMap[itemId]; item != nil {
		z.Logger.Println("[DEBUG] Item", item, "Deswapn")

		delete(z.ItemsMap, itemId)
		deswapnEvent := AquireEvent(EVENT_DESPAWN, item.Id)
		z.broadcastZone(deswapnEvent)
	}
}

func (z *Zone) onMobReswapn(e *Event) {
	mobId := e.Data[0].(int)
	if mob := z.MobsMap[mobId]; mob != nil {
		z.Logger.Println("[DEBUG] Mob", mob, "Reswapn")

		mob.IsDead = false
		mob.RecoveryHP()

		swapnEvent := AquireEvent(EVENT_SWAPN, mob.GetState()...)
		z.broadcastZone(swapnEvent)
	}
}

func (z *Zone) onMobCalm(e *Event) {
	mobId := e.Data[0].(int)
	if mob := z.MobsMap[mobId]; mob != nil {
		z.Logger.Println("[DEBUG] Mob", mob, "Calm Down")

		mob.RecoveryHP()
		for k := range mob.Haters {
			delete(mob.Haters, k)
		}
		mob.TargetId = 0
		if mob.X != mob.OriginX || mob.Y != mob.OriginY {
			mob.X, mob.Y = mob.OriginX, mob.OriginY
			moveEvent := AquireEvent(EVENT_MOVE, mob.Id, mob.X, mob.Y)
			z.broadcastZone(moveEvent)
		}
		mob.TargetId = 0
	}
}

func (z *Zone) getEntityById(id int) EntityInterface {
	entityClass := z.EntitiesIndex[id]
	switch entityClass {
	case "npc":
		return z.NPCMap[id]
	case "player":
		return z.PlayersMap[id]
	case "item":
		return z.ItemsMap[id]
	case "mob":
		return z.MobsMap[id]
	default:
		z.Logger.Println("[ERROR] No Such Entity id:", id, "IN ZONE", z.ID)
		return nil
	}
}




