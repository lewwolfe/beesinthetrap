package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/lewwolfe/beesinthetrap/internal/config"
)

func main() {
	// load in .env file for config options
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := config.LoadConfig()
	fmt.Println(config.QueenBeeAmount)
}
