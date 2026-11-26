package ir

import (
	pathpkg "path"
	"strings"
)

// ResolveImport resolves a requested module path against the indexed module
// identities without touching the filesystem.
func (t ModuleTable) ResolveImport(current ModuleID, requested string, config ModuleResolutionConfig) (ModuleID, bool) {
	if !current.Valid() || requested == "" || !config.Valid() {
		return 0, false
	}
	currentModule, ok := t.Find(current)
	if !ok {
		return 0, false
	}
	candidates := make([]string, 0, len(config.SearchPaths)+2)
	if strings.HasPrefix(requested, "/") {
		for _, root := range config.SearchPaths {
			candidates = append(candidates, pathpkg.Join(root, strings.TrimPrefix(requested, "/")))
		}
	} else {
		currentDir := pathpkg.Dir(currentModule.Name)
		if currentDir != "." {
			candidates = append(candidates, pathpkg.Join(currentDir, requested))
		}
		if config.ProgramRoot != "" {
			candidates = append(candidates, pathpkg.Join(config.ProgramRoot, requested))
		} else {
			candidates = append(candidates, requested)
		}
		for _, root := range config.SearchPaths {
			candidates = append(candidates, pathpkg.Join(root, requested))
		}
	}
	seen := make(map[string]struct{}, len(candidates))
	for _, candidate := range candidates {
		if _, exists := seen[candidate]; exists {
			continue
		}
		seen[candidate] = struct{}{}
		for _, module := range t {
			if module.Name == candidate {
				return module.ID, true
			}
		}
	}
	return 0, false
}
