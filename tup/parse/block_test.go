package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
)

func TestBlock(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.Block
		wantErr bool
	}{
		{
			name:    "empty block",
			input:   "{}",
			want:    nil,
			wantErr: true,
		},
		{
			name:  "block with single identifier",
			input: "{ x }",
			want: ast.NewBlock(
				ast.NewBlockBody(
					[]ast.Statement{},
					ast.NewIdentifier("x", nil, 0, 1),
				),
			),
		},
		{
			name:  "block with it expression",
			input: "{ it }",
			want: ast.NewBlock(
				ast.NewBlockBody(
					[]ast.Statement{},
					ast.NewItExpression(nil, 0, 2),
				),
			),
		},
		{
			name:  "block with assignment",
			input: "{ x = 1; x + 1 }",
			want: ast.NewBlock(
				ast.NewBlockBody(
					[]ast.Statement{
						ast.NewAssignment(
							ast.NewOrdinalAssignmentLHS([]*ast.Identifier{ast.NewIdentifier("x", nil, 0, 1)}, nil),
							false,
							ast.NewDecimalLiteral("1", 1, nil, 0, 1),
						),
					},
					ast.NewAddSubExpression(ast.NewIdentifier("x", nil, 0, 1), ast.OpAdd, ast.NewDecimalLiteral("1", 1, nil, 0, 1)),
				),
			),
		},
		{
			name:  "block with multiple assignments",
			input: "{ x = 1; y = 2; y + 1 }",
			want: ast.NewBlock(
				ast.NewBlockBody(
					[]ast.Statement{
						ast.NewAssignment(
							ast.NewOrdinalAssignmentLHS([]*ast.Identifier{ast.NewIdentifier("x", nil, 0, 1)}, nil),
							false,
							ast.NewDecimalLiteral("1", 1, nil, 0, 1),
						),
						ast.NewAssignment(
							ast.NewOrdinalAssignmentLHS([]*ast.Identifier{ast.NewIdentifier("y", nil, 0, 1)}, nil),
							false,
							ast.NewDecimalLiteral("2", 2, nil, 0, 1),
						),
					},
					ast.NewAddSubExpression(ast.NewIdentifier("y", nil, 0, 1), ast.OpAdd, ast.NewDecimalLiteral("1", 1, nil, 0, 1)),
				),
			),
		},
		{
			name:  "block with bare break before final expression",
			input: "{ break\nx }",
			want: ast.NewBlock(
				ast.NewBlockBody(
					[]ast.Statement{
						ast.NewBreakExpression(nil),
					},
					ast.NewIdentifier("x", nil, 0, 1),
				),
			),
		},
		{
			name:    "block ending with assignment",
			input:   "{ x = 1 }",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "block ending with multiple assignments",
			input:   "{ x = 1; y = 2 }",
			want:    nil,
			wantErr: true,
		},
		{
			name:  "block with expression statement before final expression",
			input: "{ foo(1); bar(2) }",
			want: ast.NewBlock(
				ast.NewBlockBody(
					[]ast.Statement{
						ast.NewFunctionCall(
							ast.NewFunctionIdentifier("foo", nil, 0, 3),
							nil,
							ast.NewFunctionArguments(
								ast.NewArguments([]*ast.Argument{
									ast.NewArgument(ast.NewDecimalLiteral("1", 1, nil, 0, 1), false),
								}),
								nil,
								false,
							),
							nil,
						),
					},
					ast.NewFunctionCall(
						ast.NewFunctionIdentifier("bar", nil, 0, 3),
						nil,
						ast.NewFunctionArguments(
							ast.NewArguments([]*ast.Argument{
								ast.NewArgument(ast.NewDecimalLiteral("2", 2, nil, 0, 1), false),
							}),
							nil,
							false,
						),
						nil,
					),
				),
			),
		},
		{
			name:  "block with expression statement separated by newline",
			input: "{\n  foo(1)\n  bar(2)\n}",
			want: ast.NewBlock(
				ast.NewBlockBody(
					[]ast.Statement{
						ast.NewFunctionCall(
							ast.NewFunctionIdentifier("foo", nil, 0, 3),
							nil,
							ast.NewFunctionArguments(
								ast.NewArguments([]*ast.Argument{
									ast.NewArgument(ast.NewDecimalLiteral("1", 1, nil, 0, 1), false),
								}),
								nil,
								false,
							),
							nil,
						),
					},
					ast.NewFunctionCall(
						ast.NewFunctionIdentifier("bar", nil, 0, 3),
						nil,
						ast.NewFunctionArguments(
							ast.NewArguments([]*ast.Argument{
								ast.NewArgument(ast.NewDecimalLiteral("2", 2, nil, 0, 1), false),
							}),
							nil,
							false,
						),
						nil,
					),
				),
			),
		},
		{
			name:  "block with if expression statement before final expression",
			input: "{ if ready { value }; fallback }",
			want: ast.NewBlock(
				ast.NewBlockBody(
					[]ast.Statement{
						ast.NewIfExpression(
							[]ast.Node{
								ast.NewIdentifier("ready", nil, 0, 5),
							},
							[]*ast.Block{
								ast.NewBlock(
									ast.NewBlockBody(
										[]ast.Statement{},
										ast.NewIdentifier("value", nil, 0, 5),
									),
								),
							},
							false,
						),
					},
					ast.NewIdentifier("fallback", nil, 0, 8),
				),
			),
		},
		{
			name:  "block with if expression statement before final expression separated by newline",
			input: "{\n  if ready { value }\n  fallback\n}",
			want: ast.NewBlock(
				ast.NewBlockBody(
					[]ast.Statement{
						ast.NewIfExpression(
							[]ast.Node{
								ast.NewIdentifier("ready", nil, 0, 5),
							},
							[]*ast.Block{
								ast.NewBlock(
									ast.NewBlockBody(
										[]ast.Statement{},
										ast.NewIdentifier("value", nil, 0, 5),
									),
								),
							},
							false,
						),
					},
					ast.NewIdentifier("fallback", nil, 0, 8),
				),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"Block", Block, StringerCheck[*ast.Block])
		})
	}
}
