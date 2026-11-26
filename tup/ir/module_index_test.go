package ir

import "testing"

func TestBuildModuleIndexSharesOrdinaryAndPrologLoadingPath(t *testing.T) {
	files, modules, invalid := BuildModuleIndex([]string{
		"core/math.tup",
		"core/math--linux.tup",
		"core/math--darwin.tup",
		"app/main.tup",
		"bad.txt",
	}, "linux", "amd64")
	if len(invalid) != 1 || invalid[0] != "bad.txt" {
		t.Fatalf("unexpected invalid files: %v", invalid)
	}
	if !files.Valid() || !modules.Valid() {
		t.Fatalf("shared index must produce valid tables: files=%#v modules=%#v", files, modules)
	}
	if len(files) != 3 || len(modules) != 2 {
		t.Fatalf("unexpected target-selected index: files=%d modules=%d", len(files), len(modules))
	}
	if modules[0].Path != "app/main" || modules[1].Path != "core/math" {
		t.Fatalf("unexpected module paths: %#v", modules)
	}
}
