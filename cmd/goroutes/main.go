package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Ranger0806/go-routes/internal/app"
)

const version = "0.1.0"

func main() {
	os.Exit(run())
}

func run() int {
	inputDirectory := flag.String(
		"input",
		"",
		"directory containing source BAT files",
	)

	outputPath := flag.String(
		"output",
		"routes.bat",
		"path to generated BAT file",
	)

	showVersion := flag.Bool(
		"version",
		false,
		"show application version",
	)

	flag.Usage = printUsage
	flag.Parse()

	if *showVersion {
		fmt.Printf("go-routes %s\n", version)
		return 0
	}

	if *inputDirectory == "" {
		fmt.Fprintln(
			os.Stderr,
			"error: --input is required",
		)
		fmt.Fprintln(os.Stderr)

		printUsage()

		return 2
	}

	absoluteOutputPath, err := filepath.Abs(*outputPath)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"error: resolve output path: %v\n",
			err,
		)

		return 1
	}

	stats, err := app.Build(
		*inputDirectory,
		absoluteOutputPath,
	)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"error: %v\n",
			err,
		)

		return 1
	}

	fmt.Println("Go Routes")
	fmt.Println()
	fmt.Printf(
		"Parsed routes:    %d\n",
		stats.ParsedRoutes,
	)
	fmt.Printf(
		"Optimized routes: %d\n",
		stats.OptimizedRoutes,
	)
	fmt.Printf(
		"Removed routes:   %d\n",
		stats.RemovedRoutes,
	)
	fmt.Println()
	fmt.Printf(
		"Created: %s\n",
		absoluteOutputPath,
	)

	return 0
}

func printUsage() {
	writer := flag.CommandLine.Output()

	fmt.Fprintln(
		writer,
		"Usage:",
	)

	fmt.Fprintln(
		writer,
		"  goroutes --input <directory> [--output <file>]",
	)

	fmt.Fprintln(writer)
	fmt.Fprintln(writer, "Options:")

	flag.PrintDefaults()
}
