package route

import (
	"net/netip"
	"testing"
)

func TestNewRouteMasksNetwork(t *testing.T) {
	network := netip.MustParsePrefix("192.168.1.15/24")

	route := NewRoute(network, "YouTube")

	expected := netip.MustParsePrefix("192.168.1.0/24")

	if route.Network != expected {
		t.Errorf(
			"expected network %s, got %s",
			expected,
			route.Network,
		)
	}
}

func TestNewRouteAddsDescription(t *testing.T) {
	network := netip.MustParsePrefix("192.168.1.0/24")

	route := NewRoute(network, "YouTube")

	if len(route.Descriptions) != 1 {
		t.Fatalf(
			"expected 1 description, got %d",
			len(route.Descriptions),
		)
	}

	if route.Descriptions[0] != "YouTube" {
		t.Errorf(
			"expected description YouTube, got %s",
			route.Descriptions[0],
		)
	}
}

func TestAddDescriptionDoesNotAddDuplicate(t *testing.T) {
	network := netip.MustParsePrefix("192.168.1.0/24")
	route := NewRoute(network, "YouTube")

	route.AddDescription("YouTube")

	if len(route.Descriptions) != 1 {
		t.Errorf(
			"expected 1 description, got %d",
			len(route.Descriptions),
		)
	}
}

func TestAddDescriptionIgnoresEmptyDescription(t *testing.T) {
	network := netip.MustParsePrefix("192.168.1.0/24")
	route := NewRoute(network, "")

	if len(route.Descriptions) != 0 {
		t.Errorf(
			"expected no descriptions, got %v",
			route.Descriptions,
		)
	}
}
