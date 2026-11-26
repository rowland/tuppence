package ir

import "testing"

func TestImportBindingValidity(t *testing.T) {
	binding := ImportBinding{LocalName: "io", Path: "std/io", Target: ModuleID(2), Origin: testOrigin()}
	if !binding.Valid() {
		t.Fatal("complete import binding must be valid")
	}
}

func TestImportBindingRejectsMissingFields(t *testing.T) {
	for name, binding := range map[string]ImportBinding{
		"missing local name": {Path: "io", Target: ModuleID(2), Origin: testOrigin()},
		"missing path":       {LocalName: "io", Target: ModuleID(2), Origin: testOrigin()},
		"missing target":     {LocalName: "io", Path: "io", Origin: testOrigin()},
		"missing origin":     {LocalName: "io", Target: ModuleID(2)},
	} {
		if binding.Valid() {
			t.Errorf("%s import binding must be invalid", name)
		}
	}
}
