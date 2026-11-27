package ir

import "strings"

// CompilationTarget combines target-specific file selection with module lookup.
type CompilationTarget struct {
	OS           string
	Architecture string
	ModuleRoots  ModuleResolutionConfig
	PrologPaths  []string
	PrimaryPath  string
}

// Valid reports whether the target has the information required for loading.
func (t CompilationTarget) Valid() bool {
	return strings.TrimSpace(t.OS) != "" && strings.TrimSpace(t.Architecture) != "" && t.ModuleRoots.Valid()
}

// BuildModuleIndexForTarget builds the shared source/module index for a target.
func BuildModuleIndexForTarget(filenames []string, target CompilationTarget) (SourceFileTable, ModuleTable, []string) {
	if !target.Valid() {
		return nil, nil, append([]string(nil), filenames...)
	}
	return BuildModuleIndex(filenames, target.OS, target.Architecture)
}
