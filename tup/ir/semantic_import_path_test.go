package ir

import "testing"

func TestIndexedProgramValidatesImportPathAgainstTarget(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{
			{ID: ModuleID(1), Name: "app/main", Files: []FileID{FileID(1)}, Imports: []ImportBinding{{LocalName: "foo", Path: "foo", Target: ModuleID(2), Origin: testOrigin()}}, Origin: testOrigin()},
			{ID: ModuleID(2), Name: "app/foo", Files: []FileID{FileID(2)}, Origin: Origin{Primary: Span{File: FileID(2), Start: 0, End: 1}}},
		},
	}
	if !program.Valid() {
		t.Fatal("import target matching its resolved path must be valid")
	}

	program.Modules[0].Imports[0].Target = ModuleID(1)
	if program.Valid() {
		t.Fatal("import target differing from its resolved path must be invalid")
	}
}
