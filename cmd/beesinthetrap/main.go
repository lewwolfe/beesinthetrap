package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/lewwolfe/beesinthetrap/internal/cli"
	"github.com/lewwolfe/beesinthetrap/internal/config"
	"github.com/lewwolfe/beesinthetrap/internal/game"
)

func main() {
	// load in .env file for config options
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config.LoadConfig()
	game := game.NewGame(cfg)

	gameCLI := cli.NewGameCLI(game)
	gameCLI.Start()
}
