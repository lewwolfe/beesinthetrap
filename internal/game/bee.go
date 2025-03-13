package game

import "math/rand/v2"

type Bee struct {
	beeType      string
	hp           int
	attackDamage int
	hitDamage    int
	missChance   float64
}

func (b *Bee) Attack() int {
	if rand.Float64() < b.missChance {
		return 0
	}
	return b.attackDamage
}

func (b *Bee) Hit() int {
	b.hp -= b.hitDamage
	return b.hitDamage
}

func (b *Bee) IsDead() bool {
	return b.hp <= 0
}

func (b *Bee) GetBeeType() string {
	return b.beeType
}

func (b *Bee) GetHP() int {
	return b.hp
}
