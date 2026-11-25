package ir

import "testing"

func TestIndexedProgramRejectsPrivateCrossModuleReferences(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{
			{ID: ModuleID(1), Name: "main", Imports: []ImportBinding{{LocalName: "lib", Target: ModuleID(2), Origin: testOrigin()}}, Origin: testOrigin()},
			{ID: ModuleID(2), Name: "lib", Origin: testOrigin()},
		},
		Scopes: ScopeTable{
			{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule, Origin: testOrigin()},
			{ID: ScopeID(2), Module: ModuleID(2), Kind: ScopeModule, Origin: testOrigin()},
		},
		Declarations: DeclarationTable{
			{ID: DeclID(1), Scope: ScopeID(1), Name: "main", Kind: DeclValue, Origin: testOrigin()},
			{ID: DeclID(2), Scope: ScopeID(2), Name: "secret", Kind: DeclValue, Origin: testOrigin()},
		},
		Expressions: ExpressionTable{{ID: ExprID(1), Scope: ScopeID(1), Declaration: DeclID(2), Origin: testOrigin()}},
	}
	if program.Valid() {
		t.Fatal("cross-module reference to a private declaration must be invalid")
	}
}

func TestIndexedProgramAllowsImportedExportedReferences(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{
			{ID: ModuleID(1), Name: "main", Imports: []ImportBinding{{LocalName: "lib", Target: ModuleID(2), Origin: testOrigin()}}, Origin: testOrigin()},
			{ID: ModuleID(2), Name: "lib", Origin: testOrigin()},
		},
		Scopes: ScopeTable{
			{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule, Origin: testOrigin()},
			{ID: ScopeID(2), Module: ModuleID(2), Kind: ScopeModule, Origin: testOrigin()},
		},
		Declarations: DeclarationTable{
			{ID: DeclID(1), Scope: ScopeID(1), Name: "main", Kind: DeclValue, Origin: testOrigin()},
			{ID: DeclID(2), Scope: ScopeID(2), Name: "answer", Kind: DeclValue, Exported: true, Origin: testOrigin()},
		},
		Expressions: ExpressionTable{{ID: ExprID(1), Scope: ScopeID(1), Declaration: DeclID(2), Origin: testOrigin()}},
	}
	if !program.Valid() {
		t.Fatal("imported exported reference must be valid")
	}
}
