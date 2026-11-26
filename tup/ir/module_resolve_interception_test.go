package ir

import "testing"

func TestModuleTableResolveImportDistinguishesInterceptedAndAbsoluteLibraryPaths(t *testing.T) {
	table := ModuleTable{
		{ID: ModuleID(1), Name: "main", Path: "app/main", Files: []FileID{FileID(1)}, Origin: testOrigin()},
		{ID: ModuleID(2), Name: "local foo", Path: "app/foo", Files: []FileID{FileID(2)}, Origin: Origin{Primary: Span{File: FileID(2), Start: 0, End: 1}}},
		{ID: ModuleID(3), Name: "stdlib foo", Path: "/stdlib/foo", Files: []FileID{FileID(3)}, Origin: Origin{Primary: Span{File: FileID(3), Start: 0, End: 1}}},
	}
	config := ModuleResolutionConfig{SearchPaths: []string{"/stdlib"}}
	local, ok := table.ResolveImport(ModuleID(1), "foo", config)
	if !ok || local != ModuleID(2) {
		t.Fatalf("relative import must resolve to intercepted local module: got %d, %t", local, ok)
	}
	stdlib, ok := table.ResolveImport(ModuleID(1), "/foo", config)
	if !ok || stdlib != ModuleID(3) {
		t.Fatalf("absolute import must resolve to direct library module: got %d, %t", stdlib, ok)
	}
	if local == stdlib {
		t.Fatal("intercepted and absolute imports must retain distinct module identities")
	}
}
