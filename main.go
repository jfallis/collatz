package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/jfallis/collatz/pkg/calculation"
	"github.com/jfallis/collatz/pkg/collatz"
	"github.com/jfallis/collatz/pkg/collatz/extension"
	"github.com/jfallis/collatz/pkg/domain"

	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

const (
	errMargin     = 15
	cpuMultiplier = 100
)

func main() {
	handler := pterm.NewSlogHandler(&pterm.DefaultLogger)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	_ = pterm.DefaultBigText.WithLetters(
		putils.LettersFromStringWithStyle("Collatz", pterm.FgGreen.ToStyle()),
	).Render()

	pterm.DefaultHeader.Println("The Collatz Conjecture is the simplest math problem no one can solve \n",
		"- it is easy enough for almost anyone to understand but notoriously difficult to solve.")

	if len(os.Args) != 3 {
		printErrMsg("invalid command; example: ./collatz seed 9663 or ./collatz bruteforce 100000")
		return
	}

	s, argErr := strconv.ParseUint(os.Args[2], calculation.Base, calculation.Bit)
	if argErr != nil {
		printErrMsg(argErr.Error())
		return
	}

	switch os.Args[1] {
	case "bruteforce":
		runtime.GOMAXPROCS(runtime.NumCPU() * cpuMultiplier)
		if bErr := extension.Bruteforce(s, extension.DefaultBruteforceRoutineBatchLimit, true); bErr != nil {
			printErrMsg(bErr.Error())
			return
		}

		pterm.DefaultBasicText.Println("The" + pterm.LightMagenta(" cake ") + "is a lie.")

		return
	case "seed":
		collatzConjecture(s)
	}
}

func collatzConjecture(n uint64) {
	c := collatz.New(n)
	if err := c.Calculate(); err != nil {
		printErrMsg(err.Error())
		pterm.Println()

		return
	}

	bulletListItems := []pterm.BulletListItem{
		{Level: 0, Text: pterm.LightMagenta("Total steps:") + fmt.Sprintf(" %d", len(c.Steps()))},
		{Level: 0, Text: pterm.LightMagenta("Collatz sequence:") + fmt.Sprintf(" %s", func(steps []uint64) string {
			s := make([]string, len(steps))
			for x, step := range steps {
				s[x] = strconv.FormatUint(step, 10)
			}

			return strings.Join(s, ", ")
		}(c.Steps()))},
		{Level: 0, Text: pterm.LightMagenta("Success:") + fmt.Sprintf(" %t", c.Success())},
	}

	_ = pterm.DefaultBulletList.WithItems(bulletListItems).Render()
	pterm.Println()

	resp := domain.Response{
		HailStones: c.Steps(),
	}

	pterm.DefaultSection.Println("Graphs")

	pterm.DefaultSection.WithLevel(2).Println("Hail stones")
	pterm.DefaultBasicText.Println(buildCharts[uint64](resp.HailStones))

	jsonDump, err := json.Marshal(&resp)
	if err != nil {
		printErrMsg(err.Error())
		return
	}

	pterm.Println()
	pterm.Printf("JSON data blob: %s\n", jsonDump)
}

func buildCharts[V uint64 | float64](data []V) string {
	bars := make([]pterm.Bar, len(data))
	for i, p := range data {
		bars[i] = pterm.Bar{
			Label: fmt.Sprintf("%d: %v", i+1, p),
			Value: int(p),
		}
	}

	barChart, _ := pterm.DefaultBarChart.WithHorizontal().WithBars(bars).Srender()

	return barChart
}

func printErrMsg(str string) {
	pterm.DefaultHeader.WithMargin(errMargin).
		WithBackgroundStyle(pterm.NewStyle(pterm.BgRed)).
		WithTextStyle(pterm.NewStyle(pterm.FgLightWhite)).
		Println(str)
}
