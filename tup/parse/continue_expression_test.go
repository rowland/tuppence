package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
)

func TestContinueExpression(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.ContinueExpression
		wantErr bool
	}{
		{
			name:  "bare continue",
			input: "continue",
			want:  ast.NewContinueExpression(nil),
		},
		{
			name:  "continue with value",
			input: "continue i",
			want: ast.NewContinueExpression(
				ast.NewIdentifier("i", nil, 0, 1),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"ContinueExpression", ContinueExpression, StringerCheck[*ast.ContinueExpression])
		})
	}
}
