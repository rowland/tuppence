package ir

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/source"
)

func TestNewTypeReferenceFromASTPreservesNameAndSpan(t *testing.T) {
	parsed := ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", source.NewSource([]byte("Int"), "main.tup"), 2, 3), nil, 2, 3)
	reference := NewTypeReferenceFromAST(TypeReferenceID(1), ScopeID(2), FileID(7), parsed)
	if !reference.Valid() || reference.Name != "Int" || reference.Origin.Primary.File != FileID(7) || reference.Origin.Primary.Start != 2 || reference.Origin.Primary.End != 5 {
		t.Fatalf("AST type reference conversion lost source data: %#v", reference)
	}
}
