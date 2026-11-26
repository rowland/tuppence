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

func TestModuleRejectsOriginOutsideFiles(t *testing.T) {
	module := Module{
		ID:     ModuleID(1),
		Name:   "main",
		Files:  []FileID{FileID(1)},
		Origin: Origin{Primary: Span{File: FileID(2), Start: 0, End: 1}},
	}
	if module.Valid() {
		t.Fatal("module origin outside its files must be invalid")
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

func TestModuleRejectsDuplicateImports(t *testing.T) {
	module := Module{
		ID:   ModuleID(1),
		Name: "main",
		Imports: []ImportBinding{
			{LocalName: "io", Path: "io", Target: ModuleID(2), Origin: testOrigin()},
			{LocalName: "io", Path: "other", Target: ModuleID(3), Origin: testOrigin()},
		},
		Origin: testOrigin(),
	}
	if module.Valid() {
		t.Fatal("module with duplicate import names must be invalid")
	}
}

func TestModuleRejectsMissingFields(t *testing.T) {
	for name, module := range map[string]Module{
		"missing id":     {Name: "main", Files: []FileID{FileID(1)}, Origin: testOrigin()},
		"missing name":   {ID: ModuleID(1), Origin: testOrigin()},
		"missing origin": {ID: ModuleID(1), Name: "main"},
	} {
		if module.Valid() {
			t.Errorf("%s module must be invalid", name)
		}
	}
}
