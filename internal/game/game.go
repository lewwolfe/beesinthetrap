package game

import (
	"sync"

	"github.com/lewwolfe/beesinthetrap/internal/config"
)

type GameEngine struct {
	config     *config.Config
	player     *Player
	hive       []Bee
	InputChan  chan string
	OutputChan chan string
	wg         sync.WaitGroup
}

func NewGame(cfg *config.Config) *GameEngine {
	ge := &GameEngine{
		config:     cfg,
		player:     &Player{hp: cfg.PlayerHealth},
		InputChan:  make(chan string),
		OutputChan: make(chan string),
	}

	// Spawn worker bees
	for i := 0; i < cfg.WorkerBeeAmount; i++ {
		ge.hive = append(ge.hive, Bee{
			Type:          "Worker",
			HP:            cfg.WorkerBeeHealth,
			AttackDamage:  cfg.WorkerBeeAttackDamage,
			DefenceDamage: cfg.WorkerBeeDefenseDamage,
			MissChance:    cfg.BeeMissChance,
		})
	}

	// Spawn drone bees
	for i := 0; i < cfg.DroneBeeAmount; i++ {
		ge.hive = append(ge.hive, Bee{
			Type:          "Drone",
			HP:            cfg.DroneBeeHealth,
			AttackDamage:  cfg.DroneBeeAttackDamage,
			DefenceDamage: cfg.DroneBeeDefenseDamage,
			MissChance:    cfg.BeeMissChance,
		})
	}

	// Spawn Queen bee(s)
	for i := 0; i < cfg.QueenBeeAmount; i++ {
		ge.hive = append(ge.hive, Bee{
			Type:          "Queen",
			HP:            cfg.QueenBeeHealth,
			AttackDamage:  cfg.QueenBeeAttackDamage,
			DefenceDamage: cfg.QueenBeeDefenseDamage,
			MissChance:    cfg.BeeMissChance,
		})
	}

	return ge
}

func (ge *GameEngine) HasGameFinished() bool {
	if ge.player.IsDead() {
		// Bee win logic
		return true
	} else if len(ge.hive) == 0 {
		// Player win logic
		return true
	}

	return false
}

func (ge *GameEngine) PlayerTurn() {
	// TODO: Implement player attacking a bee
	// Wait for input chan to recieve a player input
	// Let the player Attack() to get damage
	// Run Damage() on bee
	// Check if bee is dead, if so remove from hive
	// Send out response from game to cli
}

func (ge *GameEngine) BeeTurn() {
	// TODO: Implement bees attacking the player
	// Let the bee Attack() to get damage
	// Run Damage() on player
	// Send out response from game to cli
}
