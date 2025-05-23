package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// RabbitMQConfig contains configuration data for connecting to RabbitMQ
type RabbitMQConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	VHost    string `json:"vhost"`
}

// LoadRabbitMQConfig loads the RabbitMQ configuration from a JSON file
func LoadRabbitMQConfig(filename string) (*RabbitMQConfig, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open config file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := &RabbitMQConfig{}
	err = decoder.Decode(config)
	if err != nil {
		return nil, fmt.Errorf("could not decode config JSON: %v", err)
	}

	return config, nil
}

func InitRabbitMQConfig(pathConf string) {
	rabbitmqConfig, err := LoadRabbitMQConfig(pathConf)
	if err != nil {
		log.Fatalf("Error loading RabbitMQ config: %v", err)
	}
	log.Printf("RabbitMQ config: %+v", rabbitmqConfig)
}
