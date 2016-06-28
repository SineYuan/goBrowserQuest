package bqs

import (
	"time"
	"math/rand"
)



func DmgFormula(weaponLevel, armorLevel int) int {
		dealt := weaponLevel * random(5, 10)
		absorbed := armorLevel * random(1, 3)
		dmg := dealt - absorbed

		if dmg <= 0 {
			return random(0, 3)
		} else {
			return dmg
		}
}

func HpFormula(armorLevel int) int {
	var hp = 80 + ((armorLevel - 1) * 30)
	return hp
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max - min) + min
}