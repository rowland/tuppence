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
			name:  "union type rhs",
			input: "Key = Int | String",
			want: ast.NewTypeDeclaration(
				ast.NewTypeDeclarationLHS(nil, ast.NewTypeIdentifier("Key", nil, 0, 3), nil),
				ast.NewUnionType([]ast.UnionMemberType{
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("String", nil, 0, 6), nil, 0, 6),
				}),
			),
		},
		{
			name:  "named tuple union type rhs",
			input: "ComplexKey = Int | String | ComplexTuple(primary: Int, secondary: String)",
			want: ast.NewTypeDeclaration(
				ast.NewTypeDeclarationLHS(nil, ast.NewTypeIdentifier("ComplexKey", nil, 0, 10), nil),
				ast.NewUnionType([]ast.UnionMemberType{
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("String", nil, 0, 6), nil, 0, 6),
					ast.NewNamedTuple(
						ast.NewTypeIdentifier("ComplexTuple", nil, 0, 12),
						ast.NewTupleType([]ast.TupleTypeMemberNode{
							ast.NewLabeledTupleTypeMember(
								nil,
								ast.NewIdentifier("primary", nil, 0, 7),
								ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
							),
							ast.NewLabeledTupleTypeMember(
								nil,
								ast.NewIdentifier("secondary", nil, 0, 9),
								ast.NewTypeReference(nil, ast.NewTypeIdentifier("String", nil, 0, 6), nil, 0, 6),
							),
						}),
					),
				}),
			),
		},
		{
			name: "union declaration rhs",
			input: "Result[a] = union(\n    Ok()\n    Err(a)\n)",
			want: ast.NewTypeDeclaration(
				ast.NewTypeDeclarationLHS(
					nil,
					ast.NewTypeIdentifier("Result", nil, 0, 6),
					ast.NewTypeParameters([]*ast.TypeParameter{
						ast.NewTypeParameter(ast.NewIdentifier("a", nil, 0, 1)),
					}),
				),
				ast.NewUnionDeclaration(ast.UnionMembers{
					ast.NewUnionMemberDeclaration(
						nil,
						ast.NewNamedTuple(
							ast.NewTypeIdentifier("Ok", nil, 0, 2),
							ast.NewTupleType(nil),
						),
					),
					ast.NewUnionMemberDeclaration(
						nil,
						ast.NewNamedTuple(
							ast.NewTypeIdentifier("Err", nil, 0, 3),
							ast.NewTupleType([]ast.TupleTypeMemberNode{
								ast.NewTupleTypeMember(nil, ast.NewIdentifier("a", nil, 0, 1)),
							}),
						),
					),
				}),
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

func TestUnionMemberNoAnnotations(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    ast.UnionDeclarationMemberType
		wantErr bool
	}{
		{
			name:  "type reference member",
			input: "Card",
			want:  ast.NewTypeReference(nil, ast.NewTypeIdentifier("Card", nil, 0, 4), nil, 0, 4),
		},
		{
			name:  "generic member",
			input: "Result[Int, String]",
			want: ast.NewGenericType(
				ast.NewTypeReference(nil, ast.NewTypeIdentifier("Result", nil, 0, 6), nil, 0, 6),
				ast.NewTypeArgumentList([]*ast.TypeArgument{
					ast.NewTypeArgument(ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3)),
					ast.NewTypeArgument(ast.NewTypeReference(nil, ast.NewTypeIdentifier("String", nil, 0, 6), nil, 0, 6)),
				}),
			),
		},
		{
			name:    "annotated introduced member is not union_member_no_annotations",
			input:   "@flag\nOk()",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"UnionMemberNoAnnotations", UnionMemberNoAnnotations, StringerCheck[ast.UnionDeclarationMemberType])
		})
	}
}

func TestUnionMemberDeclaration(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.UnionMemberDeclaration
		wantErr bool
	}{
		{
			name:  "introduced named tuple member",
			input: "@flag\nOk()",
			want: ast.NewUnionMemberDeclaration(
				[]ast.Annotation{ast.NewSimpleAnnotation("flag")},
				ast.NewNamedTuple(
					ast.NewTypeIdentifier("Ok", nil, 0, 2),
					ast.NewTupleType(nil),
				),
			),
		},
		{
			name:  "unannotated introduced named tuple member",
			input: "Err(a)",
			want: ast.NewUnionMemberDeclaration(
				nil,
				ast.NewNamedTuple(
					ast.NewTypeIdentifier("Err", nil, 0, 3),
					ast.NewTupleType([]ast.TupleTypeMemberNode{
						ast.NewTupleTypeMember(nil, ast.NewIdentifier("a", nil, 0, 1)),
					}),
				),
			),
		},
		{
			name:  "existing type member",
			input: "Result[Int]",
			want: ast.NewUnionMemberDeclaration(
				nil,
				ast.NewGenericType(
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("Result", nil, 0, 6), nil, 0, 6),
					ast.NewTypeArgumentList([]*ast.TypeArgument{
						ast.NewTypeArgument(ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3)),
					}),
				),
			),
		},
		{
			name:    "annotations require introduced named tuple members",
			input:   "@flag\nCard",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"UnionMemberDeclaration", UnionMemberDeclaration, StringerCheck[*ast.UnionMemberDeclaration])
		})
	}
}

