package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
)

func TestTypeQualifiedDeclaration(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.TypeQualifiedDeclaration
		wantErr bool
	}{
		{
			name:  "type qualified declaration",
			input: "Person.defaultName = name",
			want: ast.NewTypeQualifiedDeclaration(
				ast.NewTypeIdentifier("Person", nil, 0, 6),
				ast.NewAssignment(
					ast.NewOrdinalAssignmentLHS([]*ast.Identifier{
						ast.NewIdentifier("defaultName", nil, 0, 11),
					}, nil),
					ast.Immutable,
					ast.NewIdentifier("name", nil, 0, 4),
				),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"TypeQualifiedDeclaration", TypeQualifiedDeclaration, StringerCheck[*ast.TypeQualifiedDeclaration])
		})
	}
}
