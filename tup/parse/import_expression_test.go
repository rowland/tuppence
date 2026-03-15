package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
)

func TestImportExpression(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.ImportExpression
		wantErr bool
	}{
		{
			name:  "basic import",
			input: `import("io")`,
			want: ast.NewImportExpression(
				ast.NewStringLiteral(`"io"`, "io", nil, 0, 4),
			),
		},
		{
			name:  "qualified module path",
			input: `import("foo/bar")`,
			want: ast.NewImportExpression(
				ast.NewStringLiteral(`"foo/bar"`, "foo/bar", nil, 0, 9),
			),
		},
		{
			name:    "missing string literal",
			input:   `import(io)`,
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"ImportExpression", ImportExpression, StringerCheck[*ast.ImportExpression])
		})
	}
}
