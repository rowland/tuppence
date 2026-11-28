package ir

import "testing"

func TestIndexedProgramRejectsFileClaimedByMultipleModules(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{
			{ID: ModuleID(1), Name: "first", Files: []FileID{FileID(1)}, Origin: testOrigin()},
			{ID: ModuleID(2), Name: "second", Files: []FileID{FileID(1)}, Origin: testOrigin()},
		},
	}
	if program.Valid() {
		t.Fatal("a source file claimed by multiple modules must be invalid")
	}
}
