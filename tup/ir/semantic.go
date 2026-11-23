package ir

// IndexedProgram is the first semantic snapshot: source identities have been
// collected into validated lexical tables, but expressions are not analyzed yet.
type IndexedProgram struct {
	Scopes       ScopeTable
	Declarations DeclarationTable
	Bindings     BindingTable
}

// Valid reports whether every indexed table is internally valid.
func (p IndexedProgram) Valid() bool {
	return p.Scopes.Valid() && p.Declarations.Valid() && p.Bindings.Valid()
}
