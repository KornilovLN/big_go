# Усовершенствование проекта с использованием Redis
    В текущей архитектуре Redis уже настроен в docker-compose.yml, но не используется в логике приложения.
    Модифицируем проект, чтобы включить Redis в рабочий процесс согласно рекомендуемой архитектуре из документа big_go-var2.md.

## План модификации:
* **Изменить коллектор, чтобы он сохранял данные в Redis вместо прямой отправки пользователям**
* **Настроить Redis Pub/Sub для уведомлений пользователей**
* **Модифицировать пользовательские приложения для получения данных из Redis**

### Шаг 1: Создадим пакет для работы с Redis (redis.go)
```GO
package redis

import (
	"big_go/config"
	"big_go/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisRepository представляет репозиторий для работы с Redis
type RedisRepository struct {
	client *redis.Client
}

// NewRedisRepository создает новый экземпляр репозитория Redis
func NewRedisRepository(cfg config.RedisConfig) (*RedisRepository, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Проверка соединения
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к Redis: %v", err)
	}

	return &RedisRepository{
		client: client,
	}, nil
}

// SaveSensorData сохраняет данные датчиков в Redis с TTL
func (r *RedisRepository) SaveSensorData(ctx context.Context, data models.SensorData) error {
	// Создаем ключ в формате user:{recipient}:data:{timestamp}
	key := fmt.Sprintf("user:%s:data:%d", data.Meta.Recipient, data.Meta.Timestamp.UnixNano())

	// Сериализуем данные в JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("ошибка сериализации данных: %v", err)
	}

	// Сохраняем данные в Redis с TTL 24 часа
	err = r.client.Set(ctx, key, jsonData, 24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("ошибка сохранения данных в Redis: %v", err)
	}

	log.Printf("Данные сохранены в Redis с ключом: %s", key)
	return nil
}

// PublishNotification отправляет уведомление через Redis Pub/Sub
func (r *RedisRepository) PublishNotification(ctx context.Context, recipient string, data models.SensorData) error {
	// Создаем канал для уведомлений в формате user:{recipient}:notifications
	channel := fmt.Sprintf("user:%s:notifications", recipient)

	// Создаем уведомление с минимальной информацией
	notification := map[string]interface{}{
		"message":   "Новые данные датчиков доступны",
		"timestamp": data.Meta.Timestamp,
		"post_id":   data.Meta.PostID,
		"address":   data.Meta.Address,
		"key":       fmt.Sprintf("user:%s:data:%d", data.Meta.Recipient, data.Meta.Timestamp.UnixNano()),
	}

	// Сериализуем уведомление в JSON
	jsonNotification, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("ошибка сериализации уведомления: %v", err)
	}

	// Публикуем уведомление
	err = r.client.Publish(ctx, channel, jsonNotification).Err()
	if err != nil {
		return fmt.Errorf("ошибка публикации уведомления: %v", err)
	}

	log.Printf("Уведомление отправлено в канал: %s", channel)
	return nil
}

// GetSensorData получает данные датчиков из Redis по ключу
func (r *RedisRepository) GetSensorData(ctx context.Context, key string) (models.SensorData, error) {
	var data models.SensorData

	// Получаем данные из Redis
	jsonData, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return data, fmt.Errorf("данные не найдены по ключу: %s", key)
		}
		return data, fmt.Errorf("ошибка получения данных из Redis: %v", err)
	}

	// Десериализуем данные из JSON
	err = json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		return data, fmt.Errorf("ошибка десериализации данных: %v", err)
	}

	return data, nil
}

// GetLatestSensorData получает последние данные датчиков для указанного пользователя
func (r *RedisRepository) GetLatestSensorData(ctx context.Context, recipient string, limit int) ([]models.SensorData, error) {
	var result []models.SensorData

	// Получаем все ключи для данного пользователя
	pattern := fmt.Sprintf("user:%s:data:*", recipient)
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения ключей из Redis: %v", err)
	}

	// Если ключей нет, возвращаем пустой массив
	if len(keys) == 0 {
		return result, nil
	}

	// Ограничиваем количество ключей, если их больше чем limit
	if len(keys) > limit && limit > 0 {
		// Сортируем ключи по времени (последняя часть ключа - timestamp)
		// Это упрощенная сортировка, в реальном проекте может потребоваться более сложная логика
		keys = keys[len(keys)-limit:]
	}

	// Получаем данные для каждого ключа
	for _, key := range keys {
		data, err := r.GetSensorData(ctx, key)
		if err != nil {
			log.Printf("Ошибка получения данных для ключа %s: %v", key, err)
			continue
		}
		result = append(result, data)
	}

	return result, nil
}

// Close закрывает соединение с Redis
func (r *RedisRepository) Close() error {
	return r.client.Close()
}
```

