package ir

// IndexedProgram is the first semantic snapshot: source identities have been
// collected into validated lexical tables, but expressions are not analyzed yet.
type IndexedProgram struct {
	Scopes       ScopeTable
	Declarations DeclarationTable
	Bindings     BindingTable
	Diagnostics  []Diagnostic
}

// Valid reports whether every indexed table and diagnostic is valid.
func (p IndexedProgram) Valid() bool {
	if !p.Scopes.Valid() || !p.Declarations.Valid() || !p.Bindings.Valid() {
		return false
	}
	for _, diagnostic := range p.Diagnostics {
		if !diagnostic.Valid() {
			return false
		}
	}
	return true
}
