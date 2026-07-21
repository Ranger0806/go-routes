package optimizer

import (
	"net/netip"
	"slices"
	"testing"

	"github.com/Ranger0806/go-routes/internal/route"
)

func TestMergeAdjacentMergesSiblingNetworks(t *testing.T) {
	routes := []*route.Route{
		route.NewRoute(
			netip.MustParsePrefix("10.0.0.0/25"),
			"first",
		),
		route.NewRoute(
			netip.MustParsePrefix("10.0.0.128/25"),
			"second",
		),
	}

	result := MergeAdjacent(routes)

	if len(result) != 1 {
		t.Fatalf("expected 1 route, got %d", len(result))
	}

	expectedNetwork := netip.MustParsePrefix("10.0.0.0/24")

	if result[0].Network != expectedNetwork {
		t.Errorf(
			"expected network %s, got %s",
			expectedNetwork,
			result[0].Network,
		)
	}

	expectedDescriptions := []string{"first", "second"}

	if !slices.Equal(result[0].Descriptions, expectedDescriptions) {
		t.Errorf(
			"expected descriptions %v, got %v",
			expectedDescriptions,
			result[0].Descriptions,
		)
	}
}

func TestMergeAdjacentMergesRepeatedly(t *testing.T) {
	routes := []*route.Route{
		route.NewRoute(
			netip.MustParsePrefix("10.0.0.0/26"),
			"part 1",
		),
		route.NewRoute(
			netip.MustParsePrefix("10.0.0.64/26"),
			"part 2",
		),
		route.NewRoute(
			netip.MustParsePrefix("10.0.0.128/26"),
			"part 3",
		),
		route.NewRoute(
			netip.MustParsePrefix("10.0.0.192/26"),
			"part 4",
		),
	}

	result := MergeAdjacent(routes)

	if len(result) != 1 {
		t.Fatalf("expected 1 route, got %d", len(result))
	}

	expectedNetwork := netip.MustParsePrefix("10.0.0.0/24")

	if result[0].Network != expectedNetwork {
		t.Errorf(
			"expected network %s, got %s",
			expectedNetwork,
			result[0].Network,
		)
	}

	expectedDescriptions := []string{
		"part 1",
		"part 2",
		"part 3",
		"part 4",
	}

	if !slices.Equal(result[0].Descriptions, expectedDescriptions) {
		t.Errorf(
			"expected descriptions %v, got %v",
			expectedDescriptions,
			result[0].Descriptions,
		)
	}
}

func TestMergeAdjacentKeepsNonSiblingNetworks(t *testing.T) {
	routes := []*route.Route{
		route.NewRoute(
			netip.MustParsePrefix("10.0.0.0/25"),
			"first",
		),
		route.NewRoute(
			netip.MustParsePrefix("10.0.1.0/25"),
			"second",
		),
	}

	result := MergeAdjacent(routes)

	if len(result) != 2 {
		t.Fatalf(
			"expected unrelated networks to remain separate, got %d routes",
			len(result),
		)
	}
}

func TestMergeAdjacentKeepsDifferentSizedNetworks(t *testing.T) {
	routes := []*route.Route{
		route.NewRoute(
			netip.MustParsePrefix("10.0.0.0/25"),
			"first",
		),
		route.NewRoute(
			netip.MustParsePrefix("10.0.0.128/26"),
			"second",
		),
	}

	result := MergeAdjacent(routes)

	if len(result) != 2 {
		t.Fatalf(
			"expected differently sized networks to remain separate, got %d",
			len(result),
		)
	}
}

func TestMergeAdjacentWorksWithUnsortedInput(t *testing.T) {
	routes := []*route.Route{
		route.NewRoute(
			netip.MustParsePrefix("10.0.0.128/25"),
			"second",
		),
		route.NewRoute(
			netip.MustParsePrefix("192.168.0.0/24"),
			"unrelated",
		),
		route.NewRoute(
			netip.MustParsePrefix("10.0.0.0/25"),
			"first",
		),
	}

	result := MergeAdjacent(routes)

	if len(result) != 2 {
		t.Fatalf("expected 2 routes, got %d", len(result))
	}

	expectedMergedNetwork := netip.MustParsePrefix("10.0.0.0/24")

	if result[0].Network != expectedMergedNetwork {
		t.Errorf(
			"expected first network %s, got %s",
			expectedMergedNetwork,
			result[0].Network,
		)
	}
}

func TestMergeAdjacentDoesNotModifyInput(t *testing.T) {
	first := route.NewRoute(
		netip.MustParsePrefix("10.0.0.0/25"),
		"first",
	)

	second := route.NewRoute(
		netip.MustParsePrefix("10.0.0.128/25"),
		"second",
	)

	MergeAdjacent([]*route.Route{first, second})

	if first.Network != netip.MustParsePrefix("10.0.0.0/25") {
		t.Errorf("first input network was modified: %s", first.Network)
	}

	if second.Network != netip.MustParsePrefix("10.0.0.128/25") {
		t.Errorf("second input network was modified: %s", second.Network)
	}

	if !slices.Equal(first.Descriptions, []string{"first"}) {
		t.Errorf(
			"first input descriptions were modified: %v",
			first.Descriptions,
		)
	}

	if !slices.Equal(second.Descriptions, []string{"second"}) {
		t.Errorf(
			"second input descriptions were modified: %v",
			second.Descriptions,
		)
	}
}
