package ir

import "testing"

func TestIndexedProgramFindsCallCandidatesWithNamedArguments(t *testing.T) {
	program := IndexedProgram{
		Declarations: DeclarationTable{
			{ID: DeclID(1), Scope: ScopeID(1), Name: "parse", Kind: DeclFunction, Origin: testOrigin()},
			{ID: DeclID(2), Scope: ScopeID(1), Name: "parse", Kind: DeclFunction, Origin: testOrigin()},
		},
		Functions: FunctionTable{
			{ID: FunctionID(1), Declaration: DeclID(1), Parameters: []FunctionParameter{{Name: "value", Type: TypeID(10)}, {Name: "base", Type: TypeID(11)}}, Effect: EffectPure, Origin: testOrigin()},
			{ID: FunctionID(2), Declaration: DeclID(2), Parameters: []FunctionParameter{{Name: "value", Type: TypeID(10)}, {Name: "base", Type: TypeID(12)}}, Effect: EffectPure, Origin: testOrigin()},
		},
	}
	call := CallExpr{ID: ExprID(1), Scope: ScopeID(1), Name: "parse", Arguments: []CallArgument{{Type: TypeID(10)}, {Name: "base", Type: TypeID(11)}}, Origin: testOrigin()}
	matches := program.FindCallCandidates(call)
	if len(matches) != 1 || matches[0] != FunctionID(1) {
		t.Fatal("call candidate lookup must match positional and named argument types")
	}
}

func TestIndexedProgramFindsCurriedCallCandidates(t *testing.T) {
	program := IndexedProgram{
		Declarations: DeclarationTable{{ID: DeclID(1), Scope: ScopeID(1), Name: "parse", Kind: DeclFunction, Origin: testOrigin()}},
		Functions:    FunctionTable{{ID: FunctionID(1), Declaration: DeclID(1), Parameters: []FunctionParameter{{Name: "value", Type: TypeID(10)}, {Name: "base", Type: TypeID(11)}}, Effect: EffectPure, Origin: testOrigin()}},
	}
	call := CallExpr{ID: ExprID(1), Scope: ScopeID(1), Name: "parse", Arguments: []CallArgument{{Type: TypeID(10)}}, Curry: true, Origin: testOrigin()}
	matches := program.FindCallCandidates(call)
	if len(matches) != 1 || matches[0] != FunctionID(1) {
		t.Fatal("curried call candidate lookup must allow partial application")
	}
}

func TestIndexedProgramFindsFunctionsBySignatureKey(t *testing.T) {
	program := IndexedProgram{
		Declarations: DeclarationTable{
			{ID: DeclID(1), Scope: ScopeID(1), Name: "parse", Kind: DeclFunction, Origin: testOrigin()},
			{ID: DeclID(2), Scope: ScopeID(1), Name: "parse", Kind: DeclFunction, Origin: testOrigin()},
		},
		Functions: FunctionTable{
			{ID: FunctionID(1), Declaration: DeclID(1), Parameters: []FunctionParameter{{Name: "value", Type: TypeID(10)}}, Effect: EffectPure, Origin: testOrigin()},
			{ID: FunctionID(2), Declaration: DeclID(2), Parameters: []FunctionParameter{{Name: "value", Type: TypeID(11)}}, Effect: EffectPure, Origin: testOrigin()},
		},
	}
	key, ok := program.FunctionKey(FunctionID(2))
	if !ok || key.Scope != ScopeID(1) || key.Name != "parse" || len(key.Parameters) != 1 || key.Parameters[0].Type != TypeID(11) {
		t.Fatal("function key must preserve scope, name, and parameter types")
	}
	matches := program.FindFunctions(key)
	if len(matches) != 1 || matches[0] != FunctionID(2) {
		t.Fatal("function lookup must return the matching overload")
	}
}

