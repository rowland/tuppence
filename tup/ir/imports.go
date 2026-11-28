package ir

import "strings"

// ImportBinding names an imported module within the importing module's scope.
type ImportBinding struct {
	LocalName string
	Path      string
	Target    ModuleID
	Origin    Origin
}

// Valid reports whether the import binding has a name, target, and provenance.
func (i ImportBinding) Valid() bool {
	return strings.TrimSpace(i.LocalName) != "" && strings.TrimSpace(i.Path) != "" && i.Target.Valid() && i.Origin.Valid()
}
