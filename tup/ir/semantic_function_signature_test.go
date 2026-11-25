package ir

import "testing"

func TestIndexedProgramAllowsFunctionTypesDeclaredLater(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{{ID: ModuleID(1), Name: "main", Files: []FileID{FileID(1)}, Origin: testOrigin()}},
		Scopes:  ScopeTable{{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule, Origin: testOrigin()}},
		Declarations: DeclarationTable{
			{ID: DeclID(1), Scope: ScopeID(1), Name: "compute", Kind: DeclFunction, Origin: testOrigin()},
			{ID: DeclID(2), Scope: ScopeID(1), Name: "Item", Kind: DeclType, Origin: testOrigin()},
		},
		Types:     TypeTable{{ID: TypeID(1), Name: "Item", Kind: TypeNominal, Declaration: DeclID(2), HomeModule: ModuleID(1), Origin: testOrigin()}},
		Functions: FunctionTable{{ID: FunctionID(1), Declaration: DeclID(1), Parameters: []TypeID{TypeID(1)}, Result: TypeID(1), Effect: EffectPure, Origin: testOrigin()}},
	}
	if !program.Valid() {
		t.Fatal("function may reference types declared later in the snapshot")
	}
}
