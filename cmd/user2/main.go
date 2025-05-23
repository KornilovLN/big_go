package main

import (
	"big_go/internal/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Канал для хранения последних полученных данных
	var latestData []models.SensorData

	// Обработчик для получения данных от коллектора
	r.POST("/data", func(c *gin.Context) {
		var data models.SensorData
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		log.Printf("User2 получил данные: %+v", data)
		// Подробное логирование полученных данных
		log.Printf("User1 получил данные:")
		log.Printf("  Метаданные:")
		log.Printf("    Получатель: %s", data.Meta.Recipient)
		log.Printf("    ID поста: %d", data.Meta.PostID)
		log.Printf("    Адрес: %d", data.Meta.Address)
		log.Printf("    Временная метка: %s", data.Meta.Timestamp.Format(time.RFC3339))
		log.Printf("  Данные измерений:")
		log.Printf("    Температура: %.2f °C", data.Data.Temperature)
		log.Printf("    Давление: %.2f мм.рт.ст.", data.Data.Pressure)
		log.Printf("    Влажность: %.2f %%", data.Data.Humidity)

		// Добавление данных в список последних данных (максимум 100 записей)
		latestData = append(latestData, data)
		if len(latestData) > 100 {
			latestData = latestData[1:]
		}

		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	// Обработчик для отображения последних данных
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "User2 Dashboard",
			"data":  latestData,
		})
	})

	// Загрузка HTML шаблонов
	r.LoadHTMLGlob("internal/templates/*.html")

	// Запуск сервера
	log.Println("User2 сервис запущен на порту 8083")
	if err := r.Run(":8083"); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
