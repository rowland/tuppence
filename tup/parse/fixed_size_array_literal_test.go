package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

func TestFixedSizeArrayLiteral(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.FixedSizeArrayLiteral
		wantErr bool
	}{
		{
			name:  "explicit elements",
			input: "[3]Int[1, 2, 3]",
			want: ast.NewFixedSizeArrayLiteral(
				ast.NewFixedSizeArrayType(
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
					ast.NewDecimalLiteral("3", 3, nil, 0, 1),
				),
				[]ast.Expression{
					ast.NewDecimalLiteral("1", 1, nil, 0, 1),
					ast.NewDecimalLiteral("2", 2, nil, 0, 1),
					ast.NewDecimalLiteral("3", 3, nil, 0, 1),
				},
				nil,
			),
		},
		{
			name:  "implicit block parameter initializer",
			input: "[8]Int { it }",
			want: ast.NewFixedSizeArrayLiteral(
				ast.NewFixedSizeArrayType(
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
					ast.NewDecimalLiteral("8", 8, nil, 0, 1),
				),
				nil,
				ast.NewFunctionBlock(
					nil,
					ast.NewBlockBody(
						[]ast.Statement{},
						ast.NewIdentifier("it", nil, 0, 2),
					),
				),
			),
		},
		{
			name:  "identifier size with explicit elements",
			input: "[n]Int[1, 2, 3]",
			want: ast.NewFixedSizeArrayLiteral(
				ast.NewFixedSizeArrayType(
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
					ast.NewIdentifier("n", nil, 0, 1),
				),
				[]ast.Expression{
					ast.NewDecimalLiteral("1", 1, nil, 0, 1),
					ast.NewDecimalLiteral("2", 2, nil, 0, 1),
					ast.NewDecimalLiteral("3", 3, nil, 0, 1),
				},
				nil,
			),
		},
		{
			name:  "explicit block parameter initializer",
			input: "[8]Int { |index| index }",
			want: ast.NewFixedSizeArrayLiteral(
				ast.NewFixedSizeArrayType(
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
					ast.NewDecimalLiteral("8", 8, nil, 0, 1),
				),
				nil,
				ast.NewFunctionBlock(
					ast.NewBlockParameters(
						ast.NewOrdinalAssignmentLHS(
							[]*ast.Identifier{
								ast.NewIdentifier("index", nil, 0, 5),
							},
							nil,
						),
					),
					ast.NewBlockBody(
						[]ast.Statement{},
						ast.NewIdentifier("index", nil, 0, 5),
					),
				),
			),
		},
		{
			name:  "identifier size with explicit block parameter initializer",
			input: "[n]Int { |index| index }",
			want: ast.NewFixedSizeArrayLiteral(
				ast.NewFixedSizeArrayType(
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
					ast.NewIdentifier("n", nil, 0, 1),
				),
				nil,
				ast.NewFunctionBlock(
					ast.NewBlockParameters(
						ast.NewOrdinalAssignmentLHS(
							[]*ast.Identifier{
								ast.NewIdentifier("index", nil, 0, 5),
							},
							nil,
						),
					),
					ast.NewBlockBody(
						[]ast.Statement{},
						ast.NewIdentifier("index", nil, 0, 5),
					),
				),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokens, err := tok.Tokenize([]byte(test.input), "test.tup")
			if err != nil {
				t.Fatalf("Tokenize(%q) error = %v", test.input, err)
			}

			got, remainder, err := FixedSizeArrayLiteral(tokens)
			if test.wantErr {
				if err == nil {
					t.Fatalf("FixedSizeArrayLiteral(%q): want error, got nil", test.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("FixedSizeArrayLiteral(%q): got error %v, want nil", test.input, err)
			}
			if got == nil {
				t.Fatalf("FixedSizeArrayLiteral(%q) = nil, want not nil", test.input)
			}
			if got.String() != test.want.String() {
				t.Fatalf("FixedSizeArrayLiteral(%q) = %q, want %q", test.input, got.String(), test.want.String())
			}
			if len(remainder) != 1 {
				t.Fatalf("FixedSizeArrayLiteral(%q) remainder = %v, want EOF only", test.input, remainder)
			}
		})
	}
}
