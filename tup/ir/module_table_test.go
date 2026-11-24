package ir

import "testing"

func TestModuleTableFindsModulesByID(t *testing.T) {
	table := ModuleTable{
		{ID: ModuleID(1), Name: "main", Origin: testOrigin()},
		{ID: ModuleID(2), Name: "lib", Origin: testOrigin()},
	}
	if !table.Valid() {
		t.Fatal("complete module table must be valid")
	}
	module, ok := table.Find(ModuleID(2))
	if !ok || module.Name != "lib" {
		t.Fatal("module table must find the requested module")
	}
	if _, ok := table.Find(ModuleID(9)); ok {
		t.Fatal("module table must not find an absent module")
	}
}

func TestModuleTableRejectsInvalidOrDuplicateModules(t *testing.T) {
	invalid := ModuleTable{{ID: ModuleID(1), Name: "main"}}
	if invalid.Valid() {
		t.Fatal("module table with an invalid module must be invalid")
	}

	duplicate := ModuleTable{
		{ID: ModuleID(1), Name: "main", Origin: testOrigin()},
		{ID: ModuleID(1), Name: "other", Origin: testOrigin()},
	}
	if duplicate.Valid() {
		t.Fatal("module table with duplicate IDs must be invalid")
	}
}
