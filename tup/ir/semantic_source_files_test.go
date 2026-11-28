package ir

import "testing"

func TestIndexedProgramValidatesModuleSourceFileMembership(t *testing.T) {
	program := IndexedProgram{
		Modules:     ModuleTable{{ID: ModuleID(1), Name: "math", Files: []FileID{FileID(1)}, Origin: testOrigin()}},
		SourceFiles: SourceFileTable{{ID: FileID(1), Path: "math.tup", Module: ModuleFileName{ModuleName: "math"}}},
	}
	if !program.Valid() {
		t.Fatal("module owning a matching indexed source file must be valid")
	}
	program.SourceFiles[0].Module.ModuleName = "other"
	if program.Valid() {
		t.Fatal("module owning a source file from another module must be invalid")
	}
}
