package ir

import "testing"

func TestIndexedProgramDumpIsDeterministic(t *testing.T) {
	program := IndexedProgram{
		Modules: ModuleTable{
			{ID: ModuleID(1), Name: "main", Origin: testOrigin()},
		},
		Scopes: ScopeTable{
			{ID: ScopeID(1), Module: ModuleID(1), Kind: ScopeModule, Entries: []ScopeEntry{{Name: "answer", Declaration: DeclID(1)}}},
		},
		Declarations: DeclarationTable{
			{ID: DeclID(1), Scope: ScopeID(1), Name: "answer", Kind: DeclValue, Origin: testOrigin()},
		},
		Bindings: BindingTable{
			{ID: BindingID(1), Scope: ScopeID(1), Name: "answer", TypeSlot: TypeVarID(1), Value: ExprID(1), Origin: testOrigin()},
		},
	}
	want := "modules:\n  1 name=main\ntypes:\nscopes:\n  1 module=1 parent=0 kind=1 entries=[answer:1]\ndeclarations:\n  1 scope=1 kind=1 name=answer\nbindings:\n  1 scope=1 previous=0 name=answer type=1 value=1 mutable=false\nfunctions:\naliases:\ndiagnostics:\n"
	if got := program.Dump(); got != want {
		t.Fatalf("unexpected dump:\n got: %q\nwant: %q", got, want)
	}
}
