package optimizer

import (
	"net/netip"

	"github.com/Ranger0806/go-routes/internal/route"
)

func Deduplicate(routes []*route.Route) []*route.Route {
	result := make([]*route.Route, 0)
	seen := make(map[netip.Prefix]*route.Route)

	for _, currentRoute := range routes {
		if currentRoute == nil {
			continue
		}

		network := currentRoute.Network.Masked()

		existingRoute, exists := seen[network]
		if exists {
			for _, description := range currentRoute.Descriptions {
				existingRoute.AddDescription(description)
			}

			continue
		}

		deduplicatedRoute := route.NewRoute(network, "")

		for _, description := range currentRoute.Descriptions {
			deduplicatedRoute.AddDescription(description)
		}

		seen[network] = deduplicatedRoute
		result = append(result, deduplicatedRoute)
	}

	return result
}
