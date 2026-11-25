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

// FindCallCandidates returns exact-type overload candidates for a direct call.
// It preserves function-table order and does not rank candidates.
func (p IndexedProgram) FindCallCandidates(call CallExpr) []FunctionID {
	if call.Name == "" {
		return nil
	}
	var matches []FunctionID
	for _, function := range p.Functions {
		declaration, ok := p.Declarations.Find(function.Declaration)
		if !ok || declaration.Scope != call.Scope || declaration.Name != call.Name {
			continue
		}
		if callMatchesFunction(call, function) {
			matches = append(matches, function.ID)
		}
	}
	return matches
}

func callMatchesFunction(call CallExpr, function Function) bool {
	if !call.Curry && len(call.Arguments) != len(function.Parameters) {
		return false
	}
	if len(call.Arguments) > len(function.Parameters) {
		return false
	}
	matched := make([]bool, len(function.Parameters))
	positional := 0
	for _, argument := range call.Arguments {
		index := positional
		if argument.Name != "" {
			index = -1
			for i, parameter := range function.Parameters {
				if parameter.Name == argument.Name {
					index = i
					break
				}
			}
		} else {
			positional++
		}
		if index < 0 || index >= len(function.Parameters) || matched[index] || function.Parameters[index].Type != argument.Type {
			return false
		}
		matched[index] = true
	}
	if !call.Curry {
		for _, present := range matched {
			if !present {
				return false
			}
		}
	}
	return true
}
