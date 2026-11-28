package ir

import "testing"

func TestIndexedProgramAllowsEnabledPrologAliasInModuleScope(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{
			{ID: ModuleID(1), Name: "main", Files: []FileID{FileID(1)}, Origin: testOrigin()},
			{ID: ModuleID(2), Name: "core", Files: []FileID{FileID(2)}, Origin: Origin{Primary: Span{File: FileID(2), Start: 0, End: 1}}},
		},
		Scopes: ScopeTable{
			{ID: ScopeID(0x80000001), Module: ModuleID(1), Kind: ScopeProlog, Origin: testOrigin()},
			{ID: ScopeID(1), Parent: ScopeID(0x80000001), Module: ModuleID(1), Kind: ScopeModule, Initial: true, Origin: testOrigin()},
			{ID: ScopeID(2), Module: ModuleID(2), Kind: ScopeModule, Origin: Origin{Primary: Span{File: FileID(2), Start: 0, End: 1}}},
		},
		Declarations: DeclarationTable{
			{ID: DeclID(1), Scope: ScopeID(2), Name: "Int", Kind: DeclType, Exported: true, Origin: Origin{Primary: Span{File: FileID(2), Start: 0, End: 1}}},
		},
		ImportAliases: ImportAliasTable{{Scope: ScopeID(0x80000001), LocalName: "Int", Module: ModuleID(2), Declaration: DeclID(1), Source: ImportAliasSourceProlog, Origin: testOrigin()}},
		Prolog:        PrologConfig{Enabled: true, Modules: []ModuleID{ModuleID(2)}, Architecture: "amd64"},
	}
	if !program.Valid() {
		t.Fatal("enabled prolog alias in a module scope must be valid")
	}
}

func TestIndexedProgramRejectsPrologAliasWhenDisabled(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{
			{ID: ModuleID(1), Name: "main", Files: []FileID{FileID(1)}, Origin: testOrigin()},
			{ID: ModuleID(2), Name: "core", Files: []FileID{FileID(2)}, Origin: Origin{Primary: Span{File: FileID(2), Start: 0, End: 1}}},
		},
		Scopes:        ScopeTable{{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule, Origin: testOrigin()}},
		Declarations:  DeclarationTable{{ID: DeclID(1), Scope: ScopeID(1), Name: "Int", Kind: DeclType, Exported: true, Origin: testOrigin()}},
		ImportAliases: ImportAliasTable{{Scope: ScopeID(0x80000001), LocalName: "Int", Module: ModuleID(2), Declaration: DeclID(1), Source: ImportAliasSourceProlog, Origin: testOrigin()}},
	}
	if program.Valid() {
		t.Fatal("prolog alias must be invalid when the prolog is disabled")
	}
}
