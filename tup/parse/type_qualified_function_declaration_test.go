package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
)

func TestTypeQualifiedFunctionDeclaration(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.TypeQualifiedFunctionDeclaration
		wantErr bool
	}{
		{
			name:  "type qualified function declaration",
			input: `Person.greet = fn(name: String) String { "hello" }`,
			want: ast.NewTypeQualifiedFunctionDeclaration(
				ast.NewTypeIdentifier("Person", nil, 0, 6),
				ast.NewFunctionDeclaration(
					nil,
					ast.NewFunctionDeclarationLHS(
						ast.NewFunctionIdentifier("greet", nil, 0, 5),
						nil,
					),
					ast.NewFunctionDeclarationType(
						false,
						[]ast.FunctionTypeParameter{
							ast.NewLabeledParameter(
								nil,
								ast.NewIdentifier("name", nil, 0, 4),
								ast.NewTypeReference(nil, ast.NewTypeIdentifier("String", nil, 0, 6), nil, 0, 6),
							),
						},
						ast.NewTypeReference(nil, ast.NewTypeIdentifier("String", nil, 0, 6), nil, 0, 6),
						false,
					),
					ast.NewBlock(
						ast.NewBlockBody(
							nil,
							ast.NewStringLiteral(`"hello"`, "hello", nil, 0, 7),
						),
					),
				),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"TypeQualifiedFunctionDeclaration", TypeQualifiedFunctionDeclaration, StringerCheck[*ast.TypeQualifiedFunctionDeclaration])
		})
	}
}
