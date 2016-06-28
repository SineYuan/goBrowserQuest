package bqs

import (
	"time"
)

type Mob struct {
	Character
	OriginX, OriginY        int
	ArmorLevel, WeaponLevel int
	Haters                  map[int]int
	IsDead                  bool
	ReswapnTimer            *time.Timer
	HateTimer               *time.Timer
	EventCh                 chan *Event
}

type MobProperty struct {
	Name   string
	Drops  map[string]int
	HP     int
	Armor  int
	Weapon int
}

var MobPropertiesMap = map[int]*MobProperty{
	TYPE_RAT.Id : &MobProperty{
		Drops: map[string]int{
			"flask": 40,
			"burger": 10,
			"firepotion": 5,
		},
		HP: 25,
		Armor: 1,
		Weapon: 1,
	},
	TYPE_SKELETON.Id: &MobProperty{
		Drops: map[string]int{
			"flask": 40,
			"mailarmor": 10,
			"axe": 20,
			"firepotion": 5,
		},
		HP: 110,
		Armor: 2,
		Weapon: 2,
	},
	TYPE_GOBLIN.Id: &MobProperty{
		Drops: map[string]int{
			"flask": 40,
			"mailarmor": 20,
			"axe": 20,
			"firepotion": 5,
		},
		HP: 19,
		Armor: 2,
		Weapon: 1,
	},
	TYPE_OGRE.Id: &MobProperty{
		Drops: map[string]int{
			"burger": 10,
			"flask": 50,
			"platearmor": 20,
			"morningstar": 20,
			"firepotion": 5,
		},
		HP: 200,
		Armor: 3,
		Weapon: 2,
	},
	TYPE_SPECTRE.Id: &MobProperty{
		Drops: map[string]int{
			"flask": 30,
			"redarmor": 40,
			"redsword": 30,
			"firepotion": 5,
		},
		HP: 250,
		Armor: 2,
		Weapon: 4,
	},
	TYPE_DEATHKNIGHT.Id: &MobProperty{
		Drops: map[string]int{
			"burger": 95,
			"firepotion": 5,
		},
		HP: 250,
		Armor: 3,
		Weapon: 3,
	},
	TYPE_CRAB.Id: &MobProperty{
		Drops: map[string]int{
			"flask": 50,
			"axe": 20,
			"leatherarmor": 10,
			"firepotion": 5,
		},
		HP: 60,
		Armor: 2,
		Weapon: 1,
	},
	TYPE_SNAKE.Id: &MobProperty{
		Drops: map[string]int{
			"flask": 50,
			"mailarmor": 10,
			"morningstar": 10,
			"firepotion": 5,
		},
		HP: 150,
		Armor: 3,
		Weapon: 2,
	},
	TYPE_SKELETON2.Id: &MobProperty{
		Drops: map[string]int{
			"flask": 60,
			"platearmor": 15,
			"bluesword": 15,
			"firepotion": 5,
		},
		HP: 200,
		Armor: 3,
		Weapon: 3,
	},
	TYPE_EYE.Id: &MobProperty{
		Drops: map[string]int{
			"flask": 50,
			"redarmor": 20,
			"redsword": 10,
			"firepotion": 5,
		},
		HP: 200,
		Armor: 3,
		Weapon: 3,
	},
	TYPE_BAT.Id: &MobProperty{
		Drops: map[string]int{
			"flask": 50,
			"axe": 20,
			"firepotion": 5,
		},
		HP: 80,
		Armor: 2,
		Weapon: 1,
	},
	TYPE_WIZARD.Id: &MobProperty{
		Drops: map[string]int{
			"flask": 50,
			"platearmor": 20,
			"firepotion": 5,
		},
		HP: 100,
		Armor: 2,
		Weapon: 6,
	},

	TYPE_BOSS.Id: &MobProperty{
		Drops: map[string]int{
			"goldensword": 100,
		},
		HP: 700,
		Armor: 6,
		Weapon: 7,
	},
}

