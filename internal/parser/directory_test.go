package parser

import (
	"net/netip"
	"os"
	"path/filepath"
	"testing"
)

func TestParseDirectory(t *testing.T) {
	rootDir := t.TempDir()
	nestedDir := filepath.Join(rootDir, "nested")

	if err := os.Mkdir(nestedDir, 0o700); err != nil {
		t.Fatalf("failed to create nested directory: %v", err)
	}

	youtubePath := filepath.Join(rootDir, "youtube.bat")
	telegramPath := filepath.Join(nestedDir, "telegram.BAT")
	ignoredPath := filepath.Join(rootDir, "readme.txt")

	youtubeContent := `
route add 216.239.0.0 mask 255.255.0.0 0.0.0.0
route add 142.250.0.0 mask 255.255.0.0 0.0.0.0
`

	telegramContent := `
route add 149.154.160.0 mask 255.255.240.0 0.0.0.0
`

	if err := os.WriteFile(youtubePath, []byte(youtubeContent), 0o600); err != nil {
		t.Fatalf("failed to create youtube file: %v", err)
	}

	if err := os.WriteFile(telegramPath, []byte(telegramContent), 0o600); err != nil {
		t.Fatalf("failed to create telegram file: %v", err)
	}

	if err := os.WriteFile(ignoredPath, []byte("not a BAT file"), 0o600); err != nil {
		t.Fatalf("failed to create ignored file: %v", err)
	}

	routes, err := ParseDirectory(rootDir)
	if err != nil {
		t.Fatalf("ParseDirectory returned error: %v", err)
	}

	if len(routes) != 3 {
		t.Fatalf("expected 3 routes, got %d", len(routes))
	}

	expectedNetworks := map[netip.Prefix]string{
		netip.MustParsePrefix("216.239.0.0/16"):   "youtube",
		netip.MustParsePrefix("142.250.0.0/16"):   "youtube",
		netip.MustParsePrefix("149.154.160.0/20"): "telegram",
	}

	for _, currentRoute := range routes {
		expectedDescription, exists := expectedNetworks[currentRoute.Network]
		if !exists {
			t.Errorf("unexpected network: %s", currentRoute.Network)
			continue
		}

		if len(currentRoute.Descriptions) != 1 {
			t.Errorf(
				"network %s: expected 1 description, got %d",
				currentRoute.Network,
				len(currentRoute.Descriptions),
			)
			continue
		}

		if currentRoute.Descriptions[0] != expectedDescription {
			t.Errorf(
				"network %s: expected description %q, got %q",
				currentRoute.Network,
				expectedDescription,
				currentRoute.Descriptions[0],
			)
		}
	}
}

func TestParseDirectoryReturnsErrorWhenDirectoryDoesNotExist(t *testing.T) {
	_, err := ParseDirectory("directory-that-does-not-exist")

	if err == nil {
		t.Fatal("expected directory walking error, got nil")
	}
}

func TestParseDirectoryExceptSkipsExcludedBATFile(t *testing.T) {
	rootDirectory := t.TempDir()

	sourcePath := filepath.Join(rootDirectory, "youtube.bat")
	outputPath := filepath.Join(rootDirectory, "result.bat")

	sourceContent := "route add 10.0.0.0 mask 255.255.255.0 0.0.0.0"
	outputContent := "route add 192.168.0.0 mask 255.255.255.0 0.0.0.0"

	if err := os.WriteFile(
		sourcePath,
		[]byte(sourceContent),
		0o600,
	); err != nil {
		t.Fatalf("create source BAT: %v", err)
	}

	if err := os.WriteFile(
		outputPath,
		[]byte(outputContent),
		0o600,
	); err != nil {
		t.Fatalf("create output BAT: %v", err)
	}

	routes, err := ParseDirectoryExcept(rootDirectory, outputPath)
	if err != nil {
		t.Fatalf("ParseDirectoryExcept returned error: %v", err)
	}

	if len(routes) != 1 {
		t.Fatalf("expected 1 route, got %d", len(routes))
	}

	expected := netip.MustParsePrefix("10.0.0.0/24")

	if routes[0].Network != expected {
		t.Errorf(
			"expected network %s, got %s",
			expected,
			routes[0].Network,
		)
	}
}
