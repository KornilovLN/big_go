package main

import (
	"big_go/config"
	"big_go/internal/repository/redis"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

// Структура для хранения данных для шаблона
type PageData struct {
	Title    string
	Messages []Message
}

// Структура для сообщения
type Message struct {
	ID        string
	Type      string
	Value     float64
	Timestamp time.Time
}

func main() {
	// Инициализация контекста
	ctx := context.Background()

	// Загрузка конфигурации Redis
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

	// Загрузка шаблонов
	tmpl, err := template.ParseFiles("internal/templates/base.html", "internal/templates/index.html")
	if err != nil {
		log.Fatalf("Ошибка загрузки шаблонов: %v", err)
	}

	// Обработчик для главной страницы
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Получение последних данных из Redis для User2
		data, err := redisRepo.GetLatestData(ctx, "User2", 10)
		if err != nil {
			log.Printf("Ошибка получения данных из Redis: %v", err)
			http.Error(w, "Ошибка получения данных", http.StatusInternalServerError)
			return
		}

		// Преобразование данных для шаблона
		messages := make([]Message, 0, len(data))
		for _, item := range data {
			messages = append(messages, Message{
				ID:        item.ID,
				Type:      item.Type,
				Value:     item.Value,
				Timestamp: time.Unix(item.Timestamp, 0),
			})
		}

		// Отображение шаблона
		pageData := PageData{
			Title:    "User2 Dashboard",
			Messages: messages,
		}

		err = tmpl.ExecuteTemplate(w, "base", pageData)
		if err != nil {
			log.Printf("Ошибка отображения шаблона: %v", err)
			http.Error(w, "Ошибка отображения страницы", http.StatusInternalServerError)
		}
	})

	// Обработчик для получения обновлений через SSE
	http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		// Настройка заголовков для SSE
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Создание канала для получения уведомлений о новых данных
		keyChan, err := redisRepo.Subscribe(ctx, "User2")
		if err != nil {
			log.Printf("Ошибка подписки на уведомления: %v", err)
			http.Error(w, "Ошибка подписки на уведомления", http.StatusInternalServerError)
			return
		}

		// Отправка событий клиенту
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "SSE не поддерживается", http.StatusInternalServerError)
			return
		}

		// Обработка закрытия соединения
		notify := r.Context().Done()
		go func() {
			<-notify
			log.Println("Клиент отключился")
		}()

		// Отправка событий
		for key := range keyChan {
			// Получение данных по ключу
			data, err := redisRepo.GetData(ctx, key)
			if err != nil {
				log.Printf("Ошибка получения данных по ключу %s: %v", key, err)
				continue
			}

			// Преобразование данных в JSON
			message := Message{
				ID:        data.ID,
				Type:      data.Type,
				Value:     data.Value,
				Timestamp: time.Unix(data.Timestamp, 0),
			}

			jsonData, err := json.Marshal(message)
			if err != nil {
				log.Printf("Ошибка сериализации данных: %v", err)
				continue
			}

			// Отправка события клиенту
			fmt.Fprintf(w, "data: %s\n\n", jsonData)
			flusher.Flush()
		}
	})

	// Обработчик для статических файлов
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Запуск сервера
	log.Println("User2 сервер запущен на порту 8083")
	if err := http.ListenAndServe(":8083", nil); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
