package parser

import (
	"path/filepath"
	"testing"
)

func TestDescriptionFromFilename(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected string
	}{
		{
			name:     "simple filename",
			filename: "youtube.bat",
			expected: "youtube",
		},
		{
			name:     "underscore replacement",
			filename: "youtube_description.bat",
			expected: "youtube description",
		},
		{
			name:     "hyphen replacement",
			filename: "discord-routes.bat",
			expected: "discord routes",
		},
		{
			name:     "multiple separators",
			filename: "my__custom---routes.bat",
			expected: "my custom routes",
		},
		{
			name:     "filename with current OS path",
			filename: filepath.Join("home", "anton", "routes", "telegram.bat"),
			expected: "telegram",
		},
		{
			name:     "filename without extension",
			filename: "youtube",
			expected: "youtube",
		},
		{
			name:     "spaces around words",
			filename: "  youtube_routes  .bat",
			expected: "youtube routes",
		},
		{
			name:     "uppercase preserved",
			filename: "YouTube_Premium.bat",
			expected: "YouTube Premium",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := DescriptionFromFilename(test.filename)

			if actual != test.expected {
				t.Errorf(
					"DescriptionFromFilename(%q): expected %q, got %q",
					test.filename,
					test.expected,
					actual,
				)
			}
		})
	}
}
