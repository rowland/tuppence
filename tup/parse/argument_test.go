package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/source"
	"github.com/rowland/tuppence/tup/tok"
)

func TestArgument(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.Argument
		wantErr bool
	}{
		{
			name:    "spread operator without expression",
			input:   "...",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "expression",
			input:   "x",
			want:    ast.NewArgument(ast.NewIdentifier("x", nil, 0, 1), false),
			wantErr: false,
		},
		{
			name:    "spread expression",
			input:   "...x",
			want:    ast.NewArgument(ast.NewIdentifier("x", nil, 0, 1), true),
			wantErr: false,
		},
		{
			name:    "addition expression",
			input:   "x + y",
			want:    ast.NewArgument(ast.NewAddSubExpression(ast.NewIdentifier("x", nil, 0, 1), ast.OpAdd, ast.NewIdentifier("y", nil, 0, 1)), false),
			wantErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			source := source.NewSource([]byte(test.input), "test.tup")
			tokens, err := tok.Tokenize(source.Contents, source.Filename)
			if err != nil {
				t.Errorf("Tokenize(%q) = %v", test.input, err)
			}
			arg, _, err := Argument(tokens)
			if test.wantErr {
				if err == nil {
					t.Fatalf("Argument() = nil, want error")
				}
				return
			}
			if !test.wantErr && err != nil {
				t.Fatalf("Argument() error = %v, want nil", err)
			}
			if arg.Expr.String() != test.want.Expr.String() {
				t.Errorf("Argument() = %v, want %v", arg, test.want)
			}
		})
	}
}

func TestArguments(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.Arguments
		wantErr bool
	}{
		{
			name:    "empty",
			input:   "",
			want:    ast.NewArguments([]*ast.Argument{}),
			wantErr: true,
		},
		{
			name:    "single argument",
			input:   "x",
			want:    ast.NewArguments([]*ast.Argument{ast.NewArgument(ast.NewIdentifier("x", nil, 0, 1), false)}),
			wantErr: false,
		},
		{
			name:  "multiple arguments",
			input: "x, y",
			want: ast.NewArguments([]*ast.Argument{
				ast.NewArgument(ast.NewIdentifier("x", nil, 0, 1), false),
				ast.NewArgument(ast.NewIdentifier("y", nil, 0, 1), false),
			}),
			wantErr: false,
		},
		{
			name:    "spread operator with expression",
			input:   "...x",
			want:    ast.NewArguments([]*ast.Argument{ast.NewArgument(ast.NewIdentifier("x", nil, 0, 1), true)}),
			wantErr: false,
		},
		{
			name:  "mixed arguments",
			input: "x, ...y",
			want: ast.NewArguments([]*ast.Argument{
				ast.NewArgument(ast.NewIdentifier("x", nil, 0, 1), false),
				ast.NewArgument(ast.NewIdentifier("y", nil, 0, 1), true),
			}),
			wantErr: false,
		},
		{
			name:  "multiple arguments with spread operator",
			input: "x, ...y",
			want: ast.NewArguments([]*ast.Argument{
				ast.NewArgument(ast.NewIdentifier("x", nil, 0, 1), false),
				ast.NewArgument(ast.NewIdentifier("y", nil, 0, 1), true),
			}),
			wantErr: false,
		},
		{
			name:    "multiple arguments with spread operator missing expression quits after first argument",
			input:   "x, ...",
			want:    ast.NewArguments([]*ast.Argument{ast.NewArgument(ast.NewIdentifier("x", nil, 0, 1), false)}),
			wantErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			source := source.NewSource([]byte(test.input), "test.tup")
			tokens, err := tok.Tokenize(source.Contents, source.Filename)
			if err != nil {
				t.Errorf("Tokenize(%q) = %v", test.input, err)
				return
			}
			args, _, err := Arguments(tokens)
			if test.wantErr {
				if err == nil {
					t.Fatalf("Arguments() = nil, want error")
				}
				return
			}
			if !test.wantErr && err != nil {
				t.Fatalf("Arguments() error = %v, want nil", err)
			}
			if args.String() != test.want.String() {
				t.Errorf("Arguments() = %v, want %v", args, test.want)
			}
		})
	}
}

