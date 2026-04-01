package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
)

type Point struct {
	X float64 `json:"X"`
	Y float64 `json:"Y"`
}

func loadPointFromJson(filename string) ([]Point, error) {
	data, err := os.ReadFile(filename)

	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("Файл '%s' не найден. Проверьте имя файла и путь", filename)
		}

		return nil, fmt.Errorf("Ошибка при чтении файла: %v", err)
	}

	points := []Point{}

	err = json.Unmarshal(data, &points)

	if err != nil {
		return nil, fmt.Errorf("Ошибка рспарсинга JSON: %v", err)
	}

	return points, nil
}

func diameterFromCircle(point1, point2, point3 Point) (diameter float64, err error) {
	x1 := point1.X
	y1 := point1.Y
	x2 := point2.X
	y2 := point2.Y
	x3 := point3.X
	y3 := point3.Y

	determinantA := x1*(y2-y3) - y1*(x2-x3) + x2*y3 - x3*y2

	if math.Abs(determinantA) < 1e-8 {
		return 0, fmt.Errorf("Точки лежат на одной прямой!")
	}

	square1 := x1*x1 + y1*y1
	square2 := x2*x2 + y2*y2
	square3 := x3*x3 + y3*y3

	xc := (square1*(y2-y3) + square2*(y3-y1) + square3*(y1-y2)) / (2 * determinantA)
	yc := (square1*(x3-x2) + square2*(x1-x3) + square3*(x2-x1)) / (2 * determinantA)

	radius := math.Hypot(x1-xc, y1-yc)

	return radius * 2, nil
}

func main() {
	var filename string
	fmt.Println("Введите название файлы: ")
	fmt.Scanln(&filename)

	points, err := loadPointFromJson(filename + ".json")

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Загружено точек: %d\n", len(points))
	for i, p := range points {
		fmt.Printf("Точка %d: X=%.2f, Y=%.2f\n", i+1, p.X, p.Y)
	}
}
