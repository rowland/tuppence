package ir

import "testing"

func TestModuleTableResolveImportUsesRelativeThenProgramThenSearchPaths(t *testing.T) {
	table := ModuleTable{
		{ID: ModuleID(1), Name: "app/sub", Files: []FileID{FileID(1)}, Origin: testOrigin()},
		{ID: ModuleID(2), Name: "app/foo", Files: []FileID{FileID(2)}, Origin: Origin{Primary: Span{File: FileID(2), Start: 0, End: 1}}},
		{ID: ModuleID(3), Name: "/program/foo", Files: []FileID{FileID(3)}, Origin: Origin{Primary: Span{File: FileID(3), Start: 0, End: 1}}},
		{ID: ModuleID(4), Name: "/stdlib/foo", Files: []FileID{FileID(4)}, Origin: Origin{Primary: Span{File: FileID(4), Start: 0, End: 1}}},
	}
	config := ModuleResolutionConfig{ProgramRoot: "/program", SearchPaths: []string{"/stdlib"}}
	if got, ok := table.ResolveImport(ModuleID(1), "foo", config); !ok || got != ModuleID(2) {
		t.Fatalf("relative module must win, got %d, %t", got, ok)
	}
	if got, ok := table.ResolveImport(ModuleID(1), "missing", config); ok || got != 0 {
		t.Fatalf("missing module must not resolve, got %d, %t", got, ok)
	}
	candidates := table.ResolveImportCandidates(ModuleID(1), "missing", config)
	want := []string{"app/missing", "/program/missing", "/stdlib/missing"}
	if len(candidates) != len(want) {
		t.Fatalf("unexpected candidate count: got %v want %v", candidates, want)
	}
	for i := range want {
		if candidates[i] != want[i] {
			t.Fatalf("unexpected candidate order: got %v want %v", candidates, want)
		}
	}
}

func TestModuleTableResolveImportAbsoluteUsesSearchRoots(t *testing.T) {
	table := ModuleTable{
		{ID: ModuleID(1), Name: "app", Files: []FileID{FileID(1)}, Origin: testOrigin()},
		{ID: ModuleID(2), Name: "/stdlib/foo", Files: []FileID{FileID(2)}, Origin: Origin{Primary: Span{File: FileID(2), Start: 0, End: 1}}},
	}
	config := ModuleResolutionConfig{ProgramRoot: "/program", SearchPaths: []string{"/stdlib"}}
	if got, ok := table.ResolveImport(ModuleID(1), "/foo", config); !ok || got != ModuleID(2) {
		t.Fatalf("absolute module must resolve through search roots, got %d, %t", got, ok)
	}
}
