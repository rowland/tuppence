package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
)

func TestExportAssignment(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.ExportAssignment
		wantErr bool
	}{
		{
			name:  "simple export assignment",
			input: "value: 1",
			want: ast.NewExportAssignment(*ast.NewAssignment(
				ast.NewOrdinalAssignmentLHS([]*ast.Identifier{
					ast.NewIdentifier("value", nil, 0, 5),
				}, nil),
				ast.Immutable,
				ast.NewDecimalLiteral("1", 1, nil, 0, 1),
			)),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"ExportAssignment", ExportAssignment, StringerCheck[*ast.ExportAssignment])
		})
	}
}

func TestExportFunctionDeclaration(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.ExportFunctionDeclaration
		wantErr bool
	}{
		{
			name:  "simple export function declaration",
			input: "sqr: fn(i: Int) Int { i * i }",
			want: ast.NewExportFunctionDeclaration(
				ast.NewFunctionDeclaration(
					nil,
					ast.NewFunctionDeclarationLHS(
						ast.NewFunctionIdentifier("sqr", nil, 0, 3),
						nil,
					),
					ast.NewFunctionDeclarationType(
						false,
						[]ast.FunctionTypeParameter{
							ast.NewLabeledParameter(
								ast.NewAnnotations(nil),
								ast.NewIdentifier("i", nil, 0, 1),
								ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
							),
						},
						ast.NewReturnType(
							ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
						),
						false,
					),
					ast.NewBlock(
						ast.NewBlockBody(
							[]ast.Statement{},
							ast.NewMulDivExpression(
								ast.NewIdentifier("i", nil, 0, 1),
								ast.OpMul,
								ast.NewIdentifier("i", nil, 0, 1),
							),
						),
					),
				),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"ExportFunctionDeclaration", ExportFunctionDeclaration, StringerCheck[*ast.ExportFunctionDeclaration])
		})
	}
}

func TestExportFunctionTypeDeclaration(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.ExportFunctionTypeDeclaration
		wantErr bool
	}{
		{
			name:  "simple export function type declaration",
			input: "Transformer: fn(Int) String",
			want: ast.NewExportFunctionTypeDeclaration(
				ast.NewFunctionTypeDeclaration(
					ast.NewTypeIdentifier("Transformer", nil, 0, 11),
					nil,
					ast.NewFunctionType(
						false,
						[]ast.FunctionTypeParameter{
							ast.NewParameter(
								ast.NewAnnotations(nil),
								ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
							),
						},
						ast.NewReturnType(
							ast.NewTypeReference(nil, ast.NewTypeIdentifier("String", nil, 0, 6), nil, 0, 6),
						),
					),
				),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"ExportFunctionTypeDeclaration", ExportFunctionTypeDeclaration, StringerCheck[*ast.ExportFunctionTypeDeclaration])
		})
	}
}

func TestExportTypeDeclaration(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.ExportTypeDeclaration
		wantErr bool
	}{
		{
			name:  "simple export type declaration",
			input: "Person: type(name: String)",
			want: ast.NewExportTypeDeclaration(*ast.NewTypeDeclaration(
				ast.NewTypeDeclarationLHS(
					nil,
					ast.NewTypeIdentifier("Person", nil, 0, 6),
					nil,
				),
				ast.NewTypeTuple(
					ast.NewTupleType([]ast.TupleTypeMemberNode{
						ast.NewLabeledTupleTypeMember(
							ast.NewAnnotations(nil),
							ast.NewIdentifier("name", nil, 0, 4),
							ast.NewTypeReference(nil, ast.NewTypeIdentifier("String", nil, 0, 6), nil, 0, 6),
						),
					}),
				),
			)),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"ExportTypeDeclaration", ExportTypeDeclaration, StringerCheck[*ast.ExportTypeDeclaration])
		})
	}
}

func TestExportTypeQualifiedDeclaration(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.ExportTypeQualifiedDeclaration
		wantErr bool
	}{
		{
			name:  "simple export type qualified declaration",
			input: "Person.name: \"Brent\"",
			want: ast.NewExportTypeQualifiedDeclaration(
				ast.NewTypeQualifiedDeclaration(
					ast.NewTypeIdentifier("Person", nil, 0, 6),
					ast.NewAssignment(
						ast.NewOrdinalAssignmentLHS([]*ast.Identifier{
							ast.NewIdentifier("name", nil, 0, 4),
						}, nil),
						ast.Immutable,
						ast.NewStringLiteral(`"Brent"`, "Brent", nil, 0, 7),
					),
				),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"ExportTypeQualifiedDeclaration", ExportTypeQualifiedDeclaration, StringerCheck[*ast.ExportTypeQualifiedDeclaration])
		})
	}
}

