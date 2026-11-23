package ir

import "testing"

func TestBindingTableFindsBindingsByID(t *testing.T) {
	table := BindingTable{
		{ID: BindingID(1), Scope: ScopeID(2), Name: "counter", TypeSlot: TypeVarID(3), Value: ExprID(4), Origin: testOrigin()},
		{ID: BindingID(2), Scope: ScopeID(2), Name: "counter", Previous: BindingID(1), TypeSlot: TypeVarID(3), Value: ExprID(5), Mutable: true, Origin: testOrigin()},
	}
	if !table.Valid() {
		t.Fatal("complete binding table must be valid")
	}
	binding, ok := table.Find(BindingID(2))
	if !ok || binding.Previous != BindingID(1) {
		t.Fatal("binding table must find the requested binding version")
	}
	if _, ok := table.Find(BindingID(9)); ok {
		t.Fatal("binding table must not find an absent binding")
	}
}

func TestBindingTableRejectsInvalidOrDuplicateBindings(t *testing.T) {
	invalid := BindingTable{{ID: BindingID(1), Scope: ScopeID(2), Name: "counter", TypeSlot: TypeVarID(3), Value: ExprID(4)}}
	if invalid.Valid() {
		t.Fatal("binding table with an invalid binding must be invalid")
	}

	duplicate := BindingTable{
		{ID: BindingID(1), Scope: ScopeID(2), Name: "counter", TypeSlot: TypeVarID(3), Value: ExprID(4), Origin: testOrigin()},
		{ID: BindingID(1), Scope: ScopeID(2), Name: "other", TypeSlot: TypeVarID(3), Value: ExprID(5), Origin: testOrigin()},
	}
	if duplicate.Valid() {
		t.Fatal("binding table with duplicate IDs must be invalid")
	}
}
