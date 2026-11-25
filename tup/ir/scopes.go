package ir

// ScopeKind identifies the lexical role of a scope.
type ScopeKind uint8

const (
	ScopeModule ScopeKind = iota + 1
	ScopeFunction
	ScopeBlock
	ScopeLoop
)

// Valid reports whether the scope kind is recognized.
func (k ScopeKind) Valid() bool {
	return k >= ScopeModule && k <= ScopeLoop
}

// Scope is the identity and parentage of one lexical scope. A zero Parent is
// valid for a module root; parent ownership is checked by the enclosing table.
type Scope struct {
	ID      ScopeID
	Parent  ScopeID
	Module  ModuleID
	Kind    ScopeKind
	Entries []ScopeEntry
	Origin  Origin
}

// Valid reports whether the scope has the required identity and kind.
func (s Scope) Valid() bool {
	if !s.ID.Valid() || !s.Module.Valid() || !s.Kind.Valid() || !s.Origin.Valid() {
		return false
	}
	for _, entry := range s.Entries {
		if !entry.Valid() {
			return false
		}
	}
	return true
}
