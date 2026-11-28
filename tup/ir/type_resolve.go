package ir

// ResolveTypeName resolves one visible nominal type name in lexical scope.
// Ambiguous names are rejected rather than selecting by storage or table order.
func (p IndexedProgram) ResolveTypeName(scope ScopeID, name string) (TypeID, bool) {
	matches := p.FindTypeCandidates(scope, name)
	if len(matches) != 1 {
		return 0, false
	}
	return matches[0], true
}
