package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
)

func TestReturnExpression(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.ReturnExpression
		wantErr bool
	}{
		{
			name:  "bare return",
			input: "return",
			want:  ast.NewReturnExpression(nil),
		},
		{
			name:  "return with value",
			input: "return i",
			want: ast.NewReturnExpression(
				ast.NewIdentifier("i", nil, 0, 1),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"ReturnExpression", ReturnExpression, StringerCheck[*ast.ReturnExpression])
		})
	}
}
