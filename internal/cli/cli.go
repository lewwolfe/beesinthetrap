package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/lewwolfe/beesinthetrap/internal/game"
)

type GameCLI struct {
	gameEngine *game.GameEngine
	playerName string
	autoMode   bool
	scanner    *bufio.Scanner
	gameLogs   []string
}

func NewGameCLI(gameEngine *game.GameEngine) *GameCLI {
	return &GameCLI{
		gameEngine: gameEngine,
		scanner:    bufio.NewScanner(os.Stdin),
	}
}

func (c *GameCLI) Start() {
	c.displayWelcomeBanner()
	c.promptPlayerName()
	c.promptAutoMode()
	c.runGame()
}

func (c *GameCLI) promptPlayerName() {
	fmt.Print("Enter your name, brave bee hunter: ")
	if c.scanner.Scan() {
		c.playerName = strings.TrimSpace(c.scanner.Text())
	}
	if c.playerName == "" {
		c.playerName = "Anonymous Hunter"
	}
	fmt.Printf("Welcome, %s!\n\n", c.playerName)
}

func (c *GameCLI) promptAutoMode() {
	for {
		fmt.Print("Do you want the game to run automatically? (y/n): ")
		if !c.scanner.Scan() {
			break
		}

		input := strings.ToLower(strings.TrimSpace(c.scanner.Text()))
		if input == "y" || input == "yes" {
			c.autoMode = true
			fmt.Println("Auto mode activated. Sit back and watch the bees battle!")
			break
		} else if input == "n" || input == "no" {
			c.autoMode = false
			fmt.Println("Manual mode activated. You'll need to type 'hit' to attack.")
			break
		}
		fmt.Println("Please enter 'y' or 'n'.")
	}
	fmt.Println()
}

func (c *GameCLI) runGame() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling
	c.setupSignalHandling(ctx, cancel)

	// Clear the screen and display game interface
	c.clearScreen()
	c.displayGameInterface()

	// Start all goroutines
	var wg sync.WaitGroup
	c.startGameRoutines(ctx, &wg)

	// Wait for game state event
	select {
	case <-ctx.Done():
		return
	case gameState := <-c.gameEngine.GameStateChan:
		time.Sleep(200 * time.Millisecond)
		c.displayGameOver(gameState)
		ctx.Done()
	}

	wg.Wait()
}

func (c *GameCLI) setupSignalHandling(ctx context.Context, cancel context.CancelFunc) {
	defer ctx.Done()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nGame interrupted! Shutting down...")
		cancel()
	}()
}

func (c *GameCLI) startGameRoutines(ctx context.Context, wg *sync.WaitGroup) {
	// Start game output reader
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.monitorGameOutput(ctx)
	}()

	// Start input handler if not in auto mode
	if !c.autoMode {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.handleUserInput(ctx)
		}()
	}

	// Start game engine
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.gameEngine.Start(c.autoMode, ctx)
	}()
}

func (c *GameCLI) monitorGameOutput(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-c.gameEngine.OutputChan:
			c.gameLogs = append(c.gameLogs, msg)
			c.displayMessage()
		}
	}
}

func (c *GameCLI) handleUserInput(ctx context.Context) {
	for c.scanner.Scan() {
		input := strings.TrimSpace(c.scanner.Text())
		select {
		case <-ctx.Done():
			return
		case c.gameEngine.InputChan <- input:
		}
	}

	if err := c.scanner.Err(); err != nil {
		fmt.Printf("Error reading input: %v\n", err)
	}
}
