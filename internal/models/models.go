package models

import "time"

// SensorData представляет данные от датчиков
type SensorData struct {
	Meta MetaData  `json:"meta"`
	Data DataPoint `json:"data"`
}

// MetaData содержит метаданные сообщения
type MetaData struct {
	Recipient string    `json:"recipient"` // User1 или User2
	PostID    int       `json:"post_id"`   // Номер поста (1-10)
	Address   int       `json:"address"`   // Адрес
	Timestamp time.Time `json:"timestamp"` // Временная метка
}

// DataPoint содержит данные измерений
type DataPoint struct {
	Temperature float64 `json:"temperature"` // Температура в градусах Цельсия
	Pressure    float64 `json:"pressure"`    // Давление в мм.рт.ст.
	Humidity    float64 `json:"humidity"`    // Влажность в %
}
