package ir

// ExpressionTable is the initial semantic expression table. It currently
// stores direct reference expressions; other expression forms will join it.
type ExpressionTable []ReferenceExpr

// Find returns the reference expression with the requested identity.
func (t ExpressionTable) Find(id ExprID) (ReferenceExpr, bool) {
	for _, expression := range t {
		if expression.ID == id {
			return expression, true
		}
	}
	return ReferenceExpr{}, false
}

// Valid reports whether every expression is valid and IDs are unique.
func (t ExpressionTable) Valid() bool {
	seen := make(map[ExprID]struct{}, len(t))
	for _, expression := range t {
		if !expression.Valid() {
			return false
		}
		if _, exists := seen[expression.ID]; exists {
			return false
		}
		seen[expression.ID] = struct{}{}
	}
	return true
}
