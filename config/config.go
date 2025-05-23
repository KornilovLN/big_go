// config/config.go
package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

// GeneratorConfig содержит конфигурационные данные для генератора
type GeneratorConfig struct {
	PostNumber            int `json:"post_number"`
	GenerationIntervalMin int `json:"address_numbers"`
	GenerationIntervalMax int `json:"recipient_numbers"`
}

// LoadGeneratorConfig загружает конфигурацию генератора из файла
func LoadGeneratorConfig(filename string) (*GeneratorConfig, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open config file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := &GeneratorConfig{}
	err = decoder.Decode(config)
	if err != nil {
		return nil, fmt.Errorf("could not decode config JSON: %v", err)
	}

	return config, nil
}

// InitGeneratorConfig загружает и логирует конфигурацию генератора
func InitGeneratorConfig(pathConf string) (*GeneratorConfig, error) {
	generatorConfig, err := LoadGeneratorConfig(pathConf)
	if err != nil {
		return nil, fmt.Errorf("error loading generator config: %v", err)
	}

	log.Printf("Generator config: Post Number: %d, Interval Min: %d, Interval Max: %d",
		generatorConfig.PostNumber, generatorConfig.GenerationIntervalMin, generatorConfig.GenerationIntervalMax)

	return generatorConfig, nil
}

// AppConfig содержит конфигурационные данные приложения
type AppConfig struct {
	ServerPort string `json:"server_port"`
	PageTitle  string `json:"page_title"`
	Database   struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		Name     string `json:"name"`
	} `json:"database"`
}

// LoadConfig загружает конфигурацию из файла
func LoadConfig(filename string) (*AppConfig, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := &AppConfig{}
	err = decoder.Decode(config)
	if err != nil {
		return nil, fmt.Errorf("could not decode config JSON: %v", err)
	}

	return config, nil
}

// ValidateConfig проверяет корректность конфигурации
func ValidateConfig(config *AppConfig) error {
	if config.ServerPort == "" {
		return fmt.Errorf("server port is required")
	}
	if config.PageTitle == "" {
		return fmt.Errorf("page title is required")
	}
	if config.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	// Дополнительные проверки...
	return nil
}

// InitAppConfig загружает, валидирует и логирует конфигурацию приложения
func InitAppConfig(pathConf string) (*AppConfig, string, error) {
	var logBuilder strings.Builder

	appConfig, err := LoadConfig(pathConf)
	if err != nil {
		return nil, "", fmt.Errorf("error loading config: %v", err)
	}

	err = ValidateConfig(appConfig)
	if err != nil {
		return nil, "", fmt.Errorf("invalid configuration: %v", err)
	}

	log.Printf("Настройки обращения к БД:")
	log.Printf("ServerPort: %s", appConfig.ServerPort)
	log.Printf("PageTitle: %s", appConfig.PageTitle)
	log.Printf("Database Host: %s", appConfig.Database.Host)
	log.Printf("Database Port: %d", appConfig.Database.Port)
	log.Printf("Database User: %s", appConfig.Database.User)
	// It is not recommended to log the password
	// log.Printf("Database Password: %s", appConfig.Database.Password)
	// It is not recommended to log the password
	log.Printf("Database Name: %s", appConfig.Database.Name)

	logBuilder.WriteString(fmt.Sprintf("Настройки обращения к БД:\n"))
	logBuilder.WriteString(fmt.Sprintf("ServerPort: %s\n", appConfig.ServerPort))
	logBuilder.WriteString(fmt.Sprintf("PageTitle: %s\n", appConfig.PageTitle))
	logBuilder.WriteString(fmt.Sprintf("Database Host: %s\n", appConfig.Database.Host))
	logBuilder.WriteString(fmt.Sprintf("Database Port: %d\n", appConfig.Database.Port))
	logBuilder.WriteString(fmt.Sprintf("Database User: %s\n", appConfig.Database.User))
	// It is not recommended to log the password
	// logBuilder.WriteString(fmt.Sprintf("Database Password: %s\n", appConfig.Database.Password))
	// It is not recommended to log the password
	logBuilder.WriteString(fmt.Sprintf("Database Name: %s\n", appConfig.Database.Name))

	return appConfig, logBuilder.String(), nil
}
