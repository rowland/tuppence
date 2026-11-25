package ir

import "testing"

func TestIndexedProgramRejectsTypeFromForeignModuleFile(t *testing.T) {
	program := IndexedProgram{
		Modules:      ModuleTable{{ID: ModuleID(1), Name: "main", Files: []FileID{FileID(1)}, Origin: testOrigin()}},
		Scopes:       ScopeTable{{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule, Origin: testOrigin()}},
		Declarations: DeclarationTable{{ID: DeclID(1), Scope: ScopeID(1), Name: "Answer", Kind: DeclType, Origin: testOrigin()}},
		Types:        TypeTable{{ID: TypeID(1), Name: "Answer", Kind: TypeNominal, Declaration: DeclID(1), HomeModule: ModuleID(1), Origin: Origin{Primary: Span{File: FileID(2), Start: 0, End: 6}}}},
	}
	if program.Valid() {
		t.Fatal("type from a foreign module file must be invalid")
	}
}
