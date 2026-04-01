package main

import (
	"encoding/json"
	"image/color"
	"math"
	"math/rand"
	"os"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

type Circle struct {
	X, R float64
}

type CircleData struct {
	CenterX float64 `json:"center_x"`
	Radius  float64 `json:"radius"`
	Points  []Point `json:"top_points"`
}

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func generateCircles(n int, minX, maxX, minR, maxR float64) []Circle {
	rand.Seed(time.Now().UnixNano())
	var circles []Circle
	attempts := 0
	const maxAttempts = 20000

	for len(circles) < n && attempts < maxAttempts {
		attempts++
		r := minR + rand.Float64()*(maxR-minR)
		x := minX + rand.Float64()*(maxX-minX)

		ok := true
		for _, c := range circles {
			dx := math.Abs(x - c.X)
			minDist := 2 * math.Sqrt(r*c.R)
			if dx < minDist-1e-6 {
				ok = false
				break
			}
		}
		if ok {
			circles = append(circles, Circle{X: x, R: r})
		}
	}
	return circles
}

func circleToXYClosed(c Circle, points int) plotter.XYs {
	xy := make(plotter.XYs, points+1)
	for i := 0; i < points; i++ {
		t := 2 * math.Pi * float64(i) / float64(points)
		xy[i].X = c.X + c.R*math.Cos(t)
		xy[i].Y = c.R + c.R*math.Sin(t)
	}
	xy[points] = xy[0]
	return xy
}

func generateTopPoints(c Circle, count int) []Point {
	points := make([]Point, count)
	for i := 0; i < count; i++ {
		t := rand.Float64() * math.Pi
		points[i].X = c.X + c.R*math.Cos(t)
		points[i].Y = c.R + c.R*math.Sin(t)
	}
	return points
}

func exportToJSON(circles []Circle, filename string) error {
	var data []CircleData
	for _, c := range circles {
		cd := CircleData{
			CenterX: c.X,
			Radius:  c.R,
			Points:  generateTopPoints(c, 3),
		}
		data = append(data, cd)
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, jsonData, 0644)
}

func main() {
	const (
		count = 8
		minX  = 0.0
		maxX  = 100.0
		minR  = 3.0
		maxR  = 12.0
	)

	circles := generateCircles(count, minX, maxX, minR, maxR)

	exportToJSON(circles, "circles.json")

	yMax := 0.0
	for _, c := range circles {
		if y := 2 * c.R; y > yMax {
			yMax = y
		}
	}
	xMin, xMax := minX-5, maxX+5
	yMin, yMax := -1.0, yMax+5

	dataWidth := xMax - xMin
	dataHeight := yMax - yMin
	aspectRatio := dataWidth / dataHeight

	baseHeight := 10 * vg.Centimeter
	width := baseHeight * vg.Length(aspectRatio)

	p := plot.New()
	p.Title.Text = "Окружности"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	p.Add(plotter.NewGrid())

	p.X.Min, p.X.Max = xMin, xMax
	p.Y.Min, p.Y.Max = yMin, yMax

	colors := []color.Color{
		color.RGBA{R: 255, G: 100, A: 255},
		color.RGBA{G: 255, B: 50, A: 255},
		color.RGBA{R: 50, G: 150, B: 255, A: 255},
		color.RGBA{R: 200, G: 50, B: 200, A: 255},
		color.RGBA{R: 50, G: 200, B: 100, A: 255},
	}

	for i, c := range circles {
		xy := circleToXYClosed(c, 120)
		line, err := plotter.NewLine(xy)
		if err != nil {
			panic(err)
		}
		line.Color = colors[i%len(colors)]
		line.Width = vg.Points(2.0)
		p.Add(line)

		topPoints := generateTopPoints(c, 3)
		xyPoints := make(plotter.XYs, len(topPoints))
		for j, tp := range topPoints {
			xyPoints[j].X = tp.X
			xyPoints[j].Y = tp.Y
		}

		scatter, err := plotter.NewScatter(xyPoints)
		if err != nil {
			panic(err)
		}
		scatter.Color = color.RGBA{R: 255, A: 255}
		scatter.GlyphStyle.Radius = vg.Points(5)
		p.Add(scatter)
	}

	p.Save(width, baseHeight, "circles_with_points.png")
}
