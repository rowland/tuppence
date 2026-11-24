package ir

import "testing"

func TestTupleSchemaValidity(t *testing.T) {
	schema := TupleSchema{
		ID: TupleSchemaID(1),
		Members: []TupleMember{
			{Name: "left", LogicalSlot: 0, Type: TypeID(1)},
			{Name: "right", LogicalSlot: 1, Type: TypeID(2)},
		},
		Origin: testOrigin(),
	}
	if !schema.Valid() {
		t.Fatal("complete tuple schema must be valid")
	}
}

func TestTupleSchemaAllowsOrdinalMembers(t *testing.T) {
	schema := TupleSchema{
		ID:      TupleSchemaID(1),
		Members: []TupleMember{{LogicalSlot: 0, Type: TypeID(1)}},
		Origin:  testOrigin(),
	}
	if !schema.Valid() {
		t.Fatal("tuple schema with ordinal members must be valid")
	}
}

func TestTupleSchemaRejectsInvalidMembers(t *testing.T) {
	for name, schema := range map[string]TupleSchema{
		"missing id":     {Members: []TupleMember{{Type: TypeID(1)}}, Origin: testOrigin()},
		"missing type":   {ID: TupleSchemaID(1), Members: []TupleMember{{}}, Origin: testOrigin()},
		"missing origin": {ID: TupleSchemaID(1), Members: []TupleMember{{Type: TypeID(1)}}},
		"duplicate slot": {ID: TupleSchemaID(1), Members: []TupleMember{{LogicalSlot: 0, Type: TypeID(1)}, {LogicalSlot: 0, Type: TypeID(2)}}, Origin: testOrigin()},
	} {
		if schema.Valid() {
			t.Errorf("%s tuple schema must be invalid", name)
		}
	}
}
