package ir

import "testing"

func TestIndexedProgramRejectsBindingFromForeignModuleFile(t *testing.T) {
	program := IndexedProgram{
		Modules:      ModuleTable{{ID: ModuleID(1), Name: "main", Files: []FileID{FileID(1)}, Origin: testOrigin()}},
		Scopes:       ScopeTable{{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule, Origin: testOrigin()}},
		Declarations: DeclarationTable{{ID: DeclID(1), Scope: ScopeID(1), Name: "answer", Kind: DeclValue, Origin: testOrigin()}},
		Expressions:  ExpressionTable{{ID: ExprID(1), Scope: ScopeID(1), Declaration: DeclID(1), Origin: testOrigin()}},
		Bindings:     BindingTable{{ID: BindingID(1), Scope: ScopeID(1), Name: "answer", TypeSlot: TypeVarID(1), Value: ExprID(1), Origin: Origin{Primary: Span{File: FileID(2), Start: 0, End: 6}}}},
	}
	if program.Valid() {
		t.Fatal("binding from a foreign module file must be invalid")
	}
}
