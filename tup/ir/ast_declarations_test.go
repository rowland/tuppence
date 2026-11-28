package ir

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
)

func TestInstallCollectedTypeDeclarationsCreatesAliases(t *testing.T) {
	lhs := ast.NewTypeDeclarationLHS(nil, ast.NewTypeIdentifier("Alias", nil, 0, 5), nil)
	rhs := ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 8, 3), nil, 8, 3)
	collected := CollectTypeDeclarations(ModuleID(1), ScopeID(1), FileID(1), DeclID(0), []*ast.TypeDeclaration{ast.NewTypeDeclaration(lhs, rhs)})
	prolog := ScopeID(0x80000001)
	program := IndexedProgram{
		Scopes: ScopeTable{
			{ID: prolog, Module: ModuleID(1), Kind: ScopeProlog, Origin: testOrigin()},
			{ID: ScopeID(1), Parent: prolog, Module: ModuleID(1), Kind: ScopeModule, Initial: true, Origin: testOrigin()},
			{ID: ScopeID(2), Module: ModuleID(2), Kind: ScopeModule, Origin: testOrigin()},
		},
		ImportAliases: ImportAliasTable{{Scope: ScopeID(1), LocalName: "Int", Module: ModuleID(2), Declaration: DeclID(20), Origin: testOrigin()}},
		Declarations:  DeclarationTable{{ID: DeclID(20), Scope: ScopeID(2), Name: "Int", Kind: DeclType, Exported: true, Origin: testOrigin()}},
		Types:         TypeTable{{ID: TypeID(20), Name: "Int", Kind: TypeNominal, Declaration: DeclID(20), HomeModule: ModuleID(2), Origin: testOrigin()}},
	}
	if !program.InstallCollectedTypeDeclarations(ModuleID(1), collected) || len(program.Aliases) != 1 || program.Aliases[0].Declaration != DeclID(1) || program.Aliases[0].Target != TypeID(20) {
		t.Fatalf("resolved type declaration must create an alias: %#v", program)
	}
}

func TestCollectTypeDeclarationsBuildsScopeEntries(t *testing.T) {
	first := ast.NewTypeDeclarationLHS(nil, ast.NewTypeIdentifier("Int", nil, 4, 3), nil)
	second := ast.NewTypeDeclarationLHS(nil, ast.NewTypeIdentifier("Answer", nil, 12, 6), nil)
	rhs := ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 20, 3), nil, 20, 3)
	collected := CollectTypeDeclarations(ModuleID(1), ScopeID(1), FileID(7), DeclID(10), []*ast.TypeDeclaration{
		ast.NewTypeDeclaration(first, rhs),
		ast.NewTypeDeclaration(second, rhs),
	})
	if len(collected.Diagnostics) != 0 || len(collected.Declarations) != 2 || len(collected.Entries) != 2 || len(collected.References) != 2 {
		t.Fatalf("type declaration collection returned unexpected results: %#v", collected)
	}
	if collected.Declarations[0].ID != DeclID(11) || collected.Declarations[0].Name != "Int" || collected.Declarations[0].Scope != ScopeID(1) || collected.Entries[1].Declaration != DeclID(12) {
		t.Fatalf("collected declarations lost semantic identity: %#v", collected)
	}
	if collected.Declarations[1].Origin.Primary.File != FileID(7) || collected.Declarations[1].Origin.Primary.Start != 12 || collected.Declarations[1].Origin.Primary.End != 18 {
		t.Fatalf("collected declaration lost source span: %#v", collected.Declarations[1])
	}
	newType := CollectTypeDeclarations(ModuleID(1), ScopeID(1), FileID(7), DeclID(20), []*ast.TypeDeclaration{ast.NewTypeDeclaration(first, ast.NewNilableType(nil))})
	if len(newType.Types) != 0 || len(newType.Diagnostics) != 1 || newType.Diagnostics[0].Code != "TYPE_BODY_UNSUPPORTED" {
		t.Fatalf("structured declaration bodies must be reported explicitly: %#v", newType)
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