func TestIndexedProgramMatchesLiteralThroughTypeConstructor(t *testing.T) {
	program := IndexedProgram{
		Constructors: TypeConstructorTable{{Type: TypeID(20), Storage: StorageI32, Origin: testOrigin()}},
		Declarations: DeclarationTable{{ID: DeclID(1), Scope: ScopeID(1), Name: "parse", Kind: DeclFunction, Origin: testOrigin()}},
		Functions:    FunctionTable{{ID: FunctionID(1), Declaration: DeclID(1), Parameters: []FunctionParameter{{Name: "value", Type: TypeID(20)}}, Effect: EffectPure, Origin: testOrigin()}},
	}
	call := CallExpr{ID: ExprID(1), Scope: ScopeID(1), Name: "parse", Arguments: []CallArgument{{Literal: LiteralInteger, LiteralTemplate: true}}, Origin: testOrigin()}
	matches := program.FindCallCandidates(call)
	if len(matches) != 1 || matches[0] != FunctionID(1) {
		t.Fatal("literal argument must match a parameter with an admissible storage constructor")
	}
}

func TestIndexedProgramMarksAmbiguousLiteralCandidates(t *testing.T) {
	program := IndexedProgram{
		Constructors: TypeConstructorTable{
			{Type: TypeID(20), Storage: StorageI32, Origin: testOrigin()},
			{Type: TypeID(21), Storage: StorageI64, Origin: testOrigin()},
		},
		Declarations: DeclarationTable{
			{ID: DeclID(1), Scope: ScopeID(1), Name: "take", Kind: DeclFunction, Origin: testOrigin()},
			{ID: DeclID(2), Scope: ScopeID(1), Name: "take", Kind: DeclFunction, Origin: testOrigin()},
		},
		Functions: FunctionTable{
			{ID: FunctionID(1), Declaration: DeclID(1), Parameters: []FunctionParameter{{Name: "value", Type: TypeID(20)}}, Effect: EffectPure, Origin: testOrigin()},
			{ID: FunctionID(2), Declaration: DeclID(2), Parameters: []FunctionParameter{{Name: "value", Type: TypeID(21)}}, Effect: EffectPure, Origin: testOrigin()},
		},
	}
	call := CallExpr{ID: ExprID(1), Scope: ScopeID(1), Name: "take", Arguments: []CallArgument{{Literal: LiteralInteger}}, Origin: testOrigin()}
	result := program.FindCallCandidatesWithResolution(call)
	if len(result.Functions) != 2 || !result.Call.Arguments[0].LiteralAmbiguous {
		t.Fatal("multiple literal candidates must mark the literal ambiguous")
	}
	if call.Arguments[0].LiteralAmbiguous {
		t.Fatal("candidate discovery must not mutate the original call")
	}
}

func TestIndexedProgramFindsNonLossyCastFunctions(t *testing.T) {
	program := IndexedProgram{
		Declarations: DeclarationTable{
			{ID: DeclID(1), Scope: ScopeID(1), Name: "cast", Kind: DeclFunction, Origin: testOrigin()},
			{ID: DeclID(2), Scope: ScopeID(1), Name: "not_cast", Kind: DeclFunction, Origin: testOrigin()},
			{ID: DeclID(3), Scope: ScopeID(2), Name: "cast", Kind: DeclFunction, Origin: testOrigin()},
		},
		Functions: FunctionTable{
			{ID: FunctionID(1), Declaration: DeclID(1), Parameters: []FunctionParameter{{Name: "value", Type: TypeID(10)}}, Result: TypeID(11), Effect: EffectPure, Cast: true, Origin: testOrigin()},
			{ID: FunctionID(2), Declaration: DeclID(2), Parameters: []FunctionParameter{{Name: "value", Type: TypeID(10)}}, Result: TypeID(11), Effect: EffectPure, Origin: testOrigin()},
			{ID: FunctionID(3), Declaration: DeclID(3), Parameters: []FunctionParameter{{Name: "value", Type: TypeID(10)}}, Result: TypeID(11), Effect: EffectPure, Cast: true, Origin: testOrigin()},
		},
	}
	matches := program.FindCastFunctions(ScopeID(1), TypeID(10), TypeID(11))
	if len(matches) != 1 || matches[0] != FunctionID(1) {
		t.Fatal("cast lookup must return only marked casts in the requested scope")
	}
}

