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

func TestArrayTypeRejectsDynamicLength(t *testing.T) {
	array := Type{ID: TypeID(1), Kind: TypeArray, Element: TypeID(2), Length: 4, Origin: testOrigin()}
	if array.Valid() {
		t.Fatal("dynamic array with a fixed length must be invalid")
	}
}

func TestTupleTypeValidity(t *testing.T) {
	tuple := Type{ID: TypeID(1), Kind: TypeTuple, Elements: []TypeID{TypeID(2), TypeID(3)}, Origin: testOrigin()}
	if !tuple.Valid() {
		t.Fatal("complete tuple type must be valid")
	}
}

func TestTupleTypeAllowsEmptyTuple(t *testing.T) {
	tuple := Type{ID: TypeID(1), Kind: TypeTuple, Elements: []TypeID{}, Origin: testOrigin()}
	if !tuple.Valid() {
		t.Fatal("empty tuple type must be valid")
	}
}

func TestUnionTypeValidity(t *testing.T) {
	union := Type{ID: TypeID(1), Kind: TypeUnion, Elements: []TypeID{TypeID(2), TypeID(3)}, Origin: testOrigin()}
	if !union.Valid() {
		t.Fatal("complete union type must be valid")
	}
}

func TestUnionTypeRejectsEmptyAlternatives(t *testing.T) {
	union := Type{ID: TypeID(1), Kind: TypeUnion, Elements: []TypeID{}, Origin: testOrigin()}
	if union.Valid() {
		t.Fatal("union without alternatives must be invalid")
	}
}

func TestArrayTypeRequiresElement(t *testing.T) {
	array := Type{ID: TypeID(1), Kind: TypeArray, Origin: testOrigin()}
	if array.Valid() {
		t.Fatal("array without an element type must be invalid")
	}
}

func TestStructuralTypesRejectNominalMetadata(t *testing.T) {
	cases := []Type{
		{ID: TypeID(1), Name: "Array", Kind: TypeArray, Element: TypeID(2), Origin: testOrigin()},
		{ID: TypeID(1), Name: "Tuple", Kind: TypeTuple, Elements: []TypeID{TypeID(2)}, Origin: testOrigin()},
		{ID: TypeID(1), Name: "Union", Kind: TypeUnion, Elements: []TypeID{TypeID(2)}, Origin: testOrigin()},
	}
	for _, typ := range cases {
		if typ.Valid() {
			t.Errorf("structural type with a nominal name must be invalid: %#v", typ)
		}
	}
}

func TestStructuralTypesRejectIrrelevantMetadata(t *testing.T) {
	tuple := Type{ID: TypeID(1), Kind: TypeTuple, Elements: []TypeID{TypeID(2)}, FixedLength: true, Origin: testOrigin()}
	if tuple.Valid() {
		t.Fatal("tuple with array length metadata must be invalid")
	}
	union := Type{ID: TypeID(1), Kind: TypeUnion, Elements: []TypeID{TypeID(2)}, Length: 1, Origin: testOrigin()}
	if union.Valid() {
		t.Fatal("union with array length metadata must be invalid")
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
		StorageF16, StorageF32, StorageF64, StorageV128, StoragePTR, StorageBytes,
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

func TestLiteralStoragePriority(t *testing.T) {
	got := LiteralStoragePriority(LiteralInteger)
	want := []StorageKind{
		StorageI8, StorageI16, StorageI32, StorageI64,
		StorageU8, StorageU16, StorageU32, StorageU64,
		StorageF16, StorageF32, StorageF64, StorageBytes,
	}
	if len(got) != len(want) {
		t.Fatalf("integer literal priority length = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("integer literal priority[%d] = %d, want %d", i, got[i], want[i])
		}
	}
	if got, ok := SelectLiteralStorage(LiteralInteger, []StorageKind{StorageF32, StorageI32}); !ok || got != StorageI32 {
		t.Fatalf("literal storage selection = %d, %t; want I32, true", got, ok)
	}
}

func TestAmbiguousLiteralDefaults(t *testing.T) {
	if got, ok := DefaultLiteralKind(LiteralInteger, false); !ok || got != LiteralDefaultInt {
		t.Fatal("ambiguous in-range integer literal must default to Int")
	}
	if got, ok := DefaultLiteralKind(LiteralInteger, true); !ok || got != LiteralDefaultFloat {
		t.Fatal("out-of-range integer literal must default to Float")
	}
	if got, ok := DefaultLiteralKind(LiteralFloat, false); !ok || got != LiteralDefaultFloat {
		t.Fatal("float literal must default to Float")
	}
}
