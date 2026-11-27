package ir

import "testing"

func TestIndexedProgramDumpIsDeterministic(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{
			{ID: ModuleID(1), Name: "main", Files: []FileID{FileID(1)}, Origin: testOrigin()},
		},
		Scopes: ScopeTable{
			{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule, Origin: testOrigin(), Entries: []ScopeEntry{{Name: "answer", Declaration: DeclID(1)}}},
		},
		Declarations: DeclarationTable{
			{ID: DeclID(1), Scope: ScopeID(1), Name: "answer", Kind: DeclValue, Origin: testOrigin()},
		},
		Bindings: BindingTable{
			{ID: BindingID(1), Scope: ScopeID(1), Name: "answer", TypeSlot: TypeVarID(1), Value: ExprID(1), Origin: testOrigin()},
		},
	}
	want := "modules:\n  1 name=main path= files=[1] imports=[]\nsource_files:\nprimary_module:\n  id=0\ntypes:\nconstructors:\nscopes:\n  1 module=1 parent=0 kind=1 initial=false origin_file=1 entries=[answer:1]\ndeclarations:\n  1 scope=1 kind=1 name=answer owner_type=0 exported=false\nbindings:\n  1 scope=1 previous=0 name=answer type=1 value=1 mutable=false\nfunctions:\naliases:\nexpressions:\ncalls:\nimport_aliases:\nliteral_defaults:\n  int=0 float=0\nmodule_resolution:\n  program_root= search_paths=[]\nprolog:\n  enabled=false modules=[] architecture=\ndiagnostics:\n"
	if got := program.Dump(); got != want {
		t.Fatalf("unexpected dump:\n got: %q\nwant: %q", got, want)
	}
}
