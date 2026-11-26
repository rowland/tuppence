package ir

import "testing"

func TestIndexedProgramRejectsPrivateCrossModuleTypeAliases(t *testing.T) {
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
			{ID: DeclID(1), Scope: ScopeID(1), Name: "Alias", Kind: DeclType, Origin: testOrigin()},
			{ID: DeclID(2), Scope: ScopeID(2), Name: "Private", Kind: DeclType, Origin: Origin{Primary: Span{File: FileID(2), Start: 0, End: 1}}},
		},
		Types:   TypeTable{{ID: TypeID(1), Name: "Private", Kind: TypeNominal, Declaration: DeclID(2), HomeModule: ModuleID(2), Origin: Origin{Primary: Span{File: FileID(2), Start: 0, End: 1}}}},
		Aliases: TypeAliasTable{{Declaration: DeclID(1), Target: TypeID(1), Origin: testOrigin()}},
	}
	if program.Valid() {
		t.Fatal("cross-module alias to a private type must be invalid")
	}
}

func TestIndexedProgramAllowsImportedExportedTypeAliases(t *testing.T) {
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
			{ID: DeclID(1), Scope: ScopeID(1), Name: "Alias", Kind: DeclType, Origin: testOrigin()},
			{ID: DeclID(2), Scope: ScopeID(2), Name: "Public", Kind: DeclType, Exported: true, Origin: Origin{Primary: Span{File: FileID(2), Start: 0, End: 1}}},
		},
		Types:   TypeTable{{ID: TypeID(1), Name: "Public", Kind: TypeNominal, Declaration: DeclID(2), HomeModule: ModuleID(2), Origin: Origin{Primary: Span{File: FileID(2), Start: 0, End: 1}}}},
		Aliases: TypeAliasTable{{Declaration: DeclID(1), Target: TypeID(1), Origin: testOrigin()}},
	}
	if !program.Valid() {
		t.Fatal("imported exported type alias must be valid")
	}
}