func NewMob(mobCount int, kind EntityType, x, y int) *Mob {
	mob := new(Mob)
	mob.Id = 70000 + kind.Id * 100 + mobCount
	mob.X, mob.Y = x, y
	mob.OriginX, mob.OriginY = x, y
	mob.Kind = kind
	mob.Haters = make(map[int]int)
	if prop := GetMobProperty(kind); prop != nil {
		mob.HP = prop.HP
		mob.ArmorLevel, mob.WeaponLevel = prop.Armor, prop.Weapon
		mob.ReswapnTimer = time.NewTimer(0)
		<-mob.ReswapnTimer.C
		mob.HateTimer = time.NewTimer(0)
		<-mob.HateTimer.C
		mob.startTimerService()
		return mob
	}
	return nil
}

func (m *Mob) GetState() []interface{} {
	state := m.Character.GetState()
	return state
}

func (m *Mob) receiveDamage(dmg int) {
	m.HP -= dmg;
}

func (m *Mob) AddHate(playerId, dmg int) {
	m.Haters[playerId] += dmg
}

func (m *Mob) ChooseMobTarget() int {
	var max int
	var maxPid int
	for pid, hate := range m.Haters {
		if hate > max {
			max = hate
			maxPid = pid
		}
	}
	if max <= 0 {
		return -1
	}
	return maxPid
}

const DEFAULT_MOB_RESWAPN_TIME = 45 * time.Second
const DEFAULT_MOB_CALM_TIME = 5 * time.Second

func (mob *Mob) ReswapnLater(eventCh chan *Event) {

	mob.ReswapnTimer.Reset(DEFAULT_MOB_RESWAPN_TIME)
}

func (mob *Mob) ResetHateLater(eventCh chan *Event) {
	mob.HateTimer.Reset(DEFAULT_MOB_CALM_TIME)
}

func (mob *Mob) RecoveryHP() {
	if prop := GetMobProperty(mob.Kind); prop != nil {
		mob.HP = prop.HP
	}
}

func (mob *Mob) returnToSpawningPosition() {
}

func (mob *Mob) ChaseTarget(zoneId string, mp *Map, targetX, targetY int) {
	zid := mp.GetGroupIdFromPosition(targetX, targetY)
	if zoneId != zid {
		mob.X, mob.Y = targetX, targetY
	} else {
		pointsAround := make([][2]int, 0)
		for _, p := range [][2]int{
			[2]int{targetX, targetY + 1},
			[2]int{targetX + 1, targetY},
			[2]int{targetX, targetY - 1},
			[2]int{targetX - 1, targetY},
		} {
			if mp.IsValidPosition(p[0], p[1]) && zoneId == mp.GetGroupIdFromPosition(p[0], p[1]) {
				pointsAround = append(pointsAround, p)
			}
		}
		minLen := 999999
		minIndex := 0
		for i, p := range pointsAround {
			pathLength := (mob.X - p[0]) * (mob.X - p[0]) + (mob.Y - p[1]) * (mob.Y - p[1])
			if pathLength <= minLen {
				minLen = pathLength
				minIndex = i
			}
		}
		mob.X, mob.Y = pointsAround[minIndex][0], pointsAround[minIndex][1]
	}
}

func (m *Mob) DropItem() *Item {
	if prop := GetMobProperty(m.Kind); prop != nil {
		seed := random(0, 99)
		acc := 0
		itemType := ""
		for k, v := range prop.Drops {
			if seed >= acc && seed < acc + v {
				itemType = k
			}
			acc += v
		}
		if kind := GetEntityTypeByString(itemType); kind.Id != 0 {
			item := NewItem(11000000 + m.OriginX * 1000 + m.OriginY, kind, m.X, m.Y)
			return item
		}
	}
	return nil
}

func (mob *Mob) startTimerService() {
	go func() {
		for {
			select {
			case <-mob.ReswapnTimer.C:
				e := AquireEvent(EVENT_MOB_RESWAPN, mob.Id)
				mob.EventCh <- e
			case <-mob.HateTimer.C:
				e := AquireEvent(EVENT_MOB_CALM, mob.Id)
				mob.EventCh <- e
			}
		}
	}()
}

func (mob *Mob) SetEventCh(evenCh chan *Event) {
	mob.EventCh = evenCh
}

func GetMobProperty(kind EntityType) *MobProperty {
	return MobPropertiesMap[kind.Id]
}
