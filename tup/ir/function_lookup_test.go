package ir

import "testing"

func TestIndexedProgramFindsFunctionsBySignatureKey(t *testing.T) {
	program := IndexedProgram{
		Declarations: DeclarationTable{
			{ID: DeclID(1), Scope: ScopeID(1), Name: "parse", Kind: DeclFunction, Origin: testOrigin()},
			{ID: DeclID(2), Scope: ScopeID(1), Name: "parse", Kind: DeclFunction, Origin: testOrigin()},
		},
		Functions: FunctionTable{
			{ID: FunctionID(1), Declaration: DeclID(1), Parameters: []TypeID{TypeID(10)}, Effect: EffectPure, Origin: testOrigin()},
			{ID: FunctionID(2), Declaration: DeclID(2), Parameters: []TypeID{TypeID(11)}, Effect: EffectPure, Origin: testOrigin()},
		},
	}
	key, ok := program.FunctionKey(FunctionID(2))
	if !ok || key.Scope != ScopeID(1) || key.Name != "parse" || len(key.Parameters) != 1 || key.Parameters[0] != TypeID(11) {
		t.Fatal("function key must preserve scope, name, and parameter types")
	}
	matches := program.FindFunctions(key)
	if len(matches) != 1 || matches[0] != FunctionID(2) {
		t.Fatal("function lookup must return the matching overload")
	}
}
