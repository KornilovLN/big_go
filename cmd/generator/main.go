package main

import (
	"big_go/config"
	"big_go/internal/services/generator"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/streadway/amqp"
)

func main() {
	// Инициализация конфигурации RabbitMQ
	rabbitConfig, err := config.LoadRabbitMQConfig("config_rabbitmq.json")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации RabbitMQ: %v", err)
	}

	// Подключение к RabbitMQ
	conn, err := amqp.Dial("amqp://" + rabbitConfig.User + ":" + rabbitConfig.Password +
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

	// Инициализация генератора данных
	gen := generator.New()

	// Запуск генерации данных
	for {
		// Генерация случайного интервала от 1 до 5 секунд
		interval := time.Duration(rand.Intn(4)+1) * time.Second

		// Генерация данных
		data := gen.GenerateData()

		// Сериализация данных в JSON
		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Printf("Ошибка сериализации данных: %v", err)
			continue
		}

		// Публикация сообщения
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        jsonData,
			})
		if err != nil {
			log.Printf("Ошибка публикации сообщения: %v", err)
		} else {
			log.Printf("Отправлено сообщение: %s", string(jsonData))
		}

		// Ожидание перед следующей генерацией
		time.Sleep(interval)
	}
}