func TestUnionDeclaration(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.UnionDeclaration
		wantErr bool
	}{
		{
			name: "introduced generic union declaration",
			input: "union(\nOk()\nErr(a)\n)",
			want: ast.NewUnionDeclaration(ast.UnionMembers{
				ast.NewUnionMemberDeclaration(
					nil,
					ast.NewNamedTuple(
						ast.NewTypeIdentifier("Ok", nil, 0, 2),
						ast.NewTupleType(nil),
					),
				),
				ast.NewUnionMemberDeclaration(
					nil,
					ast.NewNamedTuple(
						ast.NewTypeIdentifier("Err", nil, 0, 3),
						ast.NewTupleType([]ast.TupleTypeMemberNode{
							ast.NewTupleTypeMember(nil, ast.NewIdentifier("a", nil, 0, 1)),
						}),
					),
				),
			}),
		},
		{
			name: "union declaration with existing members",
			input: "union(\nResult[Int]\nCard\n)",
			want: ast.NewUnionDeclaration(ast.UnionMembers{
				ast.NewUnionMemberDeclaration(
					nil,
					ast.NewGenericType(
						ast.NewTypeReference(nil, ast.NewTypeIdentifier("Result", nil, 0, 6), nil, 0, 6),
						ast.NewTypeArgumentList([]*ast.TypeArgument{
							ast.NewTypeArgument(ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3)),
						}),
					),
				),
				ast.NewUnionMemberDeclaration(
					nil,
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("Card", nil, 0, 4), nil, 0, 4),
				),
			}),
		},
		{
			name:    "missing trailing eol before close paren",
			input:   "union(\nOk())",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"UnionDeclaration", UnionDeclaration, StringerCheck[*ast.UnionDeclaration])
		})
	}
}

func TestUnionDeclarationWithError(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.UnionDeclarationWithError
		wantErr bool
	}{
		{
			name: "introduced members with error",
			input: "union(\nOk()\nErr(a)\nerror\n)",
			want: ast.NewUnionDeclarationWithError(ast.UnionMembers{
				ast.NewUnionMemberDeclaration(
					nil,
					ast.NewNamedTuple(
						ast.NewTypeIdentifier("Ok", nil, 0, 2),
						ast.NewTupleType(nil),
					),
				),
				ast.NewUnionMemberDeclaration(
					nil,
					ast.NewNamedTuple(
						ast.NewTypeIdentifier("Err", nil, 0, 3),
						ast.NewTupleType([]ast.TupleTypeMemberNode{
							ast.NewTupleTypeMember(nil, ast.NewIdentifier("a", nil, 0, 1)),
						}),
					),
				),
			}),
		},
		{
			name: "existing members with error",
			input: "union(\nResult[Int]\nCard\nerror\n)",
			want: ast.NewUnionDeclarationWithError(ast.UnionMembers{
				ast.NewUnionMemberDeclaration(
					nil,
					ast.NewGenericType(
						ast.NewTypeReference(nil, ast.NewTypeIdentifier("Result", nil, 0, 6), nil, 0, 6),
						ast.NewTypeArgumentList([]*ast.TypeArgument{
							ast.NewTypeArgument(ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3)),
						}),
					),
				),
				ast.NewUnionMemberDeclaration(
					nil,
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("Card", nil, 0, 4), nil, 0, 4),
				),
			}),
		},
		{
			name:    "missing error line",
			input:   "union(\nOk()\nErr(a)\n)",
			wantErr: true,
		},
		{
			name:    "error requires at least one real member",
			input:   "union(\nerror\n)",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"UnionDeclarationWithError", UnionDeclarationWithError, StringerCheck[*ast.UnionDeclarationWithError])
		})
	}
}

