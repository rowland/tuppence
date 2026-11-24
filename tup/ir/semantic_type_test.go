package ir

import "testing"

func TestIndexedProgramRejectsNominalTypeWithWrongDeclarationKind(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{
			{ID: ModuleID(1), Name: "bar", Origin: testOrigin()},
		},
		Types: TypeTable{
			{ID: TypeID(1), Name: "T", Kind: TypeNominal, Declaration: DeclID(1), HomeModule: ModuleID(1), Origin: testOrigin()},
		},
		Scopes: ScopeTable{
			{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule},
		},
		Declarations: DeclarationTable{
			{ID: DeclID(1), Scope: ScopeID(1), Name: "T", Kind: DeclValue, Origin: testOrigin()},
		},
	}
	if program.Valid() {
		t.Fatal("nominal type backed by a value declaration must be invalid")
	}
}

func TestIndexedProgramRejectsNominalTypeWithWrongHomeModule(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{
			{ID: ModuleID(1), Name: "bar", Origin: testOrigin()},
			{ID: ModuleID(2), Name: "other", Origin: testOrigin()},
		},
		Types: TypeTable{
			{ID: TypeID(1), Name: "T", Kind: TypeNominal, Declaration: DeclID(1), HomeModule: ModuleID(2), Origin: testOrigin()},
		},
		Scopes: ScopeTable{
			{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule},
		},
		Declarations: DeclarationTable{
			{ID: DeclID(1), Scope: ScopeID(1), Name: "T", Kind: DeclType, Origin: testOrigin()},
		},
	}
	if program.Valid() {
		t.Fatal("nominal type with a mismatched home module must be invalid")
	}
}
