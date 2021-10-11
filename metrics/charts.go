package metrics

func BuildMetrics(c []float64) [][]float64 {
	resp := make([][]float64, len(c))
	for i, x := range c {
		resp[i] = []float64{float64(i), x}
	}

	return resp
}
