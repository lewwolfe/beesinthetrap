package game

type Bee struct {
	Type          string
	HP            int
	AttackDamage  int
	DefenceDamage int
	MissChance    float64
}

func (b *Bee) Attack(player *Player) {
}

func (b *Bee) TakeDamage(damage int) {
	b.HP -= damage
}

func (b *Bee) IsDead() bool {
	return b.HP <= 0
}

func (b *Bee) GetType() string {
	return b.Type
}

func (b *Bee) GetHP() int {
	return b.HP
}
