package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep   = 0.65  // средняя длина шага.
	mInKm     = 1000  // количество метров в километре.
	minInH    = 60    // количество минут в часе.
	kmhInMsec = 0.278 // коэффициент для преобразования км/ч в м/с.
	cmInM     = 100   // количество сантиметров в метре.
)

var (
	ErrInvalidData = errors.New("Не корректный формат данных")
)

func parseTraining(data string) (int, string, time.Duration, error) {
	// ваш код ниже
	dataParts := strings.Split(data, ",") // Разделяем входящие данные, исходя по формату: "3456,Ходьба,3h00m"
	if len(dataParts) != 3 {
		return 0, "", 0, ErrInvalidData
	}

	steps, err := strconv.Atoi(strings.TrimSpace(dataParts[0])) // Преобразуем первый элемент в целое число (количество шагов)
	if err != nil {
		return 0, "", 0, err
	}

	activity := strings.TrimSpace(dataParts[1]) // Берем значение второго элемента - нашей активности

	duration, err := time.ParseDuration(strings.TrimSpace(dataParts[2])) // Преобразуем третий элемент в time (продолжительность ходьбы)
	if err != nil {
		return 0, "", 0, err
	}

	return steps, activity, duration, nil
}

// distance возвращает дистанцию(в километрах), которую преодолел пользователь за время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий (число шагов при ходьбе и беге).
func distance(steps int) float64 {
	// ваш код ниже
	return float64(steps) * lenStep / float64(mInKm)
}

// meanSpeed возвращает значение средней скорости движения во время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий(число шагов при ходьбе и беге).
// duration time.Duration — длительность тренировки.
func meanSpeed(steps int, duration time.Duration) float64 {
	// ваш код ниже
	if duration <= 0 {
		return 0
	}

	return distance(steps) / duration.Hours()
}

// ShowTrainingInfo возвращает строку с информацией о тренировке.
//
// Параметры:
//
// data string - строка с данными.
// weight, height float64 — вес и рост пользователя.
func TrainingInfo(data string, weight, height float64) string {
	// ваш код ниже
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		return err.Error()
	}

	distanceInKm := distance(steps)
	averageSpeed := meanSpeed(steps, duration)
	calories := 0.0

	switch strings.ToLower(activity) { // Для каждого из видов тренировки рассчитать дистанцию, среднюю скорость и калории.
	case "ходьба":
		calories = WalkingSpentCalories(steps, weight, height, duration)
	case "бег":
		calories = RunningSpentCalories(steps, weight, duration)
	default:
		return "неизвестный тип тренировки"
	}

	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %0.2f ч.\nДистанция: %0.2f км.\nСкорость: %0.2f км/ч\nСожгли калорий: %0.2f", activity, duration.Hours(), distanceInKm, averageSpeed, calories)

}

// Константы для расчета калорий, расходуемых при беге.
const (
	runningCaloriesMeanSpeedMultiplier = 18.0 // множитель средней скорости.
	runningCaloriesMeanSpeedShift      = 20.0 // среднее количество сжигаемых калорий при беге.
)

// RunningSpentCalories возвращает количество потраченных калорий при беге.
//
// Параметры:
//
// steps int - количество шагов.
// weight float64 — вес пользователя.
// duration time.Duration — длительность тренировки.
func RunningSpentCalories(steps int, weight float64, duration time.Duration) float64 {
	// ваш код здесь
	averageSpeed := meanSpeed(steps, duration)
	if averageSpeed < 0 {
		averageSpeed = 0
	}

	// Количество сжигаемых калорий при беге.
	return ((runningCaloriesMeanSpeedMultiplier * averageSpeed) - runningCaloriesMeanSpeedShift) * weight
}

// Константы для расчета калорий, расходуемых при ходьбе.
const (
	walkingCaloriesWeightMultiplier = 0.035 // множитель массы тела.
	walkingSpeedHeightMultiplier    = 0.029 // множитель роста.
)

// WalkingSpentCalories возвращает количество потраченных калорий при ходьбе.
//
// Параметры:
//
// steps int - количество шагов.
// duration time.Duration — длительность тренировки.
// weight float64 — вес пользователя.
// height float64 — рост пользователя.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) float64 {
	// ваш код здесь
	averageSpeed := meanSpeed(steps, duration)
	if averageSpeed < 0 {
		averageSpeed = 0
	}

	return ((walkingCaloriesWeightMultiplier * weight) + (averageSpeed*averageSpeed/height)*walkingSpeedHeightMultiplier) * duration.Hours() * minInH
}
