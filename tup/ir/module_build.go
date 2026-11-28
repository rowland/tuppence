package ir

import (
	pathpkg "path"
	"sort"
)

// BuildModulesFromSourceFiles creates module records from indexed source files.
func BuildModulesFromSourceFiles(files SourceFileTable) ModuleTable {
	groups := make(map[string]*Module)
	for _, file := range files {
		modulePath := pathpkg.Join(pathpkg.Dir(file.Path), file.Module.ModuleName)
		if existing, ok := groups[modulePath]; ok {
			existing.Files = append(existing.Files, file.ID)
			continue
		}
		groups[modulePath] = &Module{
			ID:     ModuleID(len(groups) + 1),
			Name:   file.Module.ModuleName,
			Path:   modulePath,
			Files:  []FileID{file.ID},
			Origin: Origin{Primary: Span{File: file.ID, Start: 0, End: 0}},
		}
	}
	paths := make([]string, 0, len(groups))
	for modulePath := range groups {
		paths = append(paths, modulePath)
	}
	sort.Strings(paths)
	result := make(ModuleTable, 0, len(paths))
	for i, modulePath := range paths {
		module := *groups[modulePath]
		module.ID = ModuleID(i + 1)
		result = append(result, module)
	}
	return result
}
