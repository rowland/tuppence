package ir

import "testing"

func TestTypeValidity(t *testing.T) {
	primitive := Type{ID: TypeID(1), Name: "int", Kind: TypePrimitive, Origin: testOrigin()}
	if !primitive.Valid() {
		t.Fatal("complete primitive type must be valid")
	}

	nominal := Type{ID: TypeID(2), Name: "User", Kind: TypeNominal, Declaration: DeclID(3), HomeModule: ModuleID(4), Origin: testOrigin()}
	if !nominal.Valid() {
		t.Fatal("complete nominal type must be valid")
	}
}

func TestTypeRejectsMissingFields(t *testing.T) {
	for name, typ := range map[string]Type{
		"missing id":     {Name: "int", Kind: TypePrimitive, Origin: testOrigin()},
		"missing name":   {ID: TypeID(1), Kind: TypePrimitive, Origin: testOrigin()},
		"missing kind":   {ID: TypeID(1), Name: "int", Origin: testOrigin()},
		"missing origin": {ID: TypeID(1), Name: "int", Kind: TypePrimitive},
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
