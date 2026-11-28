package ir

import "testing"

func TestBuildModulesFromSourceFilesGroupsFilesAndKeepsVersionsDistinct(t *testing.T) {
	files, invalid := IndexSourceFiles([]string{
		"bar/foo42/foo.tup",
		"bar/foo42/foo.extra.tup",
		"bar/foo43/foo.tup",
	}, "linux", "amd64")
	if len(invalid) != 0 {
		t.Fatalf("unexpected invalid files: %v", invalid)
	}
	modules := BuildModulesFromSourceFiles(files)
	if !modules.Valid() || len(modules) != 2 {
		t.Fatalf("unexpected modules: %#v", modules)
	}
	if modules[0].Name != "foo" || modules[0].Path != "bar/foo42/foo" {
		t.Fatalf("unexpected first module: %#v", modules[0])
	}
	if len(modules[0].Files) != 2 {
		t.Fatalf("files with the same module prefix must be grouped: %#v", modules[0].Files)
	}
}
