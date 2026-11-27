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
	bestDistance := int(^uint(0) >> 1)
	for _, function := range p.Functions {
		declaration, ok := p.Declarations.Find(function.Declaration)
		visible := false
		if ok {
			if call.Qualifier == "" {
				var distance int
				distance, visible = p.functionVisibilityDistance(call.Scope, declaration)
				if visible && distance < bestDistance {
					matches = matches[:0]
					bestDistance = distance
				}
				visible = visible && distance == bestDistance
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
	current, ok := p.Scopes.Find(scope)
	for ok {
		for _, alias := range p.ImportAliases {
			if alias.Scope == current.ID && alias.Declaration == declaration {
				return true
			}
		}
		if !current.Parent.Valid() {
			break
		}
		current, ok = p.Scopes.Find(current.Parent)
	}
	return false
}

func (p IndexedProgram) functionVisibilityDistance(callScope ScopeID, declaration Declaration) (int, bool) {
	caller, callerOK := p.Scopes.Find(callScope)
	declScope, declOK := p.Scopes.Find(declaration.Scope)
	if !callerOK || !declOK {
		return 0, declaration.Scope == callScope
	}
	if caller.Module != declScope.Module {
		for distance := 0; ; distance++ {
			for _, alias := range p.ImportAliases {
				if alias.Scope == caller.ID && alias.Declaration == declaration.ID {
					return distance, true
				}
			}
			if !caller.Parent.Valid() {
				return 0, false
			}
			next, ok := p.Scopes.Find(caller.Parent)
			if !ok {
				return 0, false
			}
			caller = next
		}
	}
	for distance := 0; ; distance++ {
		if caller.ID == declScope.ID {
			return distance, true
		}
		for _, alias := range p.ImportAliases {
			if alias.Scope == caller.ID && alias.Declaration == declaration.ID {
				return distance, true
			}
		}
		if !caller.Parent.Valid() {
			return 0, false
		}
		next, ok := p.Scopes.Find(caller.Parent)
		if !ok {
			return 0, false
		}
		caller = next
	}
}

func (p IndexedProgram) functionVisible(scope ScopeID, declaration Declaration) bool {
	_, ok := p.functionVisibilityDistance(scope, declaration)
	return ok
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

// FindTypeCandidates returns visible nominal types with the nearest lexical
// scope taking precedence over outer scopes.
func (p IndexedProgram) FindTypeCandidates(scope ScopeID, name string) []TypeID {
	if name == "" {
		return nil
	}
	bestDistance := int(^uint(0) >> 1)
	var matches []TypeID
	for _, typ := range p.Types {
		if typ.Name != name || !typ.Declaration.Valid() {
			continue
		}
		declaration, ok := p.Declarations.Find(typ.Declaration)
		if !ok {
			continue
		}
		distance, visible := p.functionVisibilityDistance(scope, declaration)
		if !visible {
			continue
		}
		if distance < bestDistance {
			bestDistance = distance
			matches = matches[:0]
		}
		if distance == bestDistance {
			matches = append(matches, typ.ID)
		}
	}
	return matches
}
func (p IndexedProgram) typeVisible(scope ScopeID, typeID TypeID) bool {
	typ, ok := p.Types.Find(typeID)
	if !ok {
		return true
	}
	if _, ok := p.Declarations.Find(typ.Declaration); !ok {
		return false
	}
	if typ.Kind == TypeNominal && typ.Name != "" {
		resolved, resolvedOK := p.ResolveTypeName(scope, typ.Name)
		return resolvedOK && resolved == typeID
	}
	return true
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
