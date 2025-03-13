package game

import (
	"math/rand"
	"testing"
)

// TestPlayerAttack tests basic attack functionality of a Player
func TestPlayerAttack(t *testing.T) {
	source := rand.NewSource(42)
	rng := rand.New(source)

	// Test player always hits (missChance = 0)
	alwaysHitPlayer := &Player{
		hp:         20,
		missChance: 0,
	}

	hit := alwaysHitPlayer.Attack(rng)
	if !hit {
		t.Error("Player with missChance 0 should always hit")
	}

	alwaysMissPlayer := &Player{
		hp:         20,
		missChance: 1,
	}

	hit = alwaysMissPlayer.Attack(rng)
	if hit {
		t.Error("Player with missChance 1 should always miss")
	}
}

// TestPlayerSting tests that a player takes the right amount of damage
func TestPlayerSting(t *testing.T) {
	player := &Player{
		hp: 20,
	}

	player.Sting(5)
	if player.GetHP() != 15 {
		t.Errorf("Expected HP 15 after sting, got %d", player.GetHP())
	}

	player.Sting(7)
	player.Sting(3)
	if player.GetHP() != 5 {
		t.Errorf("Expected HP 5 after multiple stings, got %d", player.GetHP())
	}
}

// TestPlayerIsDead tests that a player correctly reports when it's dead
func TestPlayerIsDead(t *testing.T) {
	alivePlayer := &Player{
		hp: 10,
	}

	if alivePlayer.IsDead() {
		t.Error("Player with positive HP should not be dead")
	}

	// Test player with exactly 0 HP
	deadPlayer := &Player{
		hp: 0,
	}

	if !deadPlayer.IsDead() {
		t.Error("Player with 0 HP should be dead")
	}

	// Test player with negative HP
	veryDeadPlayer := &Player{
		hp: -5,
	}

	if !veryDeadPlayer.IsDead() {
		t.Error("Player with negative HP should be dead")
	}
}
