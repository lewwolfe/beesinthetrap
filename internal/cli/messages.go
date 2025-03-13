package cli

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func (c *GameCLI) displayWelcomeBanner() {
	fmt.Println("===================================================")
	fmt.Println("Welcome to Bees in the Trap!")
	fmt.Println("===================================================")
}

func (c *GameCLI) displayGameInterface() {
	fmt.Println("===================================================")
	fmt.Printf("Player: %s\n", c.playerName)
	fmt.Printf("Health: %d/%d\n", c.gameEngine.GetPlayer().GetHP(), c.gameEngine.Config.PlayerHealth)
	fmt.Printf("Bees remaining: %d\n", len(c.gameEngine.GetHive()))
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

func (c *GameCLI) displayGameOver() {
	c.clearScreen()
	fmt.Println("===================================================")
	fmt.Println("                     GAME OVER                    ")
	fmt.Println("===================================================")

	player := c.gameEngine.GetPlayer()
	if player.IsDead() {
		fmt.Printf("Sorry %s, you were defeated by the hive!\n", c.playerName)
	} else if len(c.gameEngine.GetHive()) == 0 {
		fmt.Printf("Congratulations %s! You defeated the hive!\n", c.playerName)
	}

	fmt.Printf("\nFinal Stats for %s:\n", c.playerName)
	fmt.Printf("Health remaining: %d/%d\n", player.GetHP(), c.gameEngine.Config.PlayerHealth)
	fmt.Printf("Bees remaining: %d\n\n", len(c.gameEngine.GetHive()))
	fmt.Printf("Bee Stings: %d\n", c.gameEngine.BeeStings)
	fmt.Printf("Player Hits: %d\n", c.gameEngine.PlayerHits)

	if len(c.gameEngine.GetHive()) > 0 {
		beeTypes := make(map[string]int)
		for _, bee := range c.gameEngine.GetHive() {
			beeTypes[bee.GetBeeType()]++
		}

		fmt.Println("\nRemaining bees:")
		for beeType, count := range beeTypes {
			fmt.Printf("- %s: %d\n", beeType, count)
		}
	}

	fmt.Println("\n===================================================")
	fmt.Println("Thanks for playing!")
	fmt.Println("Press Enter to exit...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func (c *GameCLI) clearScreen() {
	fmt.Print("\033[H\033[2J")
}