func TestExportTypeQualifiedFunctionDeclaration(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.ExportTypeQualifiedFunctionDeclaration
		wantErr bool
	}{
		{
			name:  "simple export type qualified function declaration",
			input: "Person.greet: fn() String { \"Hello\" }",
			want: ast.NewExportTypeQualifiedFunctionDeclaration(
				ast.NewTypeQualifiedFunctionDeclaration(
					ast.NewTypeIdentifier("Person", nil, 0, 6),
					ast.NewFunctionDeclaration(
						nil,
						ast.NewFunctionDeclarationLHS(
							ast.NewFunctionIdentifier("greet", nil, 0, 5),
							nil,
						),
						ast.NewFunctionDeclarationType(
							false,
							[]ast.FunctionTypeParameter{},
							ast.NewReturnType(
								ast.NewTypeReference(nil, ast.NewTypeIdentifier("String", nil, 0, 6), nil, 0, 6),
							),
							false,
						),
						ast.NewBlock(
							ast.NewBlockBody(
								[]ast.Statement{},
								ast.NewStringLiteral(`"Hello"`, "Hello", nil, 0, 7),
							),
						),
					),
				),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"ExportTypeQualifiedFunctionDeclaration", ExportTypeQualifiedFunctionDeclaration, StringerCheck[*ast.ExportTypeQualifiedFunctionDeclaration])
		})
	}
}

func TestExportDeclaration(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    ast.ExportDeclaration
		wantErr bool
	}{
		{
			name:  "export assignment",
			input: "value: 1",
			want: ast.NewExportAssignment(*ast.NewAssignment(
				ast.NewOrdinalAssignmentLHS([]*ast.Identifier{
					ast.NewIdentifier("value", nil, 0, 5),
				}, nil),
				ast.Immutable,
				ast.NewDecimalLiteral("1", 1, nil, 0, 1),
			)),
		},
		{
			name:  "export function declaration",
			input: "sqr: fn(i: Int) Int { i * i }",
			want: ast.NewExportFunctionDeclaration(
				ast.NewFunctionDeclaration(
					nil,
					ast.NewFunctionDeclarationLHS(
						ast.NewFunctionIdentifier("sqr", nil, 0, 3),
						nil,
					),
					ast.NewFunctionDeclarationType(
						false,
						[]ast.FunctionTypeParameter{
							ast.NewLabeledParameter(
								ast.NewAnnotations(nil),
								ast.NewIdentifier("i", nil, 0, 1),
								ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
							),
						},
						ast.NewReturnType(
							ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
						),
						false,
					),
					ast.NewBlock(
						ast.NewBlockBody(
							[]ast.Statement{},
							ast.NewMulDivExpression(
								ast.NewIdentifier("i", nil, 0, 1),
								ast.OpMul,
								ast.NewIdentifier("i", nil, 0, 1),
							),
						),
					),
				),
			),
		},
		{
			name:  "export function type declaration",
			input: "Transformer: fn(Int) String",
			want: ast.NewExportFunctionTypeDeclaration(
				ast.NewFunctionTypeDeclaration(
					ast.NewTypeIdentifier("Transformer", nil, 0, 11),
					nil,
					ast.NewFunctionType(
						false,
						[]ast.FunctionTypeParameter{
							ast.NewParameter(
								ast.NewAnnotations(nil),
								ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
							),
						},
						ast.NewReturnType(
							ast.NewTypeReference(nil, ast.NewTypeIdentifier("String", nil, 0, 6), nil, 0, 6),
						),
					),
				),
			),
		},
		{
			name:  "export type declaration",
			input: "Person: type(name: String)",
			want: ast.NewExportTypeDeclaration(*ast.NewTypeDeclaration(
				ast.NewTypeDeclarationLHS(
					nil,
					ast.NewTypeIdentifier("Person", nil, 0, 6),
					nil,
				),
				ast.NewTypeTuple(
					ast.NewTupleType([]ast.TupleTypeMemberNode{
						ast.NewLabeledTupleTypeMember(
							ast.NewAnnotations(nil),
							ast.NewIdentifier("name", nil, 0, 4),
							ast.NewTypeReference(nil, ast.NewTypeIdentifier("String", nil, 0, 6), nil, 0, 6),
						),
					}),
				),
			)),
		},
		{
			name:  "export type qualified declaration",
			input: "Person.name: \"Brent\"",
			want: ast.NewExportTypeQualifiedDeclaration(
				ast.NewTypeQualifiedDeclaration(
					ast.NewTypeIdentifier("Person", nil, 0, 6),
					ast.NewAssignment(
						ast.NewOrdinalAssignmentLHS([]*ast.Identifier{
							ast.NewIdentifier("name", nil, 0, 4),
						}, nil),
						ast.Immutable,
						ast.NewStringLiteral(`"Brent"`, "Brent", nil, 0, 7),
					),
				),
			),
		},
		{
			name:  "export type qualified function declaration",
			input: "Person.greet: fn() String { \"Hello\" }",
			want: ast.NewExportTypeQualifiedFunctionDeclaration(
				ast.NewTypeQualifiedFunctionDeclaration(
					ast.NewTypeIdentifier("Person", nil, 0, 6),
					ast.NewFunctionDeclaration(
						nil,
						ast.NewFunctionDeclarationLHS(
							ast.NewFunctionIdentifier("greet", nil, 0, 5),
							nil,
						),
						ast.NewFunctionDeclarationType(
							false,
							[]ast.FunctionTypeParameter{},
							ast.NewReturnType(
								ast.NewTypeReference(nil, ast.NewTypeIdentifier("String", nil, 0, 6), nil, 0, 6),
							),
							false,
						),
						ast.NewBlock(
							ast.NewBlockBody(
								[]ast.Statement{},
								ast.NewStringLiteral(`"Hello"`, "Hello", nil, 0, 7),
							),
						),
					),
				),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"ExportDeclaration", ExportDeclaration, StringerCheck[ast.ExportDeclaration])
		})
	}
}
