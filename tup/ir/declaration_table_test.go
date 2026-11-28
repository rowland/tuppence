package ir

import "testing"

func TestDeclarationTableFindsDeclarationsByID(t *testing.T) {
	table := DeclarationTable{
		{ID: DeclID(1), Scope: ScopeID(2), Name: "answer", Kind: DeclValue, Origin: testOrigin()},
		{ID: DeclID(2), Scope: ScopeID(2), Name: "compute", Kind: DeclFunction, Origin: testOrigin()},
	}
	if !table.Valid() {
		t.Fatal("complete declaration table must be valid")
	}
	declaration, ok := table.Find(DeclID(2))
	if !ok || declaration.Name != "compute" {
		t.Fatal("declaration table must find the requested declaration")
	}
	if _, ok := table.Find(DeclID(9)); ok {
		t.Fatal("declaration table must not find an absent declaration")
	}
}

func TestDeclarationTableRejectsInvalidOrDuplicateDeclarations(t *testing.T) {
	invalid := DeclarationTable{{ID: DeclID(1), Scope: ScopeID(2), Name: "answer", Kind: DeclValue}}
	if invalid.Valid() {
		t.Fatal("declaration table with an invalid declaration must be invalid")
	}

	duplicate := DeclarationTable{
		{ID: DeclID(1), Scope: ScopeID(2), Name: "answer", Kind: DeclValue, Origin: testOrigin()},
		{ID: DeclID(1), Scope: ScopeID(2), Name: "other", Kind: DeclValue, Origin: testOrigin()},
	}
	if duplicate.Valid() {
		t.Fatal("declaration table with duplicate IDs must be invalid")
	}
}

func testOrigin() Origin {
	return Origin{Primary: Span{File: FileID(1), Start: 0, End: 1}}
}
