package optimizer

import (
	"net/netip"
	"slices"
	"testing"

	"github.com/Ranger0806/go-routes/internal/route"
)

func TestRemoveCoveredRemovesNestedNetwork(t *testing.T) {
	wideRoute := route.NewRoute(
		netip.MustParsePrefix("216.239.0.0/16"),
		"youtube",
	)

	nestedRoute := route.NewRoute(
		netip.MustParsePrefix("216.239.32.0/19"),
		"google",
	)

	result := RemoveCovered([]*route.Route{
		wideRoute,
		nestedRoute,
	})

	if len(result) != 1 {
		t.Fatalf("expected 1 route, got %d", len(result))
	}

	expectedNetwork := netip.MustParsePrefix("216.239.0.0/16")

	if result[0].Network != expectedNetwork {
		t.Errorf(
			"expected network %s, got %s",
			expectedNetwork,
			result[0].Network,
		)
	}

	expectedDescriptions := []string{"youtube", "google"}

	if !slices.Equal(result[0].Descriptions, expectedDescriptions) {
		t.Errorf(
			"expected descriptions %v, got %v",
			expectedDescriptions,
			result[0].Descriptions,
		)
	}
}

func TestRemoveCoveredWorksWhenNestedNetworkComesFirst(t *testing.T) {
	nestedRoute := route.NewRoute(
		netip.MustParsePrefix("216.239.32.0/19"),
		"google",
	)

	wideRoute := route.NewRoute(
		netip.MustParsePrefix("216.239.0.0/16"),
		"youtube",
	)

	result := RemoveCovered([]*route.Route{
		nestedRoute,
		wideRoute,
	})

	if len(result) != 1 {
		t.Fatalf("expected 1 route, got %d", len(result))
	}

	expectedNetwork := netip.MustParsePrefix("216.239.0.0/16")

	if result[0].Network != expectedNetwork {
		t.Errorf(
			"expected network %s, got %s",
			expectedNetwork,
			result[0].Network,
		)
	}

	expectedDescriptions := []string{"youtube", "google"}

	if !slices.Equal(result[0].Descriptions, expectedDescriptions) {
		t.Errorf(
			"expected descriptions %v, got %v",
			expectedDescriptions,
			result[0].Descriptions,
		)
	}
}

func TestRemoveCoveredCollapsesNetworkChain(t *testing.T) {
	routes := []*route.Route{
		route.NewRoute(
			netip.MustParsePrefix("10.20.30.0/24"),
			"service",
		),
		route.NewRoute(
			netip.MustParsePrefix("10.20.0.0/16"),
			"platform",
		),
		route.NewRoute(
			netip.MustParsePrefix("10.0.0.0/8"),
			"company",
		),
	}

	result := RemoveCovered(routes)

	if len(result) != 1 {
		t.Fatalf("expected 1 route, got %d", len(result))
	}

	expectedNetwork := netip.MustParsePrefix("10.0.0.0/8")

	if result[0].Network != expectedNetwork {
		t.Errorf(
			"expected network %s, got %s",
			expectedNetwork,
			result[0].Network,
		)
	}

	expectedDescriptions := []string{
		"company",
		"platform",
		"service",
	}

	if !slices.Equal(result[0].Descriptions, expectedDescriptions) {
		t.Errorf(
			"expected descriptions %v, got %v",
			expectedDescriptions,
			result[0].Descriptions,
		)
	}
}

func TestRemoveCoveredKeepsUnrelatedNetworks(t *testing.T) {
	routes := []*route.Route{
		route.NewRoute(
			netip.MustParsePrefix("216.239.0.0/16"),
			"youtube",
		),
		route.NewRoute(
			netip.MustParsePrefix("149.154.160.0/20"),
			"telegram",
		),
	}

	result := RemoveCovered(routes)

	if len(result) != 2 {
		t.Fatalf("expected 2 routes, got %d", len(result))
	}
}

func TestRemoveCoveredDoesNotMergeAdjacentNetworks(t *testing.T) {
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

	result := RemoveCovered(routes)

	if len(result) != 2 {
		t.Fatalf(
			"expected adjacent networks to remain separate, got %d routes",
			len(result),
		)
	}
}

func TestRemoveCoveredAlsoDeduplicatesExactNetworks(t *testing.T) {
	routes := []*route.Route{
		route.NewRoute(
			netip.MustParsePrefix("216.239.0.0/16"),
			"youtube",
		),
		route.NewRoute(
			netip.MustParsePrefix("216.239.0.0/16"),
			"google",
		),
		route.NewRoute(
			netip.MustParsePrefix("216.239.32.0/19"),
			"video",
		),
	}

	result := RemoveCovered(routes)

	if len(result) != 1 {
		t.Fatalf("expected 1 route, got %d", len(result))
	}

	expectedDescriptions := []string{
		"youtube",
		"google",
		"video",
	}

	if !slices.Equal(result[0].Descriptions, expectedDescriptions) {
		t.Errorf(
			"expected descriptions %v, got %v",
			expectedDescriptions,
			result[0].Descriptions,
		)
	}
}

func TestRemoveCoveredDoesNotModifyInput(t *testing.T) {
	wideRoute := route.NewRoute(
		netip.MustParsePrefix("216.239.0.0/16"),
		"youtube",
	)

	nestedRoute := route.NewRoute(
		netip.MustParsePrefix("216.239.32.0/19"),
		"google",
	)

	RemoveCovered([]*route.Route{
		wideRoute,
		nestedRoute,
	})

	if !slices.Equal(wideRoute.Descriptions, []string{"youtube"}) {
		t.Errorf(
			"wide input route was modified: %v",
			wideRoute.Descriptions,
		)
	}

	if !slices.Equal(nestedRoute.Descriptions, []string{"google"}) {
		t.Errorf(
			"nested input route was modified: %v",
			nestedRoute.Descriptions,
		)
	}
}
