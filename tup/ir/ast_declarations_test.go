package ir

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
)

func TestCollectTypeDeclarationsBuildsScopeEntries(t *testing.T) {
	first := ast.NewTypeDeclarationLHS(nil, ast.NewTypeIdentifier("Int", nil, 4, 3), nil)
	second := ast.NewTypeDeclarationLHS(nil, ast.NewTypeIdentifier("Answer", nil, 12, 6), nil)
	rhs := ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 20, 3), nil, 20, 3)
	collected := CollectTypeDeclarations(ScopeID(1), FileID(7), DeclID(10), []*ast.TypeDeclaration{
		ast.NewTypeDeclaration(first, rhs),
		ast.NewTypeDeclaration(second, rhs),
	})
	if len(collected.Diagnostics) != 0 || len(collected.Declarations) != 2 || len(collected.Entries) != 2 {
		t.Fatalf("type declaration collection returned unexpected results: %#v", collected)
	}
	if collected.Declarations[0].ID != DeclID(11) || collected.Declarations[0].Name != "Int" || collected.Declarations[0].Scope != ScopeID(1) || collected.Entries[1].Declaration != DeclID(12) {
		t.Fatalf("collected declarations lost semantic identity: %#v", collected)
	}
	if collected.Declarations[1].Origin.Primary.File != FileID(7) || collected.Declarations[1].Origin.Primary.Start != 12 || collected.Declarations[1].Origin.Primary.End != 18 {
		t.Fatalf("collected declaration lost source span: %#v", collected.Declarations[1])
	}
	program := IndexedProgram{Scopes: ScopeTable{
		{ID: ScopeID(0x80000001), Module: ModuleID(1), Kind: ScopeProlog, Origin: testOrigin()},
		{ID: ScopeID(1), Parent: ScopeID(0x80000001), Module: ModuleID(1), Kind: ScopeModule, Initial: true, Origin: testOrigin()},
	}}
	for i := range collected.Declarations {
		collected.Declarations[i].Origin = testOrigin()
	}
	if !program.InstallCollectedTypeDeclarations(ModuleID(1), collected) || len(program.Declarations) != 2 || len(program.Scopes[1].Entries) != 2 {
		t.Fatal("collected declarations must be installed in the module initial scope")
	}
	if program.InstallCollectedTypeDeclarations(ModuleID(1), collected) {
		t.Fatal("declaration installation must reject duplicate names and IDs")
	}
}
