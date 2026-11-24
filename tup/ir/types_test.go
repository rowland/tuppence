package ir

import "testing"

func TestTypeValidity(t *testing.T) {
	for _, kind := range []TypeKind{TypePrimitive, TypeNominal} {
		typ := Type{ID: TypeID(1), Name: "int", Kind: kind, Origin: testOrigin()}
		if !typ.Valid() {
			t.Errorf("complete type of kind %d must be valid", kind)
		}
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
