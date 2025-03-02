package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/go1fl-4-sprint-final/internal/spentcalories"
)

var (
	StepLength     = 0.65 // длина шага в метрах
	ErrInvalidData = errors.New("Не корректный формат данных")
)

func parsePackage(data string) (int, time.Duration, error) {
	// ваш код ниже
	dataParts := strings.Split(data, ",") // Разделяем входящие данные, исходя по формату: "678,0h50m"
	if len(dataParts) != 2 {
		return 0, 0, ErrInvalidData
	}

	steps, err := strconv.Atoi(strings.TrimSpace(dataParts[0])) // Преобразуем первый элемент в целое число (количество шагов)
	if err != nil {
		return 0, 0, err
	}

	duration, err := time.ParseDuration(strings.TrimSpace(dataParts[1])) // Преобразуем второй элемент в time (продолжительность ходьбы)
	if err != nil {
		return 0, 0, err
	}

	return steps, duration, nil
}

// DayActionInfo обрабатывает входящий пакет, который передаётся в
// виде строки в параметре data. Параметр storage содержит пакеты за текущий день.
// Если время пакета относится к новым суткам, storage предварительно
// очищается.
// Если пакет валидный, он добавляется в слайс storage, который возвращает
// функция. Если пакет невалидный, storage возвращается без изменений.
func DayActionInfo(data string, weight, height float64) string {
	// ваш код ниже
	steps, duration, err := parsePackage(data)
	if err != nil {
		return ""
	}

	if steps <= 0 { // Если количество шагов не больше 0, то возвращать пустую строку
		return ""
	}

	distanceKm := (float64(steps) * StepLength) / 1000 // Вычисляем дистанцию (в километрах)

	caloriesBurned := spentcalories.WalkingSpentCalories(steps, weight, height, duration) // Вычисляем потребленные калории

	result := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.", steps, distanceKm, caloriesBurned)

	return result
}
