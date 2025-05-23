// pkg/utils/utils.go
package utils

import (
	"fmt"
	"time"
)

// FormatDate возвращает строковое представление даты в заданном формате
func FormatDate(t time.Time, layout string) string {
	return t.Format(layout)
}

// Contains проверяет, содержится ли строка в срезе строк
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// ErrorHandler обрабатывает ошибки и выводит их в стандартный вывод
func ErrorHandler(err error) {
	if err != nil {
		fmt.Println("Error:", err)
	}
}
