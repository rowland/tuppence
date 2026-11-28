package ir

import "testing"

func TestBuildIndexedProgramRecordsTypeReferenceDiagnostics(t *testing.T) {
	target := CompilationTarget{OS: "linux", Architecture: "amd64", PrimaryPath: "app/main"}
	program, invalid := BuildIndexedProgramWithTypeReferences([]string{"app/main.tup"}, target, PrologConfig{}, []TypeReference{{ID: TypeReferenceID(1), Scope: ScopeID(1), Name: "Missing", Origin: Origin{Primary: Span{File: FileID(1), Start: 0, End: 7}}}})
	if len(invalid) != 0 || len(program.TypeReferences) != 0 || len(program.Diagnostics) != 1 || program.Diagnostics[0].Code != "TYPE_NOT_FOUND" {
		t.Fatalf("builder must record type-resolution diagnostics: program=%#v invalid=%v", program, invalid)
	}
}

func TestBuildIndexedProgramRequiresPrimaryPath(t *testing.T) {
	program, invalid := BuildIndexedProgram([]string{"app/main.tup"}, CompilationTarget{OS: "linux", Architecture: "amd64"}, PrologConfig{})
	if program.Modules != nil || len(invalid) != 1 || invalid[0] != "primary module path is required" {
		t.Fatalf("missing primary module must be rejected: program=%#v invalid=%v", program, invalid)
	}
}

func TestBuildIndexedProgramRejectsMissingPrologModule(t *testing.T) {
	target := CompilationTarget{OS: "linux", Architecture: "amd64", PrimaryPath: "app/main", PrologPaths: []string{"core/missing"}}
	program, invalid := BuildIndexedProgram([]string{"app/main.tup"}, target, PrologConfig{Enabled: true})
	if len(invalid) != 1 || invalid[0] != "core/missing" || program.Valid() {
		t.Fatalf("missing prolog path must be rejected: program=%#v invalid=%v", program, invalid)
	}
	if len(program.Diagnostics) != 1 || program.Diagnostics[0].Code != "PROLOG_MODULE_NOT_FOUND" {
		t.Fatalf("missing prolog path must produce a diagnostic: %#v", program.Diagnostics)
	}
}

func TestBuildIndexedProgramSelectsConfiguredPrologModules(t *testing.T) {
	target := CompilationTarget{OS: "linux", Architecture: "amd64", PrimaryPath: "core/math", PrologPaths: []string{"core/math"}}
	program, invalid := BuildIndexedProgram([]string{"core/math.tup", "app/main.tup"}, target, PrologConfig{Enabled: true})
	if len(invalid) != 0 || !program.Valid() || len(program.Prolog.Modules) != 1 {
		t.Fatalf("configured prolog module must be selected: %#v invalid=%v", program, invalid)
	}
	if program.Prolog.Architecture != "amd64" {
		t.Fatal("prolog architecture must come from the compilation target")
	}
}

func TestBuildIndexedProgramUsesOneTargetConfiguration(t *testing.T) {
	target := CompilationTarget{OS: "linux", Architecture: "amd64", PrimaryPath: "app/main", ModuleRoots: ModuleResolutionConfig{ProgramRoot: "/program", SearchPaths: []string{"/stdlib"}}}
	program, invalid := BuildIndexedProgram([]string{"app/main.tup", "app/main--linux.tup", "app/main--darwin.tup"}, target, PrologConfig{})
	if len(invalid) != 0 || !program.Valid() {
		t.Fatalf("constructed program must be valid: %#v invalid=%v", program, invalid)
	}
	if program.ModuleResolution.ProgramRoot != "/program" || program.PrimaryModule != ModuleID(1) || len(program.Scopes) != 2 || program.Scopes[0].Initial || program.Scopes[0].Kind != ScopeProlog || !program.Scopes[1].Initial || len(program.SourceFiles) != 2 || len(program.Modules) != 1 {
		t.Fatalf("target configuration was not preserved: %#v", program)
	}
}

func TestBuildIndexedProgramSelectsTargetSpecificPrologFiles(t *testing.T) {
	linux := CompilationTarget{OS: "linux", Architecture: "amd64", PrimaryPath: "core/math", PrologPaths: []string{"core/math"}}
	program, invalid := BuildIndexedProgram([]string{"core/math--linux.tup"}, linux, PrologConfig{Enabled: true})
	if len(invalid) != 0 || !program.Valid() || len(program.SourceFiles) != 1 {
		t.Fatalf("linux prolog file must be selected: program=%#v invalid=%v", program, invalid)
	}
	darwin := CompilationTarget{OS: "darwin", Architecture: "arm64", PrimaryPath: "app/main", PrologPaths: []string{"core/math"}}
	program, invalid = BuildIndexedProgram([]string{"core/math--linux.tup", "app/main.tup"}, darwin, PrologConfig{Enabled: true})
	if len(invalid) != 1 || invalid[0] != "core/math" || len(program.Diagnostics) != 1 || program.Diagnostics[0].Code != "PROLOG_MODULE_NOT_FOUND" {
		t.Fatalf("non-matching prolog file must be rejected: program=%#v invalid=%v", program, invalid)
	}
}
