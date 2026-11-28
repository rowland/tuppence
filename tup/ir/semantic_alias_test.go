package ir

import "testing"

func TestIndexedProgramRejectsAliasWithWrongDeclarationKind(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{{ID: ModuleID(1), Name: "main", Files: []FileID{FileID(1)}, Origin: testOrigin()}},
		Types:   TypeTable{{ID: TypeID(1), Name: "Target", Kind: TypeNominal, Declaration: DeclID(1), HomeModule: ModuleID(1), Origin: testOrigin()}},
		Scopes:  ScopeTable{{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule, Origin: testOrigin()}},
		Declarations: DeclarationTable{
			{ID: DeclID(1), Scope: ScopeID(1), Name: "Target", Kind: DeclType, Origin: testOrigin()},
			{ID: DeclID(2), Scope: ScopeID(1), Name: "Alias", Kind: DeclValue, Origin: testOrigin()},
		},
		Aliases: TypeAliasTable{{Declaration: DeclID(2), Target: TypeID(1), Origin: testOrigin()}},
	}
	if program.Valid() {
		t.Fatal("alias backed by a non-type declaration must be invalid")
	}
}

func TestIndexedProgramRejectsDanglingAliasTarget(t *testing.T) {
	program := IndexedProgram{
		Modules:      ModuleTable{{ID: ModuleID(1), Name: "main", Files: []FileID{FileID(1)}, Origin: testOrigin()}},
		Scopes:       ScopeTable{{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule, Origin: testOrigin()}},
		Declarations: DeclarationTable{{ID: DeclID(1), Scope: ScopeID(1), Name: "Alias", Kind: DeclType, Origin: testOrigin()}},
		Aliases:      TypeAliasTable{{Declaration: DeclID(1), Target: TypeID(9), Origin: testOrigin()}},
	}
	if program.Valid() {
		t.Fatal("alias with an absent target must be invalid")
	}
}

func TestIndexedProgramRejectsAliasAlsoRepresentedAsType(t *testing.T) {
	program := IndexedProgram{
		Modules:      ModuleTable{{ID: ModuleID(1), Name: "main", Files: []FileID{FileID(1)}, Origin: testOrigin()}},
		Types:        TypeTable{{ID: TypeID(1), Name: "Alias", Kind: TypeNominal, Declaration: DeclID(1), HomeModule: ModuleID(1), Origin: testOrigin()}},
		Scopes:       ScopeTable{{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule, Origin: testOrigin()}},
		Declarations: DeclarationTable{{ID: DeclID(1), Scope: ScopeID(1), Name: "Alias", Kind: DeclType, Origin: testOrigin()}},
		Aliases:      TypeAliasTable{{Declaration: DeclID(1), Target: TypeID(1), Origin: testOrigin()}},
	}
	if program.Valid() {
		t.Fatal("declaration represented as both alias and type must be invalid")
	}
}
