package ir

import "testing"

func TestIndexedProgramRejectsDanglingBindingExpressions(t *testing.T) {
	program := IndexedProgram{
		Modules:  ModuleTable{{ID: ModuleID(1), Name: "main", Origin: testOrigin()}},
		Scopes:   ScopeTable{{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule}},
		Bindings: BindingTable{{ID: BindingID(1), Scope: ScopeID(1), Name: "answer", TypeSlot: TypeVarID(1), Value: ExprID(9), Origin: testOrigin()}},
	}
	if program.Valid() {
		t.Fatal("binding with an absent expression must be invalid")
	}
}
