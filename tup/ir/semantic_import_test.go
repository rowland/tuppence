package ir

import "testing"

func TestIndexedProgramRejectsDanglingModuleImports(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{
			{ID: ModuleID(1), Name: "main", Files: []FileID{FileID(1)}, Imports: []ImportBinding{{LocalName: "io", Path: "io", Target: ModuleID(9), Origin: testOrigin()}}, Origin: testOrigin()},
		},
	}
	if program.Valid() {
		t.Fatal("import referring to an absent module must be invalid")
	}
}

func TestIndexedProgramRejectsImportNameCollision(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{
			{ID: ModuleID(1), Name: "main", Files: []FileID{FileID(1)}, Imports: []ImportBinding{{LocalName: "io", Target: ModuleID(2), Origin: testOrigin()}}, Origin: testOrigin()},
			{ID: ModuleID(2), Name: "io", Origin: testOrigin()},
		},
		Scopes: ScopeTable{
			{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule, Origin: testOrigin(), Entries: []ScopeEntry{{Name: "io", Declaration: DeclID(1)}}},
		},
		Declarations: DeclarationTable{{ID: DeclID(1), Scope: ScopeID(1), Name: "io", Kind: DeclValue, Origin: testOrigin()}},
	}
	if program.Valid() {
		t.Fatal("import colliding with a module declaration must be invalid")
	}
}
