package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
)

func TestBreakExpression(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.BreakExpression
		wantErr bool
	}{
		{
			name:  "bare break",
			input: "break",
			want:  ast.NewBreakExpression(nil),
		},
		{
			name:  "break with value",
			input: "break i",
			want: ast.NewBreakExpression(
				ast.NewIdentifier("i", nil, 0, 1),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"BreakExpression", BreakExpression, StringerCheck[*ast.BreakExpression])
		})
	}
}