func TestLabeledArgument(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.LabeledArgument
		wantErr bool
	}{
		{
			name:    "empty",
			input:   "",
			want:    nil,
			wantErr: true,
		},
		{
			name:  "labeled argument",
			input: "x: y",
			want: ast.NewLabeledArgument(
				ast.NewIdentifier("x", nil, 0, 1),
				ast.NewArgument(ast.NewIdentifier("y", nil, 0, 1), false),
			),
			wantErr: false,
		},
		{
			name:  "labeled argument with number",
			input: "x: 1",
			want: ast.NewLabeledArgument(
				ast.NewIdentifier("x", nil, 0, 1),
				ast.NewArgument(ast.NewDecimalLiteral("1", 1, nil, 0, 1), false),
			),
			wantErr: false,
		},
		{
			name:  "labeled argument with spread operator",
			input: "x: ...y",
			want: ast.NewLabeledArgument(
				ast.NewIdentifier("x", nil, 0, 1),
				ast.NewArgument(ast.NewIdentifier("y", nil, 0, 1), true),
			),
			wantErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			source := source.NewSource([]byte(test.input), "test.tup")
			tokens, err := tok.Tokenize(source.Contents, source.Filename)
			if err != nil {
				t.Errorf("Tokenize(%q) = %v", test.input, err)
				return
			}
			arg, _, err := LabeledArgument(tokens)
			if test.wantErr {
				if err == nil {
					t.Fatalf("LabeledArgument(): err == nil, want error")
				}
				return
			}
			if !test.wantErr && err != nil {
				t.Fatalf("LabeledArgument(): err == %v, want nil", err)
			}
			if arg.Identifier.String() != test.want.Identifier.String() {
				t.Errorf("LabeledArgument().Identifier = %v, want %v", arg.Identifier, test.want.Identifier)
			}
			if arg.Argument.Expr.String() != test.want.Argument.Expr.String() {
				t.Errorf("LabeledArgument().Argument.Expr = %v, want %v", arg.Argument.Expr, test.want.Argument.Expr)
			}
		})
	}
}

func TestLabeledArguments(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.LabeledArguments
		wantErr bool
	}{
		{
			name:    "empty",
			input:   "",
			want:    nil,
			wantErr: true,
		},
		{
			name:  "single labeled argument",
			input: "x: y",
			want: ast.NewLabeledArguments(
				[]*ast.LabeledArgument{
					ast.NewLabeledArgument(
						ast.NewIdentifier("x", nil, 0, 1),
						ast.NewArgument(ast.NewIdentifier("y", nil, 0, 1), false),
					),
				},
			),
			wantErr: false,
		},
		{
			name:  "labeled arguments",
			input: "x: y, z: w",
			want: ast.NewLabeledArguments(
				[]*ast.LabeledArgument{
					ast.NewLabeledArgument(
						ast.NewIdentifier("x", nil, 0, 1),
						ast.NewArgument(ast.NewIdentifier("y", nil, 0, 1), false),
					),
					ast.NewLabeledArgument(
						ast.NewIdentifier("z", nil, 0, 1),
						ast.NewArgument(ast.NewIdentifier("w", nil, 0, 1), false),
					),
				},
			),
			wantErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			source := source.NewSource([]byte(test.input), "test.tup")
			tokens, err := tok.Tokenize(source.Contents, source.Filename)
			if err != nil {
				t.Errorf("Tokenize(%q) = %v", test.input, err)
				return
			}
			args, _, err := LabeledArguments(tokens)
			if test.wantErr {
				if err == nil {
					t.Fatalf("LabeledArguments() = nil, want error")
				}
				return
			}
			if !test.wantErr && err != nil {
				t.Fatalf("LabeledArguments() error = %v, want nil", err)
			}
			if args.String() != test.want.String() {
				t.Errorf("LabeledArguments() = %v, want %v", args, test.want)
			}
		})
	}
}

