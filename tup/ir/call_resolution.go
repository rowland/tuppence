package ir

// ArgumentAdaptation identifies how a source argument reaches a formal parameter.
type ArgumentAdaptation uint8

const (
	AdaptExact ArgumentAdaptation = iota + 1
	AdaptCast
	AdaptLiteral
)

// ResolvedCallArgument records one source-to-parameter mapping after candidate
// selection. Cast identifies an existing cast function inserted by the compiler.
type ResolvedCallArgument struct {
	Parameter    int
	Source       CallArgument
	Adaptation   ArgumentAdaptation
	Cast         FunctionID
	Storage      StorageKind
	ResolvedType TypeID
}

// LiteralDefaults identifies the architecture-selected semantic defaults for
// ambiguous integer and float literals.
type LiteralDefaults struct {
	Int   TypeID
	Float TypeID
}

type ResolvedCallPlan struct {
	Function  FunctionID
	Arguments []ResolvedCallArgument
}

// ResolveCallPlan resolves argument mappings for one already-selected function.
// It inserts only a unique suitable existing cast function; it does not rank
// overloads or synthesize cast functions.
func (p IndexedProgram) ResolveCallPlan(call CallExpr, functionID FunctionID) (ResolvedCallPlan, bool) {
	function, ok := p.Functions.Find(functionID)
	if !ok || call.Name == "" {
		return ResolvedCallPlan{}, false
	}
	declaration, ok := p.Declarations.Find(function.Declaration)
	if !ok || declaration.Scope != call.Scope || declaration.Name != call.Name {
		return ResolvedCallPlan{}, false
	}
	if !call.Curry && len(call.Arguments) != len(function.Parameters) || len(call.Arguments) > len(function.Parameters) {
		return ResolvedCallPlan{}, false
	}

	matched := make([]bool, len(function.Parameters))
	positional := 0
	plan := ResolvedCallPlan{Function: functionID, Arguments: make([]ResolvedCallArgument, 0, len(call.Arguments))}
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
		if index < 0 || index >= len(function.Parameters) || matched[index] {
			return ResolvedCallPlan{}, false
		}

		parameter := function.Parameters[index]
		resolved := ResolvedCallArgument{Parameter: index, Source: argument}
		if argument.Type.Valid() && argument.Type == parameter.Type {
			resolved.Adaptation = AdaptExact
		} else if argument.Type.Valid() {
			casts := p.FindCastFunctions(call.Scope, argument.Type, parameter.Type)
			if len(casts) != 1 {
				return ResolvedCallPlan{}, false
			}
			resolved.Adaptation = AdaptCast
			resolved.Cast = casts[0]
		} else if argument.Literal.Valid() && argument.LiteralAmbiguous {
			defaultType := p.LiteralDefaults.Int
			if defaultKind, ok := DefaultLiteralKind(argument.Literal, argument.LiteralOutOfIntegerRange); !ok {
				return ResolvedCallPlan{}, false
			} else if defaultKind == LiteralDefaultFloat {
				defaultType = p.LiteralDefaults.Float
			}
			if !defaultType.Valid() {
				return ResolvedCallPlan{}, false
			}
			resolved.ResolvedType = defaultType
			if defaultType == parameter.Type {
				resolved.Adaptation = AdaptLiteral
			} else {
				casts := p.FindCastFunctions(call.Scope, defaultType, parameter.Type)
				if len(casts) != 1 {
					return ResolvedCallPlan{}, false
				}
				resolved.Adaptation = AdaptCast
				resolved.Cast = casts[0]
			}
		} else if argument.Literal.Valid() {
			var admissible []StorageKind
			for _, constructor := range p.Constructors {
				if constructor.Type == parameter.Type {
					admissible = append(admissible, constructor.Storage)
				}
			}
			storage, ok := SelectLiteralStorage(argument.Literal, admissible)
			if !ok {
				return ResolvedCallPlan{}, false
			}
			resolved.Adaptation = AdaptLiteral
			resolved.Storage = storage
		} else {
			return ResolvedCallPlan{}, false
		}
		matched[index] = true
		plan.Arguments = append(plan.Arguments, resolved)
	}
	if !call.Curry {
		for _, present := range matched {
			if !present {
				return ResolvedCallPlan{}, false
			}
		}
	}
	return plan, true
}
