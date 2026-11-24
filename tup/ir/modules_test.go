package ir

import "testing"

func TestModuleValidity(t *testing.T) {
	module := Module{
		ID:     ModuleID(1),
		Name:   "main",
		Files:  []FileID{FileID(1), FileID(2)},
		Origin: testOrigin(),
	}
	if !module.Valid() {
		t.Fatal("complete module must be valid")
	}
}

func TestModuleRejectsDuplicateOrInvalidFiles(t *testing.T) {
	for name, module := range map[string]Module{
		"invalid file":   {ID: ModuleID(1), Name: "main", Files: []FileID{FileID(0)}, Origin: testOrigin()},
		"duplicate file": {ID: ModuleID(1), Name: "main", Files: []FileID{FileID(1), FileID(1)}, Origin: testOrigin()},
	} {
		if module.Valid() {
			t.Errorf("%s module must be invalid", name)
		}
	}
}

func TestModuleRejectsMissingFields(t *testing.T) {
	for name, module := range map[string]Module{
		"missing id":     {Name: "main", Origin: testOrigin()},
		"missing name":   {ID: ModuleID(1), Origin: testOrigin()},
		"missing origin": {ID: ModuleID(1), Name: "main"},
	} {
		if module.Valid() {
			t.Errorf("%s module must be invalid", name)
		}
	}
}
