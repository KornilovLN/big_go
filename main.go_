// main.go
package main

import (
	"big_go/config"
	"big_go/internal/routes"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	var initLogs string

	//--- Initialize application configuration ------------------------------
	appConfig, logs, err := config.InitAppConfig("config_go.json")
	if err != nil {
		log.Fatalf("Error initializing app config: %v", err)
	}

	initLogs = logs

	//--- конфигурации для БД -----------------------------------------------

	// Initialize configurations for each service with the correct path
	config.InitPostgresConfig("../config_postgresql.json")
	config.InitRedisConfig("../config_redis.json")
	config.InitRabbitMQConfig("../config_rabbitmq.json")

	//--- Создаем новый роутер Gin ------------------------------------------

	r := gin.Default()

	// Set up the /init_logs route using the routes package
	routes.SetupInitLogsRoute(r, appConfig.PageTitle, initLogs)

	// Настройка загрузки HTML-шаблонов
	r.LoadHTMLGlob("internal/templates/*")

	// Подключаем маршруты
	routes.SetupRoutes(r, appConfig.PageTitle)

	// Запускаем сервер на порту appConfig.ServerPort==8080
	log.Printf("Server will start at port: %s\n", appConfig.ServerPort)
	r.Run(fmt.Sprintf(":%s", appConfig.ServerPort))
}

// logConfig выводит заполненные поля структуры AppConfig в лог
func logConfig(config *config.AppConfig) {
	log.Printf("ServerPort: %s", config.ServerPort)
	log.Printf("Database Host: %s", config.Database.Host)
	log.Printf("Database Port: %d", config.Database.Port)
	log.Printf("Database User: %s", config.Database.User)
	// Не рекомендуется выводить пароль в лог
	log.Printf("Database Name: %s", config.Database.Name)
}
