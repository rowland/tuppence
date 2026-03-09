package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
)

// function_call = function_identifier [ function_parameter_types ] "(" [ function_arguments ] ")" [ function_block ] .

func TestFunctionCall(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.FunctionCall
		wantErr bool
	}{
		{
			name:    "no parameters",
			input:   "foo",
			want:    nil,
			wantErr: true,
		},
		{
			name:  "single parameter",
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
			name:  "one positional and one labeled parameter",
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
			name:  "one positional and one labeled parameter with partial application",
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
			name:  "two positional parameters with type parameters",
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
			name:  "one positional parameter with one labeled parameter and partial application",
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
			name:  "one positional parameter and function block",
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
							ast.NewArgument(ast.NewDecimalLiteral("1", 1, nil, 4, 1), false),
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
							ast.NewIdentifier("it", nil, 9, 2),
							ast.OpAdd,
							ast.NewDecimalLiteral("1", 1, nil, 14, 1),
						),
					),
				),
			),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"FunctionCall", FunctionCall, StringerCheck[*ast.FunctionCall])
		})
	}
}

// function_parameter_types = "[" local_type_reference { "," local_type_reference } "]" .

func TestFunctionParameterTypes(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.FunctionParameterTypes
		wantErr bool
	}{
		{
			name:    "no parameters",
			input:   "",
			want:    nil,
			wantErr: true,
		},
		{
			name:  "one parameter",
			input: "[Int]",
			want: ast.NewFunctionParameterTypes(
				// parameters
				[]ast.LocalTypeReference{
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
				},
			),
		},
		{
			name:  "two parameters",
			input: "[Int, Int]",
			want: ast.NewFunctionParameterTypes(
				// parameters
				[]ast.LocalTypeReference{
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
				},
			),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"FunctionParameterTypes", FunctionParameterTypes, StringerCheck[*ast.FunctionParameterTypes])
		})
	}
}

// function_arguments = ( arguments_body [ partial_application ]
//                      | "*"
// 	                    )
// 	                    [ "," ] .

func TestFunctionArguments(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.FunctionArguments
		wantErr bool
	}{
		{
			name:  "no arguments",
			input: "",
			want: ast.NewFunctionArguments(
				// args
				ast.NewArguments([]*ast.Argument{}),
				// labeledArgs
				nil,
				// partialApplication
				false,
			),
			wantErr: false,
		},
		{
			name:  "one positional argument",
			input: "1",
			want: ast.NewFunctionArguments(
				// args
				ast.NewArguments([]*ast.Argument{
					ast.NewArgument(ast.NewDecimalLiteral("1", 1, nil, 0, 1), false),
				}),
				// labeledArgs
				nil,
				// partialApplication
				false,
			),
		},
		{
			name:  "one positional argument and one labeled argument",
			input: "1, x: 2",
			want: ast.NewFunctionArguments(
				// args
				ast.NewArguments([]*ast.Argument{
					ast.NewArgument(ast.NewDecimalLiteral("1", 1, nil, 0, 1), false),
				}),
				// labeledArgs
				ast.NewLabeledArguments([]*ast.LabeledArgument{
					ast.NewLabeledArgument(ast.NewIdentifier("x", nil, 0, 1), ast.NewArgument(ast.NewDecimalLiteral("2", 2, nil, 0, 1), false)),
				}),
				// partialApplication
				false,
			),
		},
		{
			name:  "one positional argument and one labeled argument with partial application",
			input: "1, x: 2, *",
			want: ast.NewFunctionArguments(
				// args
				ast.NewArguments([]*ast.Argument{
					ast.NewArgument(ast.NewDecimalLiteral("1", 1, nil, 0, 1), false),
				}),
				// labeledArgs
				ast.NewLabeledArguments([]*ast.LabeledArgument{
					ast.NewLabeledArgument(ast.NewIdentifier("x", nil, 0, 1), ast.NewArgument(ast.NewDecimalLiteral("2", 2, nil, 0, 1), false)),
				}),
				// partialApplication
				true,
			),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"FunctionArguments", FunctionArguments, StringerCheck[*ast.FunctionArguments])
		})
	}
}

// function_block = "{" [ block_parameters ] block_body "}" .

func TestFunctionBlock(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.FunctionBlock
		wantErr bool
	}{
		{
			name:    "empty",
			input:   "",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "one parameter but no block body",
			input:   "{ |x| }",
			want:    nil,
			wantErr: true,
		},
		{
			name:  "constant expression",
			input: "{ 1 }",
			want: ast.NewFunctionBlock(
				// parameters
				nil,
				// body
				ast.NewBlockBody(
					// statements
					[]ast.Statement{},
					// expression
					ast.NewDecimalLiteral("1", 1, nil, 2, 1),
				),
			),
			wantErr: false,
		},
		{
			name:  "identifier expression",
			input: "{ x }",
			want: ast.NewFunctionBlock(
				// parameters
				nil,
				// body
				ast.NewBlockBody(
					// statements
					[]ast.Statement{},
					// expression
					ast.NewIdentifier("x", nil, 2, 1),
				),
			),
			wantErr: false,
		},
		{
			name:  "add expression",
			input: "{ x + 1 }",
			want: ast.NewFunctionBlock(
				// parameters
				nil,
				// body
				ast.NewBlockBody(
					// statements
					[]ast.Statement{},
					// expression
					ast.NewAddSubExpression(
						ast.NewIdentifier("x", nil, 2, 1),
						ast.OpAdd,
						ast.NewDecimalLiteral("1", 1, nil, 6, 1),
					),
				),
			),
			wantErr: false,
		},
		// {
		// 	name:  "one parameter and block body",
		// 	input: "{ |x| x + 1 }",
		// 	want: ast.NewFunctionBlock(
		// 		// parameters
		// 		ast.NewBlockParameters([]ast.Node{
		// 			ast.NewIdentifier("x", nil, 0, 1),
		// 		}),
		// 		// body
		// 		ast.NewBlockBody(
		// 			// statements
		// 			[]ast.Statement{},
		// 			// expression
		// 			ast.NewAddSubExpression(
		// 				ast.NewIdentifier("x", nil, 0, 1),
		// 				ast.OpAdd,
		// 				ast.NewDecimalLiteral("1", 1, nil, 0, 1),
		// 			),
		// 		),
		// 	),
		// 	wantErr: false,
		// },
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"FunctionBlock", FunctionBlock, StringerCheck[*ast.FunctionBlock])
		})
	}
}
