package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildCreatesOptimizedBATFile(t *testing.T) {
	inputDirectory := t.TempDir()
	outputDirectory := t.TempDir()

	outputPath := filepath.Join(outputDirectory, "routes.bat")

	youtubeContent := strings.Join([]string{
		"route add 10.0.0.0 mask 255.255.255.128 0.0.0.0",
		"route add 10.0.0.128 mask 255.255.255.128 0.0.0.0",
	}, "\r\n")

	googleContent := strings.Join([]string{
		"route add 10.0.0.0 mask 255.255.255.128 0.0.0.0",
	}, "\r\n")

	youtubePath := filepath.Join(inputDirectory, "youtube.bat")
	googlePath := filepath.Join(inputDirectory, "google.bat")

	if err := os.WriteFile(
		youtubePath,
		[]byte(youtubeContent),
		0o600,
	); err != nil {
		t.Fatalf("failed to create youtube BAT file: %v", err)
	}

	if err := os.WriteFile(
		googlePath,
		[]byte(googleContent),
		0o600,
	); err != nil {
		t.Fatalf("failed to create google BAT file: %v", err)
	}

	stats, err := Build(inputDirectory, outputPath)
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	if stats.ParsedRoutes != 3 {
		t.Errorf(
			"expected 3 parsed routes, got %d",
			stats.ParsedRoutes,
		)
	}

	if stats.OptimizedRoutes != 1 {
		t.Errorf(
			"expected 1 optimized route, got %d",
			stats.OptimizedRoutes,
		)
	}

	if stats.RemovedRoutes != 2 {
		t.Errorf(
			"expected 2 removed routes, got %d",
			stats.RemovedRoutes,
		)
	}

	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read generated BAT file: %v", err)
	}

	expected := "route add 10.0.0.0 mask 255.255.255.0 0.0.0.0 & rem google, youtube\r\n"

	if string(content) != expected {
		t.Errorf(
			"unexpected generated content:\nexpected: %q\ngot:      %q",
			expected,
			string(content),
		)
	}
}

func TestBuildReturnsErrorForEmptyInputDirectory(t *testing.T) {
	_, err := Build("", "routes.bat")

	if err == nil {
		t.Fatal("expected error for empty input directory")
	}
}

func TestBuildReturnsErrorForEmptyOutputPath(t *testing.T) {
	_, err := Build(t.TempDir(), "")

	if err == nil {
		t.Fatal("expected error for empty output path")
	}
}

func TestBuildReturnsErrorWhenInputDirectoryDoesNotExist(t *testing.T) {
	outputPath := filepath.Join(t.TempDir(), "routes.bat")

	_, err := Build(
		"directory-that-does-not-exist",
		outputPath,
	)

	if err == nil {
		t.Fatal("expected error for nonexistent input directory")
	}
}

func TestBuildReturnsErrorWhenNoRoutesFound(t *testing.T) {
	inputDirectory := t.TempDir()
	outputPath := filepath.Join(t.TempDir(), "routes.bat")

	readmePath := filepath.Join(inputDirectory, "README.txt")

	if err := os.WriteFile(
		readmePath,
		[]byte("no routes here"),
		0o600,
	); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	_, err := Build(inputDirectory, outputPath)

	if err == nil {
		t.Fatal("expected error when no routes were found")
	}
}

func TestBuildIgnoresPreviousOutputFile(t *testing.T) {
	inputDirectory := t.TempDir()
	outputPath := filepath.Join(
		inputDirectory,
		"generated-routes.bat",
	)

	sourcePath := filepath.Join(
		inputDirectory,
		"youtube.bat",
	)

	sourceContent := strings.Join([]string{
		"route add 10.0.0.0 mask 255.255.255.128 0.0.0.0",
		"route add 10.0.0.128 mask 255.255.255.128 0.0.0.0",
	}, "\r\n")

	if err := os.WriteFile(
		sourcePath,
		[]byte(sourceContent),
		0o600,
	); err != nil {
		t.Fatalf("create source BAT: %v", err)
	}

	firstStats, err := Build(inputDirectory, outputPath)
	if err != nil {
		t.Fatalf("first Build returned error: %v", err)
	}

	secondStats, err := Build(inputDirectory, outputPath)
	if err != nil {
		t.Fatalf("second Build returned error: %v", err)
	}

	if firstStats != secondStats {
		t.Errorf(
			"expected equal stats, first=%+v second=%+v",
			firstStats,
			secondStats,
		)
	}

	if secondStats.ParsedRoutes != 2 {
		t.Errorf(
			"expected 2 parsed source routes, got %d",
			secondStats.ParsedRoutes,
		)
	}

	if secondStats.OptimizedRoutes != 1 {
		t.Errorf(
			"expected 1 optimized route, got %d",
			secondStats.OptimizedRoutes,
		)
	}
}
