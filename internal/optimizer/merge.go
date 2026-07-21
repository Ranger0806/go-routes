package optimizer

import (
	"net/netip"
	"sort"

	"github.com/Ranger0806/go-routes/internal/route"
)

func MergeAdjacent(routes []*route.Route) []*route.Route {
	current := RemoveCovered(routes)

	for {
		next, merged := mergeAdjacentOnce(current)

		if !merged {
			return next
		}

		current = next
	}
}

func mergeAdjacentOnce(routes []*route.Route) ([]*route.Route, bool) {
	sort.SliceStable(routes, func(i, j int) bool {
		left := routes[i].Network
		right := routes[j].Network

		if left.Addr() != right.Addr() {
			return left.Addr().Less(right.Addr())
		}

		return left.Bits() < right.Bits()
	})

	result := make([]*route.Route, 0, len(routes))
	mergedAny := false

	for i := 0; i < len(routes); {
		if i+1 < len(routes) {
			parent, canMerge := commonParent(
				routes[i].Network,
				routes[i+1].Network,
			)

			if canMerge {
				mergedRoute := route.NewRoute(parent, "")

				for _, description := range routes[i].Descriptions {
					mergedRoute.AddDescription(description)
				}

				for _, description := range routes[i+1].Descriptions {
					mergedRoute.AddDescription(description)
				}

				result = append(result, mergedRoute)

				mergedAny = true
				i += 2
				continue
			}
		}

		result = append(result, routes[i])
		i++
	}

	return result, mergedAny
}

func commonParent(left, right netip.Prefix) (netip.Prefix, bool) {
	left = left.Masked()
	right = right.Masked()

	if left.Addr().BitLen() != right.Addr().BitLen() {
		return netip.Prefix{}, false
	}

	if left.Bits() != right.Bits() {
		return netip.Prefix{}, false
	}

	if left.Bits() == 0 {
		return netip.Prefix{}, false
	}

	if left == right {
		return netip.Prefix{}, false
	}

	parentBits := left.Bits() - 1

	leftParent := netip.PrefixFrom(
		left.Addr(),
		parentBits,
	).Masked()

	rightParent := netip.PrefixFrom(
		right.Addr(),
		parentBits,
	).Masked()

	if leftParent != rightParent {
		return netip.Prefix{}, false
	}

	return leftParent, true
}
