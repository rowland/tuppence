package ir

import "testing"

func TestTypeConstructorTableRejectsDuplicatePairs(t *testing.T) {
	table := TypeConstructorTable{
		{Type: TypeID(1), Storage: StorageI32, Origin: testOrigin()},
		{Type: TypeID(1), Storage: StorageI32, Origin: testOrigin()},
	}
	if table.Valid() {
		t.Fatal("duplicate type/storage constructors must be invalid")
	}
}

func TestTypeConstructorTableFindsConstructor(t *testing.T) {
	table := TypeConstructorTable{{Type: TypeID(1), Storage: StorageI32, Origin: testOrigin()}}
	constructor, ok := table.Find(TypeID(1), StorageI32)
	if !ok || constructor.Type != TypeID(1) {
		t.Fatal("type constructor table must find a type/storage pair")
	}
}
