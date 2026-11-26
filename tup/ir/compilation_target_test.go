package ir

import "testing"

func TestCompilationTargetSharesSelectionConfiguration(t *testing.T) {
	target := CompilationTarget{OS: "linux", Architecture: "amd64", ModuleRoots: ModuleResolutionConfig{SearchPaths: []string{"/stdlib"}}}
	if !target.Valid() {
		t.Fatal("complete compilation target must be valid")
	}
	files, modules, invalid := BuildModuleIndexForTarget([]string{"app/main.tup", "app/main--linux.tup", "app/main--darwin.tup"}, target)
	if len(invalid) != 0 || !files.Valid() || !modules.Valid() || len(files) != 2 {
		t.Fatalf("unexpected target index: files=%#v modules=%#v invalid=%v", files, modules, invalid)
	}
}

func TestBuildModuleIndexForTargetRejectsIncompleteTarget(t *testing.T) {
	files, modules, invalid := BuildModuleIndexForTarget([]string{"app/main.tup"}, CompilationTarget{OS: "linux"})
	if files != nil || modules != nil || len(invalid) != 1 || invalid[0] != "app/main.tup" {
		t.Fatalf("incomplete target must reject loading: files=%#v modules=%#v invalid=%v", files, modules, invalid)
	}
}
