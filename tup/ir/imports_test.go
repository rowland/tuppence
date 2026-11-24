package ir

import "testing"

func TestImportBindingValidity(t *testing.T) {
	binding := ImportBinding{LocalName: "io", Target: ModuleID(2), Origin: testOrigin()}
	if !binding.Valid() {
		t.Fatal("complete import binding must be valid")
	}
}

func TestImportBindingRejectsMissingFields(t *testing.T) {
	for name, binding := range map[string]ImportBinding{
		"missing local name": {Target: ModuleID(2), Origin: testOrigin()},
		"missing target":     {LocalName: "io", Origin: testOrigin()},
		"missing origin":     {LocalName: "io", Target: ModuleID(2)},
	} {
		if binding.Valid() {
			t.Errorf("%s import binding must be invalid", name)
		}
	}
}
