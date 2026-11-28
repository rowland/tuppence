package ir

// TypeAlias is a declaration-level name for an existing semantic type. It
// does not introduce a new TypeID.
type TypeAlias struct {
	Declaration DeclID
	Target      TypeID
	Origin      Origin
}

// Valid reports whether the alias has its declaration, target, and provenance.
func (a TypeAlias) Valid() bool {
	return a.Declaration.Valid() && a.Target.Valid() && a.Origin.Valid()
}

// TypeAliasTable is an ordered collection indexed by the alias declaration.
type TypeAliasTable []TypeAlias

// Find returns the alias attached to the requested declaration.
func (t TypeAliasTable) Find(id DeclID) (TypeAlias, bool) {
	for _, alias := range t {
		if alias.Declaration == id {
			return alias, true
		}
	}
	return TypeAlias{}, false
}

// Valid reports whether every alias is valid and declarations are unique.
func (t TypeAliasTable) Valid() bool {
	seen := make(map[DeclID]struct{}, len(t))
	for _, alias := range t {
		if !alias.Valid() {
			return false
		}
		if _, exists := seen[alias.Declaration]; exists {
			return false
		}
		seen[alias.Declaration] = struct{}{}
	}
	return true
}
