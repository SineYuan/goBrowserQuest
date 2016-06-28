package bqs

type Character struct {
	Entity
	Orientation int
	MaxHP       int
	HP          int
	TargetId    int
}

func (c *Character) ResetMaxHP(maxHitPoints int) {
	c.MaxHP = maxHitPoints
	c.HP = c.MaxHP
}

func (c *Character) HasFullHealth() bool {
	return c.HP == c.MaxHP
}

func (c *Character) RegenHealthBy(value int) {
	if c.HP < c.MaxHP {
		if c.HP + value <= c.MaxHP {
			c.HP = c.HP + value
		} else {
			c.HP = c.MaxHP
		}
	}
}

func (c *Character) SetTarget(targetId int) {
	c.TargetId = targetId
}
