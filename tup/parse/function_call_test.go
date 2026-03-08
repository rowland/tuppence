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
			source := source.NewSource([]byte(test.input), "test.tup")
			tokens, err := tok.Tokenize(source.Contents, source.Filename)
			if err != nil {
				t.Errorf("Tokenize(%q) = %v", test.input, err)
			}
			got, _, err := FunctionParameterTypes(tokens)
			if test.wantErr {
				if err == nil {
					t.Errorf("FunctionParameterTypes(%q): want error", test.input)
				}
				return
			}
			if !test.wantErr && err != nil {
				t.Fatalf("FunctionParameterTypes(%q): got error %v, want nil", test.input, err)
			}
			if got == nil {
				t.Fatalf("FunctionParameterTypes(%q): got nil, want %v", test.input, test.want)
			}
			if got.String() != test.want.String() {
				t.Errorf("FunctionParameterTypes(%q) = %v, want %v", test.input, got.String(), test.want.String())
			}
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
			source := source.NewSource([]byte(test.input), "test.tup")
			tokens, err := tok.Tokenize(source.Contents, source.Filename)
			if err != nil {
				t.Errorf("Tokenize(%q) = %v", test.input, err)
			}
			got, _, err := FunctionArguments(tokens)
			if test.wantErr {
				if err == nil {
					t.Errorf("FunctionArguments(%q): want error", test.input)
				}
				return
			}
			if !test.wantErr && err != nil {
				t.Fatalf("FunctionArguments(%q): got error %v, want nil", test.input, err)
			}
			if got == nil {
				t.Fatalf("FunctionArguments(%q): got nil, want %v", test.input, test.want)
			}
			if got.String() != test.want.String() {
				t.Errorf("FunctionArguments(%q) = %v, want %v", test.input, got.String(), test.want.String())
			}
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
			source := source.NewSource([]byte(test.input), "test.tup")
			tokens, err := tok.Tokenize(source.Contents, source.Filename)
			if err != nil {
				t.Errorf("Tokenize(%q) = %v", test.input, err)
			}
			got, _, err := FunctionBlock(tokens)
			if test.wantErr {
				if err == nil {
					t.Errorf("FunctionBlock(%q): want error", test.input)
				}
				return
			}
			if !test.wantErr && err != nil {
				t.Fatalf("FunctionBlock(%q): got error %v, want nil", test.input, err)
			}
			if got == nil {
				t.Fatalf("FunctionBlock(%q): got nil, want %v", test.input, test.want)
			}
			if got.String() != test.want.String() {
				t.Errorf("FunctionBlock(%q) = %v, want %v", test.input, got.String(), test.want.String())
			}
		})
	}
}
