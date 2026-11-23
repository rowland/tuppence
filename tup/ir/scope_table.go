package ir

// ScopeTable is an ordered collection of scopes indexed by identity.
type ScopeTable []Scope

// Find returns the scope with the requested identity.
func (t ScopeTable) Find(id ScopeID) (Scope, bool) {
	for _, scope := range t {
		if scope.ID == id {
			return scope, true
		}
	}
	return Scope{}, false
}

// Valid reports whether every scope is valid and IDs are unique.
func (t ScopeTable) Valid() bool {
	seen := make(map[ScopeID]struct{}, len(t))
	for _, scope := range t {
		if !scope.Valid() {
			return false
		}
		if _, exists := seen[scope.ID]; exists {
			return false
		}
		seen[scope.ID] = struct{}{}
	}
	return true
}
