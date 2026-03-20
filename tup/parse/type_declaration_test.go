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
			name:  "tuple type rhs",
			input: "Person = type(name: String, age: Int)",
			want: ast.NewTypeDeclaration(
				ast.NewTypeDeclarationLHS(nil, ast.NewTypeIdentifier("Person", nil, 0, 6), nil),
				ast.NewTupleType([]ast.TupleTypeMemberNode{
					ast.NewLabeledTupleTypeMember(
						nil,
						ast.NewIdentifier("name", nil, 0, 4),
						ast.NewTypeReference(nil, ast.NewTypeIdentifier("String", nil, 0, 6), nil, 0, 6),
					),
					ast.NewLabeledTupleTypeMember(
						nil,
						ast.NewIdentifier("age", nil, 0, 3),
						ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
					),
				}),
			),
		},
		{
			name:  "nested tuple type rhs",
			input: "Nested = type(id: Int, data: (name: String, value: Float))",
			want: ast.NewTypeDeclaration(
				ast.NewTypeDeclarationLHS(nil, ast.NewTypeIdentifier("Nested", nil, 0, 6), nil),
				ast.NewTupleType([]ast.TupleTypeMemberNode{
					ast.NewLabeledTupleTypeMember(
						nil,
						ast.NewIdentifier("id", nil, 0, 2),
						ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
					),
					ast.NewLabeledTupleTypeMember(
						nil,
						ast.NewIdentifier("data", nil, 0, 4),
						ast.NewTupleType([]ast.TupleTypeMemberNode{
							ast.NewLabeledTupleTypeMember(
								nil,
								ast.NewIdentifier("name", nil, 0, 4),
								ast.NewTypeReference(nil, ast.NewTypeIdentifier("String", nil, 0, 6), nil, 0, 6),
							),
							ast.NewLabeledTupleTypeMember(
								nil,
								ast.NewIdentifier("value", nil, 0, 5),
								ast.NewTypeReference(nil, ast.NewTypeIdentifier("Float", nil, 0, 5), nil, 0, 5),
							),
						}),
					),
				}),
			),
		},
		{
			name:  "dynamic array rhs",
			input: "Bytes = []Byte",
			want: ast.NewTypeDeclaration(
				ast.NewTypeDeclarationLHS(nil, ast.NewTypeIdentifier("Bytes", nil, 0, 5), nil),
				ast.NewDynamicArrayType(
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("Byte", nil, 0, 4), nil, 0, 4),
				),
			),
		},
		{
			name:  "nested dynamic array rhs",
			input: "Grid = [][]Int",
			want: ast.NewTypeDeclaration(
				ast.NewTypeDeclarationLHS(nil, ast.NewTypeIdentifier("Grid", nil, 0, 4), nil),
				ast.NewDynamicArrayType(
					ast.NewDynamicArrayType(
						ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
					),
				),
			),
		},
		{
			name:  "error tuple rhs",
			input: "HttpError = error(code: Int, message: String)",
			want: ast.NewTypeDeclaration(
				ast.NewTypeDeclarationLHS(nil, ast.NewTypeIdentifier("HttpError", nil, 0, 9), nil),
				ast.NewErrorTuple(
					ast.NewTupleType([]ast.TupleTypeMemberNode{
						ast.NewLabeledTupleTypeMember(
							nil,
							ast.NewIdentifier("code", nil, 0, 4),
							ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
						),
						ast.NewLabeledTupleTypeMember(
							nil,
							ast.NewIdentifier("message", nil, 0, 7),
							ast.NewTypeReference(nil, ast.NewTypeIdentifier("String", nil, 0, 6), nil, 0, 6),
						),
					}),
				),
			),
		},
		{
			name:  "single member error tuple rhs",
			input: "BogusCard = error(Card)",
			want: ast.NewTypeDeclaration(
				ast.NewTypeDeclarationLHS(nil, ast.NewTypeIdentifier("BogusCard", nil, 0, 9), nil),
				ast.NewErrorTuple(
					ast.NewTupleType([]ast.TupleTypeMemberNode{
						ast.NewTupleTypeMember(
							nil,
							ast.NewTypeReference(nil, ast.NewTypeIdentifier("Card", nil, 0, 4), nil, 0, 4),
						),
					}),
				),
			),
		},
		{
			name:  "fixed size array rhs",
			input: "IPv4 = [4]Byte",
			want: ast.NewTypeDeclaration(
				ast.NewTypeDeclarationLHS(nil, ast.NewTypeIdentifier("IPv4", nil, 0, 4), nil),
				ast.NewFixedSizeArrayType(
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("Byte", nil, 0, 4), nil, 0, 4),
					ast.NewDecimalLiteral("4", 4, nil, 0, 1),
				),
			),
		},
		{
			name:  "nested fixed size array rhs",
			input: "Matrix = [3][3]Int",
			want: ast.NewTypeDeclaration(
				ast.NewTypeDeclarationLHS(nil, ast.NewTypeIdentifier("Matrix", nil, 0, 6), nil),
				ast.NewFixedSizeArrayType(
					ast.NewFixedSizeArrayType(
						ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
						ast.NewDecimalLiteral("3", 3, nil, 0, 1),
					),
					ast.NewDecimalLiteral("3", 3, nil, 0, 1),
				),
			),
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

func TestDynamicArray(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    ast.TypeDeclarationRHS
		wantErr bool
	}{
		{
			name:  "simple dynamic array",
			input: "[]Byte",
			want: ast.NewDynamicArrayType(
				ast.NewTypeReference(nil, ast.NewTypeIdentifier("Byte", nil, 0, 4), nil, 0, 4),
			),
		},
		{
			name:  "nested dynamic array",
			input: "[][]Int",
			want: ast.NewDynamicArrayType(
				ast.NewDynamicArrayType(
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
				),
			),
		},
		{
			name:    "missing element type",
			input:   "[]",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"DynamicArray", DynamicArray, StringerCheck[ast.TypeDeclarationRHS])
		})
	}
}

