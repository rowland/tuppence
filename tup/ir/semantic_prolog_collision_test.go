package ir

import "testing"

func TestIndexedProgramRejectsPrologAliasCollidingWithDeclaration(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{
			{ID: ModuleID(1), Name: "main", Files: []FileID{FileID(1)}, Origin: testOrigin()},
			{ID: ModuleID(2), Name: "core", Files: []FileID{FileID(2)}, Origin: Origin{Primary: Span{File: FileID(2), Start: 0, End: 1}}},
		},
		Scopes: ScopeTable{
			{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule, Initial: true, Entries: []ScopeEntry{{Name: "Int", Declaration: DeclID(2)}}, Origin: testOrigin()},
			{ID: ScopeID(2), Module: ModuleID(2), Kind: ScopeModule, Origin: Origin{Primary: Span{File: FileID(2), Start: 0, End: 1}}},
		},
		Declarations: DeclarationTable{
			{ID: DeclID(2), Scope: ScopeID(1), Name: "Int", Kind: DeclType, Origin: testOrigin()},
			{ID: DeclID(3), Scope: ScopeID(2), Name: "Int", Kind: DeclType, Exported: true, Origin: Origin{Primary: Span{File: FileID(2), Start: 0, End: 1}}},
		},
		ImportAliases: ImportAliasTable{{Scope: ScopeID(1), LocalName: "Int", Module: ModuleID(2), Declaration: DeclID(3), Source: ImportAliasSourceProlog, Origin: testOrigin()}},
		Prolog:        PrologConfig{Enabled: true, Modules: []ModuleID{ModuleID(2)}, Architecture: "amd64"},
	}
	if program.Valid() {
		t.Fatal("prolog alias colliding with a declaration must be invalid")
	}
}
