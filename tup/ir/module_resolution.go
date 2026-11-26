package ir

import "strings"

// ModuleResolutionConfig defines the module roots used after relative lookup.
type ModuleResolutionConfig struct {
	ProgramRoot string
	SearchPaths []string
}

// Valid reports whether configured roots are absolute and unique.
func (c ModuleResolutionConfig) Valid() bool {
	if c.ProgramRoot != "" && !strings.HasPrefix(c.ProgramRoot, "/") {
		return false
	}
	seen := make(map[string]struct{}, len(c.SearchPaths))
	for _, path := range c.SearchPaths {
		if path == "" || !strings.HasPrefix(path, "/") {
			return false
		}
		if _, exists := seen[path]; exists {
			return false
		}
		seen[path] = struct{}{}
	}
	return true
}
