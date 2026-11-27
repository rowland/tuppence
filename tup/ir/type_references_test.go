package ir

import "testing"

func TestIndexedProgramResolvesTypeReferencesInOrder(t *testing.T) {
	prolog := ScopeID(0x80000001)
	program := IndexedProgram{
		Scopes: ScopeTable{
			{ID: prolog, Module: ModuleID(1), Kind: ScopeProlog, Origin: testOrigin()},
			{ID: ScopeID(1), Parent: prolog, Module: ModuleID(1), Kind: ScopeModule, Initial: true, Origin: testOrigin()},
			{ID: ScopeID(2), Module: ModuleID(2), Kind: ScopeModule, Origin: testOrigin()},
		},
		ImportAliases: ImportAliasTable{{Scope: prolog, LocalName: "Int", Module: ModuleID(2), Declaration: DeclID(1), Origin: testOrigin()}},
		Declarations: DeclarationTable{
			{ID: DeclID(1), Scope: ScopeID(2), Name: "Int", Kind: DeclType, Exported: true, Origin: testOrigin()},
			{ID: DeclID(2), Scope: ScopeID(1), Name: "Ambig", Kind: DeclType, Origin: testOrigin()},
			{ID: DeclID(3), Scope: ScopeID(1), Name: "Ambig", Kind: DeclType, Origin: testOrigin()},
		},
		Types: TypeTable{
			{ID: TypeID(10), Name: "Int", Kind: TypeNominal, Declaration: DeclID(1), HomeModule: ModuleID(2), Origin: testOrigin()},
			{ID: TypeID(11), Name: "Ambig", Kind: TypeNominal, Declaration: DeclID(2), HomeModule: ModuleID(1), Origin: testOrigin()},
			{ID: TypeID(12), Name: "Ambig", Kind: TypeNominal, Declaration: DeclID(3), HomeModule: ModuleID(1), Origin: testOrigin()},
		},
	}
	references := []TypeReference{
		{ID: TypeReferenceID(1), Scope: ScopeID(1), Name: "Int", Origin: testOrigin()},
		{ID: TypeReferenceID(2), Scope: ScopeID(1), Name: "Missing", Origin: testOrigin()},
		{ID: TypeReferenceID(3), Scope: ScopeID(1), Name: "Ambig", Origin: testOrigin()},
	}
	resolved, diagnostics := program.ResolveTypeReferences(references)
	if len(resolved) != 1 || resolved[0].Reference.ID != TypeReferenceID(1) || resolved[0].Type != TypeID(10) {
		t.Fatal("type reference resolution must preserve the successful reference")
	}
	if len(diagnostics) != 2 || diagnostics[0].Code != "TYPE_NOT_FOUND" || diagnostics[1].Code != "TYPE_AMBIGUOUS" {
		t.Fatal("type reference diagnostics must preserve source order")
	}
}
