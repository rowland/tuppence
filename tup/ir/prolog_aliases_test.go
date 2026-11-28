package ir

import "testing"

func TestIndexedProgramAddsPrologAliasOnlyToInitialScope(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{
			{ID: ModuleID(1), Name: "main", Files: []FileID{FileID(1)}, Origin: testOrigin()},
			{ID: ModuleID(2), Name: "core", Files: []FileID{FileID(2)}, Origin: Origin{Primary: Span{File: FileID(2), Start: 0, End: 1}}},
		},
		Scopes: ScopeTable{
			{ID: ScopeID(0x80000001), Module: ModuleID(1), Kind: ScopeProlog, Origin: testOrigin()},
			{ID: ScopeID(1), Parent: ScopeID(0x80000001), Module: ModuleID(1), Kind: ScopeModule, Initial: true, Origin: testOrigin()},
			{ID: ScopeID(2), Module: ModuleID(1), Kind: ScopeFunction, Origin: testOrigin()},
		},
		Prolog: PrologConfig{Enabled: true, Modules: []ModuleID{ModuleID(2)}, Architecture: "amd64"},
	}
	alias := ImportAlias{Scope: ScopeID(0x80000001), LocalName: "Int", Module: ModuleID(2), Declaration: DeclID(1), Origin: testOrigin()}
	if !program.AddPrologAlias(alias) || len(program.ImportAliases) != 1 || program.ImportAliases[0].Source != ImportAliasSourceProlog {
		t.Fatal("prolog alias must be installed in an initial scope")
	}
	alias.Scope = ScopeID(2)
	if program.AddPrologAlias(alias) {
		t.Fatal("prolog alias must not be installed in a non-initial scope")
	}
	if program.AddPrologAliases([]ImportAlias{{Scope: ScopeID(1), LocalName: "Float", Module: ModuleID(2), Declaration: DeclID(2), Origin: testOrigin()}, alias}) || len(program.ImportAliases) != 1 {
		t.Fatal("prolog alias batches must be atomic")
	}
}
