package ir

import "testing"

func TestIndexedProgramAllowsArrayElementDeclaredLater(t *testing.T) {
	program := IndexedProgram{
		Modules:      ModuleTable{{ID: ModuleID(1), Name: "main", Files: []FileID{FileID(1)}, Origin: testOrigin()}},
		Scopes:       ScopeTable{{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule, Origin: testOrigin()}},
		Declarations: DeclarationTable{{ID: DeclID(1), Scope: ScopeID(1), Name: "Item", Kind: DeclType, Origin: testOrigin()}},
		Types: TypeTable{
			{ID: TypeID(1), Kind: TypeArray, Element: TypeID(2), Origin: testOrigin()},
			{ID: TypeID(2), Name: "Item", Kind: TypeNominal, Declaration: DeclID(1), HomeModule: ModuleID(1), Origin: testOrigin()},
		},
	}
	if !program.Valid() {
		t.Fatal("array may reference a type declared later in the snapshot")
	}
}
