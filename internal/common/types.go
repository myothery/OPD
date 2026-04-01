package common

type Point struct {
	X float64 `json:"X"`
	Y float64 `json:"Y"`
}

type Circle struct {
	X, R float64
}

type CircleData struct {
	CenterX float64 `json:"CENTER_X"`
	Radius  float64 `json:"RADIUS"`
	Points  []Point `json:"TOP_POINTS"`
}
