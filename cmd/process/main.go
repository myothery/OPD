//go:generate go run ../generator/graphg.go

package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/se1lzor/OPD/internal/common"
)

func loadCirclesFromJSON(filename string) ([]common.CircleData, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("Файл '%s' не найден. Проверьте имя файла и путь", filename)
		}
		return nil, fmt.Errorf("Ошибка при чтении файла: %v", err)
	}

	var circles []common.CircleData
	err = json.Unmarshal(data, &circles)
	if err != nil {
		return nil, fmt.Errorf("Ошибка парсинга JSON: %v", err)
	}

	return circles, nil
}

func diameterFromPoints(point1, point2, point3 common.Point) (float64, error) {
	x1 := point1.X
	y1 := point1.Y
	x2 := point2.X
	y2 := point2.Y
	x3 := point3.X
	y3 := point3.Y

	determinantA := x1*(y2-y3) - y1*(x2-x3) + x2*y3 - x3*y2

	if math.Abs(determinantA) < 1e-8 {
		return 0, fmt.Errorf("Точки лежат на одной прямой")
	}

	square1 := x1*x1 + y1*y1
	square2 := x2*x2 + y2*y2
	square3 := x3*x3 + y3*y3

	xc := (square1*(y2-y3) + square2*(y3-y1) + square3*(y1-y2)) / (2 * determinantA)
	yc := (square1*(x3-x2) + square2*(x1-x3) + square3*(x2-x1)) / (2 * determinantA)

	radius := math.Hypot(x1-xc, y1-yc)
	return radius * 2, nil
}

func roundDiameter(diameter float64) float64 {
	// Округляем до 2 знаков после запятой для группировки
	return math.Round(diameter*100) / 100
}

func main() {
	var filename string
	fmt.Println("Введите название файла: ")
	fmt.Scanln(&filename)

	// Загружаем данные из JSON
	circles, err := loadCirclesFromJSON(filename + ".json")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Загружено окружностей: %d\n\n", len(circles))

	// Словарь для подсчета диаметров
	diameterCount := make(map[float64]int)

	// Обрабатываем каждую окружность
	for i, circle := range circles {
		if len(circle.Points) < 3 {
			fmt.Printf("Окружность %d: недостаточно точек (только %d)\n", i+1, len(circle.Points))
			continue
		}

		// Вычисляем диаметр по трем точкам
		diameter, err := diameterFromPoints(
			circle.Points[0],
			circle.Points[1],
			circle.Points[2],
		)
		if err != nil {
			fmt.Printf("Окружность %d: ошибка - %v\n", i+1, err)
			continue
		}

		// Округляем диаметр для группировки
		roundedDiameter := roundDiameter(diameter)
		diameterCount[roundedDiameter]++

		fmt.Printf("Окружность %d: центр X=%.2f, радиус=%.2f, диаметр=%.2f\n",
			i+1, circle.CenterX, circle.Radius, diameter)
	}

	// Выводим результаты в отсортированном порядке
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("Результаты группировки по диаметрам:")
	fmt.Println(strings.Repeat("=", 50))

	// Сортируем диаметры
	var diameters []float64
	for d := range diameterCount {
		diameters = append(diameters, d)
	}
	sort.Float64s(diameters)

	// Выводим статистику
	for _, d := range diameters {
		fmt.Printf("Количество труб диаметра %.2f: %d\n", d, diameterCount[d])
	}

	// Итоговая статистика
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("Всего обработано окружностей: %d\n", len(circles))
	fmt.Printf("Уникальных диаметров: %d\n", len(diameterCount))
}
