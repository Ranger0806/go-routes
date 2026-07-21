package parser

import (
	"path/filepath"
	"strings"
)

func DescriptionFromFilename(filename string) string {
	fileName := filepath.Base(filename)
	ext := filepath.Ext(fileName)
	description := strings.TrimSuffix(fileName, ext)
	description = strings.ReplaceAll(description, "_", " ")
	description = strings.ReplaceAll(description, "-", " ")
	description = strings.TrimSpace(description)
	description = strings.Join(strings.Fields(description), " ")
	return description
}
