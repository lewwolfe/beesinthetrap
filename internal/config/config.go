package config

import (
	"os"
	"strconv"
)

type Config struct {
	PlayerHealth          int
	LogSize               int
	AutoRunSpeed          int
	RandomSeed            int64
	PlayerMissChance      float64
	BeeMissChance         float64
	QueenBeeAmount        int
	QueenBeeHealth        int
	QueenBeeAttackDamage  int
	QueenBeeHitDamage     int
	WorkerBeeAmount       int
	WorkerBeeHealth       int
	WorkerBeeAttackDamage int
	WorkerBeeHitDamage    int
	DroneBeeAmount        int
	DroneBeeHealth        int
	DroneBeeAttackDamage  int
	DroneBeeHitDamage     int
}

func LoadConfig() *Config {
	config := new(Config)

	// Player
	config.PlayerHealth = getEnvAsInt("PLAYER_HEALTH", 100)

	//Game Options
	config.RandomSeed = int64(getEnvAsInt("RANDOM_SEED", 0))
	config.PlayerMissChance = getEnvAsFloat("PLAYER_MISS_CHANCE", 0.1)
	config.BeeMissChance = getEnvAsFloat("BEE_MISS_CHANCE", 0.2)
	config.LogSize = getEnvAsInt("LOG_SIZE", 10)
	config.AutoRunSpeed = getEnvAsInt("AUTO_RUN_SPEED", 1)

	// Queen Bee
	config.QueenBeeAmount = getEnvAsInt("QUEEN_BEE_AMOUNT", 1)
	config.QueenBeeHealth = getEnvAsInt("QUEEN_BEE_HEALTH", 100)
	config.QueenBeeAttackDamage = getEnvAsInt("QUEEN_BEE_ATTACK_DAMAGE", 10)
	config.QueenBeeHitDamage = getEnvAsInt("QUEEN_BEE_DEFENSE_DAMAGE", 10)

	// Worker Bee
	config.WorkerBeeAmount = getEnvAsInt("WORKER_BEE_AMOUNT", 5)
	config.WorkerBeeHealth = getEnvAsInt("WORKER_BEE_HEALTH", 75)
	config.WorkerBeeAttackDamage = getEnvAsInt("WORKER_BEE_ATTACK_DAMAGE", 5)
	config.WorkerBeeHitDamage = getEnvAsInt("WORKER_BEE_DEFENSE_DAMAGE", 25)

	// Drone Bee
	config.DroneBeeAmount = getEnvAsInt("DRONE_BEE_AMOUNT", 25)
	config.DroneBeeHealth = getEnvAsInt("DRONE_BEE_HEALTH", 60)
	config.DroneBeeAttackDamage = getEnvAsInt("DRONE_BEE_ATTACK_DAMAGE", 1)
	config.DroneBeeHitDamage = getEnvAsInt("DRONE_BEE_DEFENSE_DAMAGE", 30)

	return config
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
