package ir

import "testing"

func TestTypeAliasTableFindsAliasesByDeclaration(t *testing.T) {
	table := TypeAliasTable{
		{Declaration: DeclID(1), Target: TypeID(2), Origin: testOrigin()},
		{Declaration: DeclID(3), Target: TypeID(4), Origin: testOrigin()},
	}
	if !table.Valid() {
		t.Fatal("complete type alias table must be valid")
	}
	alias, ok := table.Find(DeclID(3))
	if !ok || alias.Target != TypeID(4) {
		t.Fatal("alias table must find the requested alias")
	}
	if _, ok := table.Find(DeclID(9)); ok {
		t.Fatal("alias table must not find an absent alias")
	}
}

func TestTypeAliasTableRejectsDuplicateDeclarations(t *testing.T) {
	table := TypeAliasTable{
		{Declaration: DeclID(1), Target: TypeID(2), Origin: testOrigin()},
		{Declaration: DeclID(1), Target: TypeID(3), Origin: testOrigin()},
	}
	if table.Valid() {
		t.Fatal("alias table with duplicate declarations must be invalid")
	}
}
