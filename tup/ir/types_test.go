package ir

import "testing"

func TestTypeValidity(t *testing.T) {
	nominal := Type{ID: TypeID(1), Name: "User", Kind: TypeNominal, Declaration: DeclID(3), HomeModule: ModuleID(4), Origin: testOrigin()}
	if !nominal.Valid() {
		t.Fatal("complete nominal type must be valid")
	}
}

func TestArrayTypeValidity(t *testing.T) {
	array := Type{ID: TypeID(1), Kind: TypeArray, Element: TypeID(2), FixedLength: true, Length: 4, Origin: testOrigin()}
	if !array.Valid() {
		t.Fatal("complete array type must be valid")
	}
}

func TestArrayTypeRequiresElement(t *testing.T) {
	array := Type{ID: TypeID(1), Kind: TypeArray, Origin: testOrigin()}
	if array.Valid() {
		t.Fatal("array without an element type must be invalid")
	}
}
func TestTypeRejectsMissingFields(t *testing.T) {
	for name, typ := range map[string]Type{
		"missing id":     {Name: "User", Kind: TypeNominal, Declaration: DeclID(1), HomeModule: ModuleID(1), Origin: testOrigin()},
		"missing name":   {ID: TypeID(1), Kind: TypeNominal, Declaration: DeclID(1), HomeModule: ModuleID(1), Origin: testOrigin()},
		"missing kind":   {ID: TypeID(1), Name: "User", Declaration: DeclID(1), HomeModule: ModuleID(1), Origin: testOrigin()},
		"missing origin": {ID: TypeID(1), Name: "User", Kind: TypeNominal, Declaration: DeclID(1), HomeModule: ModuleID(1)},
	} {
		if typ.Valid() {
			t.Errorf("%s type must be invalid", name)
		}
	}
}

func TestNominalTypeRequiresDeclarationAndHomeModule(t *testing.T) {
	for name, typ := range map[string]Type{
		"missing declaration": {ID: TypeID(1), Name: "User", Kind: TypeNominal, HomeModule: ModuleID(1), Origin: testOrigin()},
		"missing home module": {ID: TypeID(1), Name: "User", Kind: TypeNominal, Declaration: DeclID(1), Origin: testOrigin()},
	} {
		if typ.Valid() {
			t.Errorf("%s nominal type must be invalid", name)
		}
	}
}

func TestStorageKinds(t *testing.T) {
	for _, kind := range []StorageKind{
		StorageI8, StorageI16, StorageI32, StorageI64,
		StorageU8, StorageU16, StorageU32, StorageU64,
		StorageF16, StorageF32, StorageF64, StorageV128, StoragePTR,
	} {
		if !kind.Valid() {
			t.Errorf("storage kind %d must be valid", kind)
		}
	}
	if !(StorageType{Kind: StoragePTR}).Valid() {
		t.Fatal("pointer storage type must be valid")
	}
	if (StorageType{}).Valid() {
		t.Fatal("empty storage type must be invalid")
	}
}
