package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
)

func TestFunctionCallContext(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.FunctionCallContext
		wantErr bool
	}{
		{
			name:  "bare processor",
			input: "json",
			want: ast.NewFunctionCallContext(
				ast.NewFunctionIdentifier("json", nil, 0, 4),
				nil,
			),
		},
		{
			name:  "processor with explicit arguments",
			input: "mustache(context)",
			want: ast.NewFunctionCallContext(
				ast.NewFunctionIdentifier("mustache", nil, 0, 8),
				ast.NewFunctionArguments(
					ast.NewArguments([]*ast.Argument{
						ast.NewArgument(ast.NewIdentifier("context", nil, 0, 7), false),
					}),
					nil,
					false,
				),
			),
		},
		{
			name:  "processor with labeled arguments",
			input: "html(trim: true, escape: false)",
			want: ast.NewFunctionCallContext(
				ast.NewFunctionIdentifier("html", nil, 0, 4),
				ast.NewFunctionArguments(
					nil,
					ast.NewLabeledArguments([]*ast.LabeledArgument{
						ast.NewLabeledArgument(
							ast.NewIdentifier("trim", nil, 0, 4),
							ast.NewArgument(ast.NewBooleanLiteral("true", true, nil, 0, 4), false),
						),
						ast.NewLabeledArgument(
							ast.NewIdentifier("escape", nil, 0, 6),
							ast.NewArgument(ast.NewBooleanLiteral("false", false, nil, 0, 5), false),
						),
					}),
					false,
				),
			),
		},
		{
			name:  "processor with empty argument list",
			input: "processor()",
			want: ast.NewFunctionCallContext(
				ast.NewFunctionIdentifier("processor", nil, 0, 9),
				ast.NewFunctionArguments(nil, nil, false),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"FunctionCallContext", FunctionCallContext, StringerCheck[*ast.FunctionCallContext])
		})
	}
}
