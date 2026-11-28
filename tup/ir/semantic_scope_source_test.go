package ir

import "testing"

func TestIndexedProgramRejectsScopeFromForeignModuleFile(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{{ID: ModuleID(1), Name: "main", Files: []FileID{FileID(1)}, Origin: testOrigin()}},
		Scopes:  ScopeTable{{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule, Origin: Origin{Primary: Span{File: FileID(2), Start: 0, End: 1}}}},
	}
	if program.Valid() {
		t.Fatal("scope from a foreign module file must be invalid")
	}
}
