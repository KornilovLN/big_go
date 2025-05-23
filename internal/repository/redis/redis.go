package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"big_go/config"
	"big_go/internal/models"

	"github.com/go-redis/redis/v8"
)

// Repository представляет репозиторий для работы с Redis
type Repository struct {
	client *redis.Client
}

// New создает новый репозиторий Redis
func New(config *config.RedisConfig) (*Repository, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	// Проверка соединения
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &Repository{
		client: client,
	}, nil
}

// SaveData сохраняет данные в Redis
func (r *Repository) SaveData(ctx context.Context, data *models.SensorData) error {
	// Создаем ключ в формате "sensor:{type}:{id}:{timestamp}"
	key := fmt.Sprintf("sensor:%s:%s:%d", data.Type, data.ID, data.Timestamp)

	// Сериализуем данные в JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	// Сохраняем данные в Redis с TTL 24 часа
	err = r.client.Set(ctx, key, jsonData, 24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to save data to Redis: %w", err)
	}

	// Публикуем уведомление о новых данных
	channel := fmt.Sprintf("new_data:%s", data.Recipient)
	err = r.client.Publish(ctx, channel, key).Err()
	if err != nil {
		log.Printf("Warning: failed to publish notification: %v", err)
	}

	return nil
}

// GetData получает данные из Redis по ключу
func (r *Repository) GetData(ctx context.Context, key string) (*models.SensorData, error) {
	// Получаем данные из Redis
	jsonData, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("data not found")
		}
		return nil, fmt.Errorf("failed to get data from Redis: %w", err)
	}

	// Десериализуем данные из JSON
	var data models.SensorData
	err = json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return &data, nil
}

// GetLatestData получает последние данные для указанного получателя
func (r *Repository) GetLatestData(ctx context.Context, recipient string, limit int) ([]*models.SensorData, error) {
	// Ищем ключи, соответствующие шаблону
	pattern := fmt.Sprintf("sensor:*:*:*")
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get keys from Redis: %w", err)
	}

	var result []*models.SensorData

	// Получаем данные для каждого ключа
	for _, key := range keys {
		data, err := r.GetData(ctx, key)
		if err != nil {
			log.Printf("Warning: failed to get data for key %s: %v", key, err)
			continue
		}

		// Фильтруем по получателю
		if data.Recipient == recipient {
			result = append(result, data)
		}

		// Если достигли лимита, прекращаем
		if len(result) >= limit {
			break
		}
	}

	return result, nil
}

// Subscribe подписывается на уведомления о новых данных
func (r *Repository) Subscribe(ctx context.Context, recipient string) (<-chan string, error) {
	// Создаем канал для подписки
	pubsub := r.client.Subscribe(ctx, fmt.Sprintf("new_data:%s", recipient))

	// Создаем канал для передачи ключей
	keyChan := make(chan string, 100)

	// Запускаем горутину для обработки сообщений
	go func() {
		defer pubsub.Close()
		defer close(keyChan)

		for {
			msg, err := pubsub.ReceiveMessage(ctx)
			if err != nil {
				log.Printf("Error receiving message: %v", err)
				return
			}

			keyChan <- msg.Payload
		}
	}()

	return keyChan, nil
}

// Close закрывает соединение с Redis
func (r *Repository) Close() error {
	return r.client.Close()
}
