package ir

import "testing"

func TestScopeTableFindsScopesByID(t *testing.T) {
	table := ScopeTable{
		{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule, Origin: testOrigin()},
		{ID: ScopeID(2), Parent: ScopeID(1), Module: ModuleID(1), Kind: ScopeBlock, Origin: testOrigin()},
	}
	if !table.Valid() {
		t.Fatal("complete scope table must be valid")
	}
	scope, ok := table.Find(ScopeID(2))
	if !ok || scope.Parent != ScopeID(1) {
		t.Fatal("scope table must find the requested scope")
	}
	if _, ok := table.Find(ScopeID(9)); ok {
		t.Fatal("scope table must not find an absent scope")
	}
}

func TestScopeTableRejectsInvalidOrDuplicateScopes(t *testing.T) {
	invalid := ScopeTable{{ID: ScopeID(1), Module: ModuleID(1)}}
	if invalid.Valid() {
		t.Fatal("scope table with an invalid scope must be invalid")
	}

	duplicate := ScopeTable{
		{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule, Origin: testOrigin()},
		{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeBlock, Origin: testOrigin()},
	}
	if duplicate.Valid() {
		t.Fatal("scope table with duplicate IDs must be invalid")
	}
}
