package ir

import "testing"

func TestIndexedProgramRejectsDeclarationFromForeignModuleFile(t *testing.T) {
	program := IndexedProgram{
		Modules:      ModuleTable{{ID: ModuleID(1), Name: "main", Files: []FileID{FileID(1)}, Origin: testOrigin()}},
		Scopes:       ScopeTable{{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule}},
		Declarations: DeclarationTable{{ID: DeclID(1), Scope: ScopeID(1), Name: "answer", Kind: DeclValue, Origin: Origin{Primary: Span{File: FileID(2), Start: 0, End: 6}}}},
	}
	if program.Valid() {
		t.Fatal("declaration from a foreign module file must be invalid")
	}
}
