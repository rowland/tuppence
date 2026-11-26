package ir

import "strings"

// ImportAlias binds one exported member of an imported module into a local scope.
type ImportAlias struct {
	Scope       ScopeID
	LocalName   string
	Module      ModuleID
	Declaration DeclID
	Source      ImportAliasSource
	Origin      Origin
}

// Valid reports whether the alias has valid identity and provenance.
func (a ImportAlias) Valid() bool {
	return a.Scope.Valid() && strings.TrimSpace(a.LocalName) != "" && a.Module.Valid() && a.Declaration.Valid() && a.Source.Valid() && a.Origin.Valid()
}
