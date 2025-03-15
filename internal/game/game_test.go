package game_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/lewwolfe/beesinthetrap/internal/config"
	"github.com/lewwolfe/beesinthetrap/internal/game"
)

func TestNewGame(t *testing.T) {
	// Test with default config
	cfg := &config.Config{
		PlayerHealth:          100,
		PlayerMissChance:      0.2,
		WorkerBeeAmount:       10,
		WorkerBeeHealth:       10,
		WorkerBeeAttackDamage: 5,
		WorkerBeeHitDamage:    8,
		DroneBeeAmount:        5,
		DroneBeeHealth:        20,
		DroneBeeAttackDamage:  7,
		DroneBeeHitDamage:     10,
		QueenBeeAmount:        1,
		QueenBeeHealth:        50,
		QueenBeeAttackDamage:  12,
		QueenBeeHitDamage:     15,
		BeeMissChance:         0.3,
		RandomSeed:            12345,
	}

	ge := game.NewGame(cfg)

	// Verify game engine was created with proper config
	if ge == nil {
		t.Fatalf("NewGame() returned nil")
	}

	// Verify hive size matches config
	expectedBeeCount := cfg.WorkerBeeAmount + cfg.DroneBeeAmount + cfg.QueenBeeAmount
	actualBeeCount := len(ge.GetHive())
	if actualBeeCount != expectedBeeCount {
		t.Errorf("Expected hive size %d, got %d", expectedBeeCount, actualBeeCount)
	}

	// Check custom settings are loaded
	cfg = &config.Config{
		PlayerHealth:     10,
		PlayerMissChance: 0.2,
		WorkerBeeAmount:  1,
		BeeMissChance:    0.3,
		RandomSeed:       12345,
	}

	ge = game.NewGame(cfg)

	// Check custom health works
	if ge.GetPlayer().GetHP() != 10 {
		t.Fatalf("Expected player health %d, got %d", 10, ge.GetPlayer().GetHP())
	}

	// Check custom hive size
	expectedBeeCount = cfg.WorkerBeeAmount + cfg.DroneBeeAmount + cfg.QueenBeeAmount
	actualBeeCount = len(ge.GetHive())

	if actualBeeCount != expectedBeeCount {
		t.Errorf("Expected custom hive size %d, got %d", expectedBeeCount, actualBeeCount)
	}
}

func TestHasGameFinished(t *testing.T) {
	tests := []struct {
		name           string
		playerHealth   int
		hiveSize       int
		expectedResult bool
	}{
		{
			name:           "Game in progress",
			playerHealth:   10,
			hiveSize:       5,
			expectedResult: false,
		},
		{
			name:           "Player dead",
			playerHealth:   0,
			hiveSize:       5,
			expectedResult: true,
		},
		{
			name:           "All bees dead",
			playerHealth:   10,
			hiveSize:       0,
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				PlayerHealth:    tt.playerHealth,
				WorkerBeeAmount: tt.hiveSize,
				WorkerBeeHealth: 1,
				RandomSeed:      12345,
			}

			ge := game.NewGame(cfg)

			if tt.playerHealth == 0 {
				ge.GetPlayer().Sting(100)
			}

			if tt.hiveSize == 0 {
				ge.ClearHive()
			}

			result := ge.IsGameFinished()
			if result != tt.expectedResult {
				t.Errorf("IsGameFinished() = %v, want %v", result, tt.expectedResult)
			}
		})
	}
}

func TestTakePlayerTurn(t *testing.T) {
	cfg := &config.Config{
		PlayerHealth:          100,
		PlayerMissChance:      0,
		WorkerBeeAmount:       1,
		WorkerBeeHealth:       10,
		WorkerBeeAttackDamage: 5,
		WorkerBeeHitDamage:    8,
		BeeMissChance:         0,
		RandomSeed:            42,
	}

	ge := game.NewGame(cfg)

	// Start a goroutine to read from OutputChan
	outputMessages := make([]string, 0)
	done := make(chan bool)
	go func() {
		for {
			select {
			case msg := <-ge.OutputChan:
				outputMessages = append(outputMessages, msg)
			case <-done:
				return
			}
		}
	}()

	ge.TakePlayerTurn()
	done <- true

	// Verify player hit a bee
	if len(outputMessages) == 0 {
		t.Fatalf("No output messages received")
	}

	// Check if the first message indicates a hit
	if !strings.Contains(outputMessages[0], "Direct Hit!") {
		t.Errorf("Expected player to hit, got message: %s", outputMessages[0])
	}

	// Validate damage was dealt to the bee
	if ge.GetHive()[0].GetHP() != 2 { // 10 - 8 = 2
		t.Errorf("Expected bee health to be 2, got %d", ge.GetHive()[0].GetHP())
	}
}

