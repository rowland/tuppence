package ir

// ReferenceExpr represents a direct semantic reference to an existing
// declaration. Resolution may additionally identify the function or value
// identity; neither implies a copy or allocation.
type ReferenceExpr struct {
	ID          ExprID
	Declaration DeclID
	Function    FunctionID
	Value       ValueID
	Origin      Origin
}

// Valid reports whether the reference has an expression identity, declaration,
// provenance, and at most one resolved runtime-facing identity.
func (r ReferenceExpr) Valid() bool {
	if !r.ID.Valid() || !r.Declaration.Valid() || !r.Origin.Valid() {
		return false
	}
	return !(r.Function.Valid() && r.Value.Valid())
}
