package game

import (
	"context"
	"fmt"
	"math/rand/v2"

	"github.com/lewwolfe/beesinthetrap/internal/config"
)

type GameEngine struct {
	config     *config.Config
	player     *Player
	hive       []Bee
	playerTurn bool
	playerHits int
	beeStings  int
	InputChan  chan string
	OutputChan chan string
}

func NewGame(cfg *config.Config) *GameEngine {
	ge := &GameEngine{
		config:     cfg,
		player:     &Player{hp: cfg.PlayerHealth, missChance: cfg.PlayerMissChance},
		playerTurn: true,
		playerHits: 0,
		beeStings:  0,
		InputChan:  make(chan string),
		OutputChan: make(chan string),
	}

	// Spawn worker bees
	for i := 0; i < cfg.WorkerBeeAmount; i++ {
		ge.hive = append(ge.hive, Bee{
			beeType:      "Worker",
			hp:           cfg.WorkerBeeHealth,
			attackDamage: cfg.WorkerBeeAttackDamage,
			hitDamage:    cfg.WorkerBeeHitDamage,
			missChance:   cfg.BeeMissChance,
		})
	}

	// Spawn drone bees
	for i := 0; i < cfg.DroneBeeAmount; i++ {
		ge.hive = append(ge.hive, Bee{
			beeType:      "Drone",
			hp:           cfg.DroneBeeHealth,
			attackDamage: cfg.DroneBeeAttackDamage,
			hitDamage:    cfg.DroneBeeHitDamage,
			missChance:   cfg.BeeMissChance,
		})
	}

	// Spawn Queen bee(s)
	for i := 0; i < cfg.QueenBeeAmount; i++ {
		ge.hive = append(ge.hive, Bee{
			beeType:      "Queen",
			hp:           cfg.QueenBeeHealth,
			attackDamage: cfg.QueenBeeAttackDamage,
			hitDamage:    cfg.QueenBeeHitDamage,
			missChance:   cfg.BeeMissChance,
		})
	}

	return ge
}

// Start runs the game loop, handling turns and input
func (ge *GameEngine) Start(auto bool, ctx context.Context) {
	for !ge.HasGameFinished() {
		if ge.playerTurn {
			if !auto {
				ge.waitForPlayerAction()
			}
			ge.TakePlayerTurn()
		} else {
			ge.TakeBeeTurn()
		}

		// Toggle turn
		ge.playerTurn = !ge.playerTurn
	}
}

func (ge *GameEngine) waitForPlayerAction() {
	for {
		ge.OutputChan <- "Type 'hit' to attack..."

		// Wait for valid input
		input := <-ge.InputChan
		if input == "hit" {
			return
		}

		ge.OutputChan <- fmt.Sprintf("Invalid command! '%s'", input)
	}
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

// Wait for input chan to recieve a player input
func (ge *GameEngine) TakePlayerTurn() {

	// Let the player Attack() to see if they miss
	if !ge.player.Attack() {
		ge.OutputChan <- "Miss! You just missed the hive, better luck next time!"
		return
	}

	//Select a random bee from the hive and damage it
	ge.playerHits++
	beePos := rand.IntN(len(ge.hive))
	bee := ge.hive[beePos]

	// Deal damage to the bee
	beeDamage := bee.Hit()
	ge.OutputChan <- fmt.Sprintf("Direct Hit! You dealt %d damage to a %s Bee.", beeDamage, bee.beeType)

	// Check if bee is dead and which type of bee to update the hive
	if bee.IsDead() && bee.beeType == "Queen" {
		ge.OutputChan <- "The Queen Bee is dead, and the entire hive collapses!"
		ge.hive = []Bee{}
	} else if bee.IsDead() {
		ge.OutputChan <- fmt.Sprintf("You killed a %s!", bee.beeType)
		ge.hive = append(ge.hive[:beePos], ge.hive[beePos+1:]...)
	}
}

func (ge *GameEngine) TakeBeeTurn() {
	// Select random bee from the hive
	beePos := rand.IntN(len(ge.hive))
	bee := ge.hive[beePos]

	// Let the bee Attack() to get damage
	damage := bee.Attack()
	if damage == 0 {
		ge.OutputChan <- fmt.Sprintf("Buzz! That was close! The %s Bee just missed you!", bee.beeType)
		return
	}

	// Run Sting() on player
	ge.player.Sting(damage)

	// Send out response from game to cli
	ge.OutputChan <- fmt.Sprintf("Ouch! A %s Bee stung you for %d damage!", bee.beeType, damage)
}
