package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// RedisConfig contains configuration data for connecting to Redis
type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

// LoadRedisConfig loads the Redis configuration from a JSON file
func LoadRedisConfig(filename string) (*RedisConfig, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open config file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := &RedisConfig{}
	err = decoder.Decode(config)
	if err != nil {
		return nil, fmt.Errorf("could not decode config JSON: %v", err)
	}

	return config, nil
}

func InitRedisConfig(pathConf string) {
	redisConfig, err := LoadRedisConfig(pathConf)
	if err != nil {
		log.Fatalf("Error loading Redis config: %v", err)
	}
	log.Printf("Redis config: %+v", redisConfig)
}
