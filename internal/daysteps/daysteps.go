package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	splitStrings := strings.Split(data, ",")
	if len(splitStrings) != 2 {
		return 0, time.Duration(0), errors.New("Строка не распарсилась")
	}
	steps, err := strconv.Atoi(splitStrings[0])
	if err != nil {
		return 0, time.Duration(0.0), err
	}

	if steps <= 0 {
		return 0, time.Duration(0), errors.New("Количество шагов пустое")
	}

	duration, err := time.ParseDuration(splitStrings[1])
	if err != nil {
		return 0, time.Duration(0), err
	}
	if duration <= time.Duration(0.0) {
		return 0, time.Duration(0), errors.New("Продолжительность меньше или равна нулю")
	}

	return steps, duration, nil

}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		return ""
	}
	if steps == 0 {
		return ""
	}
	distanceKm := (float64(steps) * stepLength) / mInKm
	spentcalories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		return ""
	}

	resp := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distanceKm, spentcalories)

	return resp

}
