package ir

import "strings"

// Module identifies one semantic compilation unit.
type Module struct {
	ID      ModuleID
	Name    string
	Files   []FileID
	Imports []ImportBinding
	Origin  Origin
}

func moduleHasFile(files []FileID, file FileID) bool {
	for _, candidate := range files {
		if candidate == file {
			return true
		}
	}
	return false
}

// Valid reports whether the module has the required identity, source data, and unique files/imports.
func (m Module) Valid() bool {
	if !m.ID.Valid() || strings.TrimSpace(m.Name) == "" || !m.Origin.Valid() || len(m.Files) == 0 {
		return false
	}
	if !moduleHasFile(m.Files, m.Origin.Primary.File) {
		return false
	}
	seenFiles := make(map[FileID]struct{}, len(m.Files))
	for _, file := range m.Files {
		if !file.Valid() {
			return false
		}
		if _, exists := seenFiles[file]; exists {
			return false
		}
		seenFiles[file] = struct{}{}
	}
	seenImports := make(map[string]struct{}, len(m.Imports))
	for _, binding := range m.Imports {
		if !binding.Valid() {
			return false
		}
		if _, exists := seenImports[binding.LocalName]; exists {
			return false
		}
		seenImports[binding.LocalName] = struct{}{}
	}
	return true
}
