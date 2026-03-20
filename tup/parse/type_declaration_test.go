package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
)

func TestTypeParameter(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.TypeParameter
		wantErr bool
	}{
		{
			name:  "single type parameter",
			input: "a",
			want:  ast.NewTypeParameter(ast.NewIdentifier("a", nil, 0, 1)),
		},
		{
			name:    "type parameters use identifiers, not type identifiers",
			input:   "A",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"TypeParameter", TypeParameter, StringerCheck[*ast.TypeParameter])
		})
	}
}

func TestTypeParameters(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.TypeParameters
		wantErr bool
	}{
		{
			name:  "single parameter",
			input: "[a]",
			want: ast.NewTypeParameters([]*ast.TypeParameter{
				ast.NewTypeParameter(ast.NewIdentifier("a", nil, 0, 1)),
			}),
		},
		{
			name:  "multiple parameters",
			input: "[a, b]",
			want: ast.NewTypeParameters([]*ast.TypeParameter{
				ast.NewTypeParameter(ast.NewIdentifier("a", nil, 0, 1)),
				ast.NewTypeParameter(ast.NewIdentifier("b", nil, 0, 1)),
			}),
		},
		{
			name:    "empty type parameters are rejected",
			input:   "[]",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"TypeParameters", TypeParameters, StringerCheck[*ast.TypeParameters])
		})
	}
}

func TestNilableType(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.NilableType
		wantErr bool
	}{
		{
			name:  "nilable type reference",
			input: "?Foo",
			want: ast.NewNilableType(
				ast.NewTypeReference(nil, ast.NewTypeIdentifier("Foo", nil, 0, 3), nil, 0, 3),
			),
		},
		{
			name:  "nilable generic parameter",
			input: "?a",
			want: ast.NewNilableType(
				ast.NewIdentifier("a", nil, 0, 1),
			),
		},
		{
			name:    "missing local type reference",
			input:   "?",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"NilableType", NilableType, StringerCheck[*ast.NilableType])
		})
	}
}

func TestTypeDeclarationLHS(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.TypeDeclarationLHS
		wantErr bool
	}{
		{
			name:  "simple lhs",
			input: "Result",
			want:  ast.NewTypeDeclarationLHS(nil, ast.NewTypeIdentifier("Result", nil, 0, 6), nil),
		},
		{
			name:  "annotated generic lhs",
			input: "@serializable\nResult[a, b]",
			want: ast.NewTypeDeclarationLHS(
				[]ast.Annotation{
					ast.NewSimpleAnnotation("serializable"),
				},
				ast.NewTypeIdentifier("Result", nil, 0, 6),
				ast.NewTypeParameters([]*ast.TypeParameter{
					ast.NewTypeParameter(ast.NewIdentifier("a", nil, 0, 1)),
					ast.NewTypeParameter(ast.NewIdentifier("b", nil, 0, 1)),
				}),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"TypeDeclarationLHS", TypeDeclarationLHS, StringerCheck[*ast.TypeDeclarationLHS])
		})
	}
}

func TestTypeDeclaration(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.TypeDeclaration
		wantErr bool
	}{
		{
			name:  "type alias",
			input: "Result = foo.Bar",
			want: ast.NewTypeDeclaration(
				ast.NewTypeDeclarationLHS(nil, ast.NewTypeIdentifier("Result", nil, 0, 6), nil),
				ast.NewTypeReference(
					[]*ast.Identifier{ast.NewIdentifier("foo", nil, 0, 3)},
					ast.NewTypeIdentifier("Bar", nil, 0, 3),
					nil,
					0,
					0,
				),
			),
		},
		{
			name:  "generic nilable type declaration",
			input: "Maybe[a] = ?a",
			want: ast.NewTypeDeclaration(
				ast.NewTypeDeclarationLHS(
					nil,
					ast.NewTypeIdentifier("Maybe", nil, 0, 5),
					ast.NewTypeParameters([]*ast.TypeParameter{
						ast.NewTypeParameter(ast.NewIdentifier("a", nil, 0, 1)),
					}),
				),
				ast.NewNilableType(ast.NewIdentifier("a", nil, 0, 1)),
			),
		},
		{
			name:  "annotated type alias",
			input: "@serializable\nResult = Foo",
			want: ast.NewTypeDeclaration(
				ast.NewTypeDeclarationLHS(
					[]ast.Annotation{ast.NewSimpleAnnotation("serializable")},
					ast.NewTypeIdentifier("Result", nil, 0, 6),
					nil,
				),
				ast.NewTypeReference(nil, ast.NewTypeIdentifier("Foo", nil, 0, 3), nil, 0, 3),
			),
		},
		{
			name:    "tuple type rhs is not implemented yet",
			input:   "Person = type(name: String, age: Int)",
			wantErr: true,
		},
		{
			name:    "dynamic array rhs is not implemented yet",
			input:   "Bytes = []Byte",
			wantErr: true,
		},
		{
			name:    "fixed size array rhs is not implemented yet",
			input:   "IPv4 = [4]Byte",
			wantErr: true,
		},
		{
			name:    "missing equals",
			input:   "Result Foo",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"TypeDeclaration", TypeDeclaration, StringerCheck[*ast.TypeDeclaration])
		})
	}
}