func TestTakeBeeTurn(t *testing.T) {
	cfg := &config.Config{
		PlayerHealth:          100,
		PlayerMissChance:      0,
		WorkerBeeAmount:       1,
		WorkerBeeHealth:       10,
		WorkerBeeAttackDamage: 5,
		WorkerBeeHitDamage:    8,
		BeeMissChance:         0,
		RandomSeed:            42,
	}

	ge := game.NewGame(cfg)

	// Start a goroutine to read from OutputChan
	outputMessages := make([]string, 0)
	done := make(chan bool)
	go func() {
		for {
			select {
			case msg := <-ge.OutputChan:
				outputMessages = append(outputMessages, msg)
			case <-done:
				return
			}
		}
	}()

	// Take a bee turn
	ge.TakeBeeTurn()
	done <- true

	// Verify bee stung the player
	if len(outputMessages) == 0 {
		t.Fatalf("No output messages received")
	}

	// Check if the message indicates a sting
	if !strings.Contains(outputMessages[0], "Ouch!") {
		t.Errorf("Expected bee to sting, got message: %s", outputMessages[0])
	}

	// Validate damage was dealt to the player
	if ge.GetPlayer().GetHP() != 95 { // 100 - 5 = 95
		t.Errorf("Expected player health to be 95, got %d", ge.GetPlayer().GetHP())
	}
}

func TestQueenBeeKill(t *testing.T) {
	cfg := &config.Config{
		PlayerHealth:         100,
		PlayerMissChance:     0, // Ensure player never misses
		QueenBeeAmount:       1,
		QueenBeeHealth:       1, // One hit will kill
		QueenBeeAttackDamage: 10,
		QueenBeeHitDamage:    10,
		BeeMissChance:        0,
		RandomSeed:           12345,
	}

	ge := game.NewGame(cfg)

	// Start a goroutine to read from OutputChan
	outputMessages := make([]string, 0)
	done := make(chan bool)
	go func() {
		for {
			select {
			case msg := <-ge.OutputChan:
				outputMessages = append(outputMessages, msg)
			case <-done:
				return
			}
		}
	}()

	// Take a player turn which should kill the queen
	ge.TakePlayerTurn()

	if !ge.IsGameFinished() {
		t.Errorf("Game should have ended when Queen Bee died")
	}

	// Stop reading from OutputChan
	done <- true

	// Verify queen was killed and hive collapsed
	hiveSize := len(ge.GetHive())
	if hiveSize != 0 {
		t.Errorf("Expected hive to collapse (size 0), got size %d", hiveSize)
	}

	// Check for queen death message
	queenDeathMessageFound := false
	for _, msg := range outputMessages {
		if strings.Contains(msg, "The Queen Bee is dead, and the entire hive collapses!") {
			queenDeathMessageFound = true
			break
		}
	}

	if !queenDeathMessageFound {
		t.Errorf("Queen death message not found in output")
	}
}

func TestGameLoop(t *testing.T) {
	cfg := &config.Config{
		PlayerHealth:          10,
		PlayerMissChance:      0,
		WorkerBeeAmount:       1,
		WorkerBeeHealth:       1,
		WorkerBeeAttackDamage: 5,
		WorkerBeeHitDamage:    1,
		BeeMissChance:         0,
		RandomSeed:            12345,
	}

	ge := game.NewGame(cfg)

	// Run the game in auto mode with a cancellable context
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	go func() {
		for output := range ge.OutputChan {
			fmt.Println(output)
		}
	}()

	done := make(chan bool)
	go func() {
		ge.Start(true, ctx)
		close(done)
	}()

	// Wait for game to finish or timeout
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatalf("Game did not finish within expected time")
	}

	// Verify game is finished
	if !ge.IsGameFinished() {
		t.Errorf("Expected game to be finished")
	}
}
