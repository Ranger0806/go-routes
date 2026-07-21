package optimizer

import (
	"net/netip"
	"sort"

	"github.com/Ranger0806/go-routes/internal/route"
)

func RemoveCovered(routes []*route.Route) []*route.Route {
	deduplicated := Deduplicate(routes)

	sort.SliceStable(deduplicated, func(i, j int) bool {
		left := deduplicated[i].Network
		right := deduplicated[j].Network

		if left.Bits() != right.Bits() {
			return left.Bits() < right.Bits()
		}

		return left.Addr().Less(right.Addr())
	})

	result := make([]*route.Route, 0, len(deduplicated))

	for _, currentRoute := range deduplicated {
		coveringRoute := findCoveringRoute(result, currentRoute.Network)

		if coveringRoute != nil {
			for _, description := range currentRoute.Descriptions {
				coveringRoute.AddDescription(description)
			}

			continue
		}

		result = append(result, currentRoute)
	}

	return result
}

func findCoveringRoute(
	routes []*route.Route,
	network netip.Prefix,
) *route.Route {
	for _, currentRoute := range routes {
		if prefixContainsPrefix(currentRoute.Network, network) {
			return currentRoute
		}
	}

	return nil
}

func prefixContainsPrefix(parent, child netip.Prefix) bool {
	if parent.Addr().BitLen() != child.Addr().BitLen() {
		return false
	}

	if parent.Bits() > child.Bits() {
		return false
	}

	return parent.Contains(child.Addr())
}
