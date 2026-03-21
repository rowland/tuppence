package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
)

func TestScopedIdentifier(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.ScopedIdentifier
		wantErr bool
	}{
		{
			name:  "simple identifier",
			input: "value",
			want: ast.NewScopedIdentifier([]*ast.Identifier{
				ast.NewIdentifier("value", nil, 0, 5),
			}),
		},
		{
			name:  "module scoped identifier",
			input: "math.pi",
			want: ast.NewScopedIdentifier([]*ast.Identifier{
				ast.NewIdentifier("math", nil, 0, 4),
				ast.NewIdentifier("pi", nil, 0, 2),
			}),
		},
		{
			name:    "function identifier not allowed as scoped identifier segment",
			input:   "fail!",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"ScopedIdentifier", ScopedIdentifier, StringerCheck[*ast.ScopedIdentifier])
		})
	}
}

func TestConstant(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.Constant
		wantErr bool
	}{
		{
			name:  "literal constant",
			input: "1",
			want: ast.NewConstant(
				ast.NewDecimalLiteral("1", 1, nil, 0, 1),
			),
		},
		{
			name:  "scoped identifier constant",
			input: "math.pi",
			want: ast.NewConstant(
				ast.NewScopedIdentifier([]*ast.Identifier{
					ast.NewIdentifier("math", nil, 0, 4),
					ast.NewIdentifier("pi", nil, 0, 2),
				}),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"Constant", Constant, StringerCheck[*ast.Constant])
		})
	}
}
