package ir

import "sort"

// SelectedModuleFiles groups eligible source files by semantic module name.
type SelectedModuleFiles struct {
	ModuleName string
	Files      []string
}

// SelectModuleFiles parses and groups source filenames deterministically. Invalid
// filenames are returned separately so loaders can issue source diagnostics.
func SelectModuleFiles(filenames []string, osName, architecture string) ([]SelectedModuleFiles, []string) {
	groups := make(map[string][]string)
	var invalid []string
	for _, filename := range filenames {
		parsed, ok := ParseModuleFileName(filename)
		if !ok {
			invalid = append(invalid, filename)
			continue
		}
		if parsed.Matches(osName, architecture) {
			groups[parsed.ModuleName] = append(groups[parsed.ModuleName], filename)
		}
	}
	moduleNames := make([]string, 0, len(groups))
	for moduleName := range groups {
		moduleNames = append(moduleNames, moduleName)
	}
	sort.Strings(moduleNames)
	selected := make([]SelectedModuleFiles, 0, len(moduleNames))
	for _, moduleName := range moduleNames {
		files := groups[moduleName]
		sort.Strings(files)
		selected = append(selected, SelectedModuleFiles{ModuleName: moduleName, Files: files})
	}
	sort.Strings(invalid)
	return selected, invalid
}