func TestFixedSizeArray(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    ast.TypeDeclarationRHS
		wantErr bool
	}{
		{
			name:  "simple fixed size array",
			input: "[4]Byte",
			want: ast.NewFixedSizeArrayType(
				ast.NewTypeReference(nil, ast.NewTypeIdentifier("Byte", nil, 0, 4), nil, 0, 4),
				ast.NewDecimalLiteral("4", 4, nil, 0, 1),
			),
		},
		{
			name:  "identifier sized fixed size array",
			input: "[n]Byte",
			want: ast.NewFixedSizeArrayType(
				ast.NewTypeReference(nil, ast.NewTypeIdentifier("Byte", nil, 0, 4), nil, 0, 4),
				ast.NewIdentifier("n", nil, 0, 1),
			),
		},
		{
			name:  "nested fixed size array",
			input: "[3][3]Int",
			want: ast.NewFixedSizeArrayType(
				ast.NewFixedSizeArrayType(
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
					ast.NewDecimalLiteral("3", 3, nil, 0, 1),
				),
				ast.NewDecimalLiteral("3", 3, nil, 0, 1),
			),
		},
		{
			name:    "missing size",
			input:   "[]Byte",
			wantErr: true,
		},
		{
			name:    "missing element type",
			input:   "[4]",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"FixedSizeArray", FixedSizeArray, StringerCheck[ast.TypeDeclarationRHS])
		})
	}
}

func TestErrorTuple(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    ast.TypeDeclarationRHS
		wantErr bool
	}{
		{
			name:  "labeled error tuple",
			input: "error(code: Int, message: String)",
			want: ast.NewErrorTuple(
				ast.NewTupleType([]ast.TupleTypeMemberNode{
					ast.NewLabeledTupleTypeMember(
						nil,
						ast.NewIdentifier("code", nil, 0, 4),
						ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
					),
					ast.NewLabeledTupleTypeMember(
						nil,
						ast.NewIdentifier("message", nil, 0, 7),
						ast.NewTypeReference(nil, ast.NewTypeIdentifier("String", nil, 0, 6), nil, 0, 6),
					),
				}),
			),
		},
		{
			name:  "single member error tuple",
			input: "error(Card)",
			want: ast.NewErrorTuple(
				ast.NewTupleType([]ast.TupleTypeMemberNode{
					ast.NewTupleTypeMember(
						nil,
						ast.NewTypeReference(nil, ast.NewTypeIdentifier("Card", nil, 0, 4), nil, 0, 4),
					),
				}),
			),
		},
		{
			name:    "missing tuple type",
			input:   "error",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"ErrorTuple", ErrorTuple, StringerCheck[ast.TypeDeclarationRHS])
		})
	}
}

func TestTupleType(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.TupleType
		wantErr bool
	}{
		{
			name:  "ordinal tuple type",
			input: "(Int, String)",
			want: ast.NewTupleType([]ast.TupleTypeMemberNode{
				ast.NewTupleTypeMember(nil, ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3)),
				ast.NewTupleTypeMember(nil, ast.NewTypeReference(nil, ast.NewTypeIdentifier("String", nil, 0, 6), nil, 0, 6)),
			}),
		},
		{
			name:  "labeled tuple type",
			input: "(name: String, age: Int)",
			want: ast.NewTupleType([]ast.TupleTypeMemberNode{
				ast.NewLabeledTupleTypeMember(
					nil,
					ast.NewIdentifier("name", nil, 0, 4),
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("String", nil, 0, 6), nil, 0, 6),
				),
				ast.NewLabeledTupleTypeMember(
					nil,
					ast.NewIdentifier("age", nil, 0, 3),
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
				),
			}),
		},
		{
			name:  "nested tuple type member",
			input: "(coords: (x: Float, y: Float))",
			want: ast.NewTupleType([]ast.TupleTypeMemberNode{
				ast.NewLabeledTupleTypeMember(
					nil,
					ast.NewIdentifier("coords", nil, 0, 6),
					ast.NewTupleType([]ast.TupleTypeMemberNode{
						ast.NewLabeledTupleTypeMember(
							nil,
							ast.NewIdentifier("x", nil, 0, 1),
							ast.NewTypeReference(nil, ast.NewTypeIdentifier("Float", nil, 0, 5), nil, 0, 5),
						),
						ast.NewLabeledTupleTypeMember(
							nil,
							ast.NewIdentifier("y", nil, 0, 1),
							ast.NewTypeReference(nil, ast.NewTypeIdentifier("Float", nil, 0, 5), nil, 0, 5),
						),
					}),
				),
			}),
		},
		{
			name:  "nilable tuple type member",
			input: "(id: ?Int)",
			want: ast.NewTupleType([]ast.TupleTypeMemberNode{
				ast.NewLabeledTupleTypeMember(
					nil,
					ast.NewIdentifier("id", nil, 0, 2),
					ast.NewNilableType(ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3)),
				),
			}),
		},
		{
			name:    "mixed labeled and ordinal members are rejected",
			input:   "(name: String, Int)",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"TupleType", TupleType, StringerCheck[*ast.TupleType])
		})
	}
}
