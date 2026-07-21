package app

import (
	"fmt"
	"strings"

	"github.com/Ranger0806/go-routes/internal/generator"
	"github.com/Ranger0806/go-routes/internal/optimizer"
	"github.com/Ranger0806/go-routes/internal/parser"
)

type Stats struct {
	ParsedRoutes    int
	OptimizedRoutes int
	RemovedRoutes   int
}

func Build(inputDirectory string, outputPath string) (Stats, error) {
	var stats Stats

	if strings.TrimSpace(inputDirectory) == "" {
		return stats, fmt.Errorf("input directory is empty")
	}

	if strings.TrimSpace(outputPath) == "" {
		return stats, fmt.Errorf("output path is empty")
	}

	routes, err := parser.ParseDirectoryExcept(
		inputDirectory,
		outputPath,
	)
	if err != nil {
		return stats, fmt.Errorf(
			"parse input directory: %w",
			err,
		)
	}

	stats.ParsedRoutes = len(routes)

	if len(routes) == 0 {
		return stats, fmt.Errorf(
			"no routes found in directory %q",
			inputDirectory,
		)
	}

	optimizedRoutes := optimizer.MergeAdjacent(routes)

	stats.OptimizedRoutes = len(optimizedRoutes)
	stats.RemovedRoutes =
		stats.ParsedRoutes - stats.OptimizedRoutes

	if err := generator.WriteBAT(
		outputPath,
		optimizedRoutes,
	); err != nil {
		return stats, fmt.Errorf(
			"generate output BAT file: %w",
			err,
		)
	}

	return stats, nil
}
