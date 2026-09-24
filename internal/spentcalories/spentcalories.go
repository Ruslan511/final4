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
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {

	splitStrings := strings.Split(data, ",")
	if len(splitStrings) != 3 {
		return 0, "", time.Duration(0), errors.New("Передана некорректная строка")
	}

	duration, err := time.ParseDuration(splitStrings[2])
	if err != nil {
		return 0, "", time.Duration(0), err
	}
	if duration <= time.Duration(0) {
		return 0, "", time.Duration(0), errors.New("Продолжительность тренировки равна нулю или отрицательная")
	}

	steps, err := strconv.Atoi(splitStrings[0])
	if err != nil {
		return 0, "", time.Duration(0), err
	}
	if steps <= 0 {
		return 0, "", time.Duration(0), errors.New("Количество шагов равно нулю, или отрицательное")
	}

	return steps, splitStrings[1], duration, nil
}

func distance(steps int, height float64) float64 {
	return (height * stepLengthCoefficient * float64(steps)) / float64(mInKm)

}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {

	if duration <= 0 {
		return 0
	}

	distance := distance(steps, height)
	return distance / duration.Hours()

}

func TrainingInfo(data string, weight, height float64) (string, error) {

	steps, action, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}

	averageSpeed := meanSpeed(steps, height, duration)
	distance := distance(steps, height)

	resp := ""
	switch action {
	case "Ходьба":
		calories, err := WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
		resp = fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", action, duration.Hours(), distance, averageSpeed, calories)
	case "Бег":
		calories, err := RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
		resp = fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", action, duration.Hours(), distance, averageSpeed, calories)
	default:
		return "", errors.New("неизвестный тип тренировки")
	}
	return resp, nil

}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	if steps <= 0 || duration <= time.Duration(0) || weight <= 0 || height <= 0 {
		return 0, errors.New("Вы ничего не прошли, или не затратили время!")
	}
	averageSpeed := meanSpeed(steps, height, duration)
	return (weight * duration.Minutes() * averageSpeed) / minInH, nil

}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || duration <= time.Duration(0) || weight <= 0 || height <= 0 {
		return 0, errors.New("Вы ничего не прошли, или не затратили время!")
	}
	averageSpeed := meanSpeed(steps, height, duration)
	return ((weight * duration.Minutes() * averageSpeed) / minInH) * walkingCaloriesCoefficient, nil

}
