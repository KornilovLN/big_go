package main

import (
	"big_go/config"
	"big_go/internal/models"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/streadway/amqp"
)

func main() {
	// Инициализация конфигурации RabbitMQ
	rabbitConfig, err := config.LoadRabbitMQConfig("config_rabbitmq.json")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации RabbitMQ: %v", err)
	}

	// Подключение к RabbitMQ
	conn, err := amqp.Dial(
		"amqp://" + rabbitConfig.User + ":" + rabbitConfig.Password +
			"@" + rabbitConfig.Host + ":" + fmt.Sprintf("%d", rabbitConfig.Port) + "/" + rabbitConfig.VHost,
	)
	if err != nil {
		log.Fatalf("Ошибка подключения к RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Создание канала
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Ошибка создания канала: %v", err)
	}
	defer ch.Close()

	// Объявление очереди
	q, err := ch.QueueDeclare(
		"sensor_data", // имя очереди
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		log.Fatalf("Ошибка объявления очереди: %v", err)
	}

	// Настройка потребителя сообщений
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Ошибка регистрации потребителя: %v", err)
	}

	log.Println("Коллектор запущен. Ожидание сообщений...")

	// Обработка сообщений
	for d := range msgs {
		var sensorData models.SensorData
		err := json.Unmarshal(d.Body, &sensorData)
		if err != nil {
			log.Printf("Ошибка десериализации сообщения: %v", err)
			continue
		}

		log.Printf("Получено сообщение: %+v", sensorData)

		// Определяем, куда отправить данные
		var endpoint string
		if sensorData.Meta.Recipient == "User1" {
			endpoint = "http://user1:8082/data"
		} else {
			endpoint = "http://user2:8083/data"
		}

		// Отправляем данные соответствующему пользователю
		jsonData, err := json.Marshal(sensorData)
		if err != nil {
			log.Printf("Ошибка сериализации данных: %v", err)
			continue
		}

		resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("Ошибка отправки данных пользователю: %v", err)
			continue
		}
		resp.Body.Close()

		log.Printf("Данные успешно отправлены на %s", endpoint)
	}
}
