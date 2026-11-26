package ir

import "testing"

func TestBuildIndexedProgramUsesOneTargetConfiguration(t *testing.T) {
	target := CompilationTarget{OS: "linux", Architecture: "amd64", ModuleRoots: ModuleResolutionConfig{ProgramRoot: "/program", SearchPaths: []string{"/stdlib"}}}
	program, invalid := BuildIndexedProgram([]string{"app/main.tup", "app/main--linux.tup", "app/main--darwin.tup"}, target, PrologConfig{})
	if len(invalid) != 0 || !program.Valid() {
		t.Fatalf("constructed program must be valid: %#v invalid=%v", program, invalid)
	}
	if program.ModuleResolution.ProgramRoot != "/program" || len(program.SourceFiles) != 2 || len(program.Modules) != 1 {
		t.Fatalf("target configuration was not preserved: %#v", program)
	}
}
