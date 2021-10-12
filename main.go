package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/buger/goterm"
	"github.com/jfallis/collatz/calculation"
	"github.com/jfallis/collatz/domain"
	"github.com/jfallis/collatz/metrics"
)

const (
	headerBoxWidth   = 50
	headerLineHeight = 4
	headerPosition   = 5
	chartWidth       = 100
	chartHeight      = 20
)

func main() {
	goterm.Clear()
	introBox := goterm.NewBox(headerBoxWidth|goterm.PCT, headerLineHeight, 0)
	fmt.Fprintln(introBox,
		fmt.Sprint("The Collatz Conjecture is the simplest math problem no one can solve \n",
			"- it is easy enough for almost anyone to understand but notoriously difficult to solve."),
	)
	goterm.Println(goterm.MoveTo(introBox.String(), headerPosition|goterm.PCT, headerPosition|goterm.PCT))
	goterm.Flush()

	if len(os.Args) != 3 {
		panic(errors.New("invalid command; example: ./collatz seed 9663 or ./collatz bruteforce 100000"))
	}

	s, err := strconv.ParseUint(os.Args[2], calculation.Base, calculation.Bit)
	if err != nil {
		panic(err)
	}

	switch os.Args[1] {
	case "bruteforce":
		if num, err := calculation.Bruteforce(s); err != nil {
			goterm.Println(num)
			panic(err)
		}

		goterm.Println("The cake is a lie.")
		goterm.Flush()

		return
	case "seed":
		collatz(s)
	}
}

func collatz(n uint64) {
	c := calculation.Create(n)
	goterm.Printf("\nNumber value: %d, Number of steps: %d, Success: %t\n\n", n, c.Len(), calculation.Success(c.Results()))
	goterm.Flush()

	resp := domain.Response{
		HailStones: c.Results(),
		Charts: domain.Charts{
			Logarithm: metrics.BuildMetrics(func(x calculation.Statement) []float64 {
				y := make([]float64, x.Len())
				for i, z := range x.Results() {
					y[i] = math.Log(z)
				}

				return y
			}(c)),
			HailStones: metrics.BuildMetrics(c.Results()),
			Histogram: metrics.BuildMetrics(
				calculation.ConvertIntToFloat(calculation.CreateHistogram(
					calculation.ConvertFloatToInt(c.Results()), calculation.LeadingDigit),
				),
			),
		},
	}

	// Log
	createChart(resp.Charts.Logarithm, "Logarithm")

	// HailStones
	createChart(resp.Charts.HailStones, "Hail stones")

	// Histogram
	createChart(resp.Charts.Histogram, "Histogram")

	// json blob
	jsonDump, err := json.Marshal(&resp)
	if err != nil {
		panic(err)
	}

	goterm.Printf("JSON data blob: %s\n", jsonDump)
	goterm.Flush()
}

func createChart(values [][]float64, colName string) {
	chart := goterm.NewLineChart(chartWidth, chartHeight)
	data := new(goterm.DataTable)
	data.AddColumn("Iterations")
	data.AddColumn(colName)

	for _, num := range values {
		data.AddRow(num[0], num[1])
	}

	goterm.Println(chart.Draw(data))
	goterm.Flush()
}
