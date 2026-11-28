package ir

import "testing"

func TestImportResolutionDiagnosticsReportsUnresolvedPath(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{{ID: ModuleID(1), Name: "app", Files: []FileID{FileID(1)}, Imports: []ImportBinding{{LocalName: "foo", Path: "foo", Target: ModuleID(9), Origin: testOrigin()}}, Origin: testOrigin()}},
	}
	diagnostics := program.ImportResolutionDiagnostics()
	if len(diagnostics) != 1 || diagnostics[0].Code != "IMPORT_NOT_FOUND" {
		t.Fatalf("unexpected import diagnostics: %#v", diagnostics)
	}
	if diagnostics[0].ImportPath != "foo" {
		t.Fatal("unresolved diagnostic must preserve the requested path")
	}
}

func TestNewUnresolvedImportDiagnosticPreservesResolution(t *testing.T) {
	candidates := []string{"app/foo", "/program/foo", "/stdlib/foo"}
	diagnostic := NewUnresolvedImportDiagnostic(testOrigin(), "foo", candidates)
	if !diagnostic.Valid() {
		t.Fatal("unresolved import diagnostic must be valid")
	}
	if diagnostic.Code != "IMPORT_NOT_FOUND" || diagnostic.ImportPath != "foo" {
		t.Fatalf("unexpected diagnostic identity: %#v", diagnostic)
	}
	if got := FormatImportResolution(diagnostic); got != "import foo (tried: app/foo, /program/foo, /stdlib/foo)" {
		t.Fatalf("unexpected resolution text: %q", got)
	}
	candidates[0] = "changed"
	if diagnostic.ImportCandidates[0] != "app/foo" {
		t.Fatal("diagnostic must own its candidate list")
	}
}