### Шаг 2: Обновим конфигурацию Redis (redis.go)
```GO
package config

import (
	"encoding/json"
	"os"
)

// RedisConfig содержит конфигурацию для подключения к Redis
type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

// LoadRedisConfig загружает конфигурацию Redis из файла
func LoadRedisConfig(filePath string) (RedisConfig, error) {
	var config RedisConfig

	// Чтение файла конфигурации
	file, err := os.ReadFile(filePath)
	if err != nil {
		return config, err
	}

	// Разбор JSON
	err = json.Unmarshal(file, &config)
	if err != nil {
		return config, err
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

	return config, nil
}
```

### Шаг 3: Модифицируем коллектор (main.go)
```GO
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

	// Инициализация репозитория Redis
	redisRepo, err := redis.NewRedisRepository(redisConfig)
	if err != nil {
		log.Fatalf("Ошибка инициализации репозитория Redis: %v", err)
	}
	defer redisRepo.Close()

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

		// Сохраняем данные в Redis
		err = redisRepo.SaveSensorData(ctx, sensorData)
		if err != nil {
			log.Printf("Ошибка сохранения данных в Redis: %v", err)
			continue
		}

		// Отправляем уведомление через Redis Pub/Sub
		err = redisRepo.PublishNotification(ctx, sensorData.Meta.Recipient, sensorData)
		if err != nil {
			log.Printf("Ошибка отправки уведомления: %v", err)
			continue
		}

		log.Printf("Данные успешно обработаны и сохранены в Redis")
	}
}
```

### Шаг 4: Модифицируем пользовательские приложения (User1/main.go)
```GO
package main

import (
	"big_go/config"
	"big_go/internal/models"
	"big_go/internal/repository/redis"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	goredis "github.com/go-redis/redis/v8"
)

func main() {
	// Инициализация контекста
	ctx := context.Background()

	// Инициализация конфигурации Redis
	redisConfig, err := config.LoadRedisConfig("config_redis.json")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации Redis: %v", err)
	}

	// Инициализация репозитория Redis
	redisRepo, err := redis.NewRedisRepository(redisConfig)
	if err != nil {
		log.Fatalf("Ошибка инициализации репозитория Redis: %v", err)
	}
	defer redisRepo.Close()

	// Инициализация Gin
	r := gin.Default()

	// Канал для хранения последних полученных данных
	var latestData []models.SensorData

	// Загружаем начальные данные из Redis
	initialData, err := redisRepo.GetLatestSensorData(ctx, "User1", 100)
	if err != nil {
		log.Printf("Ошибка загрузки начальных данных из Redis: %v", err)
	} else {
		latestData = initialData
		log.Printf("Загружено %d записей из Redis", len(initialData))
	}

	// Создаем клиент Redis для подписки на уведомления
	redisClient := goredis.NewClient(&goredis.Options{
		Addr:     redisConfig.Host + ":" + fmt.Sprintf("%d", redisConfig.Port),
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})
	defer redisClient.Close()

	// Подписываемся на канал уведомлений
	pubsub := redisClient.Subscribe(ctx, "user:User1:notifications")
	defer pubsub.Close()

	// Запускаем горутину для обработки уведомлений
	go func() {
		for {
			msg, err := pubsub.ReceiveMessage(ctx)
			if err != nil {
				log.Printf("Ошибка получения сообщения из Redis Pub/Sub: %v", err)
				continue
			}

			var notification map[string]interface{}
			if err := json.Unmarshal([]byte(msg.Payload), &notification); err != nil {
				log.Printf("Ошибка десериализации уведомления: %v", err)
				continue
			}

			log.Printf("Получено уведомление: %v", notification)

			// Получаем ключ данных из уведомления
			if key, ok := notification["key"].(string); ok {
				// Получаем данные из Redis по ключу
				data, err := redisRepo.GetSensorData(ctx, key)
				if err != nil {
					log.Printf("Ошибка получения данных из Redis: %v", err)
					continue
				}

				// Добавляем данные в список последних данных
				latestData = append(latestData, data)
				if len(latestData) > 100 {
					latestData = latestData[1:]
				}

				log.Printf("Данные успешно получены из Redis и добавлены в список")
			}
		}
	}()

	// Обработчик для отображения последних данных
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "User1 Dashboard",
			"data":  latestData,
		})
	})

	// Обработчик для получения данных через API
	r.GET("/api/data", func(c *gin.Context) {
		c.JSON(http.StatusOK, latestData)
	})

	// Обработчик для обновления данных вручную
	r.GET("/refresh", func(c *gin.Context) {
		// Загружаем последние данные из Redis
		data, err := redisRepo.GetLatestSensorData(ctx, "User1", 100)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		latestData = data
		c.Redirect(http.StatusFound, "/")
	})

	// Загрузка HTML шаблонов
	r.LoadHTMLGlob("internal/templates/*.html")

	// Запуск сервера
	log.Println("User1 сервис запущен на порту 8082")
	if err := r.Run(":8082"); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
```