func TestIndexedProgramResolvesInsertedCast(t *testing.T) {
	program := IndexedProgram{
		Declarations: DeclarationTable{
			{ID: DeclID(1), Scope: ScopeID(1), Name: "use_int64", Kind: DeclFunction, Origin: testOrigin()},
			{ID: DeclID(2), Scope: ScopeID(1), Name: "cast", Kind: DeclFunction, Origin: testOrigin()},
		},
		Functions: FunctionTable{
			{ID: FunctionID(1), Declaration: DeclID(1), Parameters: []FunctionParameter{{Name: "value", Type: TypeID(11)}}, Effect: EffectPure, Origin: testOrigin()},
			{ID: FunctionID(2), Declaration: DeclID(2), Parameters: []FunctionParameter{{Name: "value", Type: TypeID(10)}}, Result: TypeID(11), Effect: EffectPure, Cast: true, Origin: testOrigin()},
		},
	}
	call := CallExpr{ID: ExprID(1), Scope: ScopeID(1), Name: "use_int64", Arguments: []CallArgument{{Type: TypeID(10)}}, Origin: testOrigin()}
	plan, ok := program.ResolveCallPlan(call, FunctionID(1))
	if !ok || len(plan.Arguments) != 1 || plan.Arguments[0].Adaptation != AdaptCast || plan.Arguments[0].Cast != FunctionID(2) {
		t.Fatal("resolved call plan must record the inserted cast function")
	}
}

func TestIndexedProgramRejectsCallWithoutUniqueCast(t *testing.T) {
	program := IndexedProgram{
		Declarations: DeclarationTable{{ID: DeclID(1), Scope: ScopeID(1), Name: "use_int64", Kind: DeclFunction, Origin: testOrigin()}},
		Functions:    FunctionTable{{ID: FunctionID(1), Declaration: DeclID(1), Parameters: []FunctionParameter{{Name: "value", Type: TypeID(11)}}, Effect: EffectPure, Origin: testOrigin()}},
	}
	call := CallExpr{ID: ExprID(1), Scope: ScopeID(1), Name: "use_int64", Arguments: []CallArgument{{Type: TypeID(10)}}, Origin: testOrigin()}
	if _, ok := program.ResolveCallPlan(call, FunctionID(1)); ok {
		t.Fatal("call plan must reject a typed mismatch without a unique cast function")
	}
}

func TestIndexedProgramPrefersExactCandidateTier(t *testing.T) {
	program := IndexedProgram{
		Declarations: DeclarationTable{
			{ID: DeclID(1), Scope: ScopeID(1), Name: "take", Kind: DeclFunction, Origin: testOrigin()},
			{ID: DeclID(2), Scope: ScopeID(1), Name: "take", Kind: DeclFunction, Origin: testOrigin()},
			{ID: DeclID(3), Scope: ScopeID(1), Name: "take", Kind: DeclFunction, Origin: testOrigin()},
		},
		Functions: FunctionTable{
			{ID: FunctionID(1), Declaration: DeclID(1), Parameters: []FunctionParameter{{Name: "value", Type: TypeID(10)}}, Effect: EffectPure, Generic: true, Origin: testOrigin()},
			{ID: FunctionID(2), Declaration: DeclID(2), Parameters: []FunctionParameter{{Name: "value", Type: TypeID(10)}}, Effect: EffectPure, Any: true, Origin: testOrigin()},
			{ID: FunctionID(3), Declaration: DeclID(3), Parameters: []FunctionParameter{{Name: "value", Type: TypeID(10)}}, Effect: EffectPure, Origin: testOrigin()},
		},
	}
	call := CallExpr{ID: ExprID(1), Scope: ScopeID(1), Name: "take", Arguments: []CallArgument{{Type: TypeID(10)}}, Origin: testOrigin()}
	result := program.FindCallCandidatesWithResolution(call)
	if len(result.Functions) != 1 || result.Functions[0] != FunctionID(3) {
		t.Fatal("exact candidates must displace generic and any candidates")
	}
}

