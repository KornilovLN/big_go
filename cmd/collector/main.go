package main

import (
	"big_go/config"
	"big_go/internal/models"
	"big_go/internal/repository/redis"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	// Инициализация контекста
	ctx := context.Background()

	// Инициализация конфигурации RabbitMQ
	rabbitConfig, err := config.LoadRabbitMQConfig("config_rabbitmq.json")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации RabbitMQ: %v", err)
	}

	// Инициализация конфигурации Redis
	redisConfig, err := config.LoadRedisConfig("config_redis.json")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации Redis: %v", err)
	}

	// Подключение к Redis
	redisRepo, err := redis.New(redisConfig)
	if err != nil {
		log.Fatalf("Ошибка подключения к Redis: %v", err)
	}
	defer redisRepo.Close()

	// Подключение к RabbitMQ
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/",
		rabbitConfig.User, rabbitConfig.Password, rabbitConfig.Host, rabbitConfig.Port))
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

	// Получение сообщений
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
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
	for msg := range msgs {
		var data models.SensorData
		err := json.Unmarshal(msg.Body, &data)
		if err != nil {
			log.Printf("Ошибка разбора сообщения: %v", err)
			msg.Nack(false, true) // Отклонить сообщение и вернуть в очередь
			continue
		}

		// Сохранение данных в Redis
		err = redisRepo.SaveData(ctx, &data)
		if err != nil {
			log.Printf("Ошибка сохранения данных в Redis: %v", err)
			msg.Nack(false, true) // Отклонить сообщение и вернуть в очередь
			continue
		}

		log.Printf("Получены данные от сенсора %s (тип: %s): %.2f", data.ID, data.Type, data.Value)
		msg.Ack(false) // Подтвердить обработку сообщения
	}
}
