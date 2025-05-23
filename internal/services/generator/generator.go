package generator

import (
	"big_go/internal/models"
	"fmt"
	"math/rand"
	"time"
)

// Generator представляет генератор данных
type Generator struct{}

// New создает новый генератор данных
func New() *Generator {
	// Инициализация генератора случайных чисел
	rand.Seed(time.Now().UnixNano())
	return &Generator{}
}

// GenerateData генерирует случайные данные
func (g *Generator) GenerateData() *models.SensorData {
	// Типы сенсоров
	sensorTypes := []string{"temperature", "pressure", "humidity"}

	// Получатели данных
	recipients := []string{"User1", "User2"}

	// Выбор случайного типа сенсора
	sensorType := sensorTypes[rand.Intn(len(sensorTypes))]

	// Выбор случайного получателя
	recipient := recipients[rand.Intn(len(recipients))]

	// Генерация случайного значения в зависимости от типа сенсора
	var value float64
	switch sensorType {
	case "temperature":
		value = 15 + rand.Float64()*25 // от 15 до 40 градусов
	case "pressure":
		value = 740 + rand.Float64()*40 // от 740 до 780 мм.рт.ст.
	case "humidity":
		value = 30 + rand.Float64()*60 // от 30% до 90%
	}

	// Создание данных сенсора
	sensorData := &models.SensorData{
		ID:        fmt.Sprintf("%s-%d", sensorType, rand.Intn(100)),
		Type:      sensorType,
		Value:     value,
		Timestamp: time.Now().Unix(),
		Recipient: recipient,
		CreatedAt: time.Now(),
	}

	return sensorData
}
