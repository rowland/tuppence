package ir

import "testing"

func TestIndexedProgramRejectsSelfImports(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{{
			ID:      ModuleID(1),
			Name:    "main",
			Files:   []FileID{FileID(1)},
			Imports: []ImportBinding{{LocalName: "self", Path: "self", Target: ModuleID(1), Origin: testOrigin()}},
			Origin:  testOrigin(),
		}},
	}
	if program.Valid() {
		t.Fatal("module importing itself must be invalid")
	}
}
