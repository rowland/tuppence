package ir

import "strings"

// Binding is one source-level value version. Previous links same-scope mutable
// versions; a first binding leaves Previous invalid.
type Binding struct {
	ID       BindingID
	Scope    ScopeID
	Name     string
	Previous BindingID
	TypeSlot TypeVarID
	Value    ExprID
	Mutable  bool
	Origin   Origin
}

// Valid reports whether the binding has the required identity and source data.
func (b Binding) Valid() bool {
	return b.ID.Valid() && b.Scope.Valid() && strings.TrimSpace(b.Name) != "" && b.TypeSlot.Valid() && b.Value.Valid() && b.Origin.Valid()
}
