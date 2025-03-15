package game

import "math/rand"

type Bee struct {
	beeType      BeeType
	hp           int
	attackDamage int
	hitDamage    int
	missChance   float64
}

type BeeType int

// Enum(ish) representation of bee types
const (
	QueenBee BeeType = iota
	WorkerBee
	DroneBee
)

func (bt BeeType) String() string {
	return [...]string{"Queen", "Worker", "Drone"}[bt]
}

func (b *Bee) Attack(rng *rand.Rand) int {
	if rng.Float64() > b.missChance {
		return b.attackDamage
	}
	return 0

}

func (b *Bee) Hit() int {
	b.hp -= b.hitDamage
	return b.hitDamage
}

func (b *Bee) IsDead() bool {
	return b.hp <= 0
}

func (b *Bee) GetBeeType() BeeType {
	return b.beeType
}

func (b *Bee) GetHP() int {
	return b.hp
}
