// config/postgresql.go
package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func LoadPostgresConfig(filename string) (*PostgresConfig, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open config file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := &PostgresConfig{}
	err = decoder.Decode(config)
	if err != nil {
		return nil, fmt.Errorf("could not decode config JSON: %v", err)
	}

	return config, nil
}

func InitPostgresConfig(pathConf string) {
	postgresConfig, err := LoadPostgresConfig(pathConf)
	if err != nil {
		log.Fatalf("Error loading PostgreSQL config: %v", err)
	}
	log.Printf("PostgreSQL config: %+v", postgresConfig)
}
