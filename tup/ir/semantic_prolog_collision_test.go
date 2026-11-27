package ir

import "testing"

func TestIndexedProgramAllowsNestedShadowingOfPrologAlias(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{
			{ID: ModuleID(1), Name: "main", Files: []FileID{FileID(1)}, Origin: testOrigin()},
			{ID: ModuleID(2), Name: "core", Files: []FileID{FileID(2)}, Origin: Origin{Primary: Span{File: FileID(2), Start: 0, End: 1}}},
		},
		Scopes: ScopeTable{
			{ID: ScopeID(0x80000001), Module: ModuleID(1), Kind: ScopeProlog, Origin: testOrigin()},
			{ID: ScopeID(1), Parent: ScopeID(0x80000001), Module: ModuleID(1), Kind: ScopeModule, Initial: true, Origin: testOrigin()},
			{ID: ScopeID(2), Parent: ScopeID(1), Module: ModuleID(1), Kind: ScopeBlock, Entries: []ScopeEntry{{Name: "Int", Declaration: DeclID(2)}}, Origin: testOrigin()},
			{ID: ScopeID(3), Module: ModuleID(2), Kind: ScopeModule, Origin: Origin{Primary: Span{File: FileID(2), Start: 0, End: 1}}},
		},
		Declarations: DeclarationTable{
			{ID: DeclID(2), Scope: ScopeID(2), Name: "Int", Kind: DeclType, Origin: testOrigin()},
			{ID: DeclID(3), Scope: ScopeID(3), Name: "Int", Kind: DeclType, Exported: true, Origin: Origin{Primary: Span{File: FileID(2), Start: 0, End: 1}}},
		},
		ImportAliases: ImportAliasTable{{Scope: ScopeID(0x80000001), LocalName: "Int", Module: ModuleID(2), Declaration: DeclID(3), Source: ImportAliasSourceProlog, Origin: testOrigin()}},
		Prolog:        PrologConfig{Enabled: true, Modules: []ModuleID{ModuleID(2)}, Architecture: "amd64"},
	}
	if !program.Valid() {
		t.Fatal("nested user declaration may shadow a prolog alias")
	}
}

func TestIndexedProgramAllowsModuleShadowingOfPrologAlias(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{
			{ID: ModuleID(1), Name: "main", Files: []FileID{FileID(1)}, Origin: testOrigin()},
			{ID: ModuleID(2), Name: "core", Files: []FileID{FileID(2)}, Origin: Origin{Primary: Span{File: FileID(2), Start: 0, End: 1}}},
		},
		Scopes: ScopeTable{
			{ID: ScopeID(0x80000001), Module: ModuleID(1), Kind: ScopeProlog, Origin: testOrigin()},
			{ID: ScopeID(1), Parent: ScopeID(0x80000001), Module: ModuleID(1), Kind: ScopeModule, Initial: true, Entries: []ScopeEntry{{Name: "Int", Declaration: DeclID(2)}}, Origin: testOrigin()},
			{ID: ScopeID(2), Module: ModuleID(2), Kind: ScopeModule, Origin: Origin{Primary: Span{File: FileID(2), Start: 0, End: 1}}},
		},
		Declarations: DeclarationTable{
			{ID: DeclID(2), Scope: ScopeID(1), Name: "Int", Kind: DeclType, Origin: testOrigin()},
			{ID: DeclID(3), Scope: ScopeID(2), Name: "Int", Kind: DeclType, Exported: true, Origin: Origin{Primary: Span{File: FileID(2), Start: 0, End: 1}}},
		},
		ImportAliases: ImportAliasTable{{Scope: ScopeID(0x80000001), LocalName: "Int", Module: ModuleID(2), Declaration: DeclID(3), Source: ImportAliasSourceProlog, Origin: testOrigin()}},
		Prolog:        PrologConfig{Enabled: true, Modules: []ModuleID{ModuleID(2)}, Architecture: "amd64"},
	}
	if !program.Valid() {
		t.Fatal("module declaration may shadow a prolog alias")
	}
}
