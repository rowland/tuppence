package ir

import "testing"

func TestImportAliasTableRejectsDuplicateLocalNames(t *testing.T) {
	table := ImportAliasTable{
		{Scope: ScopeID(1), LocalName: "Foo", Module: ModuleID(2), Declaration: DeclID(3), Origin: testOrigin()},
		{Scope: ScopeID(1), LocalName: "Foo", Module: ModuleID(2), Declaration: DeclID(4), Origin: testOrigin()},
	}
	if table.Valid() {
		t.Fatal("duplicate imported aliases in one scope must be invalid")
	}
}

func TestImportAliasTableAllowsSameNameInDifferentScopes(t *testing.T) {
	table := ImportAliasTable{
		{Scope: ScopeID(1), LocalName: "Foo", Module: ModuleID(2), Declaration: DeclID(3), Origin: testOrigin()},
		{Scope: ScopeID(2), LocalName: "Foo", Module: ModuleID(2), Declaration: DeclID(3), Origin: testOrigin()},
	}
	if !table.Valid() {
		t.Fatal("same imported alias may be used in different scopes")
	}
}
