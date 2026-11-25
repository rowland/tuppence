package ir

import "strings"

// CallExpr records a call site before overload resolution selects a function.
type CallExpr struct {
	ID        ExprID
	Scope     ScopeID
	Name      string
	Arguments []TypeID
	Origin    Origin
}

// Valid reports whether the call site has identity, scope, name, and provenance.
func (c CallExpr) Valid() bool {
	if !c.ID.Valid() || !c.Scope.Valid() || strings.TrimSpace(c.Name) == "" || !c.Origin.Valid() {
		return false
	}
	for _, argument := range c.Arguments {
		if !argument.Valid() {
			return false
		}
	}
	return true
}

// CallTable is an ordered collection of unresolved call sites.
type CallTable []CallExpr

// Find returns the call site with the requested identity.
func (t CallTable) Find(id ExprID) (CallExpr, bool) {
	for _, call := range t {
		if call.ID == id {
			return call, true
		}
	}
	return CallExpr{}, false
}

// Valid reports whether every call site is valid and IDs are unique.
func (t CallTable) Valid() bool {
	seen := make(map[ExprID]struct{}, len(t))
	for _, call := range t {
		if !call.Valid() {
			return false
		}
		if _, exists := seen[call.ID]; exists {
			return false
		}
		seen[call.ID] = struct{}{}
	}
	return true
}
