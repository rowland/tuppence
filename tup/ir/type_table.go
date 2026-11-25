package ir

// TypeTable is an ordered collection of types indexed by identity.
type TypeTable []Type

// Find returns the type with the requested identity.
func (t TypeTable) Find(id TypeID) (Type, bool) {
	for _, typ := range t {
		if typ.ID == id {
			return typ, true
		}
	}
	return Type{}, false
}

// Valid reports whether every type is valid and IDs are unique.
func (t TypeTable) Valid() bool {
	seen := make(map[TypeID]struct{}, len(t))
	seenDeclarations := make(map[DeclID]struct{}, len(t))
	for _, typ := range t {
		if !typ.Valid() {
			return false
		}
		if _, exists := seen[typ.ID]; exists {
			return false
		}
		if typ.Kind == TypeNominal {
			if _, exists := seenDeclarations[typ.Declaration]; exists {
				return false
			}
			seenDeclarations[typ.Declaration] = struct{}{}
		}
		seen[typ.ID] = struct{}{}
	}
	return true
}
