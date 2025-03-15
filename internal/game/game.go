package game

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/lewwolfe/beesinthetrap/internal/config"
)

type GameState int

const (
	Running GameState = iota
	PlayerWin
	PlayerLose
)

type GameEngine struct {
	Config        *config.Config
	player        *Player
	hive          []*Bee
	playerTurn    bool
	PlayerHits    int
	BeeStings     int
	InputChan     chan string
	OutputChan    chan string
	GameStateChan chan GameState
	rng           *rand.Rand
}

func NewGame(cfg *config.Config) *GameEngine {
	//Input a random seed for randomness, this allows for repetable games for testing
	seed := cfg.RandomSeed
	if seed == 0 {
		seed = time.Now().Unix()
	}
	source := rand.NewSource(seed)
	rng := rand.New(source)

	ge := &GameEngine{
		Config:        cfg,
		player:        &Player{hp: cfg.PlayerHealth, missChance: cfg.PlayerMissChance},
		playerTurn:    true,
		PlayerHits:    0,
		BeeStings:     0,
		InputChan:     make(chan string),
		OutputChan:    make(chan string),
		GameStateChan: make(chan GameState, 1),
		rng:           rng,
	}

	// Spawn worker bees
	for i := 0; i < cfg.WorkerBeeAmount; i++ {
		ge.hive = append(ge.hive, &Bee{
			beeType:      WorkerBee,
			hp:           cfg.WorkerBeeHealth,
			attackDamage: cfg.WorkerBeeAttackDamage,
			hitDamage:    cfg.WorkerBeeHitDamage,
			missChance:   cfg.BeeMissChance,
		})
	}

	// Spawn drone bees
	for i := 0; i < cfg.DroneBeeAmount; i++ {
		ge.hive = append(ge.hive, &Bee{
			beeType:      DroneBee,
			hp:           cfg.DroneBeeHealth,
			attackDamage: cfg.DroneBeeAttackDamage,
			hitDamage:    cfg.DroneBeeHitDamage,
			missChance:   cfg.BeeMissChance,
		})
	}

	// Spawn Queen bee(s)
	for i := 0; i < cfg.QueenBeeAmount; i++ {
		ge.hive = append(ge.hive, &Bee{
			beeType:      QueenBee,
			hp:           cfg.QueenBeeHealth,
			attackDamage: cfg.QueenBeeAttackDamage,
			hitDamage:    cfg.QueenBeeHitDamage,
			missChance:   cfg.BeeMissChance,
		})
	}

	return ge
}

func (ge *GameEngine) IsPlayerTurn() bool {
	return ge.playerTurn
}

func (ge *GameEngine) GetHive() []*Bee {
	return ge.hive
}

func (ge *GameEngine) GetPlayer() *Player {
	return ge.player
}

func (ge *GameEngine) ClearHive() {
	ge.hive = []*Bee{}
}

// Start runs the game loop, handling turns and input
func (ge *GameEngine) Start(auto bool, ctx context.Context) {
	defer ctx.Done()
	// Main game loop
	for !ge.IsGameFinished() {
		select {
		case <-ctx.Done():
			return
		default:
			// Existing game logic
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

	// Send final game state messages
	if ge.player.IsDead() {
		ge.OutputChan <- "💀 You have been defeated by the hive!"
		ge.GameStateChan <- PlayerLose
	} else {
		ge.OutputChan <- "🏆 Congratulations! You've destroyed the entire hive!"
		ge.GameStateChan <- PlayerWin
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

func (ge *GameEngine) IsGameFinished() bool {
	if ge.player.IsDead() || len(ge.hive) == 0 {

		return true
	}
	return false
}

// Wait for input chan to recieve a player input
func (ge *GameEngine) TakePlayerTurn() {
	// Let the player Attack() to see if they miss
	if !ge.player.Attack(ge.rng) {
		ge.OutputChan <- "❌ Miss! You just missed the hive, better luck next time!"
		return
	}

	//Select a random bee from the hive and damage it
	ge.PlayerHits++
	beePos := ge.rng.Intn(len(ge.hive))
	bee := ge.hive[beePos]

	// Deal damage to the bee
	beeDamage := bee.Hit()
	ge.OutputChan <- fmt.Sprintf("🧑 Direct Hit! You dealt %d damage to a %s Bee.", beeDamage, bee.beeType)

	// Check if bee is dead and which type of bee to update the hive
	if bee.IsDead() && bee.beeType == QueenBee {
		ge.OutputChan <- "🎉 The Queen Bee is dead, and the entire hive collapses!"
		ge.ClearHive()
	} else if bee.IsDead() {
		ge.OutputChan <- fmt.Sprintf("💀 You killed a %s!", bee.beeType)

		// Swap with the last element and shrink the slice
		ge.hive[beePos] = ge.hive[len(ge.hive)-1]
		ge.hive = ge.hive[:len(ge.hive)-1]
	}
}

func (ge *GameEngine) TakeBeeTurn() {
	// Select random bee from the hive
	beePos := ge.rng.Intn(len(ge.hive))
	bee := ge.hive[beePos]

	// Let the bee Attack() to get damage
	damage := bee.Attack(ge.rng)
	if damage == 0 {
		ge.OutputChan <- fmt.Sprintf("❌ Buzz! That was close! The %s Bee just missed you!", bee.beeType)
		return
	}

	// Run Sting() on player
	ge.player.Sting(damage)
	ge.BeeStings++

	// Send out response from game to cli
	ge.OutputChan <- fmt.Sprintf("🐝 Ouch! A %s Bee stung you for %d damage!", bee.beeType, damage)
}
