package game

import "math/rand/v2"

type Player struct {
	hp         int
	missChance float64
}

func (p *Player) Attack() bool {
	return rand.Float64() < p.missChance
}

func (p *Player) Sting(damage int) {
	p.hp -= damage
}

func (p *Player) IsDead() bool {
	return p.hp <= 0
}

func (p *Player) Gethp() int {
	return p.hp
}
