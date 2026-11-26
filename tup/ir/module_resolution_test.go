package ir

import "testing"

func TestModuleResolutionConfigValidity(t *testing.T) {
	if !(ModuleResolutionConfig{}).Valid() {
		t.Fatal("empty module resolution configuration must be valid")
	}
	if (ModuleResolutionConfig{ProgramRoot: "relative"}).Valid() {
		t.Fatal("program root must be absolute")
	}
	if !(ModuleResolutionConfig{ProgramRoot: "/program", SearchPaths: []string{"/stdlib", "/vendor"}}).Valid() {
		t.Fatal("absolute module roots must be valid")
	}
	if (ModuleResolutionConfig{SearchPaths: []string{"/stdlib", "/stdlib"}}).Valid() {
		t.Fatal("duplicate module search paths must be invalid")
	}
}
