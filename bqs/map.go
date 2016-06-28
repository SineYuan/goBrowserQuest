package bqs
//package main

import (
	"io/ioutil"
	"fmt"
	"encoding/json"
	"math/rand"
)

type Map struct {
	Width, Height  int
	TileSize       int
	Collisions     []int
	Tilesize       int
	Grid           [][]bool
	CheckPoints    []CheckPoint
	ChestAreas     []ChestArea
	StaticChests   []Chest
	MobAreas       []MobArea `json:"roamingAreas"`
	StaticEntities map[string]string

	GroupWidth     int
	GroupHeight    int
	ZoneWidth      int
	ZoneHeight     int

	CheckPointsMap map[int]*CheckPoint
	StartingAreas  []*CheckPoint
}

const zoneWidth = 28
const zoneHeight = 12

type CheckPoint struct {
	Id   int
	X, Y int
	W, H int
	S    int
}

type Chest struct {
	X, Y  int
	Items []int `json:"i"`
}

func LoadMap(filename string) (*Map, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var m = Map{}
	err = json.Unmarshal(bytes, &m)
	if err != nil {
		return nil, err
	}

	err = m.InitMap()
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (m Map) String() string {
	formatString := "Map.Collisions: %v\n"
	formatString += "Map.CheckPoints: %v\n"
	formatString += "Map.StaticChests: %v\n"
	formatString += "Map.ChestAreas: %v\n"
	formatString += "Map.MobAreas: %v\n"
	formatString += "Map.StaticEntities: %v\n"
	formatString += "Map.Width: %v Map.Height: %v\n"
	formatString += "Map.TileSize: %v\n"

	return fmt.Sprintf(formatString,
		m.Collisions,
		m.CheckPoints,
		m.StaticChests,
		m.ChestAreas,
		m.MobAreas,
		m.StaticEntities,
		m.Width, m.Height, m.TileSize)
}

func (m *Map) InitMap() error {
	m.ZoneWidth = zoneWidth
	m.ZoneHeight = zoneHeight
	m.GroupWidth = m.Width / m.ZoneWidth
	m.GroupHeight = m.Height / m.ZoneHeight
	m.generateCollisionGrid()
	m.initCheckPoints()
	return nil
}

func (m *Map) generateCollisionGrid() {
	tileIndex := 0
	collisionsSet := make(map[int]bool)
	for i := range m.Collisions {
		collisionsSet[i] = true
	}
	m.Grid = make([][]bool, m.Height)
	for j, i := 0, 0; i < m.Height; i++ {

		m.Grid[i] = make([]bool, m.Width)
		for j = 0; j < m.Width; j++ {
			if _, ok := collisionsSet[tileIndex]; ok {
				m.Grid[i][j] = true
			} else {
				m.Grid[i][j] = false
			}
			tileIndex += 1;
		}
	}
}

func (m *Map) initCheckPoints() {
	m.CheckPointsMap = make(map[int]*CheckPoint)
	for i, cp := range m.CheckPoints {
		m.CheckPointsMap[cp.Id] = &m.CheckPoints[i]
		if cp.S == 1 {
			m.StartingAreas = append(m.StartingAreas, &m.CheckPoints[i])
		}
	}
}

func (m *Map) IsOutOfBounds(x, y int) bool {
	return x <= 0 || x >= m.Width || y <= 0 || y >= m.Height;
}

func (m *Map) IsColliding(x, y int) bool {
	return m.Grid[y][x]
}

func (m *Map) IsValidPosition(x, y int) bool {
	if !m.IsOutOfBounds(x, y) && !m.IsColliding(x, y) {
		return true
	}
	return false
}

func (m *Map) GetRandomStartingPosition() (int, int) {
	n := rand.Int() % len(m.StartingAreas)
	cp := m.StartingAreas[n]
	return cp.GetRandomPosition()
}

func (m *Map) GetGroupIdFromPosition(x, y int) string {
	zx, zy := (x - 1) / m.ZoneWidth, (y - 1) / m.ZoneHeight
	return fmt.Sprintf("%v-%v", zx, zy)
}

func (m *Map) TileIndexToGridPosition(tileNum int) (x, y int) {
	if tileNum != 0 {
		if tileNum % m.Width == 0 {
			x = m.Width
		} else {
			x = (tileNum % m.Width)
		}
	}

	y = (tileNum - 1) / m.Width

	return
}

func (cp *CheckPoint) GetRandomPosition() (x, y int) {
	x = cp.X + rand.Int() % cp.W
	y = cp.Y + rand.Int() % cp.H
	return
}