### Шаг 5: Модифицируем пользовательское приложение User2 (User2/main.go)
```GO
package main

import (
	"big_go/config"
	"big_go/internal/models"
	"big_go/internal/repository/redis"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	goredis "github.com/go-redis/redis/v8"
)

func main() {
	// Инициализация контекста
	ctx := context.Background()

	// Инициализация конфигурации Redis
	redisConfig, err := config.LoadRedisConfig("config_redis.json")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации Redis: %v", err)
	}

	// Инициализация репозитория Redis
	redisRepo, err := redis.NewRedisRepository(redisConfig)
	if err != nil {
		log.Fatalf("Ошибка инициализации репозитория Redis: %v", err)
	}
	defer redisRepo.Close()

	// Инициализация Gin
	r := gin.Default()

	// Канал для хранения последних полученных данных
	var latestData []models.SensorData

	// Загружаем начальные данные из Redis
	initialData, err := redisRepo.GetLatestSensorData(ctx, "User2", 100)
	if err != nil {
		log.Printf("Ошибка загрузки начальных данных из Redis: %v", err)
	} else {
		latestData = initialData
		log.Printf("Загружено %d записей из Redis", len(initialData))
	}

	// Создаем клиент Redis для подписки на уведомления
	redisClient := goredis.NewClient(&goredis.Options{
		Addr:     redisConfig.Host + ":" + fmt.Sprintf("%d", redisConfig.Port),
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})
	defer redisClient.Close()

	// Подписываемся на канал уведомлений
	pubsub := redisClient.Subscribe(ctx, "user:User2:notifications")
	defer pubsub.Close()

	// Запускаем горутину для обработки уведомлений
	go func() {
		for {
			msg, err := pubsub.ReceiveMessage(ctx)
			if err != nil {
				log.Printf("Ошибка получения сообщения из Redis Pub/Sub: %v", err)
				continue
			}

			var notification map[string]interface{}
			if err := json.Unmarshal([]byte(msg.Payload), &notification); err != nil {
				log.Printf("Ошибка десериализации уведомления: %v", err)
				continue
			}

			log.Printf("Получено уведомление: %v", notification)

			// Получаем ключ данных из уведомления
			if key, ok := notification["key"].(string); ok {
				// Получаем данные из Redis по ключу
				data, err := redisRepo.GetSensorData(ctx, key)
				if err != nil {
					log.Printf("Ошибка получения данных из Redis: %v", err)
					continue
				}

				// Добавляем данные в список последних данных
				latestData = append(latestData, data)
				if len(latestData) > 100 {
					latestData = latestData[1:]
				}

				log.Printf("Данные успешно получены из Redis и добавлены в список")
			}
		}
	}()

	// Обработчик для отображения последних данных
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "User2 Dashboard",
			"data":  latestData,
		})
	})

	// Обработчик для получения данных через API
	r.GET("/api/data", func(c *gin.Context) {
		c.JSON(http.StatusOK, latestData)
	})

	// Обработчик для обновления данных вручную
	r.GET("/refresh", func(c *gin.Context) {
		// Загружаем последние данные из Redis
		data, err := redisRepo.GetLatestSensorData(ctx, "User2", 100)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		latestData = data
		c.Redirect(http.StatusFound, "/")
	})

	// Загрузка HTML шаблонов
	r.LoadHTMLGlob("internal/templates/*.html")

	// Запуск сервера
	log.Println("User2 сервис запущен на порту 8083")
	if err := r.Run(":8083"); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
```

