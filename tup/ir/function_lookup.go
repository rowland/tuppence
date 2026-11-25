package ir

// FunctionKey identifies an overload by lexical scope, base name, and ordered parameter types.
type FunctionKey struct {
	Scope      ScopeID
	Name       string
	Parameters []FunctionParameter
}

// FunctionKey returns the overload key for a function declaration.
func (p IndexedProgram) FunctionKey(id FunctionID) (FunctionKey, bool) {
	function, ok := p.Functions.Find(id)
	if !ok {
		return FunctionKey{}, false
	}
	declaration, ok := p.Declarations.Find(function.Declaration)
	if !ok {
		return FunctionKey{}, false
	}
	return FunctionKey{Scope: declaration.Scope, Name: declaration.Name, Parameters: function.Parameters}, true
}

// FindFunctions returns matching overloads in stored order.
func (p IndexedProgram) FindFunctions(key FunctionKey) []FunctionID {
	var matches []FunctionID
	for _, function := range p.Functions {
		candidate, ok := p.FunctionKey(function.ID)
		if ok && candidate.Scope == key.Scope && candidate.Name == key.Name && sameParameterTypes(candidate.Parameters, key.Parameters) {
			matches = append(matches, function.ID)
		}
	}
	return matches
}
