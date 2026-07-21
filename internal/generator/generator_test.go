package generator

import (
	"net/netip"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Ranger0806/go-routes/internal/route"
)

func TestFormatRoute(t *testing.T) {
	currentRoute := route.NewRoute(
		netip.MustParsePrefix("216.239.15.20/16"),
		"youtube",
	)
	currentRoute.AddDescription("google")

	actual, err := FormatRoute(currentRoute)
	if err != nil {
		t.Fatalf("FormatRoute returned error: %v", err)
	}

	expected := "route add 216.239.0.0 mask 255.255.0.0 0.0.0.0 & rem youtube, google"

	if actual != expected {
		t.Errorf(
			"expected %q, got %q",
			expected,
			actual,
		)
	}
}

func TestFormatRouteWithoutDescription(t *testing.T) {
	currentRoute := route.NewRoute(
		netip.MustParsePrefix("10.0.0.0/24"),
		"",
	)

	actual, err := FormatRoute(currentRoute)
	if err != nil {
		t.Fatalf("FormatRoute returned error: %v", err)
	}

	expected := "route add 10.0.0.0 mask 255.255.255.0 0.0.0.0"

	if actual != expected {
		t.Errorf(
			"expected %q, got %q",
			expected,
			actual,
		)
	}
}

func TestFormatRouteSanitizesDescription(t *testing.T) {
	currentRoute := route.NewRoute(
		netip.MustParsePrefix("10.0.0.0/24"),
		"youtube & echo broken",
	)

	actual, err := FormatRoute(currentRoute)
	if err != nil {
		t.Fatalf("FormatRoute returned error: %v", err)
	}

	expected := "route add 10.0.0.0 mask 255.255.255.0 0.0.0.0 & rem youtube echo broken"

	if actual != expected {
		t.Errorf(
			"expected %q, got %q",
			expected,
			actual,
		)
	}
}

func TestWriteBATSortsRoutesAndUsesCRLF(t *testing.T) {
	outputPath := filepath.Join(t.TempDir(), "routes.bat")

	routes := []*route.Route{
		route.NewRoute(
			netip.MustParsePrefix("216.239.0.0/16"),
			"youtube",
		),
		route.NewRoute(
			netip.MustParsePrefix("10.0.0.0/24"),
			"internal",
		),
		route.NewRoute(
			netip.MustParsePrefix("142.250.0.0/16"),
			"google",
		),
	}

	err := WriteBAT(outputPath, routes)
	if err != nil {
		t.Fatalf("WriteBAT returned error: %v", err)
	}

	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read generated BAT file: %v", err)
	}

	expectedLines := []string{
		"route add 10.0.0.0 mask 255.255.255.0 0.0.0.0 & rem internal",
		"route add 142.250.0.0 mask 255.255.0.0 0.0.0.0 & rem google",
		"route add 216.239.0.0 mask 255.255.0.0 0.0.0.0 & rem youtube",
	}

	expected := strings.Join(expectedLines, "\r\n") + "\r\n"

	if string(content) != expected {
		t.Errorf(
			"unexpected BAT content:\nexpected: %q\ngot:      %q",
			expected,
			string(content),
		)
	}
}

func TestFormatRouteReturnsErrorForNilRoute(t *testing.T) {
	_, err := FormatRoute(nil)

	if err == nil {
		t.Fatal("expected error for nil route")
	}
}

func TestWriteBATReturnsErrorForEmptyPath(t *testing.T) {
	err := WriteBAT("", nil)

	if err == nil {
		t.Fatal("expected error for empty output path")
	}
}
