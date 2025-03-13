package game

type Player struct {
	hp int
}

func (p *Player) getHealth() int {
	return p.hp
}

func (p *Player) IsDead() bool {
	return p.hp <= 0
}