### Шаг 6: Обновим index.html шаблон для отображения данных с информацией о Redis
```html
<!DOCTYPE html>
<html>
<head>
    <title>{{ .title }}</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 20px;
        }
        h1 {
            color: #333;
        }
        .info-box {
            background-color: #f8f9fa;
            border: 1px solid #ddd;
            padding: 10px;
            margin-bottom: 20px;
            border-radius: 4px;
        }
        .refresh-btn {
            background-color: #4CAF50;
            color: white;
            padding: 10px 15px;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            margin-bottom: 20px;
        }
        .refresh-btn:hover {
            background-color: #45a049;
        }
        table {
            width: 100%;
            border-collapse: collapse;
            margin-top: 20px;
        }
        th, td {
            border: 1px solid #ddd;
            padding: 8px;
            text-align: left;
        }
        th {
            background-color: #f2f2f2;
        }
        tr:nth-child(even) {
            background-color: #f9f9f9;
        }
    </style>
</head>
<body>
    <h1>{{ .title }}</h1>
    
    <div class="info-box">
        <p><strong>Информация:</strong> Данные поступают из Redis через систему уведомлений Pub/Sub. 
        Коллектор сохраняет данные в Redis и отправляет уведомления о новых данных.</p>
    </div>
    
    <a href="/refresh" class="refresh-btn">Обновить данные</a>
    
    <h2>Последние полученные данные</h2>
    
    <table>
        <tr>
            <th>Время</th>
            <th>Пост ID</th>
            <th>Адрес</th>
            <th>Температура (°C)</th>
            <th>Давление (мм.рт.ст.)</th>
            <th>Влажность (%)</th>
        </tr>
        {{ range .data }}
        <tr>
            <td>{{ .Meta.Timestamp.Format "2006-01-02 15:04:05" }}</td>
            <td>{{ .Meta.PostID }}</td>
            <td>{{ .Meta.Address }}</td>
            <td>{{ printf "%.2f" .Data.Temperature }}</td>
            <td>{{ printf "%.2f" .Data.Pressure }}</td>
            <td>{{ printf "%.2f" .Data.Humidity }}</td>
        </tr>
        {{ else }}
        <tr>
            <td colspan="6" style="text-align: center;">Нет данных</td>
        </tr>
        {{ end }}
    </table>
    
    <script>
        // Автоматическое обновление страницы каждые 10 секунд
        setTimeout(function() {
            location.reload();
        }, 10000);
    </script>
</body>
</html>
```

### Шаг 7: Создадим пример файла конфигурации Redis (config_redis.json)
```json
{
  "host": "redis",
  "port": 6379,
  "password": "",
  "db": 0
}
```

### Шаг 8: Обновим go.mod для добавления зависимости go-redis
* **Для работы с Redis нам потребуется добавить зависимость go-redis. Выполните следующую команду:**
```bash
go get github.com/go-redis/redis/v8
```
* Это обновит ваш файл go.mod, добавив необходимую зависимость.

### Итоги модификации:
#### Архитектурные изменения:
* Коллектор теперь сохраняет данные в Redis вместо прямой отправки пользователям
* Используется Redis Pub/Sub для уведомления пользователей о новых данных
* Пользовательские приложения получают данные из Redis по уведомлениям

#### Преимущества новой архитектуры:
* Отказоустойчивость: Если пользовательские приложения временно недоступны, данные сохраняются в Redis и могут быть получены позже
* Масштабируемость: Можно легко добавлять новых пользователей без изменения логики коллектора
* Производительность: Redis обеспечивает высокую скорость чтения/записи данных
* Разделение ответственности: Коллектор отвечает только за сбор и сохранение данных, а не за их доставку

#### Дополнительные функции:
* Возможность ручного обновления данных через кнопку "Обновить данные"
* API-эндпоинт для получения данных в формате JSON
* Улучшенный пользовательский интерфейс с информационным блоком
 
### Шаг 9: Обновим README.md с описанием новой архитектуры

