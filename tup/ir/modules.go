package ir

import "strings"

// Module identifies one semantic compilation unit.
type Module struct {
	ID     ModuleID
	Name   string
	Origin Origin
}

// Valid reports whether the module has the required identity and source data.
func (m Module) Valid() bool {
	return m.ID.Valid() && strings.TrimSpace(m.Name) != "" && m.Origin.Valid()
}
