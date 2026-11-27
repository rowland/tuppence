package ir

import "testing"

func TestIndexedProgramTypeLookupFallsThroughToPrologAndAllowsModuleShadowing(t *testing.T) {
	prolog := ScopeID(0x80000001)
	program := IndexedProgram{
		Scopes: ScopeTable{
			{ID: prolog, Module: ModuleID(1), Kind: ScopeProlog, Origin: testOrigin()},
			{ID: ScopeID(1), Parent: prolog, Module: ModuleID(1), Kind: ScopeModule, Initial: true, Origin: testOrigin()},
			{ID: ScopeID(2), Module: ModuleID(2), Kind: ScopeModule, Origin: testOrigin()},
		},
		ImportAliases: ImportAliasTable{{Scope: prolog, LocalName: "Int", Module: ModuleID(2), Declaration: DeclID(1), Origin: testOrigin()}},
		Declarations:  DeclarationTable{{ID: DeclID(1), Scope: ScopeID(2), Name: "Int", Kind: DeclType, Exported: true, Origin: testOrigin()}},
		Types:         TypeTable{{ID: TypeID(10), Name: "Int", Kind: TypeNominal, Declaration: DeclID(1), HomeModule: ModuleID(2), Origin: testOrigin()}},
	}
	if matches := program.FindTypeCandidates(ScopeID(1), "Int"); len(matches) != 1 || matches[0] != TypeID(10) {
		t.Fatal("type lookup must fall through from the module scope to the prolog scope")
	}

	program.Declarations = append(program.Declarations, Declaration{ID: DeclID(2), Scope: ScopeID(1), Name: "Int", Kind: DeclType, Origin: testOrigin()})
	program.Types = append(program.Types, Type{ID: TypeID(11), Name: "Int", Kind: TypeNominal, Declaration: DeclID(2), HomeModule: ModuleID(1), Origin: testOrigin()})
	if matches := program.FindTypeCandidates(ScopeID(1), "Int"); len(matches) != 1 || matches[0] != TypeID(11) {
		t.Fatal("module types must shadow prolog types")
	}
}
