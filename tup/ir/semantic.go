package ir

// IndexedProgram is the first semantic snapshot: source identities have been
// collected into validated lexical tables, but expressions are not analyzed yet.
type IndexedProgram struct {
	Modules      ModuleTable
	Types        TypeTable
	Scopes       ScopeTable
	Declarations DeclarationTable
	Bindings     BindingTable
	Functions    FunctionTable
	Aliases      TypeAliasTable
	Expressions  ExpressionTable
	Diagnostics  []Diagnostic
}

// Valid reports whether every indexed table, cross-reference, and diagnostic is valid.
func (p IndexedProgram) Valid() bool {
	if !p.Modules.Valid() || !p.Types.Valid() || !p.Scopes.Valid() || !p.Declarations.Valid() || !p.Bindings.Valid() || !p.Functions.Valid() || !p.Aliases.Valid() || !p.Expressions.Valid() {
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
	if !scopeParentsAcyclic(p.Scopes) {
		return false
	}
	for _, typ := range p.Types {
		declaration, ok := p.Declarations.Find(typ.Declaration)
		if !ok || declaration.Kind != DeclType {
			return false
		}
		scope, ok := p.Scopes.Find(declaration.Scope)
		if !ok || scope.Module != typ.HomeModule {
			return false
		}
	}
	for _, alias := range p.Aliases {
		declaration, ok := p.Declarations.Find(alias.Declaration)
		if !ok || declaration.Kind != DeclType {
			return false
		}
		if _, ok := p.Types.Find(alias.Target); !ok {
			return false
		}
	}
	for _, function := range p.Functions {
		declaration, ok := p.Declarations.Find(function.Declaration)
		if !ok || declaration.Kind != DeclFunction {
			return false
		}
	}
	for _, expression := range p.Expressions {
		declaration, ok := p.Declarations.Find(expression.Declaration)
		if !ok {
			return false
		}
		if expression.Function.Valid() {
			function, ok := p.Functions.Find(expression.Function)
			if !ok || function.Declaration != expression.Declaration || declaration.Kind != DeclFunction {
				return false
			}
		}
	}
	for _, declaration := range p.Declarations {
		if _, ok := p.Scopes.Find(declaration.Scope); !ok {
			return false
		}
	}
	for _, scope := range p.Scopes {
		seenNames := make(map[string]DeclKind, len(scope.Entries))
		for _, entry := range scope.Entries {
			declaration, ok := p.Declarations.Find(entry.Declaration)
			if !ok || declaration.Scope != scope.ID || declaration.Name != entry.Name {
				return false
			}
			if previousKind, exists := seenNames[entry.Name]; exists && (previousKind != DeclFunction || declaration.Kind != DeclFunction) {
				return false
			}
			seenNames[entry.Name] = declaration.Kind
		}
	}
	for _, binding := range p.Bindings {
		if _, ok := p.Scopes.Find(binding.Scope); !ok {
			return false
		}
		if _, ok := p.Expressions.Find(binding.Value); !ok {
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

func scopeParentsAcyclic(scopes ScopeTable) bool {
	for _, start := range scopes {
		seen := make(map[ScopeID]struct{})
		current := start
		for current.Parent.Valid() {
			if _, exists := seen[current.ID]; exists {
				return false
			}
			seen[current.ID] = struct{}{}
			next, ok := scopes.Find(current.Parent)
			if !ok {
				return false
			}
			current = next
		}
	}
	return true
}
