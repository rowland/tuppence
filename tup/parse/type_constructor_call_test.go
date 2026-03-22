package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
)

func TestTypeConstructorCall(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.TypeConstructorCall
		wantErr bool
	}{
		{
			name:    "no constructor invocation",
			input:   "Person",
			wantErr: true,
		},
		{
			name:  "plain constructor call",
			input: `Person(name: "Brent")`,
			want: ast.NewTypeConstructorCall(
				ast.NewTypeReference(nil, ast.NewTypeIdentifier("Person", nil, 0, 6), nil, 0, 6),
				nil,
				ast.NewFunctionArguments(
					nil,
					ast.NewLabeledArguments([]*ast.LabeledArgument{
						ast.NewLabeledArgument(
							ast.NewIdentifier("name", nil, 0, 4),
							ast.NewArgument(ast.NewStringLiteral(`"Brent"`, "Brent", nil, 0, 7), false),
						),
					}),
					false,
				),
				nil,
			),
		},
		{
			name:  "constructor call with selector types and function block",
			input: `Reader[[]Byte](handle) { |bytes| bytes }`,
			want: ast.NewTypeConstructorCall(
				ast.NewTypeReference(nil, ast.NewTypeIdentifier("Reader", nil, 0, 6), nil, 0, 6),
				ast.NewFunctionParameterTypes([]ast.FunctionParameterType{
					ast.NewDynamicArrayType(
						ast.NewTypeReference(nil, ast.NewTypeIdentifier("Byte", nil, 0, 4), nil, 0, 4),
					),
				}),
				ast.NewFunctionArguments(
					ast.NewArguments([]*ast.Argument{
						ast.NewArgument(ast.NewIdentifier("handle", nil, 0, 6), false),
					}),
					nil,
					false,
				),
				ast.NewFunctionBlock(
					ast.NewBlockParameters(
						ast.NewOrdinalAssignmentLHS([]*ast.Identifier{
							ast.NewIdentifier("bytes", nil, 0, 5),
						}, nil),
					),
					ast.NewBlockBody(nil, ast.NewIdentifier("bytes", nil, 0, 5)),
				),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"TypeConstructorCall", TypeConstructorCall, StringerCheck[*ast.TypeConstructorCall])
		})
	}
}
