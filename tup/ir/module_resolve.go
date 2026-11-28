package ir

import (
	pathpkg "path"
	"strings"
)

// ResolveImportCandidates returns the ordered module identities considered for
// a requested path, without touching the filesystem.
func (t ModuleTable) ResolveImportCandidates(current ModuleID, requested string, config ModuleResolutionConfig) []string {
	if !current.Valid() || requested == "" || !config.Valid() {
		return nil
	}
	currentModule, ok := t.Find(current)
	if !ok {
		return nil
	}
	candidates := make([]string, 0, len(config.SearchPaths)+2)
	if strings.HasPrefix(requested, "/") {
		for _, root := range config.SearchPaths {
			candidates = append(candidates, pathpkg.Join(root, strings.TrimPrefix(requested, "/")))
		}
	} else {
		currentDir := pathpkg.Dir(currentModule.ResolutionPath())
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
	result := candidates[:0]
	for _, candidate := range candidates {
		if _, exists := seen[candidate]; exists {
			continue
		}
		seen[candidate] = struct{}{}
		result = append(result, candidate)
	}
	return result
}

// ResolveImport resolves a requested module path against the indexed module
// identities without touching the filesystem.
func (t ModuleTable) ResolveImport(current ModuleID, requested string, config ModuleResolutionConfig) (ModuleID, bool) {
	for _, candidate := range t.ResolveImportCandidates(current, requested, config) {
		for _, module := range t {
			if module.ResolutionPath() == candidate {
				return module.ID, true
			}
		}
	}
	return 0, false
}
