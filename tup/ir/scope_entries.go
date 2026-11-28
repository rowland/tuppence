package ir

import "strings"

// ScopeEntry associates a source name with its declaration at a lexical
// position. Offset is the declaration's byte offset within its source file.
type ScopeEntry struct {
	Name        string
	Declaration DeclID
	Offset      uint32
}

// Valid reports whether the entry has a name and declaration.
func (e ScopeEntry) Valid() bool {
	return strings.TrimSpace(e.Name) != "" && e.Declaration.Valid()
}