func TestArgumentsBody(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		wantArgs        *ast.Arguments
		wantLabeledArgs *ast.LabeledArguments
		wantErr         bool
		tokensRemaining int
	}{
		{
			name:            "empty",
			input:           "",
			wantArgs:        nil,
			wantLabeledArgs: nil,
			wantErr:         true,
			tokensRemaining: 1, // EOF token
		},
		{
			name:     "labeled arguments",
			input:    "x: y, z: w",
			wantArgs: nil,
			wantLabeledArgs: ast.NewLabeledArguments([]*ast.LabeledArgument{
				ast.NewLabeledArgument(ast.NewIdentifier("x", nil, 0, 1),
					ast.NewArgument(ast.NewIdentifier("y", nil, 0, 1), false),
				),
				ast.NewLabeledArgument(ast.NewIdentifier("z", nil, 0, 1),
					ast.NewArgument(ast.NewIdentifier("w", nil, 0, 1), false),
				),
			}),
			wantErr:         false,
			tokensRemaining: 1, // EOF token
		},
		{
			name:  "positional arguments and labeled arguments",
			input: "x, y: z",
			wantArgs: ast.NewArguments([]*ast.Argument{
				ast.NewArgument(ast.NewIdentifier("x", nil, 0, 1), false),
			}),
			wantLabeledArgs: ast.NewLabeledArguments([]*ast.LabeledArgument{
				ast.NewLabeledArgument(
					ast.NewIdentifier("y", nil, 0, 1),
					ast.NewArgument(ast.NewIdentifier("z", nil, 0, 1), false),
				),
			}),
			wantErr:         false,
			tokensRemaining: 1, // EOF token
		},
		{
			name:  "positional arguments only",
			input: "x, y",
			wantArgs: ast.NewArguments([]*ast.Argument{
				ast.NewArgument(ast.NewIdentifier("x", nil, 0, 1), false),
				ast.NewArgument(ast.NewIdentifier("y", nil, 0, 1), false),
			}),
			wantLabeledArgs: nil,
			wantErr:         false,
			tokensRemaining: 1, // EOF token
		},
		{
			name:  "positional arguments with trailing comma",
			input: "x, y,",
			wantArgs: ast.NewArguments([]*ast.Argument{
				ast.NewArgument(ast.NewIdentifier("x", nil, 0, 1), false),
				ast.NewArgument(ast.NewIdentifier("y", nil, 0, 1), false),
			}),
			wantLabeledArgs: nil,
			wantErr:         false,
			tokensRemaining: 2, // trailing comma token, then EOF token
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			source := source.NewSource([]byte(test.input), "test.tup")
			tokens, err := tok.Tokenize(source.Contents, source.Filename)
			if err != nil {
				t.Errorf("Tokenize(%q) = %v", test.input, err)
				return
			}
			args, labeledArgs, remainder, err := ArgumentsBody(tokens)
			if test.wantErr {
				if err == nil {
					t.Fatalf("ArgumentsBody(): err == nil, want error")
				}
				return
			}
			if !test.wantErr && err != nil {
				t.Fatalf("ArgumentsBody(): err == %v, want nil", err)
			}
			if test.wantArgs != nil && args.String() != test.wantArgs.String() {
				t.Errorf("ArgumentsBody(): args = %v, want %v", args, test.wantArgs)
			}
			if test.wantLabeledArgs != nil && labeledArgs.String() != test.wantLabeledArgs.String() {
				t.Errorf("ArgumentsBody() = %v, want %v", labeledArgs, test.wantLabeledArgs)
			}
			if len(remainder) != test.tokensRemaining {
				t.Errorf("ArgumentsBody(): tokensRemaining = %v, want %v (%v)", len(remainder), test.tokensRemaining, tok.Types(remainder))
			}
		})
	}
}
