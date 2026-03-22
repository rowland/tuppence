package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
)

func TestMultiLineStringLiteral(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.MultiLineStringLiteral
		wantErr bool
	}{
		{
			name:  "plain multiline text",
			input: "```\nThis is\nsome text\n```",
			want: ast.NewMultiLineStringLiteral(contents(
				ast.NewStringLiteral("This is\nsome text\n", "This is\nsome text\n", nil, 0, 18),
			), nil),
		},
		{
			name:  "dedented multiline text",
			input: "```\n    This text\n    will have\n    leading whitespace\n    removed.\n```",
			want: ast.NewMultiLineStringLiteral(contents(
				ast.NewStringLiteral("This text\nwill have\nleading whitespace\nremoved.\n", "This text\nwill have\nleading whitespace\nremoved.\n", nil, 0, 48),
			), nil),
		},
		{
			name:  "blank first line with indentation baseline on next line",
			input: "```\n\n    Hello\n    World\n```",
			want: ast.NewMultiLineStringLiteral(contents(
				ast.NewStringLiteral("\nHello\nWorld\n", "\nHello\nWorld\n", nil, 0, 13),
			), nil),
		},
		{
			name:  "interpolation within one line",
			input: "```\nHello \\(name)!\n```",
			want: ast.NewMultiLineStringLiteral(contents(
				ast.NewStringLiteral("Hello ", "Hello ", nil, 0, 6),
				&ast.Interpolation{Expression: ast.NewIdentifier("name", nil, 0, 4)},
				ast.NewStringLiteral("!\n", "!\n", nil, 0, 2),
			), nil),
		},
		{
			name:  "multiple interpolations across lines",
			input: "```\nGirls: \\(girls.join(\", \"))\nBoys: \\(boys.join(\", \"))\n```",
			want: ast.NewMultiLineStringLiteral(contents(
				ast.NewStringLiteral("Girls: ", "Girls: ", nil, 0, 7),
				&ast.Interpolation{
					Expression: ast.NewFunctionCall(
						ast.NewMemberAccess(
							ast.NewIdentifier("girls", nil, 0, 5),
							ast.NewIdentifier("join", nil, 0, 4),
						),
						nil,
						ast.NewFunctionArguments(
							ast.NewArguments([]*ast.Argument{
								ast.NewArgument(ast.NewStringLiteral(`", "`, ", ", nil, 0, 4), false),
							}),
							nil,
							false,
						),
						nil,
					),
				},
				ast.NewStringLiteral("\nBoys: ", "\nBoys: ", nil, 0, 7),
				&ast.Interpolation{
					Expression: ast.NewFunctionCall(
						ast.NewMemberAccess(
							ast.NewIdentifier("boys", nil, 0, 4),
							ast.NewIdentifier("join", nil, 0, 4),
						),
						nil,
						ast.NewFunctionArguments(
							ast.NewArguments([]*ast.Argument{
								ast.NewArgument(ast.NewStringLiteral(`", "`, ", ", nil, 0, 4), false),
							}),
							nil,
							false,
						),
						nil,
					),
				},
				ast.NewStringLiteral("\n", "\n", nil, 0, 1),
			), nil),
		},
		{
			name:  "processor with no args",
			input: "```json\n  { \"a\": \"b\" }\n```",
			want: ast.NewMultiLineStringLiteral(contents(
				ast.NewStringLiteral("{ \"a\": \"b\" }\n", "{ \"a\": \"b\" }\n", nil, 0, 13),
			), ast.NewFunctionCallContext(
				ast.NewFunctionIdentifier("json", nil, 0, 4),
				nil,
			)),
		},
		{
			name:  "processor with args",
			input: "```mustache(context)\n  Hello, {{name}}\n```",
			want: ast.NewMultiLineStringLiteral(contents(
				ast.NewStringLiteral("Hello, {{name}}\n", "Hello, {{name}}\n", nil, 0, 16),
			), ast.NewFunctionCallContext(
				ast.NewFunctionIdentifier("mustache", nil, 0, 8),
				ast.NewFunctionArguments(
					ast.NewArguments([]*ast.Argument{
						ast.NewArgument(ast.NewIdentifier("context", nil, 0, 7), false),
					}),
					nil,
					false,
				),
			)),
		},
		{
			name:  "interpolation spanning lines joins those lines",
			input: "```\n    Value: \\(\n        a +\n        b\n    )\n```",
			want: ast.NewMultiLineStringLiteral(contents(
				ast.NewStringLiteral("Value: ", "Value: ", nil, 0, 7),
				&ast.Interpolation{
					Expression: ast.NewAddSubExpression(
						ast.NewIdentifier("a", nil, 0, 1),
						ast.OpAdd,
						ast.NewIdentifier("b", nil, 0, 1),
					),
				},
				ast.NewStringLiteral("\n", "\n", nil, 0, 1),
			), nil),
		},
		{
			name:  "preserves crlf line endings",
			input: "```\r\n    Hello\r\n    World\r\n```",
			want: ast.NewMultiLineStringLiteral(contents(
				ast.NewStringLiteral("Hello\r\nWorld\r\n", "Hello\r\nWorld\r\n", nil, 0, 14),
			), nil),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"MultiLineStringLiteral", MultiLineStringLiteral, StringerCheck[*ast.MultiLineStringLiteral])
		})
	}
}

func contents(parts ...ast.Node) *ast.InterpolatedStringLiteral {
	return &ast.InterpolatedStringLiteral{
		BaseNode: ast.BaseNode{Type: ast.NodeInterpolatedStringLiteral},
		Parts:    parts,
	}
}
