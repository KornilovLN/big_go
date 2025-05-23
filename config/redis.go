package config

import (
    "encoding/json"
    "os"
)

// RedisConfig содержит конфигурацию для подключения к Redis
type RedisConfig struct {
    Host     string \`json:"host"\`
    Port     int    \`json:"port"\`
    Password string \`json:"password"\`
    DB       int    \`json:"db"\`
}

// LoadRedisConfig загружает конфигурацию Redis из файла
func LoadRedisConfig(configPath string) (*RedisConfig, error) {
    file, err := os.Open(configPath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var config RedisConfig
    decoder := json.NewDecoder(file)
    err = decoder.Decode(&config)
    if err != nil {
        return nil, err
    }

    // Переопределение из переменных окружения, если они установлены
    if host := os.Getenv("REDIS_HOST"); host != "" {
        config.Host = host
    }
    if port := os.Getenv("REDIS_PORT"); port != "" {
        var portInt int
        if _, err := fmt.Sscanf(port, "%d", &portInt); err == nil {
            config.Port = portInt
        }
    }
    if password := os.Getenv("REDIS_PASSWORD"); password != "" {
        config.Password = password
    }
    if db := os.Getenv("REDIS_DB"); db != "" {
        var dbInt int
        if _, err := fmt.Sscanf(db, "%d", &dbInt); err == nil {
            config.DB = dbInt
        }
    }

    return &config, nil
}