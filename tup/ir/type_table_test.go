package ir

import "testing"

func TestTypeTableFindsTypesByID(t *testing.T) {
	table := TypeTable{
		{ID: TypeID(1), Name: "Int64", Kind: TypeNominal, Declaration: DeclID(1), HomeModule: ModuleID(1), Origin: testOrigin()},
		{ID: TypeID(2), Name: "User", Kind: TypeNominal, Declaration: DeclID(2), HomeModule: ModuleID(1), Origin: testOrigin()},
	}
	if !table.Valid() {
		t.Fatal("complete type table must be valid")
	}
	typ, ok := table.Find(TypeID(2))
	if !ok || typ.Name != "User" {
		t.Fatal("type table must find the requested type")
	}
	if _, ok := table.Find(TypeID(9)); ok {
		t.Fatal("type table must not find an absent type")
	}
}

func TestTypeTableRejectsDuplicateDeclarations(t *testing.T) {
	table := TypeTable{
		{ID: TypeID(1), Name: "first", Kind: TypeNominal, Declaration: DeclID(1), HomeModule: ModuleID(1), Origin: testOrigin()},
		{ID: TypeID(2), Name: "second", Kind: TypeNominal, Declaration: DeclID(1), HomeModule: ModuleID(1), Origin: testOrigin()},
	}
	if table.Valid() {
		t.Fatal("type table with duplicate declarations must be invalid")
	}
}
func TestTypeTableRejectsInvalidOrDuplicateTypes(t *testing.T) {
	invalid := TypeTable{{ID: TypeID(1), Name: "Int64", Kind: TypeNominal, Declaration: DeclID(1), HomeModule: ModuleID(1)}}
	if invalid.Valid() {
		t.Fatal("type table with an invalid type must be invalid")
	}

	duplicate := TypeTable{
		{ID: TypeID(1), Name: "Int64", Kind: TypeNominal, Declaration: DeclID(1), HomeModule: ModuleID(1), Origin: testOrigin()},
		{ID: TypeID(1), Name: "other", Kind: TypeNominal, Declaration: DeclID(2), HomeModule: ModuleID(1), Origin: testOrigin()},
	}
	if duplicate.Valid() {
		t.Fatal("type table with duplicate IDs must be invalid")
	}
}
