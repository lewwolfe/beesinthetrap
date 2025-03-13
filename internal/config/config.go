package config

import (
	"os"
	"strconv"
)

type Config struct {
	PlayerHealth           int
	PlayerMissChance       float64
	BeeMissChance          float64
	QueenBeeAmount         int
	QueenBeeHealth         int
	QueenBeeAttackDamage   int
	QueenBeeDefenseDamage  int
	WorkerBeeAmount        int
	WorkerBeeHealth        int
	WorkerBeeAttackDamage  int
	WorkerBeeDefenseDamage int
	DroneBeeAmount         int
	DroneBeeHealth         int
	DroneBeeAttackDamage   int
	DroneBeeDefenseDamage  int
}

func LoadConfig() *Config {
	config := new(Config)

	// Player
	config.PlayerHealth = getEnvAsInt("PLAYER_HEALTH", 100)
	config.PlayerMissChance = getEnvAsFloat("PLAYER_MISS_CHANCE", 0.1)
	config.BeeMissChance = getEnvAsFloat("BEE_MISS_CHANCE", 0.2)

	// Queen Bee
	config.QueenBeeAmount = getEnvAsInt("QUEEN_BEE_AMOUNT", 1)
	config.QueenBeeHealth = getEnvAsInt("QUEEN_BEE_HEALTH", 100)
	config.QueenBeeAttackDamage = getEnvAsInt("QUEEN_BEE_ATTACK_DAMAGE", 10)
	config.QueenBeeDefenseDamage = getEnvAsInt("QUEEN_BEE_DEFENSE_DAMAGE", 10)

	// Worker Bee
	config.WorkerBeeAmount = getEnvAsInt("WORKER_BEE_AMOUNT", 5)
	config.WorkerBeeHealth = getEnvAsInt("WORKER_BEE_HEALTH", 75)
	config.WorkerBeeAttackDamage = getEnvAsInt("WORKER_BEE_ATTACK_DAMAGE", 5)
	config.WorkerBeeDefenseDamage = getEnvAsInt("WORKER_BEE_DEFENSE_DAMAGE", 25)

	// Drone Bee
	config.DroneBeeAmount = getEnvAsInt("DRONE_BEE_AMOUNT", 25)
	config.DroneBeeHealth = getEnvAsInt("DRONE_BEE_HEALTH", 60)
	config.DroneBeeAttackDamage = getEnvAsInt("DRONE_BEE_ATTACK_DAMAGE", 1)
	config.DroneBeeDefenseDamage = getEnvAsInt("DRONE_BEE_DEFENSE_DAMAGE", 30)

	return config
}

func getEnvAsString(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultVal
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultVal
	}
	return intValue

}

func getEnvAsFloat(key string, defaultVal float64) float64 {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultVal
	}

	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return defaultVal
	}
	return floatValue

}
