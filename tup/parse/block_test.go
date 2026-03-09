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
			wantErr: false,
		},
		{
			name:  "block with it expression",
			input: "{ it }",
			want: ast.NewBlock(
				ast.NewBlockBody(
					[]ast.Statement{},
					ast.NewIdentifier("it", nil, 0, 2),
				),
			),
			wantErr: false,
		},
		// {
		// 	name:  "block with assignment",
		// 	input: "{ x = 1; x + 1 }",
		// 	want: ast.NewBlock(
		// 		ast.NewBlockBody(
		// 			[]ast.Statement{
		// 				ast.NewAssignment(
		// 					ast.NewOrdinalAssignmentLHS([]*ast.Identifier{ast.NewIdentifier("x", nil, 0, 1)}, nil),
		// 					false,
		// 					ast.NewDecimalLiteral("1", 1, nil, 0, 1),
		// 				),
		// 			},
		// 			ast.NewAddSubExpression(ast.NewIdentifier("x", nil, 0, 1), ast.OpAdd, ast.NewDecimalLiteral("1", 1, nil, 0, 1)),
		// 		),
		// 	),
		// 	wantErr: false,
		// },
		// {
		// 	name:  "block with multiple assignments",
		// 	input: "{ x = 1; y = 2; y + 1 }",
		// 	want: ast.NewBlock(
		// 		ast.NewBlockBody(
		// 			[]ast.Statement{
		// 				ast.NewAssignment(
		// 					ast.NewOrdinalAssignmentLHS([]*ast.Identifier{ast.NewIdentifier("x", nil, 0, 1)}, nil),
		// 					false,
		// 					ast.NewDecimalLiteral("1", 1, nil, 0, 1),
		// 				),
		// 				ast.NewAssignment(
		// 					ast.NewOrdinalAssignmentLHS([]*ast.Identifier{ast.NewIdentifier("y", nil, 0, 1)}, nil),
		// 					false,
		// 					ast.NewDecimalLiteral("2", 2, nil, 0, 1),
		// 				),
		// 			},
		// 			ast.NewAddSubExpression(ast.NewIdentifier("y", nil, 0, 1), ast.OpAdd, ast.NewDecimalLiteral("1", 1, nil, 0, 1)),
		// 		),
		// 	),
		// 	wantErr: false,
		// },
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.input, test.want, test.wantErr, Block)
		})
	}
}
