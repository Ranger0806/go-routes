package parser

import (
	"net/netip"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseBATFile(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "youtube_routes.bat")

	content := strings.Join([]string{
		"@echo off",
		"rem YouTube routes",
		"",
		"route add 216.239.15.20 mask 255.255.0.0 0.0.0.0",
		"ROUTE ADD 142.250.1.5 MASK 255.255.0.0 0.0.0.0 & rem old description",
	}, "\r\n")

	err := os.WriteFile(filePath, []byte(content), 0o600)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	routes, err := ParseBATFile(filePath)
	if err != nil {
		t.Fatalf("ParseBATFile returned error: %v", err)
	}

	if len(routes) != 2 {
		t.Fatalf("expected 2 routes, got %d", len(routes))
	}

	expectedNetworks := []netip.Prefix{
		netip.MustParsePrefix("216.239.0.0/16"),
		netip.MustParsePrefix("142.250.0.0/16"),
	}

	for i, expectedNetwork := range expectedNetworks {
		if routes[i].Network != expectedNetwork {
			t.Errorf(
				"route %d: expected network %s, got %s",
				i,
				expectedNetwork,
				routes[i].Network,
			)
		}

		if len(routes[i].Descriptions) != 1 {
			t.Fatalf(
				"route %d: expected 1 description, got %d",
				i,
				len(routes[i].Descriptions),
			)
		}

		if routes[i].Descriptions[0] != "youtube routes" {
			t.Errorf(
				"route %d: expected description %q, got %q",
				i,
				"youtube routes",
				routes[i].Descriptions[0],
			)
		}
	}
}

func TestParseBATFileReturnsErrorForInvalidRoute(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "broken.bat")

	content := strings.Join([]string{
		"route add 216.239.0.0 mask 255.255.0.0 0.0.0.0",
		"route add invalid-ip mask 255.255.0.0 0.0.0.0",
	}, "\n")

	err := os.WriteFile(filePath, []byte(content), 0o600)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	_, err = ParseBATFile(filePath)
	if err == nil {
		t.Fatal("expected parsing error, got nil")
	}

	if !strings.Contains(err.Error(), "line 2") {
		t.Errorf("expected error to contain line number, got: %v", err)
	}
}

func TestParseBATFileReturnsErrorForInvalidMask(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "broken_mask.bat")

	content := "route add 216.239.0.0 mask 255.0.255.0 0.0.0.0"

	err := os.WriteFile(filePath, []byte(content), 0o600)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	_, err = ParseBATFile(filePath)
	if err == nil {
		t.Fatal("expected invalid mask error, got nil")
	}
}

func TestParseBATFileReturnsErrorWhenFileDoesNotExist(t *testing.T) {
	_, err := ParseBATFile("file-that-does-not-exist.bat")

	if err == nil {
		t.Fatal("expected file opening error, got nil")
	}
}

func TestParseRouteLineIgnoresNonRouteLines(t *testing.T) {
	lines := []string{
		"",
		"rem comment",
		"@echo off",
		"echo loading routes",
		"route print",
		"route delete 216.239.0.0",
	}

	for _, line := range lines {
		parsedRoute, err := parseRouteLine(line, "test")
		if err != nil {
			t.Errorf("line %q returned unexpected error: %v", line, err)
		}

		if parsedRoute != nil {
			t.Errorf("line %q unexpectedly produced a route", line)
		}
	}
}
