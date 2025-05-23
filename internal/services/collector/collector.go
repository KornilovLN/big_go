package collector

import (
	"big_go/internal/models"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Collector представляет сервис коллектора данных
type Collector struct {
	user1Data  chan models.SensorData
	user2Data  chan models.SensorData
	httpClient *http.Client
}

// NewCollector создает новый экземпляр коллектора
func NewCollector() *Collector {
	c := &Collector{
		user1Data:  make(chan models.SensorData, 100),
		user2Data:  make(chan models.SensorData, 100),
		httpClient: &http.Client{},
	}

	// Запуск горутин для отправки данных пользователям
	go c.sendDataToUser("user1", c.user1Data)
	go c.sendDataToUser("user2", c.user2Data)

	return c
}

// ProcessData обрабатывает полученные данные и направляет их соответствующему пользователю
func (c *Collector) ProcessData(data models.SensorData) error {
	log.Printf("Обработка данных для %s от поста %d", data.Meta.Recipient, data.Meta.PostID)

	// Направление данных соответствующему пользователю
	switch data.Meta.Recipient {
	case "User1":
		c.user1Data <- data
	case "User2":
		c.user2Data <- data
	default:
		return fmt.Errorf("неизвестный получатель: %s", data.Meta.Recipient)
	}

	return nil
}

// sendDataToUser отправляет данные соответствующему пользовательскому сервису
func (c *Collector) sendDataToUser(userService string, dataChan <-chan models.SensorData) {
	var endpoint string
	if userService == "user1" {
		endpoint = "http://user1:8082/data"
	} else {
		endpoint = "http://user2:8083/data"
	}

	for data := range dataChan {
		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Printf("Ошибка сериализации данных для %s: %v", userService, err)
			continue
		}

		resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("Ошибка отправки данных для %s: %v", userService, err)
			continue
		}
		resp.Body.Close()

		log.Printf("Данные успешно отправлены %s", userService)
	}
}
