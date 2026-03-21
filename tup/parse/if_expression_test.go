package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
)

func TestCondition(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    ast.Expression
		wantErr bool
	}{
		{
			name:  "simple condition",
			input: "ready",
			want:  ast.NewIdentifier("ready", nil, 0, 5),
		},
		{
			name:  "comparison condition",
			input: "x > 0",
			want: ast.NewRelationalComparison(
				ast.NewIdentifier("x", nil, 0, 1),
				ast.OpGt,
				ast.NewDecimalLiteral("0", 0, nil, 0, 1),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"Condition", Condition, StringerCheck[ast.Expression])
		})
	}
}

func TestElseBlock(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.ElseBlock
		wantErr bool
	}{
		{
			name:  "simple else block",
			input: "else { fallback }",
			want: ast.NewElseBlock(
				ast.NewBlock(
					ast.NewBlockBody(
						[]ast.Statement{},
						ast.NewIdentifier("fallback", nil, 0, 8),
					),
				),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"ElseBlock", ElseBlock, StringerCheck[*ast.ElseBlock])
		})
	}
}

func TestIfExpression(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.IfExpression
		wantErr bool
	}{
		{
			name:  "simple if expression",
			input: "if ready { value }",
			want: ast.NewIfExpression(
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
		{
			name:  "if expression with else",
			input: "if ready { value } else { fallback }",
			want: ast.NewIfExpression(
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
					ast.NewBlock(
						ast.NewBlockBody(
							[]ast.Statement{},
							ast.NewIdentifier("fallback", nil, 0, 8),
						),
					),
				},
				true,
			),
		},
		{
			name:  "if expression with else if and else",
			input: "if x > 0 { positive } else if x < 0 { negative } else { zero }",
			want: ast.NewIfExpression(
				[]ast.Node{
					ast.NewRelationalComparison(
						ast.NewIdentifier("x", nil, 0, 1),
						ast.OpGt,
						ast.NewDecimalLiteral("0", 0, nil, 0, 1),
					),
					ast.NewRelationalComparison(
						ast.NewIdentifier("x", nil, 0, 1),
						ast.OpLt,
						ast.NewDecimalLiteral("0", 0, nil, 0, 1),
					),
				},
				[]*ast.Block{
					ast.NewBlock(
						ast.NewBlockBody(
							[]ast.Statement{},
							ast.NewIdentifier("positive", nil, 0, 8),
						),
					),
					ast.NewBlock(
						ast.NewBlockBody(
							[]ast.Statement{},
							ast.NewIdentifier("negative", nil, 0, 8),
						),
					),
					ast.NewBlock(
						ast.NewBlockBody(
							[]ast.Statement{},
							ast.NewIdentifier("zero", nil, 0, 4),
						),
					),
				},
				true,
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"IfExpression", IfExpression, StringerCheck[*ast.IfExpression])
		})
	}
}
