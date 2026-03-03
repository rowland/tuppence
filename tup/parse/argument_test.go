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
