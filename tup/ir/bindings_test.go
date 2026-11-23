package ir

import "testing"

func TestBindingValidity(t *testing.T) {
	binding := Binding{
		ID:       BindingID(1),
		Scope:    ScopeID(2),
		Name:     "counter",
		TypeSlot: TypeVarID(3),
		Value:    ExprID(4),
		Mutable:  true,
		Origin:   Origin{Primary: Span{File: FileID(1), Start: 0, End: 7}},
	}
	if !binding.Valid() {
		t.Fatal("complete first binding must be valid")
	}

	binding.Previous = BindingID(5)
	if !binding.Valid() {
		t.Fatal("binding version with a previous binding must be valid")
	}
}

func TestBindingRejectsMissingFields(t *testing.T) {
	origin := Origin{Primary: Span{File: FileID(1), Start: 0, End: 1}}
	for name, binding := range map[string]Binding{
		"missing id":        {Scope: ScopeID(2), Name: "x", TypeSlot: TypeVarID(3), Value: ExprID(4), Origin: origin},
		"missing scope":     {ID: BindingID(1), Name: "x", TypeSlot: TypeVarID(3), Value: ExprID(4), Origin: origin},
		"missing name":      {ID: BindingID(1), Scope: ScopeID(2), TypeSlot: TypeVarID(3), Value: ExprID(4), Origin: origin},
		"missing type slot": {ID: BindingID(1), Scope: ScopeID(2), Name: "x", Value: ExprID(4), Origin: origin},
		"missing value":     {ID: BindingID(1), Scope: ScopeID(2), Name: "x", TypeSlot: TypeVarID(3), Origin: origin},
		"missing origin":    {ID: BindingID(1), Scope: ScopeID(2), Name: "x", TypeSlot: TypeVarID(3), Value: ExprID(4)},
	} {
		if binding.Valid() {
			t.Errorf("%s binding must be invalid", name)
		}
	}
}
