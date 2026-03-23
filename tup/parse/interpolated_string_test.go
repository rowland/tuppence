package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
)

func TestInterpolation(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.Interpolation
		wantErr bool
	}{
		{
			name:  "identifier",
			input: `"\(name)"`,
			want: &ast.Interpolation{
				Expression: ast.NewIdentifier("name", nil, 0, 4),
			},
		},
		{
			name:  "expression",
			input: `"\(a + b)"`,
			want: &ast.Interpolation{
				Expression: ast.NewAddSubExpression(
					ast.NewIdentifier("a", nil, 0, 1),
					ast.OpAdd,
					ast.NewIdentifier("b", nil, 0, 1),
				),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"Interpolation", Interpolation, StringerCheck[*ast.Interpolation])
		})
	}
}

func TestInterpolatedStringLiteral(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.InterpolatedStringLiteral
		wantErr bool
	}{
		{
			name:  "single interpolation",
			input: `"Hello \(name)!"`,
			want: &ast.InterpolatedStringLiteral{
				Parts: []ast.InterpolatedStringPart{
					ast.NewStringLiteral("Hello ", "Hello ", nil, 0, 6),
					&ast.Interpolation{Expression: ast.NewIdentifier("name", nil, 0, 4)},
					ast.NewStringLiteral("!", "!", nil, 0, 1),
				},
			},
		},
		{
			name:  "multiple interpolations",
			input: `"\(first) \(last)"`,
			want: &ast.InterpolatedStringLiteral{
				Parts: []ast.InterpolatedStringPart{
					&ast.Interpolation{Expression: ast.NewIdentifier("first", nil, 0, 5)},
					ast.NewStringLiteral(" ", " ", nil, 0, 1),
					&ast.Interpolation{Expression: ast.NewIdentifier("last", nil, 0, 4)},
				},
			},
		},
		{
			name:  "expression interpolation",
			input: `"Sum: \(a + b)"`,
			want: &ast.InterpolatedStringLiteral{
				Parts: []ast.InterpolatedStringPart{
					ast.NewStringLiteral("Sum: ", "Sum: ", nil, 0, 5),
					&ast.Interpolation{
						Expression: ast.NewAddSubExpression(
							ast.NewIdentifier("a", nil, 0, 1),
							ast.OpAdd,
							ast.NewIdentifier("b", nil, 0, 1),
						),
					},
				},
			},
		},
		{
			name:  "adjacent interpolations",
			input: `"\(a)\(b)"`,
			want: &ast.InterpolatedStringLiteral{
				Parts: []ast.InterpolatedStringPart{
					&ast.Interpolation{Expression: ast.NewIdentifier("a", nil, 0, 1)},
					&ast.Interpolation{Expression: ast.NewIdentifier("b", nil, 0, 1)},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"InterpolatedStringLiteral", InterpolatedStringLiteral, StringerCheck[*ast.InterpolatedStringLiteral])
		})
	}
}
