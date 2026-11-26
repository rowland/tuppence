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

// FindCastFunctions returns marked cast functions in lexical scope order
// for an exact source and target type pair.
func (p IndexedProgram) FindCastFunctions(scope ScopeID, source TypeID, target TypeID) []FunctionID {
	if !scope.Valid() || !source.Valid() || !target.Valid() {
		return nil
	}
	var matches []FunctionID
	for _, function := range p.Functions {
		if !function.Cast || function.Effect != EffectPure || len(function.Parameters) != 1 || function.Parameters[0].Type != source || function.Result != target {
			continue
		}
		declaration, ok := p.Declarations.Find(function.Declaration)
		if ok && (declaration.Scope == scope || p.importAliasVisible(scope, declaration.ID)) {
			matches = append(matches, function.ID)
		}
	}
	return matches
}

// CallCandidateResult preserves the call with literal-resolution metadata and
// the viable functions discovered for it.
type CallCandidateResult struct {
	Call      CallExpr
	Functions []FunctionID
}

// FindCallCandidatesWithResolution discovers candidates and marks literal
// arguments ambiguous when more than one candidate remains.
func (p IndexedProgram) FindCallCandidatesWithResolution(call CallExpr) CallCandidateResult {
	matches := p.selectCandidateTier(p.findCallCandidates(call))
	if len(matches) > 1 {
		call.Arguments = append([]CallArgument(nil), call.Arguments...)
		for i := range call.Arguments {
			if call.Arguments[i].Literal.Valid() {
				call.Arguments[i].LiteralAmbiguous = true
			}
		}
	}
	return CallCandidateResult{Call: call, Functions: matches}
}

// FindCallCandidates returns exact-type overload candidates for a direct call.
// It preserves function-table order and does not rank candidates.
func (p IndexedProgram) FindCallCandidates(call CallExpr) []FunctionID {
	return p.findCallCandidates(call)
}

func (p IndexedProgram) findCallCandidates(call CallExpr) []FunctionID {
	if call.Name == "" {
		return nil
	}
	var matches []FunctionID
	for _, function := range p.Functions {
		declaration, ok := p.Declarations.Find(function.Declaration)
		visible := false
		if ok {
			if call.Qualifier == "" {
				visible = p.functionVisible(call.Scope, declaration)
			} else {
				visible = p.qualifiedFunctionVisible(call.Scope, call.Qualifier, declaration)
			}
		}
		if !ok || !visible || declaration.Name != call.Name {
			continue
		}
		if callMatchesFunction(p, call, function) {
			matches = append(matches, function.ID)
		}
	}
	return matches
}

func (p IndexedProgram) importAliasVisible(scope ScopeID, declaration DeclID) bool {
	for _, alias := range p.ImportAliases {
		if alias.Scope == scope && alias.Declaration == declaration {
			return true
		}
	}
	return false
}

func (p IndexedProgram) functionVisible(callScope ScopeID, declaration Declaration) bool {
	caller, callerOK := p.Scopes.Find(callScope)
	declScope, declOK := p.Scopes.Find(declaration.Scope)
	if !callerOK || !declOK {
		return declaration.Scope == callScope
	}
	if p.importAliasVisible(callScope, declaration.ID) {
		return true
	}
	if caller.Module == declScope.Module {
		for current := caller; ; {
			if current.ID == declScope.ID {
				return true
			}
			if !current.Parent.Valid() {
				return false
			}
			next, ok := p.Scopes.Find(current.Parent)
			if !ok {
				return false
			}
			current = next
		}
	}
	return false
}

func (p IndexedProgram) qualifiedFunctionVisible(callScope ScopeID, qualifier string, declaration Declaration) bool {
	caller, ok := p.Scopes.Find(callScope)
	if !ok {
		return false
	}
	declScope, ok := p.Scopes.Find(declaration.Scope)
	if !ok || !declaration.Exported || caller.Module == declScope.Module {
		return false
	}
	module, ok := p.Modules.Find(caller.Module)
	if !ok {
		return false
	}
	for _, binding := range module.Imports {
		if binding.LocalName == qualifier && binding.Target == declScope.Module {
			return true
		}
	}
	return false
}

func (p IndexedProgram) selectCandidateTier(candidates []FunctionID) []FunctionID {
	for _, tier := range []func(Function) bool{
		func(function Function) bool { return !function.Generic && !function.Any },
		func(function Function) bool { return function.Generic },
		func(function Function) bool { return function.Any },
	} {
		var selected []FunctionID
		for _, id := range candidates {
			function, ok := p.Functions.Find(id)
			if ok && tier(function) {
				selected = append(selected, id)
			}
		}
		if len(selected) > 0 {
			return selected
		}
	}
	return nil
}

func callMatchesFunction(program IndexedProgram, call CallExpr, function Function) bool {
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
		if index < 0 || index >= len(function.Parameters) || matched[index] || !program.argumentMatchesParameter(argument, function.Parameters[index], function, call.Scope) {
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

func (p IndexedProgram) typeVisible(scope ScopeID, typeID TypeID) bool {
	typ, ok := p.Types.Find(typeID)
	if !ok {
		return true
	}
	declaration, ok := p.Declarations.Find(typ.Declaration)
	if !ok {
		return false
	}
	if p.importAliasVisible(scope, declaration.ID) {
		return true
	}
	caller, callerOK := p.Scopes.Find(scope)
	declScope, declOK := p.Scopes.Find(declaration.Scope)
	if !callerOK || !declOK {
		return true
	}
	if caller.Module != declScope.Module {
		return false
	}
	for current := caller; ; {
		if current.ID == declScope.ID {
			return true
		}
		if !current.Parent.Valid() {
			return false
		}
		next, ok := p.Scopes.Find(current.Parent)
		if !ok {
			return false
		}
		current = next
	}
}

func (p IndexedProgram) argumentMatchesParameter(argument CallArgument, parameter FunctionParameter, function Function, scope ScopeID) bool {
	if function.Generic || function.Any {
		return argument.Type.Valid() || argument.Literal.Valid()
	}
	if argument.Type.Valid() {
		return argument.Type == parameter.Type
	}
	if !argument.Literal.Valid() {
		return false
	}
	if !p.typeVisible(scope, parameter.Type) {
		return false
	}
	var admissible []StorageKind
	for _, constructor := range p.Constructors {
		if constructor.Type == parameter.Type {
			admissible = append(admissible, constructor.Storage)
		}
	}
	_, ok := SelectLiteralStorage(argument.Literal, admissible)
	return ok
}
