package parser

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/Ranger0806/go-routes/internal/route"
)

func ParseDirectory(rootPath string) ([]*route.Route, error) {
	return ParseDirectoryExcept(rootPath, "")
}

func ParseDirectoryExcept(
	rootPath string,
	excludedPath string,
) ([]*route.Route, error) {
	routes := make([]*route.Route, 0)

	rootAbsolutePath, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, fmt.Errorf(
			"resolve root directory %q: %w",
			rootPath,
			err,
		)
	}

	excludedAbsolutePath := ""

	if strings.TrimSpace(excludedPath) != "" {
		excludedAbsolutePath, err = filepath.Abs(excludedPath)
		if err != nil {
			return nil, fmt.Errorf(
				"resolve excluded path %q: %w",
				excludedPath,
				err,
			)
		}

		excludedAbsolutePath = filepath.Clean(excludedAbsolutePath)
	}

	err = filepath.WalkDir(
		filepath.Clean(rootAbsolutePath),
		func(
			path string,
			entry fs.DirEntry,
			walkErr error,
		) error {
			if walkErr != nil {
				return fmt.Errorf(
					"access path %q: %w",
					path,
					walkErr,
				)
			}

			if entry.IsDir() {
				return nil
			}

			absolutePath, err := filepath.Abs(path)
			if err != nil {
				return fmt.Errorf(
					"resolve path %q: %w",
					path,
					err,
				)
			}

			absolutePath = filepath.Clean(absolutePath)

			if excludedAbsolutePath != "" &&
				absolutePath == excludedAbsolutePath {
				return nil
			}

			if !strings.EqualFold(
				filepath.Ext(entry.Name()),
				".bat",
			) {
				return nil
			}

			fileRoutes, err := ParseBATFile(path)
			if err != nil {
				return fmt.Errorf(
					"parse BAT file %q: %w",
					path,
					err,
				)
			}

			routes = append(routes, fileRoutes...)

			return nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf(
			"walk directory %q: %w",
			rootPath,
			err,
		)
	}

	return routes, nil
}
