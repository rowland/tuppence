package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
)

func TestTupleUpdateExpression(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.TupleUpdateExpression
		wantErr bool
	}{
		{
			name:  "identifier receiver with labeled tuple update",
			input: "user.(name: \"Brent\")",
			want: ast.NewTupleUpdateExpression(
				ast.NewIdentifier("user", nil, 0, 4),
				ast.NewTupleLiteral(true, []*ast.TupleMember{
					ast.NewTupleMember(
						ast.NewIdentifier("name", nil, 0, 4),
						ast.NewStringLiteral(`"Brent"`, "Brent", nil, 0, 7),
					),
				}),
			),
		},
		{
			name:  "function call receiver with multiple labeled updates",
			input: "load_user().(name: \"Brent\", active: true)",
			want: ast.NewTupleUpdateExpression(
				ast.NewFunctionCall(
					ast.NewFunctionIdentifier("load_user", nil, 0, 9),
					nil,
					ast.NewFunctionArguments(ast.NewArguments([]*ast.Argument{}), nil, false),
					nil,
				),
				ast.NewTupleLiteral(true, []*ast.TupleMember{
					ast.NewTupleMember(
						ast.NewIdentifier("name", nil, 0, 4),
						ast.NewStringLiteral(`"Brent"`, "Brent", nil, 0, 7),
					),
					ast.NewTupleMember(
						ast.NewIdentifier("active", nil, 0, 6),
						ast.NewBooleanLiteral("true", true, nil, 0, 4),
					),
				}),
			),
		},
		{
			name:  "indexed receiver",
			input: "users[0].(name: \"Brent\")",
			want: ast.NewTupleUpdateExpression(
				ast.NewIndexedAccess(
					ast.NewIdentifier("users", nil, 0, 5),
					ast.NewDecimalLiteral("0", 0, nil, 0, 0),
				),
				ast.NewTupleLiteral(true, []*ast.TupleMember{
					ast.NewTupleMember(
						ast.NewIdentifier("name", nil, 0, 4),
						ast.NewStringLiteral(`"Brent"`, "Brent", nil, 0, 7),
					),
				}),
			),
		},
		{
			name:    "unlabeled tuple update is rejected",
			input:   "user.(1, 2)",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "missing receiver",
			input:   ".(name: \"Brent\")",
			want:    nil,
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"TupleUpdateExpression", TupleUpdateExpression, StringerCheck[*ast.TupleUpdateExpression])
		})
	}
}
