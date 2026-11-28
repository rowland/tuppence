package ir

// DeclarationTable is an ordered collection of declarations indexed by identity.
type DeclarationTable []Declaration

// Find returns the declaration with the requested identity.
func (t DeclarationTable) Find(id DeclID) (Declaration, bool) {
	for _, declaration := range t {
		if declaration.ID == id {
			return declaration, true
		}
	}
	return Declaration{}, false
}

// Valid reports whether every declaration is valid and IDs are unique.
func (t DeclarationTable) Valid() bool {
	seen := make(map[DeclID]struct{}, len(t))
	for _, declaration := range t {
		if !declaration.Valid() {
			return false
		}
		if _, exists := seen[declaration.ID]; exists {
			return false
		}
		seen[declaration.ID] = struct{}{}
	}
	return true
}
