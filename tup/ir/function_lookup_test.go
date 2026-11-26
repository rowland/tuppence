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
