package ir

import "testing"

func TestIndexedProgramRejectsPrivateCrossModuleReferences(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{
			{ID: ModuleID(1), Name: "main", Files: []FileID{FileID(1)}, Imports: []ImportBinding{{LocalName: "lib", Path: "lib", Target: ModuleID(2), Origin: testOrigin()}}, Origin: testOrigin()},
			{ID: ModuleID(2), Name: "lib", Files: []FileID{FileID(2)}, Origin: Origin{Primary: Span{File: FileID(2), Start: 0, End: 1}}},
		},
		Scopes: ScopeTable{
			{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule, Origin: testOrigin()},
			{ID: ScopeID(2), Module: ModuleID(2), Kind: ScopeModule, Origin: Origin{Primary: Span{File: FileID(2), Start: 0, End: 1}}},
		},
		Declarations: DeclarationTable{
			{ID: DeclID(1), Scope: ScopeID(1), Name: "main", Kind: DeclValue, Origin: testOrigin()},
			{ID: DeclID(2), Scope: ScopeID(2), Name: "secret", Kind: DeclValue, Origin: Origin{Primary: Span{File: FileID(2), Start: 0, End: 1}}},
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
			{ID: ModuleID(1), Name: "main", Files: []FileID{FileID(1)}, Imports: []ImportBinding{{LocalName: "lib", Path: "lib", Target: ModuleID(2), Origin: testOrigin()}}, Origin: testOrigin()},
			{ID: ModuleID(2), Name: "lib", Files: []FileID{FileID(2)}, Origin: Origin{Primary: Span{File: FileID(2), Start: 0, End: 1}}},
		},
		Scopes: ScopeTable{
			{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule, Origin: testOrigin()},
			{ID: ScopeID(2), Module: ModuleID(2), Kind: ScopeModule, Origin: Origin{Primary: Span{File: FileID(2), Start: 0, End: 1}}},
		},
		Declarations: DeclarationTable{
			{ID: DeclID(1), Scope: ScopeID(1), Name: "main", Kind: DeclValue, Origin: testOrigin()},
			{ID: DeclID(2), Scope: ScopeID(2), Name: "answer", Kind: DeclValue, Exported: true, Origin: Origin{Primary: Span{File: FileID(2), Start: 0, End: 1}}},
		},
		Expressions: ExpressionTable{{ID: ExprID(1), Scope: ScopeID(1), Declaration: DeclID(2), Origin: testOrigin()}},
	}
	if !program.Valid() {
		t.Fatal("imported exported reference must be valid")
	}
}
