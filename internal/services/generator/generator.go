package generator

import (
	"big_go/internal/models"
	"math/rand"
	"time"
)

const post_numbers = 4
const address_numbers = 10
const recipient_numbers = 2

// Generator представляет генератор данных
type Generator struct {
	rand *rand.Rand
}

// NewGenerator создает новый экземпляр генератора
func NewGenerator() *Generator {
	return &Generator{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// GenerateData генерирует случайные данные датчиков
func (g *Generator) GenerateData() models.SensorData {
	// Определяем получателя (User1 или User2)
	recipient := "User1"
	if g.rand.Intn(recipient_numbers) == 1 {
		recipient = "User2"
	}

	// Генерируем случайный номер поста (1-10)
	postID := g.rand.Intn(post_numbers) + 1

	// Генерируем случайный адрес
	address := g.rand.Intn(address_numbers) + 1

	// Создаем метаданные
	meta := models.MetaData{
		Recipient: recipient,
		PostID:    postID,
		Address:   address,
		Timestamp: time.Now(),
	}

	// Генерируем данные измерений
	data := models.DataPoint{
		Temperature: 22.0 + g.rand.Float64()*3.0,   // от 20 до +25 градусов
		Pressure:    740.0 + g.rand.Float64()*40.0, // от 740 до 780 мм.рт.ст.
		Humidity:    40.0 + g.rand.Float64()*40.0,  // от 40 до 80%
	}

	return models.SensorData{
		Meta: meta,
		Data: data,
	}
}
