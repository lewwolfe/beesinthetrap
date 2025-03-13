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
		if c.scanner.Scan() {
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
	}
	fmt.Println()
}

func (c *GameCLI) runGame() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	// Handle Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nGame interrupted! Shutting down...")
		cancel()
	}()

	// Clear the screen and display game interface
	c.clearScreen()
	c.displayGameInterface()

	// Start game output reader goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.monitorGameOutput(ctx)
	}()

	// If not in auto mode, start the input handler goroutine
	if !c.autoMode {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.handleUserInput(ctx)
		}()
	}

	// Start the game engine in its own goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.gameEngine.Start(c.autoMode, ctx)
	}()

	// Monitor game state in the main thread
	gameOver := make(chan bool)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if c.gameEngine.HasGameFinished() {
					gameOver <- true
					return
				}
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	// Wait for game to finish or be cancelled
	select {
	case <-ctx.Done():
		wg.Wait()
	case <-gameOver:
		time.Sleep(500 * time.Millisecond) // Wait for final messages
		c.displayGameOver()
	}
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
	scanner := bufio.NewScanner(os.Stdin)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if scanner.Scan() {
				input := strings.TrimSpace(scanner.Text())
				select {
				case <-ctx.Done():
					return
				case c.gameEngine.InputChan <- input:
				}
			} else {
				if err := scanner.Err(); err != nil {
					fmt.Printf("Error reading input: %v\n", err)
				}
				return
			}
		}
	}
}
