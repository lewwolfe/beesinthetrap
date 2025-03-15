package game

import (
	"math/rand"
	"testing"
)

// TestBeeAttack tests basic attack functionality of a Bee
func TestBeeAttack(t *testing.T) {
	alwaysHitBee := &Bee{
		beeType:      WorkerBee,
		hp:           10,
		attackDamage: 5,
		missChance:   0,
	}

	source := rand.NewSource(42)
	rng := rand.New(source)

	// Test always hit
	damage := alwaysHitBee.Attack(rng)
	if damage != 5 {
		t.Errorf("Expected damage 5 when bee always hits, got %d", damage)
	}

	// Test bee always misses
	alwaysMissBee := &Bee{
		beeType:      WorkerBee,
		hp:           10,
		attackDamage: 5,
		missChance:   1,
	}

	damage = alwaysMissBee.Attack(rng)
	if damage != 0 {
		t.Errorf("Expected damage 0 when bee always misses, got %d", damage)
	}
}

// TestBeeHit tests that a bee takes the right amount of damage
func TestBeeHit(t *testing.T) {
	bee := &Bee{
		beeType:   WorkerBee,
		hp:        10,
		hitDamage: 3,
	}

	damage := bee.Hit()
	if damage != 3 {
		t.Errorf("Expected hit damage 3, got %d", damage)
	}

	if bee.GetHP() != 7 {
		t.Errorf("Expected HP 7 after hit, got %d", bee.GetHP())
	}
}

// TestBeeIsDead tests that a bee correctly reports when it's dead
func TestBeeIsDead(t *testing.T) {
	aliveBee := &Bee{
		beeType: WorkerBee,
		hp:      5,
	}

	if aliveBee.IsDead() {
		t.Error("Bee with positive HP should not be dead")
	}

	// Test bee with exactly 0 HP
	deadBee := &Bee{
		beeType: WorkerBee,
		hp:      0,
	}

	if !deadBee.IsDead() {
		t.Error("Bee with 0 HP should be dead")
	}

	// Test bee with negative HP
	veryDeadBee := &Bee{
		beeType: WorkerBee,
		hp:      -2,
	}

	if !veryDeadBee.IsDead() {
		t.Error("Bee with negative HP should be dead")
	}
}