func TestDynamicArray(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.DynamicArrayType
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
				"DynamicArray", DynamicArray, StringerCheck[*ast.DynamicArrayType])
		})
	}
}

func TestFixedSizeArray(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.FixedSizeArrayType
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
				"FixedSizeArray", FixedSizeArray, StringerCheck[*ast.FixedSizeArrayType])
		})
	}
}

func TestErrorTuple(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.ErrorTuple
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
				"ErrorTuple", ErrorTuple, StringerCheck[*ast.ErrorTuple])
		})
	}
}

func TestUnionType(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.UnionType
		wantErr bool
	}{
		{
			name:  "any union type",
			input: "any",
			want:  ast.NewUnionType(nil),
		},
		{
			name:  "simple union type",
			input: "Int | String",
			want: ast.NewUnionType([]ast.UnionMemberType{
				ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
				ast.NewTypeReference(nil, ast.NewTypeIdentifier("String", nil, 0, 6), nil, 0, 6),
			}),
		},
		{
			name:  "named tuple member union type",
			input: "Int | ComplexTuple(primary: Int, secondary: String)",
			want: ast.NewUnionType([]ast.UnionMemberType{
				ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
				ast.NewNamedTuple(
					ast.NewTypeIdentifier("ComplexTuple", nil, 0, 12),
					ast.NewTupleType([]ast.TupleTypeMemberNode{
						ast.NewLabeledTupleTypeMember(
							nil,
							ast.NewIdentifier("primary", nil, 0, 7),
							ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
						),
						ast.NewLabeledTupleTypeMember(
							nil,
							ast.NewIdentifier("secondary", nil, 0, 9),
							ast.NewTypeReference(nil, ast.NewTypeIdentifier("String", nil, 0, 6), nil, 0, 6),
						),
					}),
				),
			}),
		},
		{
			name:  "generic member union type",
			input: "Result[Int, String] | Card",
			want: ast.NewUnionType([]ast.UnionMemberType{
				ast.NewGenericType(
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("Result", nil, 0, 6), nil, 0, 6),
					ast.NewTypeArgumentList([]*ast.TypeArgument{
						ast.NewTypeArgument(ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3)),
						ast.NewTypeArgument(ast.NewTypeReference(nil, ast.NewTypeIdentifier("String", nil, 0, 6), nil, 0, 6)),
					}),
				),
				ast.NewTypeReference(nil, ast.NewTypeIdentifier("Card", nil, 0, 4), nil, 0, 4),
			}),
		},
		{
			name:    "single type reference is not a union",
			input:   "Int",
			wantErr: true,
		},
		{
			name:    "missing right union member",
			input:   "Int |",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"UnionType", UnionType, StringerCheck[*ast.UnionType])
		})
	}
}

func TestInlineUnion(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.InlineUnion
		wantErr bool
	}{
		{
			name:  "simple inline union",
			input: "(Int | String)",
			want: ast.NewInlineUnion(
				ast.NewUnionType([]ast.UnionMemberType{
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("String", nil, 0, 6), nil, 0, 6),
				}),
			),
		},
		{
			name:  "inline union with generic member",
			input: "(Result[Int, String] | Card)",
			want: ast.NewInlineUnion(
				ast.NewUnionType([]ast.UnionMemberType{
					ast.NewGenericType(
						ast.NewTypeReference(nil, ast.NewTypeIdentifier("Result", nil, 0, 6), nil, 0, 6),
						ast.NewTypeArgumentList([]*ast.TypeArgument{
							ast.NewTypeArgument(ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3)),
							ast.NewTypeArgument(ast.NewTypeReference(nil, ast.NewTypeIdentifier("String", nil, 0, 6), nil, 0, 6)),
						}),
					),
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("Card", nil, 0, 4), nil, 0, 4),
				}),
			),
		},
		{
			name:    "parenthesized type reference is not an inline union",
			input:   "(Int)",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"InlineUnion", InlineUnion, StringerCheck[*ast.InlineUnion])
		})
	}
}

