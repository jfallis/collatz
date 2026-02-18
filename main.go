// Package main provides the entry point for the Collatz Conjecture CLI,
// supporting seed and bruteforce modes for sequence calculation and analysis.
package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/jfallis/collatz/pkg/collatz"
	"github.com/jfallis/collatz/pkg/collatz/extension"
	"github.com/jfallis/collatz/pkg/collatz/extension/bruteforce"
)

const argCount = 2

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == "time" {
				return slog.Attr{}
			}
			return a
		},
	})))

	var (
		begin, end, size         string
		steps, logging, printAll bool
	)
	bruteforceFlags := flag.NewFlagSet("bruteforce", flag.ExitOnError)
	bruteforceFlags.StringVar(&begin, "begin", "0", "define the start number")
	bruteforceFlags.StringVar(&end, "end", "", "define the end number")
	bruteforceFlags.StringVar(&size, "size", extension.CPUBatchSize(), "set the batch size")
	bruteforceFlags.BoolVar(&printAll, "print-all", false, "print all steps, recommend logging to be enabled")
	bruteforceFlags.BoolVar(&logging, "logging", true, "enable or disable logging")
	bruteforceFlags.BoolVar(&steps, "steps", true,
		"enable or disable step collection, disable for performance improvement")

	var number string
	seedFlags := flag.NewFlagSet("seed", flag.ExitOnError)
	seedFlags.StringVar(&number, "number", "", "define seed number")
	seedFlags.BoolVar(&steps, "steps", true, "enable or disable step collection")

	if len(os.Args) < argCount {
		printUsage(nil)
		return
	}

	switch os.Args[1] {
	case "bruteforce":
		if err := bruteforceFlags.Parse(os.Args[argCount:]); err != nil || end == "" {
			printUsage(bruteforceFlags)
			return
		}

		ctx := context.Background()
		if _, err := bruteforce.Run(ctx, bruteforce.Request{
			Start:       begin,
			End:         end,
			BatchSize:   size,
			Logger:      slog.Default(),
			EnableSteps: steps,
			PrintAll:    printAll,
		}); err != nil {
			slog.Error(err.Error())
		}
	case "seed":
		if err := seedFlags.Parse(os.Args[argCount:]); err != nil || number == "" {
			printUsage(seedFlags)
			return
		}

		collatzConjecture(number, steps)
	}
}

func printUsage(flagSet *flag.FlagSet) {
	slog.Info("Collatz Conjecture")
	slog.Info("The Collatz Conjecture is the simplest math problem no one can solve " +
		"- it is easy enough for almost anyone to understand but notoriously difficult to solve.")
	slog.Info("Usage:")
	slog.Info("  collatz seed -number=9663")
	slog.Info("  collatz bruteforce -start=0 -end=100000")
	if flagSet != nil {
		slog.Info("Options:")
		flagSet.PrintDefaults()
	}
	slog.Info("Help:")
	slog.Info("  collatz seed [--help | -h]")
	slog.Info("  collatz bruteforce [--help | -h]")
}

func collatzConjecture(n string, steps bool) {
	cal := collatz.New(n)
	if err := cal.Calculate(steps); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	slog.Info(cal.String())
	if steps {
		slog.Info(fmt.Sprintf("Collatz sequence: %s", cal.Steps()))
	}
}
