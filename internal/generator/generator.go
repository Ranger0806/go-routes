package generator

import (
	"fmt"
	"net"
	"os"
	"sort"
	"strings"

	"github.com/Ranger0806/go-routes/internal/route"
)

const windowsLineBreak = "\r\n"

func FormatRoute(currentRoute *route.Route) (string, error) {
	if currentRoute == nil {
		return "", fmt.Errorf("route is nil")
	}

	network := currentRoute.Network.Masked()

	if !network.IsValid() {
		return "", fmt.Errorf("route contains invalid network")
	}

	if !network.Addr().Is4() {
		return "", fmt.Errorf(
			"network %q is not IPv4",
			network,
		)
	}

	mask := net.CIDRMask(network.Bits(), 32)
	if mask == nil {
		return "", fmt.Errorf(
			"cannot create mask for network %q",
			network,
		)
	}

	maskString := net.IP(mask).String()

	command := fmt.Sprintf(
		"route add %s mask %s 0.0.0.0",
		network.Addr(),
		maskString,
	)

	descriptions := make([]string, 0, len(currentRoute.Descriptions))

	for _, description := range currentRoute.Descriptions {
		description = sanitizeDescription(description)

		if description == "" {
			continue
		}

		descriptions = append(descriptions, description)
	}

	if len(descriptions) == 0 {
		return command, nil
	}

	return command + " & rem " + strings.Join(descriptions, ", "), nil
}

// WriteBAT создаёт итоговый BAT-файл со всеми маршрутами.
func WriteBAT(path string, routes []*route.Route) error {
	if strings.TrimSpace(path) == "" {
		return fmt.Errorf("output path is empty")
	}

	sortedRoutes := make([]*route.Route, 0, len(routes))

	for index, currentRoute := range routes {
		if currentRoute == nil {
			return fmt.Errorf("route at index %d is nil", index)
		}

		network := currentRoute.Network.Masked()

		if !network.IsValid() {
			return fmt.Errorf(
				"route at index %d contains invalid network",
				index,
			)
		}

		if !network.Addr().Is4() {
			return fmt.Errorf(
				"route at index %d is not IPv4: %s",
				index,
				network,
			)
		}

		sortedRoutes = append(sortedRoutes, currentRoute)
	}

	sort.SliceStable(sortedRoutes, func(i, j int) bool {
		left := sortedRoutes[i].Network.Masked()
		right := sortedRoutes[j].Network.Masked()

		if left.Addr() != right.Addr() {
			return left.Addr().Less(right.Addr())
		}

		return left.Bits() < right.Bits()
	})

	lines := make([]string, 0, len(sortedRoutes))

	for index, currentRoute := range sortedRoutes {
		line, err := FormatRoute(currentRoute)
		if err != nil {
			return fmt.Errorf(
				"format route at index %d: %w",
				index,
				err,
			)
		}

		lines = append(lines, line)
	}

	content := strings.Join(lines, windowsLineBreak)

	// Добавляем перенос после последней строки.
	if len(lines) > 0 {
		content += windowsLineBreak
	}

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf(
			"write BAT file %q: %w",
			path,
			err,
		)
	}

	return nil
}

// sanitizeDescription убирает символы, которые могут сломать BAT-команду.
func sanitizeDescription(description string) string {
	replacer := strings.NewReplacer(
		"\r", " ",
		"\n", " ",
		"&", " ",
		"|", " ",
		"<", " ",
		">", " ",
		"^", " ",
	)

	description = replacer.Replace(description)

	return strings.Join(strings.Fields(description), " ")
}
