package optimizer

import (
	"net/netip"
	"slices"
	"testing"

	"github.com/Ranger0806/go-routes/internal/route"
)

func TestDeduplicateMergesDescriptions(t *testing.T) {
	first := route.NewRoute(
		netip.MustParsePrefix("216.239.0.0/16"),
		"youtube",
	)

	second := route.NewRoute(
		netip.MustParsePrefix("216.239.0.0/16"),
		"google",
	)

	third := route.NewRoute(
		netip.MustParsePrefix("216.239.0.0/16"),
		"youtube",
	)

	result := Deduplicate([]*route.Route{
		first,
		second,
		third,
	})

	if len(result) != 1 {
		t.Fatalf("expected 1 route, got %d", len(result))
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

func TestDeduplicateKeepsDifferentNetworks(t *testing.T) {
	first := route.NewRoute(
		netip.MustParsePrefix("216.239.0.0/16"),
		"youtube",
	)

	second := route.NewRoute(
		netip.MustParsePrefix("142.250.0.0/16"),
		"google",
	)

	result := Deduplicate([]*route.Route{
		first,
		second,
	})

	if len(result) != 2 {
		t.Fatalf("expected 2 routes, got %d", len(result))
	}

	expectedNetworks := []netip.Prefix{
		netip.MustParsePrefix("216.239.0.0/16"),
		netip.MustParsePrefix("142.250.0.0/16"),
	}

	for index, expectedNetwork := range expectedNetworks {
		if result[index].Network != expectedNetwork {
			t.Errorf(
				"route %d: expected network %s, got %s",
				index,
				expectedNetwork,
				result[index].Network,
			)
		}
	}
}

func TestDeduplicatePreservesFirstAppearanceOrder(t *testing.T) {
	routes := []*route.Route{
		route.NewRoute(
			netip.MustParsePrefix("149.154.160.0/20"),
			"telegram",
		),
		route.NewRoute(
			netip.MustParsePrefix("216.239.0.0/16"),
			"youtube",
		),
		route.NewRoute(
			netip.MustParsePrefix("149.154.160.0/20"),
			"proxy",
		),
	}

	result := Deduplicate(routes)

	if len(result) != 2 {
		t.Fatalf("expected 2 routes, got %d", len(result))
	}

	if result[0].Network != netip.MustParsePrefix("149.154.160.0/20") {
		t.Errorf(
			"expected first network to be telegram network, got %s",
			result[0].Network,
		)
	}

	if result[1].Network != netip.MustParsePrefix("216.239.0.0/16") {
		t.Errorf(
			"expected second network to be youtube network, got %s",
			result[1].Network,
		)
	}
}

func TestDeduplicateIgnoresNilRoutes(t *testing.T) {
	validRoute := route.NewRoute(
		netip.MustParsePrefix("216.239.0.0/16"),
		"youtube",
	)

	result := Deduplicate([]*route.Route{
		nil,
		validRoute,
		nil,
	})

	if len(result) != 1 {
		t.Fatalf("expected 1 route, got %d", len(result))
	}

	if result[0].Network != validRoute.Network {
		t.Errorf(
			"expected network %s, got %s",
			validRoute.Network,
			result[0].Network,
		)
	}
}

func TestDeduplicateDoesNotModifyInputRoutes(t *testing.T) {
	first := route.NewRoute(
		netip.MustParsePrefix("216.239.0.0/16"),
		"youtube",
	)

	second := route.NewRoute(
		netip.MustParsePrefix("216.239.0.0/16"),
		"google",
	)

	Deduplicate([]*route.Route{first, second})

	if !slices.Equal(first.Descriptions, []string{"youtube"}) {
		t.Errorf(
			"first input route was modified: %v",
			first.Descriptions,
		)
	}

	if !slices.Equal(second.Descriptions, []string{"google"}) {
		t.Errorf(
			"second input route was modified: %v",
			second.Descriptions,
		)
	}
}
