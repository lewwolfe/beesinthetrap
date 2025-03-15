package cli

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/lewwolfe/beesinthetrap/internal/game"
)

func (c *GameCLI) displayWelcomeBanner() {
	fmt.Println("===================================================")
	fmt.Println("Welcome to Bees in the Trap!")
	fmt.Println("===================================================")
}

func (c *GameCLI) displayGameInterface() {
	fmt.Println("===================================================")
	fmt.Printf("Player: %s\n", c.playerName)
	fmt.Printf("Health: %d/%d\n\n", c.gameEngine.GetPlayer().GetHP(), c.gameEngine.Config.PlayerHealth)
	c.printRemainingBee()
	fmt.Println("===================================================")
	fmt.Println("GAME LOG:")
}

func (c *GameCLI) displayMessage() {
	c.clearScreen()
	c.displayGameInterface()

	// Filter log view to only show config amount
	viewLogs := c.gameLogs
	if len(c.gameLogs) > c.gameEngine.Config.LogSize {
		viewLogs = c.gameLogs[len(c.gameLogs)-c.gameEngine.Config.LogSize:]
	}

	for _, msg := range viewLogs {
		fmt.Println(msg)
	}

	time.Sleep(time.Duration(c.gameEngine.Config.AutoRunSpeed) * time.Second)
}

func (c *GameCLI) displayGameOver(state game.GameState) {
	c.clearScreen()
	fmt.Println("===================================================")
	fmt.Println("                     GAME OVER                    ")
	fmt.Println("===================================================")

	if state == game.PlayerLose {
		fmt.Printf("Sorry %s, you were defeated by the hive!\n", c.playerName)
	} else if state == game.PlayerWin {
		fmt.Printf("Congratulations %s! You defeated the hive!\n", c.playerName)
	}

	fmt.Printf("\nFinal Stats for %s:\n", c.playerName)
	fmt.Printf("Health remaining: %d/%d\n\n", c.gameEngine.GetPlayer().GetHP(), c.gameEngine.Config.PlayerHealth)
	fmt.Printf("Bee Stings: %d\n", c.gameEngine.BeeStings)
	fmt.Printf("Player Hits: %d\n\n", c.gameEngine.PlayerHits)

	if len(c.gameEngine.GetHive()) > 0 {
		c.printRemainingBee()
	}

	fmt.Println("\n===================================================")
	fmt.Println("Thanks for playing!")
	fmt.Println("Press Enter to exit...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func (c *GameCLI) printRemainingBee() {
	fmt.Println("Bees remaining:")

	// Define the fixed order of bee types
	beeOrder := []string{"Queen", "Worker", "Drone"}
	beeCount := map[string]int{}
	beeHPs := map[string][]int{}

	// Count bees and track HPs
	for _, bee := range c.gameEngine.GetHive() {
		beeType := bee.GetBeeType().String()
		beeCount[beeType]++
		beeHPs[beeType] = append(beeHPs[beeType], bee.GetHP())
	}

	// Print bees in the fixed order
	for _, beeType := range beeOrder {
		if count, exists := beeCount[beeType]; exists {
			fmt.Printf("%s: %d [", beeType, count)
			for i, hp := range beeHPs[beeType] {
				if i > 0 {
					fmt.Print(", ")
				}
				fmt.Printf("%d", hp)
			}
			fmt.Println("]")
		}
	}
}

func (c *GameCLI) clearScreen() {
	fmt.Print("\033[H\033[2J")
}
