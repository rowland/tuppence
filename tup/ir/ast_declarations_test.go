package ir

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
)

func TestCollectTypeDeclarationsRecursivelyCollectsStructuredMembers(t *testing.T) {
	ref := func(name string, offset int) *ast.TypeReference {
		return ast.NewTypeReference(nil, ast.NewTypeIdentifier(name, nil, int32(offset), int32(len(name))), nil, int32(offset), int32(len(name)))
	}
	arrayOfArrays := ast.NewDynamicArrayType(ast.NewDynamicArrayType(ref("Int", 20)))
	nestedUnion := ast.NewUnionType([]ast.UnionMemberType{ref("Bool", 30), arrayOfArrays})
	nestedTuple := ast.NewTupleType([]ast.TupleTypeMemberNode{
		ast.NewTupleTypeMember(nil, nestedUnion),
		ast.NewTupleTypeMember(nil, ast.NewDynamicArrayType(ref("String", 40))),
	})
	root := ast.NewTypeTuple(ast.NewTupleType([]ast.TupleTypeMemberNode{
		ast.NewTupleTypeMember(nil, nestedTuple),
		ast.NewTupleTypeMember(nil, ref("Float", 50)),
	}))

	collected := CollectTypeDeclarations(ModuleID(1), ScopeID(1), FileID(1), DeclID(10), []*ast.TypeDeclaration{
		ast.NewTypeDeclaration(ast.NewTypeDeclarationLHS(nil, ast.NewTypeIdentifier("Nested", nil, 0, 7), nil), root),
	})
	if len(collected.Diagnostics) != 0 {
		t.Fatalf("recursive structured collection returned diagnostics: %#v", collected.Diagnostics)
	}
	if len(collected.Structured) != 6 {
		t.Fatalf("recursive structured collection must retain every nested type: %#v", collected.Structured)
	}
	if collected.Structured[0].Kind != TypeArray || collected.Structured[1].Kind != TypeArray || collected.Structured[2].Kind != TypeUnion || collected.Structured[3].Kind != TypeArray || collected.Structured[4].Kind != TypeTuple || collected.Structured[5].Kind != TypeTuple {
		t.Fatalf("nested structured types must be collected deterministically: %#v", collected.Structured)
	}
	if got := collected.Structured[2].ElementOrder; len(got) != 2 || got[0].Reference == 0 || got[1].Type != collected.Structured[1].ID {
		t.Fatalf("union members must preserve source order: %#v", got)
	}
	if got := collected.Structured[4].ElementOrder; len(got) != 2 || got[0].Type != collected.Structured[2].ID || got[1].Type != collected.Structured[3].ID {
		t.Fatalf("nested tuple members must preserve source order: %#v", got)
	}
	if got := collected.Structured[5].ElementOrder; len(got) != 2 || got[0].Type != collected.Structured[4].ID || got[1].Reference == 0 {
		t.Fatalf("root tuple members must preserve source order: %#v", got)
	}
}

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
	arrayRHS := ast.NewDynamicArrayType(ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 8, 3), nil, 8, 3))
	arrayCollected := CollectTypeDeclarations(ModuleID(1), ScopeID(1), FileID(1), DeclID(30), []*ast.TypeDeclaration{ast.NewTypeDeclaration(ast.NewTypeDeclarationLHS(nil, ast.NewTypeIdentifier("Bytes", nil, 0, 5), nil), arrayRHS)})
	if !program.InstallCollectedTypeDeclarations(ModuleID(1), arrayCollected) || len(program.Aliases) != 2 || len(program.Types) != 2 || program.Aliases[1].Target != TypeID(31) || program.Types[1].Element != TypeID(20) {
		t.Fatalf("dynamic array declaration must resolve its element and create an alias: %#v", program)
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
	arrayRHS := ast.NewDynamicArrayType(ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 24, 3), nil, 24, 3))
	arrayType := CollectTypeDeclarations(ModuleID(1), ScopeID(1), FileID(7), DeclID(30), []*ast.TypeDeclaration{ast.NewTypeDeclaration(second, arrayRHS)})
	if len(arrayType.Diagnostics) != 0 || len(arrayType.Structured) != 1 || len(arrayType.References) != 1 || arrayType.Structured[0].ID != TypeID(31) {
		t.Fatalf("dynamic array declarations must retain pending element resolution: %#v", arrayType)
	}
	fixedRHS := ast.NewFixedSizeArrayType(ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 24, 3), nil, 24, 3), ast.NewDecimalLiteral("4", 4, nil, 28, 1))
	fixedType := CollectTypeDeclarations(ModuleID(1), ScopeID(1), FileID(7), DeclID(40), []*ast.TypeDeclaration{ast.NewTypeDeclaration(second, fixedRHS)})
	if len(fixedType.Diagnostics) != 0 || len(fixedType.Structured) != 1 || !fixedType.Structured[0].FixedLength || fixedType.Structured[0].Length != 4 {
		t.Fatalf("fixed array declarations must retain validated lengths: %#v", fixedType)
	}
	tupleRHS := ast.NewTypeTuple(ast.NewTupleType([]ast.TupleTypeMemberNode{
		ast.NewTupleTypeMember(nil, ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 32, 3), nil, 32, 3)),
		ast.NewTupleTypeMember(nil, ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 37, 3), nil, 37, 3)),
	}))
	tupleType := CollectTypeDeclarations(ModuleID(1), ScopeID(1), FileID(7), DeclID(50), []*ast.TypeDeclaration{ast.NewTypeDeclaration(second, tupleRHS)})
	if len(tupleType.Diagnostics) != 0 || len(tupleType.Structured) != 1 || tupleType.Structured[0].Kind != TypeTuple || len(tupleType.Structured[0].ElementReferences) != 2 {
		t.Fatalf("tuple declarations must retain member references: %#v", tupleType)
	}
	unionRHS := ast.NewUnionType([]ast.UnionMemberType{
		ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 42, 3), nil, 42, 3),
		ast.NewTypeReference(nil, ast.NewTypeIdentifier("Bool", nil, 47, 4), nil, 47, 4),
	})
	unionType := CollectTypeDeclarations(ModuleID(1), ScopeID(1), FileID(8), DeclID(51), []*ast.TypeDeclaration{ast.NewTypeDeclaration(second, unionRHS)})
	if len(unionType.Diagnostics) != 0 || len(unionType.Structured) != 1 || unionType.Structured[0].Kind != TypeUnion || len(unionType.Structured[0].ElementReferences) != 2 {
		t.Fatalf("union declarations must retain member references: %#v", unionType)
	}
	nestedTuple := ast.NewTypeTuple(ast.NewTupleType([]ast.TupleTypeMemberNode{
		ast.NewTupleTypeMember(nil, ast.NewDynamicArrayType(ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 55, 3), nil, 55, 3))),
	}))
	nested := CollectTypeDeclarations(ModuleID(1), ScopeID(1), FileID(9), DeclID(52), []*ast.TypeDeclaration{ast.NewTypeDeclaration(second, nestedTuple)})
	if len(nested.Diagnostics) != 0 || len(nested.Structured) != 2 || nested.Structured[0].Kind != TypeArray || nested.Structured[1].Kind != TypeTuple || len(nested.Structured[1].ElementTypes) != 1 {
		t.Fatalf("tuple declarations must retain nested structured members: %#v", nested)
	}
	nestedUnion := ast.NewUnionType([]ast.UnionMemberType{
		ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 60, 3), nil, 60, 3),
		ast.NewDynamicArrayType(ast.NewTypeReference(nil, ast.NewTypeIdentifier("Bool", nil, 65, 4), nil, 65, 4)),
	})
	nestedUnionCollected := CollectTypeDeclarations(ModuleID(1), ScopeID(1), FileID(10), DeclID(53), []*ast.TypeDeclaration{ast.NewTypeDeclaration(second, nestedUnion)})
	if len(nestedUnionCollected.Diagnostics) != 0 || len(nestedUnionCollected.Structured) != 2 || nestedUnionCollected.Structured[0].Kind != TypeArray || nestedUnionCollected.Structured[1].Kind != TypeUnion || len(nestedUnionCollected.Structured[1].ElementOrder) != 2 {
		t.Fatalf("union declarations must retain nested structured members: %#v", nestedUnionCollected)
	}
	nestedTupleValue := ast.NewTupleType([]ast.TupleTypeMemberNode{
		ast.NewTupleTypeMember(nil, ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 72, 3), nil, 72, 3)),
	})
	nestedTupleDecl := ast.NewTypeTuple(ast.NewTupleType([]ast.TupleTypeMemberNode{
		ast.NewTupleTypeMember(nil, nestedTupleValue),
	}))
	nestedTupleCollected := CollectTypeDeclarations(ModuleID(1), ScopeID(1), FileID(11), DeclID(54), []*ast.TypeDeclaration{ast.NewTypeDeclaration(second, nestedTupleDecl)})
	if len(nestedTupleCollected.Diagnostics) != 0 || len(nestedTupleCollected.Structured) != 2 || nestedTupleCollected.Structured[0].Kind != TypeTuple || nestedTupleCollected.Structured[1].Kind != TypeTuple {
		t.Fatalf("tuple declarations must retain nested tuples: %#v", nestedTupleCollected)
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
