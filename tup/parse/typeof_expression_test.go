package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
)

func TestTypeofExpression(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.TypeofExpression
		wantErr bool
	}{
		{
			name:  "identifier expression",
			input: "typeof(x)",
			want: ast.NewTypeofExpression(
				ast.NewIdentifier("x", nil, 0, 1),
			),
		},
		{
			name:  "complex expression",
			input: "typeof(1 + 2)",
			want: ast.NewTypeofExpression(
				ast.NewAddSubExpression(
					ast.NewDecimalLiteral("1", 1, nil, 0, 0),
					ast.OpAdd,
					ast.NewDecimalLiteral("2", 2, nil, 0, 0),
				),
			),
		},
		{
			name:    "missing expression",
			input:   "typeof()",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "missing closing paren",
			input:   "typeof(x",
			want:    nil,
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"TypeofExpression", TypeofExpression, StringerCheck[*ast.TypeofExpression])
		})
	}
}
