package ir

// IndexedProgram is the first semantic snapshot: source identities have been
// collected into validated lexical tables, but expressions are not analyzed yet.
type IndexedProgram struct {
	Modules      ModuleTable
	Types        TypeTable
	Constructors TypeConstructorTable
	Scopes       ScopeTable
	Declarations DeclarationTable
	Bindings     BindingTable
	Functions    FunctionTable
	Aliases      TypeAliasTable
	Expressions  ExpressionTable
	Calls        CallTable
	Diagnostics  []Diagnostic
}

// Valid reports whether every indexed table, cross-reference, and diagnostic is valid.
func (p IndexedProgram) Valid() bool {
	if !p.Modules.Valid() || !p.Types.Valid() || !p.Constructors.Valid() || !p.Scopes.Valid() || !p.Declarations.Valid() || !p.Bindings.Valid() || !p.Functions.Valid() || !p.Aliases.Valid() || !p.Expressions.Valid() || !p.Calls.Valid() {
		return false
	}
	fileOwners := make(map[FileID]ModuleID)
	for _, module := range p.Modules {
		for _, file := range module.Files {
			if owner, exists := fileOwners[file]; exists && owner != module.ID {
				return false
			}
			fileOwners[file] = module.ID
		}
	}
	for _, scope := range p.Scopes {
		if _, ok := p.Modules.Find(scope.Module); !ok {
			return false
		}
		if scope.Origin.Primary.File.Valid() && !moduleOwnsOrigin(p.Modules, scope.Module, scope.Origin) {
			return false
		}
		if scope.Parent.Valid() {
			parent, ok := p.Scopes.Find(scope.Parent)
			if !ok || parent.Module != scope.Module {
				return false
			}
		}
	}
	for _, module := range p.Modules {
		for _, importBinding := range module.Imports {
			if importBinding.Target == module.ID {
				return false
			}
			if !moduleOwnsOrigin(p.Modules, module.ID, importBinding.Origin) {
				return false
			}
			if _, ok := p.Modules.Find(importBinding.Target); !ok {
				return false
			}
			if moduleImportCollides(p.Scopes, module.ID, importBinding.LocalName) {
				return false
			}
		}
	}
	if !scopeParentsAcyclic(p.Scopes) {
		return false
	}
	for _, typ := range p.Types {
		if typ.Kind == TypeArray {
			if _, ok := p.Types.Find(typ.Element); !ok {
				return false
			}
			continue
		}
		if typ.Kind == TypeTuple || typ.Kind == TypeUnion {
			for _, element := range typ.Elements {
				if _, ok := p.Types.Find(element); !ok {
					return false
				}
			}
			continue
		}
		declaration, ok := p.Declarations.Find(typ.Declaration)
		if !ok || declaration.Kind != DeclType {
			return false
		}
		scope, ok := p.Scopes.Find(declaration.Scope)
		if !ok || scope.Module != typ.HomeModule || !moduleOwnsOrigin(p.Modules, typ.HomeModule, typ.Origin) {
			return false
		}
	}
	for _, constructor := range p.Constructors {
		typ, ok := p.Types.Find(constructor.Type)
		if !ok || typ.Kind != TypeNominal || !moduleOwnsOrigin(p.Modules, typ.HomeModule, constructor.Origin) {
			return false
		}
	}
	for _, alias := range p.Aliases {
		declaration, ok := p.Declarations.Find(alias.Declaration)
		if !ok || declaration.Kind != DeclType {
			return false
		}
		declarationScope, ok := p.Scopes.Find(declaration.Scope)
		if !ok || !moduleOwnsOrigin(p.Modules, declarationScope.Module, declaration.Origin) || !moduleOwnsOrigin(p.Modules, declarationScope.Module, alias.Origin) {
			return false
		}
		for _, typ := range p.Types {
			if typ.Declaration == alias.Declaration {
				return false
			}
		}
		target, ok := p.Types.Find(alias.Target)
		if !ok {
			return false
		}
		targetDeclaration, ok := p.Declarations.Find(target.Declaration)
		if !ok || (declarationScope.Module != target.HomeModule && (!targetDeclaration.Exported || !moduleImports(p.Modules, declarationScope.Module, target.HomeModule))) {
			return false
		}
	}
	for i, function := range p.Functions {
		declaration, ok := p.Declarations.Find(function.Declaration)
		if !ok || declaration.Kind != DeclFunction {
			return false
		}
		declarationScope, ok := p.Scopes.Find(declaration.Scope)
		if !ok || !moduleOwnsOrigin(p.Modules, declarationScope.Module, function.Origin) {
			return false
		}
		seenParameterNames := make(map[string]struct{}, len(function.Parameters))
		for _, parameter := range function.Parameters {
			if !parameter.Valid() {
				return false
			}
			if _, exists := seenParameterNames[parameter.Name]; exists {
				return false
			}
			seenParameterNames[parameter.Name] = struct{}{}
			if _, ok := p.Types.Find(parameter.Type); !ok {
				return false
			}
		}
		if function.Result.Valid() {
			if _, ok := p.Types.Find(function.Result); !ok {
				return false
			}
		}
		for _, other := range p.Functions[i+1:] {
			otherDeclaration, ok := p.Declarations.Find(other.Declaration)
			if !ok {
				return false
			}
			if declaration.Scope == otherDeclaration.Scope && declaration.Name == otherDeclaration.Name && sameParameterTypes(function.Parameters, other.Parameters) {
				return false
			}
		}
	}
	for _, expression := range p.Expressions {
		declaration, ok := p.Declarations.Find(expression.Declaration)
		if !ok {
			return false
		}
		referenceScope, ok := p.Scopes.Find(expression.Scope)
		if !ok {
			return false
		}
		declarationScope, ok := p.Scopes.Find(declaration.Scope)
		if !ok || !moduleOwnsOrigin(p.Modules, declarationScope.Module, declaration.Origin) || (referenceScope.Module != declarationScope.Module && (!declaration.Exported || !moduleImports(p.Modules, referenceScope.Module, declarationScope.Module))) {
			return false
		}
		if !moduleOwnsOrigin(p.Modules, referenceScope.Module, expression.Origin) {
			return false
		}
		if expression.Function.Valid() {
			function, ok := p.Functions.Find(expression.Function)
			if !ok || function.Declaration != expression.Declaration || declaration.Kind != DeclFunction {
				return false
			}
		}
	}
	for _, call := range p.Calls {
		scope, ok := p.Scopes.Find(call.Scope)
		if !ok || !moduleOwnsOrigin(p.Modules, scope.Module, call.Origin) {
			return false
		}
		if call.Callee.Valid() {
			if _, ok := p.Calls.Find(call.Callee); !ok {
				return false
			}
		}
		for _, argument := range call.Arguments {
			if _, ok := p.Types.Find(argument.Type); !ok {
				return false
			}
		}
	}
	for _, declaration := range p.Declarations {
		scope, ok := p.Scopes.Find(declaration.Scope)
		if !ok || !moduleOwnsOrigin(p.Modules, scope.Module, declaration.Origin) {
			return false
		}
		if declaration.OwnerType.Valid() {
			owner, ok := p.Declarations.Find(declaration.OwnerType)
			if !ok || owner.Kind != DeclType || owner.Scope != declaration.Scope {
				return false
			}
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
		scope, ok := p.Scopes.Find(binding.Scope)
		if !ok || !moduleOwnsOrigin(p.Modules, scope.Module, binding.Origin) {
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

func sameParameterTypes(left, right []FunctionParameter) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i].Type != right[i].Type {
			return false
		}
	}
	return true
}

func sameTypeIDs(left, right []TypeID) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}
	return true
}

func moduleOwnsOrigin(modules ModuleTable, moduleID ModuleID, origin Origin) bool {
	module, ok := modules.Find(moduleID)
	return ok && moduleHasFile(module.Files, origin.Primary.File)
}

func moduleImportCollides(scopes ScopeTable, module ModuleID, name string) bool {
	for _, scope := range scopes {
		if scope.Module != module || scope.Kind != ScopeModule || scope.Parent.Valid() {
			continue
		}
		for _, entry := range scope.Entries {
			if entry.Name == name {
				return true
			}
		}
	}
	return false
}

func moduleImports(modules ModuleTable, source, target ModuleID) bool {
	module, ok := modules.Find(source)
	if !ok {
		return false
	}
	for _, binding := range module.Imports {
		if binding.Target == target {
			return true
		}
	}
	return false
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
