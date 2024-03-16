package domain

type Charts struct {
	HailStones [][]float64 `json:"hailStones"`
	Logarithm  [][]float64 `json:"logarithm"`
	Histogram  [][]float64 `json:"histogram"`
}

type Response struct {
	HailStones []float64 `json:"hailStones"`
	Charts     `json:"charts"`
}
