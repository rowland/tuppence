package ir

import "testing"

func TestIndexedProgramRejectsFunctionFromForeignModuleFile(t *testing.T) {
	program := IndexedProgram{
		Modules:      ModuleTable{{ID: ModuleID(1), Name: "main", Files: []FileID{FileID(1)}, Origin: testOrigin()}},
		Scopes:       ScopeTable{{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule}},
		Declarations: DeclarationTable{{ID: DeclID(1), Scope: ScopeID(1), Name: "compute", Kind: DeclFunction, Origin: testOrigin()}},
		Functions:    FunctionTable{{ID: FunctionID(1), Declaration: DeclID(1), Effect: EffectPure, Origin: Origin{Primary: Span{File: FileID(2), Start: 0, End: 7}}}},
	}
	if program.Valid() {
		t.Fatal("function from a foreign module file must be invalid")
	}
}
