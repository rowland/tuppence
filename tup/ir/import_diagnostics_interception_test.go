package ir

import "testing"

func TestImportResolutionDiagnosticsPreservesInterceptedOrigin(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{
			{ID: ModuleID(1), Name: "main", Path: "app/main", Files: []FileID{FileID(1)}, Imports: []ImportBinding{{LocalName: "foo", Path: "foo", Target: ModuleID(3), Origin: testOrigin()}}, Origin: testOrigin()},
			{ID: ModuleID(2), Name: "local foo", Path: "app/foo", Files: []FileID{FileID(2)}, Origin: Origin{Primary: Span{File: FileID(2), Start: 0, End: 1}}},
			{ID: ModuleID(3), Name: "stdlib foo", Path: "/stdlib/foo", Files: []FileID{FileID(3)}, Origin: Origin{Primary: Span{File: FileID(3), Start: 0, End: 1}}},
		},
		ModuleResolution: ModuleResolutionConfig{SearchPaths: []string{"/stdlib"}},
	}
	diagnostics := program.ImportResolutionDiagnostics()
	if len(diagnostics) != 1 || diagnostics[0].Code != "IMPORT_TARGET_MISMATCH" {
		t.Fatalf("unexpected diagnostics: %#v", diagnostics)
	}
	diagnostic := diagnostics[0]
	if diagnostic.ImportResolvedPath != "app/foo" {
		t.Fatalf("diagnostic lost intercepted module path: %q", diagnostic.ImportResolvedPath)
	}
	if len(diagnostic.ImportCandidates) != 3 || diagnostic.ImportCandidates[0] != "app/foo" || diagnostic.ImportCandidates[2] != "/stdlib/foo" {
		t.Fatalf("diagnostic lost candidate order: %v", diagnostic.ImportCandidates)
	}
}
