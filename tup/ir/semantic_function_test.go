package ir

import "testing"

func TestIndexedProgramRejectsFunctionWithWrongDeclarationKind(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{
			{ID: ModuleID(1), Name: "main", Origin: testOrigin()},
		},
		Scopes: ScopeTable{
			{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule},
		},
		Declarations: DeclarationTable{
			{ID: DeclID(1), Scope: ScopeID(1), Name: "answer", Kind: DeclValue, Origin: testOrigin()},
		},
		Functions: FunctionTable{
			{ID: FunctionID(1), Declaration: DeclID(1), Effect: EffectPure, Origin: testOrigin()},
		},
	}
	if program.Valid() {
		t.Fatal("function backed by a non-function declaration must be invalid")
	}
}

func TestIndexedProgramRejectsDanglingFunctionDeclarations(t *testing.T) {
	program := IndexedProgram{
		Functions: FunctionTable{
			{ID: FunctionID(1), Declaration: DeclID(9), Effect: EffectPure, Origin: testOrigin()},
		},
	}
	if program.Valid() {
		t.Fatal("function referring to an absent declaration must be invalid")
	}
}
