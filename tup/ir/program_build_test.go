package ir

import "testing"

func TestBuildIndexedProgramRejectsMissingPrologModule(t *testing.T) {
	target := CompilationTarget{OS: "linux", Architecture: "amd64", PrologPaths: []string{"core/missing"}}
	program, invalid := BuildIndexedProgram([]string{"app/main.tup"}, target, PrologConfig{Enabled: true})
	if len(invalid) != 1 || invalid[0] != "core/missing" || program.Valid() {
		t.Fatalf("missing prolog path must be rejected: program=%#v invalid=%v", program, invalid)
	}
	if len(program.Diagnostics) != 1 || program.Diagnostics[0].Code != "PROLOG_MODULE_NOT_FOUND" {
		t.Fatalf("missing prolog path must produce a diagnostic: %#v", program.Diagnostics)
	}
}

func TestBuildIndexedProgramSelectsConfiguredPrologModules(t *testing.T) {
	target := CompilationTarget{OS: "linux", Architecture: "amd64", PrologPaths: []string{"core/math"}}
	program, invalid := BuildIndexedProgram([]string{"core/math.tup", "app/main.tup"}, target, PrologConfig{Enabled: true})
	if len(invalid) != 0 || !program.Valid() || len(program.Prolog.Modules) != 1 {
		t.Fatalf("configured prolog module must be selected: %#v invalid=%v", program, invalid)
	}
	if program.Prolog.Architecture != "amd64" {
		t.Fatal("prolog architecture must come from the compilation target")
	}
}

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

func TestBuildIndexedProgramSelectsTargetSpecificPrologFiles(t *testing.T) {
	linux := CompilationTarget{OS: "linux", Architecture: "amd64", PrologPaths: []string{"core/math"}}
	program, invalid := BuildIndexedProgram([]string{"core/math--linux.tup"}, linux, PrologConfig{Enabled: true})
	if len(invalid) != 0 || !program.Valid() || len(program.SourceFiles) != 1 {
		t.Fatalf("linux prolog file must be selected: program=%#v invalid=%v", program, invalid)
	}
	darwin := CompilationTarget{OS: "darwin", Architecture: "arm64", PrologPaths: []string{"core/math"}}
	program, invalid = BuildIndexedProgram([]string{"core/math--linux.tup"}, darwin, PrologConfig{Enabled: true})
	if len(invalid) != 1 || invalid[0] != "core/math" || len(program.Diagnostics) != 0 {
		t.Fatalf("non-matching prolog file must be rejected: program=%#v invalid=%v", program, invalid)
	}
}
