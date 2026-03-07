package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/source"
	"github.com/rowland/tuppence/tup/tok"
)

// function_call = function_identifier [ function_parameter_types ] "(" [ function_arguments ] ")" [ function_block ] .

func TestFunctionCall(t *testing.T) {
	tests := []struct {
		input   string
		want    *ast.FunctionCall
		wantErr bool
	}{
		{
			input:   "foo",
			want:    nil,
			wantErr: true,
		},
		{
			input: "foo(1)",
			want: ast.NewFunctionCall(
				// function
				ast.NewFunctionIdentifier("foo", nil, 0, 3),
				// parameterTypes
				nil,
				// arguments
				ast.NewFunctionArguments(
					// args
					ast.NewArguments(
						// args
						[]*ast.Argument{
							ast.NewArgument(ast.NewDecimalLiteral("1", 1, nil, 0, 1), false),
						},
					),
					// labeledArgs
					nil,
					// partialApplication
					false,
				),
				// functionBlock
				nil,
			),
		},
		{
			input: "bar(1, 2)",
			want: ast.NewFunctionCall(
				// function
				ast.NewFunctionIdentifier("bar", nil, 0, 3),
				// parameterTypes
				nil,
				// arguments
				ast.NewFunctionArguments(
					// args
					ast.NewArguments(
						// args
						[]*ast.Argument{
							ast.NewArgument(ast.NewDecimalLiteral("1", 1, nil, 0, 1), false),
							ast.NewArgument(ast.NewDecimalLiteral("2", 2, nil, 0, 1), false),
						},
					),
					// labeledArgs
					nil,
					// partialApplication
					false,
				),
				// functionBlock
				nil,
			),
		},
		{
			input: "baz(1, x: 2)",
			want: ast.NewFunctionCall(
				// function
				ast.NewFunctionIdentifier("baz", nil, 0, 3),
				// parameterTypes
				nil,
				// arguments
				ast.NewFunctionArguments(
					// args
					ast.NewArguments(
						// args
						[]*ast.Argument{
							ast.NewArgument(ast.NewDecimalLiteral("1", 1, nil, 0, 1), false),
						},
					),
					// labeledArgs
					ast.NewLabeledArguments([]*ast.LabeledArgument{
						ast.NewLabeledArgument(ast.NewIdentifier("x", nil, 0, 1), ast.NewArgument(ast.NewDecimalLiteral("2", 2, nil, 0, 1), false)),
					}),
					// partialApplication
					false,
				),
				// functionBlock
				nil,
			),
		},
		{
			input: "boom(1, 2, *)",
			want: ast.NewFunctionCall(
				// function
				ast.NewFunctionIdentifier("boom", nil, 0, 3),
				// parameterTypes
				nil,
				// arguments
				ast.NewFunctionArguments(
					// args
					ast.NewArguments(
						// args
						[]*ast.Argument{
							ast.NewArgument(ast.NewDecimalLiteral("1", 1, nil, 0, 1), false),
							ast.NewArgument(ast.NewDecimalLiteral("2", 2, nil, 0, 1), false),
						},
					),
					// labeledArgs
					nil,
					// partialApplication
					true,
				),
				// functionBlock
				nil,
			),
		},
		{
			input: "fizz[Int, Int](1, 2)",
			want: ast.NewFunctionCall(
				// function
				ast.NewFunctionIdentifier("fizz", nil, 0, 3),
				// parameterTypes
				ast.NewFunctionParameterTypes(
					// parameters
					[]ast.LocalTypeReference{
						ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
						ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
					},
				),
				// arguments
				ast.NewFunctionArguments(
					// args
					ast.NewArguments(
						// args
						[]*ast.Argument{
							ast.NewArgument(ast.NewDecimalLiteral("1", 1, nil, 0, 1), false),
							ast.NewArgument(ast.NewDecimalLiteral("2", 2, nil, 0, 1), false),
						},
					),
					// labeledArgs
					nil,
					// partialApplication
					false,
				),
				// functionBlock
				nil,
			),
		},
		{
			input: "buzz(1, x: 2, *)",
			want: ast.NewFunctionCall(
				// function
				ast.NewFunctionIdentifier("buzz", nil, 0, 3),
				// parameterTypes
				nil,
				// arguments
				ast.NewFunctionArguments(
					// args
					ast.NewArguments(
						// args
						[]*ast.Argument{
							ast.NewArgument(ast.NewDecimalLiteral("1", 1, nil, 0, 1), false),
						},
					),
					// labeledArgs
					ast.NewLabeledArguments([]*ast.LabeledArgument{
						ast.NewLabeledArgument(ast.NewIdentifier("x", nil, 0, 1), ast.NewArgument(ast.NewDecimalLiteral("2", 2, nil, 0, 1), false)),
					}),
					// partialApplication
					true,
				),
				// functionBlock
				nil,
			),
		},
		{
			input: "foo(1) { it + 1 }",
			want: ast.NewFunctionCall(
				// function
				ast.NewFunctionIdentifier("foo", nil, 0, 3),
				// parameterTypes
				nil,
				// arguments
				ast.NewFunctionArguments(
					// args
					ast.NewArguments(
						// args
						[]*ast.Argument{
							ast.NewArgument(ast.NewDecimalLiteral("1", 1, nil, 0, 1), false),
						},
					),
					// labeledArgs
					nil,
					// partialApplication
					false,
				),
				// functionBlock
				ast.NewFunctionBlock(
					// parameters
					nil,
					// body
					ast.NewBlockBody(
						// statements
						[]ast.Statement{},
						// expression
						ast.NewAddSubExpression(
							ast.NewIdentifier("it", nil, 0, 1),
							ast.OpAdd,
							ast.NewDecimalLiteral("2", 2, nil, 0, 0),
						),
					),
				),
			),
		},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			source := source.NewSource([]byte(test.input), "test.tup")
			tokens, err := tok.Tokenize(source.Contents, source.Filename)
			if err != nil {
				t.Errorf("Tokenize(%q) = %v", test.input, err)
			}
			got, _, err := FunctionCall(tokens)
			if test.wantErr {
				if err == nil {
					t.Errorf("FunctionCall(%q): want error", test.input)
				}
				return
			}
			if !test.wantErr && err != nil {
				t.Fatalf("FunctionCall(%q): got error %v, want nil", test.input, err)
			}
			if got == nil {
				t.Fatalf("FunctionCall(%q): got nil, want %v", test.input, test.want)
			}
			if got.String() != test.want.String() {
				t.Errorf("FunctionCall(%q) = %v, want %v", test.input, got.String(), test.want.String())
			}
		})
	}
}