func TestUnionWithError(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.UnionWithError
		wantErr bool
	}{
		{
			name:  "exclamation shorthand",
			input: "!Result[Int, String]",
			want: ast.NewUnionWithError(
				[]ast.UnionMemberType{
					ast.NewGenericType(
						ast.NewTypeReference(nil, ast.NewTypeIdentifier("Result", nil, 0, 6), nil, 0, 6),
						ast.NewTypeArgumentList([]*ast.TypeArgument{
							ast.NewTypeArgument(ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3)),
							ast.NewTypeArgument(ast.NewTypeReference(nil, ast.NewTypeIdentifier("String", nil, 0, 6), nil, 0, 6)),
						}),
					),
				},
				true,
			),
		},
		{
			name:  "simple union with error",
			input: "Card | error",
			want: ast.NewUnionWithError(
				[]ast.UnionMemberType{
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("Card", nil, 0, 4), nil, 0, 4),
				},
				false,
			),
		},
		{
			name:  "parenthesized union with error",
			input: "(Int | String | error)",
			want: ast.NewUnionWithError(
				[]ast.UnionMemberType{
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("String", nil, 0, 6), nil, 0, 6),
				},
				false,
			),
		},
		{
			name:    "plain union without error is rejected",
			input:   "Int | String",
			wantErr: true,
		},
		{
			name:    "parenthesized plain union is rejected",
			input:   "(Int | String)",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"UnionWithError", UnionWithError, StringerCheck[*ast.UnionWithError])
		})
	}
}

func TestTypeArgument(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.TypeArgument
		wantErr bool
	}{
		{
			name:  "type reference argument",
			input: "Int",
			want: ast.NewTypeArgument(
				ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
			),
		},
		{
			name:  "generic type argument",
			input: "List[Int]",
			want: ast.NewTypeArgument(
				ast.NewGenericType(
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("List", nil, 0, 4), nil, 0, 4),
					ast.NewTypeArgumentList([]*ast.TypeArgument{
						ast.NewTypeArgument(ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3)),
					}),
				),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"TypeArgument", TypeArgument, StringerCheck[*ast.TypeArgument])
		})
	}
}

func TestTypeArgumentList(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.TypeArgumentList
		wantErr bool
	}{
		{
			name:  "multiple type arguments",
			input: "[Int, String]",
			want: ast.NewTypeArgumentList([]*ast.TypeArgument{
				ast.NewTypeArgument(ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3)),
				ast.NewTypeArgument(ast.NewTypeReference(nil, ast.NewTypeIdentifier("String", nil, 0, 6), nil, 0, 6)),
			}),
		},
		{
			name:  "nested generic type argument",
			input: "[Result[Int, String], []Byte]",
			want: ast.NewTypeArgumentList([]*ast.TypeArgument{
				ast.NewTypeArgument(
					ast.NewGenericType(
						ast.NewTypeReference(nil, ast.NewTypeIdentifier("Result", nil, 0, 6), nil, 0, 6),
						ast.NewTypeArgumentList([]*ast.TypeArgument{
							ast.NewTypeArgument(ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3)),
							ast.NewTypeArgument(ast.NewTypeReference(nil, ast.NewTypeIdentifier("String", nil, 0, 6), nil, 0, 6)),
						}),
					),
				),
				ast.NewTypeArgument(
					ast.NewDynamicArrayType(
						ast.NewTypeReference(nil, ast.NewTypeIdentifier("Byte", nil, 0, 4), nil, 0, 4),
					),
				),
			}),
		},
		{
			name:    "empty type argument list is rejected",
			input:   "[]",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"TypeArgumentList", TypeArgumentList, StringerCheck[*ast.TypeArgumentList])
		})
	}
}

func TestGenericType(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.GenericType
		wantErr bool
	}{
		{
			name:  "simple generic type",
			input: "List[Int]",
			want: ast.NewGenericType(
				ast.NewTypeReference(nil, ast.NewTypeIdentifier("List", nil, 0, 4), nil, 0, 4),
				ast.NewTypeArgumentList([]*ast.TypeArgument{
					ast.NewTypeArgument(ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3)),
				}),
			),
		},
		{
			name:  "nested generic type",
			input: "Result[List[Int], []Byte]",
			want: ast.NewGenericType(
				ast.NewTypeReference(nil, ast.NewTypeIdentifier("Result", nil, 0, 6), nil, 0, 6),
				ast.NewTypeArgumentList([]*ast.TypeArgument{
					ast.NewTypeArgument(
						ast.NewGenericType(
							ast.NewTypeReference(nil, ast.NewTypeIdentifier("List", nil, 0, 4), nil, 0, 4),
							ast.NewTypeArgumentList([]*ast.TypeArgument{
								ast.NewTypeArgument(ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3)),
							}),
						),
					),
					ast.NewTypeArgument(
						ast.NewDynamicArrayType(
							ast.NewTypeReference(nil, ast.NewTypeIdentifier("Byte", nil, 0, 4), nil, 0, 4),
						),
					),
				}),
			),
		},
		{
			name:    "type reference without arguments is not generic",
			input:   "List",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"GenericType", GenericType, StringerCheck[*ast.GenericType])
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
