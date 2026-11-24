package ir

import "testing"

func TestIndexedProgramRejectsDanglingModuleImports(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{
			{ID: ModuleID(1), Name: "main", Imports: []ImportBinding{{LocalName: "io", Target: ModuleID(9), Origin: testOrigin()}}, Origin: testOrigin()},
		},
	}
	if program.Valid() {
		t.Fatal("import referring to an absent module must be invalid")
	}
}
