package ir

// IndexedProgram is the first semantic snapshot: source identities have been
// collected into validated lexical tables, but expressions are not analyzed yet.
type IndexedProgram struct {
	Modules      ModuleTable
	Scopes       ScopeTable
	Declarations DeclarationTable
	Bindings     BindingTable
	Diagnostics  []Diagnostic
}

// Valid reports whether every indexed table and diagnostic is valid.
func (p IndexedProgram) Valid() bool {
	if !p.Modules.Valid() || !p.Scopes.Valid() || !p.Declarations.Valid() || !p.Bindings.Valid() {
		return false
	}
	for _, scope := range p.Scopes {
		if _, ok := p.Modules.Find(scope.Module); !ok {
			return false
		}
		if scope.Parent.Valid() {
			parent, ok := p.Scopes.Find(scope.Parent)
			if !ok || parent.Module != scope.Module {
				return false
			}
		}
	}
	for _, declaration := range p.Declarations {
		if _, ok := p.Scopes.Find(declaration.Scope); !ok {
			return false
		}
	}
	for _, binding := range p.Bindings {
		if _, ok := p.Scopes.Find(binding.Scope); !ok {
			return false
		}
		if binding.Previous.Valid() {
			previous, ok := p.Bindings.Find(binding.Previous)
			if !ok || previous.Scope != binding.Scope || previous.Name != binding.Name {
				return false
			}
		}
	}
	for _, diagnostic := range p.Diagnostics {
		if !diagnostic.Valid() {
			return false
		}
	}
	return true
}
