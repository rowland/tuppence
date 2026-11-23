package ir

import "testing"

func TestScopeValidity(t *testing.T) {
	scope := Scope{
		ID:     ScopeID(1),
		Module: ModuleID(2),
		Kind:   ScopeModule,
	}
	if !scope.Valid() {
		t.Fatal("complete scope must be valid")
	}
	if scope.Parent.Valid() {
		t.Fatal("root scope should not require a parent")
	}
}

func TestScopeRejectsMissingIdentityOrKind(t *testing.T) {
	for name, invalid := range map[string]Scope{
		"missing id":     {Module: ModuleID(2), Kind: ScopeBlock},
		"missing module": {ID: ScopeID(1), Kind: ScopeBlock},
		"missing kind":   {ID: ScopeID(1), Module: ModuleID(2)},
	} {
		if invalid.Valid() {
			t.Errorf("%s scope must be invalid", name)
		}
	}
}

func TestScopeKindsHaveValidValues(t *testing.T) {
	for _, kind := range []ScopeKind{ScopeModule, ScopeFunction, ScopeBlock, ScopeLoop} {
		if !kind.Valid() {
			t.Errorf("scope kind %d must be valid", kind)
		}
	}
}
