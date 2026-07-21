package route

import "slices"

import "net/netip"

type Route struct {
	Network      netip.Prefix
	Descriptions []string
}

func NewRoute(network netip.Prefix, description string) *Route {
	route := &Route{Network: network.Masked()}
	route.AddDescription(description)
	return route
}

func (route *Route) HasDescription(description string) bool {
	return slices.Contains(route.Descriptions, description)
}

func (route *Route) AddDescription(description string) {
	if route.HasDescription(description) {
		return
	}
	if description == "" {
		return
	}
	route.Descriptions = append(route.Descriptions, description)
}
