package parser

import (
	"bufio"
	"fmt"
	"net"
	"net/netip"
	"os"
	"strings"

	"github.com/Ranger0806/go-routes/internal/route"
)

func ParseBATFile(path string) ([]*route.Route, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file %q: %w", path, err)
	}
	defer file.Close()

	description := DescriptionFromFilename(path)
	routes := make([]*route.Route, 0)

	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++

		currentRoute, err := parseRouteLine(scanner.Text(), description)
		if err != nil {
			return nil, fmt.Errorf(
				"parse file %q at line %d: %w",
				path,
				lineNumber,
				err,
			)
		}

		if currentRoute == nil {
			continue
		}

		routes = append(routes, currentRoute)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read file %q: %w", path, err)
	}

	return routes, nil
}

func parseRouteLine(line string, description string) (*route.Route, error) {
	line = strings.TrimSpace(line)

	if line == "" {
		return nil, nil
	}

	fields := strings.Fields(line)

	if len(fields) == 0 {
		return nil, nil
	}

	command := strings.TrimPrefix(fields[0], "@")

	if !strings.EqualFold(command, "route") {
		return nil, nil
	}

	if len(fields) < 2 || !strings.EqualFold(fields[1], "add") {
		return nil, nil
	}

	if len(fields) < 6 {
		return nil, fmt.Errorf("invalid route command: not enough arguments")
	}

	if !strings.EqualFold(fields[3], "mask") {
		return nil, fmt.Errorf(
			"invalid route command: expected mask, got %q",
			fields[3],
		)
	}

	address, err := netip.ParseAddr(fields[2])
	if err != nil {
		return nil, fmt.Errorf(
			"invalid destination address %q: %w",
			fields[2],
			err,
		)
	}

	if !address.Is4() {
		return nil, fmt.Errorf(
			"destination address %q is not IPv4",
			fields[2],
		)
	}

	maskIP := net.ParseIP(fields[4])
	if maskIP == nil {
		return nil, fmt.Errorf("invalid network mask %q", fields[4])
	}

	maskIPv4 := maskIP.To4()
	if maskIPv4 == nil {
		return nil, fmt.Errorf("network mask %q is not IPv4", fields[4])
	}

	ones, bits := net.IPMask(maskIPv4).Size()
	if bits != 32 {
		return nil, fmt.Errorf(
			"invalid or non-contiguous network mask %q",
			fields[4],
		)
	}

	network := netip.PrefixFrom(address, ones)

	return route.NewRoute(network, description), nil
}