func TestIndexedProgramFindsOuterLexicalFunction(t *testing.T) {
	program := IndexedProgram{
		Scopes:       ScopeTable{{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule, Origin: testOrigin()}, {ID: ScopeID(2), Parent: ScopeID(1), Module: ModuleID(1), Kind: ScopeBlock, Origin: testOrigin()}},
		Declarations: DeclarationTable{{ID: DeclID(1), Scope: ScopeID(1), Name: "take", Kind: DeclFunction, Origin: testOrigin()}},
		Functions:    FunctionTable{{ID: FunctionID(1), Declaration: DeclID(1), Parameters: []FunctionParameter{{Name: "value", Type: TypeID(10)}}, Effect: EffectPure, Origin: testOrigin()}},
	}
	call := CallExpr{ID: ExprID(1), Scope: ScopeID(2), Name: "take", Arguments: []CallArgument{{Type: TypeID(10)}}, Origin: testOrigin()}
	if matches := program.FindCallCandidates(call); len(matches) != 1 || matches[0] != FunctionID(1) {
		t.Fatal("call lookup must search outward through lexical scopes")
	}
}

func TestIndexedProgramFindsQualifiedAliasAsUnqualifiedFunction(t *testing.T) {
	program := IndexedProgram{
		ImportAliases: ImportAliasTable{{Scope: ScopeID(1), LocalName: "cast", Module: ModuleID(2), Declaration: DeclID(1), Origin: testOrigin()}},
		Scopes:        ScopeTable{{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule, Origin: testOrigin()}, {ID: ScopeID(2), Module: ModuleID(2), Kind: ScopeModule, Origin: testOrigin()}},
		Declarations:  DeclarationTable{{ID: DeclID(1), Scope: ScopeID(2), Name: "cast", Kind: DeclFunction, Exported: true, Origin: testOrigin()}},
		Functions:     FunctionTable{{ID: FunctionID(1), Declaration: DeclID(1), Parameters: []FunctionParameter{{Name: "value", Type: TypeID(10)}}, Result: TypeID(11), Effect: EffectPure, Cast: true, Origin: testOrigin()}},
	}
	call := CallExpr{ID: ExprID(1), Scope: ScopeID(1), Name: "cast", Arguments: []CallArgument{{Type: TypeID(10)}}, Origin: testOrigin()}
	if matches := program.FindCastFunctions(ScopeID(1), TypeID(10), TypeID(11)); len(matches) != 1 || matches[0] != FunctionID(1) {
		t.Fatal("local import aliases must expose cast functions to the current scope")
	}
	if matches := program.FindCallCandidates(call); len(matches) != 1 || matches[0] != FunctionID(1) {
		t.Fatal("local import aliases must expose functions to unqualified lookup")
	}
}

func TestIndexedProgramRejectsUnqualifiedImportedFunction(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{
			{ID: ModuleID(1), Name: "main", Files: []FileID{FileID(1)}, Imports: []ImportBinding{{LocalName: "core", Target: ModuleID(2), Origin: testOrigin()}}, Origin: testOrigin()},
			{ID: ModuleID(2), Name: "core", Files: []FileID{FileID(2)}, Origin: testOrigin()},
		},
		Scopes:       ScopeTable{{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule, Origin: testOrigin()}, {ID: ScopeID(2), Module: ModuleID(2), Kind: ScopeModule, Origin: testOrigin()}},
		Declarations: DeclarationTable{{ID: DeclID(1), Scope: ScopeID(2), Name: "take", Kind: DeclFunction, Exported: true, Origin: testOrigin()}},
		Functions:    FunctionTable{{ID: FunctionID(1), Declaration: DeclID(1), Parameters: []FunctionParameter{{Name: "value", Type: TypeID(10)}}, Effect: EffectPure, Origin: testOrigin()}},
	}
	call := CallExpr{ID: ExprID(1), Scope: ScopeID(1), Name: "take", Arguments: []CallArgument{{Type: TypeID(10)}}, Origin: testOrigin()}
	if matches := program.FindCallCandidates(call); len(matches) != 0 {
		t.Fatal("unqualified lookup must not expose functions from an imported module")
	}
	qualified := call
	qualified.Qualifier = "core"
	if matches := program.FindCallCandidates(qualified); len(matches) != 1 || matches[0] != FunctionID(1) {
		t.Fatal("qualified lookup must resolve exported functions through the local module alias")
	}
}
