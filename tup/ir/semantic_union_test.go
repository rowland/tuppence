package ir

import "testing"

func TestIndexedProgramAllowsUnionAlternativesDeclaredLater(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{{ID: ModuleID(1), Name: "main", Files: []FileID{FileID(1)}, Origin: testOrigin()}},
		Scopes:  ScopeTable{{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule, Origin: testOrigin()}},
		Declarations: DeclarationTable{
			{ID: DeclID(1), Scope: ScopeID(1), Name: "Left", Kind: DeclType, Origin: testOrigin()},
			{ID: DeclID(2), Scope: ScopeID(1), Name: "Right", Kind: DeclType, Origin: testOrigin()},
		},
		Types: TypeTable{
			{ID: TypeID(1), Kind: TypeUnion, Elements: []TypeID{TypeID(2), TypeID(3)}, Origin: testOrigin()},
			{ID: TypeID(2), Name: "Left", Kind: TypeNominal, Declaration: DeclID(1), HomeModule: ModuleID(1), Origin: testOrigin()},
			{ID: TypeID(3), Name: "Right", Kind: TypeNominal, Declaration: DeclID(2), HomeModule: ModuleID(1), Origin: testOrigin()},
		},
	}
	if !program.Valid() {
		t.Fatal("union may reference types declared later in the snapshot")
	}
}
