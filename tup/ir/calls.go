package ir

import "strings"

// CallArgument describes one positional or named argument in source order.
type CallArgument struct {
	Name            string
	Type            TypeID
	Literal         LiteralKind
	LiteralTemplate bool
}

// Valid reports whether the argument has either a known type or a numeric literal
// category and, when present, a name.
func (a CallArgument) Valid() bool {
	if a.Name != "" && strings.TrimSpace(a.Name) == "" {
		return false
	}
	if a.Type.Valid() == a.Literal.Valid() {
		return false
	}
	if a.LiteralTemplate && !a.Literal.Valid() {
		return false
	}
	return true
}

// CallExpr records a call site before overload resolution selects a function.
type CallExpr struct {
	ID        ExprID
	Scope     ScopeID
	Callee    ExprID
	Name      string
	Arguments []CallArgument
	Curry     bool
	Origin    Origin
}

// Valid reports whether the call site has identity, scope, name, and provenance.
func (c CallExpr) Valid() bool {
	if !c.ID.Valid() || !c.Scope.Valid() || !c.Origin.Valid() {
		return false
	}
	if strings.TrimSpace(c.Name) == "" && !c.Callee.Valid() {
		return false
	}
	if c.Callee == c.ID {
		return false
	}
	named := false
	seenNames := make(map[string]struct{})
	for _, argument := range c.Arguments {
		if !argument.Valid() {
			return false
		}
		if argument.Name == "" {
			if named {
				return false
			}
			continue
		}
		named = true
		if _, exists := seenNames[argument.Name]; exists {
			return false
		}
		seenNames[argument.Name] = struct{}{}
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
