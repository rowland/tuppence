package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

func TestPostfixTails(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		parse   func([]tok.Token) (ast.Expression, []tok.Token, error)
		want    ast.Expression
		wantErr bool
	}{
		{
			name:  "function call tail",
			input: "(x)",
			parse: func(tokens []tok.Token) (ast.Expression, []tok.Token, error) {
				return functionCallTail(ast.NewIdentifier("callback", nil, 0, 8), tokens)
			},
			want: ast.NewFunctionCall(
				ast.NewIdentifier("callback", nil, 0, 8),
				nil,
				ast.NewFunctionArguments(
					ast.NewArguments([]*ast.Argument{
						ast.NewArgument(ast.NewIdentifier("x", nil, 0, 1), false),
					}),
					nil,
					false,
				),
				nil,
			),
		},
		{
			name:  "member access tail",
			input: ".name",
			parse: func(tokens []tok.Token) (ast.Expression, []tok.Token, error) {
				return memberAccessTail(ast.NewIdentifier("user", nil, 0, 4), tokens)
			},
			want: ast.NewMemberAccess(
				ast.NewIdentifier("user", nil, 0, 4),
				ast.NewIdentifier("name", nil, 0, 4),
			),
		},
		{
			name:  "ordinal member access tail",
			input: ".0",
			parse: func(tokens []tok.Token) (ast.Expression, []tok.Token, error) {
				return memberAccessTail(ast.NewIdentifier("pair", nil, 0, 4), tokens)
			},
			want: ast.NewMemberAccess(
				ast.NewIdentifier("pair", nil, 0, 4),
				ast.NewDecimalLiteral("0", 0, nil, 0, 1),
			),
		},
		{
			name:  "tuple update tail",
			input: ".(name: \"Brent\")",
			parse: func(tokens []tok.Token) (ast.Expression, []tok.Token, error) {
				return tupleUpdateTail(ast.NewIdentifier("user", nil, 0, 4), tokens)
			},
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
			name:  "safe indexed access tail",
			input: "[1]!",
			parse: func(tokens []tok.Token) (ast.Expression, []tok.Token, error) {
				return safeIndexedAccessTail(ast.NewIdentifier("entries", nil, 0, 7), tokens)
			},
			want: ast.NewSafeIndexedAccess(
				ast.NewIdentifier("entries", nil, 0, 7),
				ast.NewDecimalLiteral("1", 1, nil, 0, 0),
			),
		},
		{
			name:  "indexed access tail",
			input: "[idx]",
			parse: func(tokens []tok.Token) (ast.Expression, []tok.Token, error) {
				return indexedAccessTail(ast.NewIdentifier("entries", nil, 0, 7), tokens)
			},
			want: ast.NewIndexedAccess(
				ast.NewIdentifier("entries", nil, 0, 7),
				ast.NewIdentifier("idx", nil, 0, 3),
			),
		},
		{
			name:  "malformed member access tail",
			input: ".",
			parse: func(tokens []tok.Token) (ast.Expression, []tok.Token, error) {
				return memberAccessTail(ast.NewIdentifier("user", nil, 0, 4), tokens)
			},
			wantErr: true,
		},
		{
			name:  "malformed indexed access tail",
			input: "[]",
			parse: func(tokens []tok.Token) (ast.Expression, []tok.Token, error) {
				return indexedAccessTail(ast.NewIdentifier("entries", nil, 0, 7), tokens)
			},
			wantErr: true,
		},
		{
			name:  "malformed tuple update tail",
			input: ".(1, 2)",
			parse: func(tokens []tok.Token) (ast.Expression, []tok.Token, error) {
				return tupleUpdateTail(ast.NewIdentifier("user", nil, 0, 4), tokens)
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokens, err := tok.Tokenize([]byte(test.input), "test.tup")
			if err != nil {
				t.Fatalf("Tokenize(%q): %v", test.input, err)
			}

			got, remainder, err := test.parse(tokens)
			if test.wantErr {
				if err == nil {
					t.Fatalf("parse(%q): want error, got nil", test.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("parse(%q): got error %v, want nil", test.input, err)
			}

			remainder = skipTrivia(remainder)
			if len(remainder) == 1 && remainder[0].Type == tok.TokEOF {
				remainder = nil
			}
			if len(remainder) != 0 {
				t.Fatalf("parse(%q) left trailing tokens: %v", test.input, tok.Types(remainder))
			}

			if got.String() != test.want.String() {
				t.Fatalf("parse(%q) = %v, want %v", test.input, got, test.want)
			}
		})
	}
}
