package main

import (
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"os"
	"runtime"
	"strings"

	"github.com/jfallis/collatz/pkg/math"

	"github.com/jfallis/collatz/pkg/collatz"
	"github.com/jfallis/collatz/pkg/collatz/extension"

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

	if len(os.Args) < 3 {
		printErrMsg("invalid command; example: ./collatz seed 9663 or ./collatz bruteforce 0 100000")
		return
	}

	sArg2 := new(big.Int)
	if _, ok := sArg2.SetString(os.Args[2], math.Base); !ok {
		printErrMsg("Error: Failed to convert string to big.Int")
		return
	}

	switch os.Args[1] {
	case "bruteforce":
		runtime.GOMAXPROCS(runtime.NumCPU() * cpuMultiplier)

		sArg3 := new(big.Int)
		if _, ok := sArg3.SetString(os.Args[3], math.Base); !ok {
			printErrMsg("Error: Failed to convert string to big.Int")
			return
		}

		if _, err := extension.Bruteforce(sArg2, sArg3, extension.DefaultBruteforceRoutineBatchLimit, true); err != nil {
			var success collatz.SuccessError
			if errors.As(err, &success) {
				pterm.DefaultBasicText.Printf("The %s is a lie.\n", pterm.LightMagenta("cake"))
			}

			printErrMsg(err.Error())
			return
		}

		return
	case "seed":
		collatzConjecture(sArg2)
	}
}

func collatzConjecture(n *big.Int) {
	c := collatz.New(n)
	if err := c.Calculate(); err != nil {
		printErrMsg(err.Error())
		pterm.Println()

		return
	}

	bulletListItems := []pterm.BulletListItem{
		{Level: 0, Text: fmt.Sprintf("%s %d", pterm.LightMagenta("Total steps:"), len(c.Steps()))},
		{Level: 0, Text: fmt.Sprintf("%s %s", pterm.LightMagenta("Collatz sequence:"), func(steps []*big.Int) string {
			s := make([]string, len(steps))
			for x, step := range steps {
				s[x] = step.String()
			}

			return strings.Join(s, ", ")
		}(c.Steps()))},
		{Level: 0, Text: fmt.Sprintf("%s %t", pterm.LightMagenta("Success:"), c.Success())},
	}

	_ = pterm.DefaultBulletList.WithItems(bulletListItems).Render()
	pterm.Println()
}

func printErrMsg(str string) {
	pterm.DefaultHeader.WithMargin(errMargin).
		WithBackgroundStyle(pterm.NewStyle(pterm.BgRed)).
		WithTextStyle(pterm.NewStyle(pterm.FgLightWhite)).
		Println(str)
}
