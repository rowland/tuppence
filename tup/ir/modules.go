package ir

import "strings"

// Module identifies one semantic compilation unit.
type Module struct {
	ID     ModuleID
	Name   string
	Files  []FileID
	Origin Origin
}

// Valid reports whether the module has the required identity, source data, and unique files.
func (m Module) Valid() bool {
	if !m.ID.Valid() || strings.TrimSpace(m.Name) == "" || !m.Origin.Valid() {
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
	return true
}
