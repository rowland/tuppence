package ir

import "testing"

func TestIndexedProgramRejectsImportFromForeignModuleFile(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{
			{ID: ModuleID(1), Name: "main", Files: []FileID{FileID(1)}, Imports: []ImportBinding{{LocalName: "lib", Path: "lib", Target: ModuleID(2), Origin: Origin{Primary: Span{File: FileID(2), Start: 0, End: 1}}}}, Origin: testOrigin()},
			{ID: ModuleID(2), Name: "lib", Files: []FileID{FileID(2)}, Origin: Origin{Primary: Span{File: FileID(2), Start: 0, End: 1}}},
		},
	}
	if program.Valid() {
		t.Fatal("import from a foreign module file must be invalid")
	}
}
