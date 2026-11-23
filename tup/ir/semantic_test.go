package ir

import "testing"

func TestIndexedProgramValidity(t *testing.T) {
	program := IndexedProgram{
		Scopes: ScopeTable{
			{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule},
		},
		Declarations: DeclarationTable{
			{ID: DeclID(1), Scope: ScopeID(1), Name: "answer", Kind: DeclValue, Origin: testOrigin()},
		},
		Bindings: BindingTable{
			{ID: BindingID(1), Scope: ScopeID(1), Name: "answer", TypeSlot: TypeVarID(1), Value: ExprID(1), Origin: testOrigin()},
		},
	}
	if !program.Valid() {
		t.Fatal("complete indexed program must be valid")
	}
}

func TestIndexedProgramRejectsInvalidTables(t *testing.T) {
	program := IndexedProgram{
		Scopes: ScopeTable{
			{ID: ScopeID(1), Module: ModuleID(1)},
		},
	}
	if program.Valid() {
		t.Fatal("indexed program with an invalid table must be invalid")
	}
}
