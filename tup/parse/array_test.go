package parse

import (
	"slices"
	"testing"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/source"
	"github.com/rowland/tuppence/tup/tok"
)

func TestSize(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		tokenTypes []tok.TokenType
		want       ast.Size
		wantErr    bool
	}{
		{
			name:       "simple",
			input:      "1",
			tokenTypes: []tok.TokenType{tok.TokDecLit, tok.TokEOF},
			want:       ast.NewDecimalLiteral("1", 1, nil, 0, 1),
		},
		{
			name:       "identifier",
			input:      "x",
			tokenTypes: []tok.TokenType{tok.TokID, tok.TokEOF},
			want:       ast.NewIdentifier("x", nil, 0, 1),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			source := source.NewSource([]byte(test.input), "test.tup")
			tokens, err := tok.Tokenize(source.Contents, source.Filename)
			if err != nil {
				t.Errorf("Tokenize(%q) = %v", test.input, err)
			}
			if !slices.Equal(tok.Types(tokens), test.tokenTypes) {
				t.Fatalf("tokenTypes: %v, want %v", tok.Types(tokens), test.tokenTypes)
			}

			got, _, err := Size(tokens)

			if test.wantErr {
				if err == nil {
					t.Errorf("Size(%q): want error", test.input)
				}
				return
			}

			if !test.wantErr && err != nil {
				t.Errorf("Size(%q): got error %v, want nil", test.input, err)
			}

			if err != nil {
				t.Errorf("Size(%q) = %v", test.input, err)
			}

			switch got := got.(type) {
			case *ast.IntegerLiteral:
				want, ok := test.want.(*ast.IntegerLiteral)
				if !ok {
					t.Fatalf("Size(%q).want: got %T, want *ast.IntegerLiteral", test.input, test.want)
				}
				if got.Value != want.Value {
					t.Fatalf("Size(%q).Value: got %v, want %v", test.input, got.Value, want.Value)
				}
				if got.IntegerValue != want.IntegerValue {
					t.Fatalf("Size(%q).IntegerValue: got %v, want %v", test.input, got.IntegerValue, want.IntegerValue)
				}
				if got.Base != want.Base {
					t.Fatalf("Size(%q).Base: got %v, want %v", test.input, got.Base, want.Base)
				}
			case *ast.Identifier:
				want, ok := test.want.(*ast.Identifier)
				if !ok {
					t.Fatalf("Size(%q).want: got %T, want *ast.Identifier", test.input, test.want)
				}
				if got.Name != want.Name {
					t.Fatalf("Size(%q).Name: got %v, want %v", test.input, got.Name, want.Name)
				}
			default:
				t.Errorf("Size(%q): got %T, want *ast.IntegerLiteral or *ast.Identifier", test.input, got)
			}
		})
	}
}

func TestArrayLiteral(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.ArrayLiteral
		wantErr bool
	}{
		{
			name:  "empty array",
			input: "[]",
			want:  ast.NewArrayLiteral(nil, nil, nil),
		},
		{
			name:  "untyped array",
			input: "[1, 2, 3]",
			want: ast.NewArrayLiteral(nil, []ast.Expression{
				ast.NewDecimalLiteral("1", 1, nil, 0, 1),
				ast.NewDecimalLiteral("2", 2, nil, 0, 1),
				ast.NewDecimalLiteral("3", 3, nil, 0, 1),
			}, nil),
		},
		{
			name:  "untyped array with trailing comma",
			input: "[1,\n2,\n3,\n]",
			want: ast.NewArrayLiteral(nil, []ast.Expression{
				ast.NewDecimalLiteral("1", 1, nil, 0, 1),
				ast.NewDecimalLiteral("2", 2, nil, 0, 1),
				ast.NewDecimalLiteral("3", 3, nil, 0, 1),
			}, nil),
		},
		{
			name:  "typed array",
			input: "Int[1, 2]",
			want: ast.NewArrayLiteral(ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3), []ast.Expression{
				ast.NewDecimalLiteral("1", 1, nil, 0, 1),
				ast.NewDecimalLiteral("2", 2, nil, 0, 1),
			}, nil),
		},
		{
			name:  "typed array with trailing comma",
			input: "Int[1,\n2,\n]",
			want: ast.NewArrayLiteral(ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3), []ast.Expression{
				ast.NewDecimalLiteral("1", 1, nil, 0, 1),
				ast.NewDecimalLiteral("2", 2, nil, 0, 1),
			}, nil),
		},
		{
			name:  "typed empty array",
			input: "String[]",
			want:  ast.NewArrayLiteral(ast.NewTypeReference(nil, ast.NewTypeIdentifier("String", nil, 0, 6), nil, 0, 6), nil, nil),
		},
		{
			name:  "array with symbol members",
			input: "[:a, :b]",
			want: ast.NewArrayLiteral(nil, []ast.Expression{
				ast.NewSymbolLiteral(":a", nil, 0, 2),
				ast.NewSymbolLiteral(":b", nil, 0, 2),
			}, nil),
		},
		{
			name:  "typed array with symbol members",
			input: "Any[:a, :b]",
			want: ast.NewArrayLiteral(ast.NewTypeReference(nil, ast.NewTypeIdentifier("Any", nil, 0, 3), nil, 0, 3), []ast.Expression{
				ast.NewSymbolLiteral(":a", nil, 0, 2),
				ast.NewSymbolLiteral(":b", nil, 0, 2),
			}, nil),
		},
		{
			name:  "namespaced typed array",
			input: "foo.Foo[a, b, c]",
			want: ast.NewArrayLiteral(
				ast.NewTypeReference(
					[]*ast.Identifier{ast.NewIdentifier("foo", nil, 0, 3)},
					ast.NewTypeIdentifier("Foo", nil, 0, 3),
					nil, 0, 0,
				),
				[]ast.Expression{
					ast.NewIdentifier("a", nil, 0, 1),
					ast.NewIdentifier("b", nil, 0, 1),
					ast.NewIdentifier("c", nil, 0, 1),
				},
				nil,
			),
		},
		{
			name:  "fixed size array literal",
			input: "[3]Int[1, 2, 3]",
			want: ast.NewArrayLiteral(
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
			name:  "identifier sized fixed size array literal",
			input: "[n]Int[1, 2, 3]",
			want: ast.NewArrayLiteral(
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
			name:  "fixed size array block initializer",
			input: "[4]Int { |i| i + 1 }",
			want: ast.NewArrayLiteral(
				ast.NewFixedSizeArrayType(
					ast.NewTypeReference(nil, ast.NewTypeIdentifier("Int", nil, 0, 3), nil, 0, 3),
					ast.NewDecimalLiteral("4", 4, nil, 0, 1),
				),
				nil,
				ast.NewFunctionBlock(
					ast.NewBlockParameters(
						ast.NewOrdinalAssignmentLHS([]*ast.Identifier{
							ast.NewIdentifier("i", nil, 0, 1),
						}, nil),
					),
					ast.NewBlockBody(
						[]ast.Statement{},
						ast.NewAddSubExpression(
							ast.NewIdentifier("i", nil, 0, 1),
							ast.OpAdd,
							ast.NewDecimalLiteral("1", 1, nil, 0, 1),
						),
					),
				),
			),
		},
		{
			name:  "named fixed size array initializer with elements",
			input: "IPv4Address[1, 2, 3, 4]",
			want: ast.NewArrayLiteral(
				ast.NewTypeReference(nil, ast.NewTypeIdentifier("IPv4Address", nil, 0, 11), nil, 0, 11),
				[]ast.Expression{
					ast.NewDecimalLiteral("1", 1, nil, 0, 1),
					ast.NewDecimalLiteral("2", 2, nil, 0, 1),
					ast.NewDecimalLiteral("3", 3, nil, 0, 1),
					ast.NewDecimalLiteral("4", 4, nil, 0, 1),
				},
				nil,
			),
		},
		{
			name:  "named fixed size array initializer with function block",
			input: "IPv4Address { |i| i + 1 }",
			want: ast.NewArrayLiteral(
				ast.NewTypeReference(nil, ast.NewTypeIdentifier("IPv4Address", nil, 0, 11), nil, 0, 11),
				nil,
				ast.NewFunctionBlock(
					ast.NewBlockParameters(
						ast.NewOrdinalAssignmentLHS([]*ast.Identifier{
							ast.NewIdentifier("i", nil, 0, 1),
						}, nil),
					),
					ast.NewBlockBody(
						[]ast.Statement{},
						ast.NewAddSubExpression(
							ast.NewIdentifier("i", nil, 0, 1),
							ast.OpAdd,
							ast.NewDecimalLiteral("1", 1, nil, 0, 1),
						),
					),
				),
			),
		},
		{
			name:  "namespaced fixed size array initializer with elements",
			input: "ipv4.Address[1, 2, 3, 4]",
			want: ast.NewArrayLiteral(
				ast.NewTypeReference(
					[]*ast.Identifier{
						ast.NewIdentifier("ipv4", nil, 0, 4),
					},
					ast.NewTypeIdentifier("Address", nil, 0, 7),
					nil, 0, 0,
				),
				[]ast.Expression{
					ast.NewDecimalLiteral("1", 1, nil, 0, 1),
					ast.NewDecimalLiteral("2", 2, nil, 0, 1),
					ast.NewDecimalLiteral("3", 3, nil, 0, 1),
					ast.NewDecimalLiteral("4", 4, nil, 0, 1),
				},
				nil,
			),
		},
		{
			name:  "namespaced fixed size array initializer with function block",
			input: "ipv4.Address { |i| i + 1 }",
			want: ast.NewArrayLiteral(
				ast.NewTypeReference(
					[]*ast.Identifier{
						ast.NewIdentifier("ipv4", nil, 0, 4),
					},
					ast.NewTypeIdentifier("Address", nil, 0, 7),
					nil, 0, 0,
				),
				nil,
				ast.NewFunctionBlock(
					ast.NewBlockParameters(
						ast.NewOrdinalAssignmentLHS([]*ast.Identifier{
							ast.NewIdentifier("i", nil, 0, 1),
						}, nil),
					),
					ast.NewBlockBody(
						[]ast.Statement{},
						ast.NewAddSubExpression(
							ast.NewIdentifier("i", nil, 0, 1),
							ast.OpAdd,
							ast.NewDecimalLiteral("1", 1, nil, 0, 1),
						),
					),
				),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr,
				"ArrayLiteral", ArrayLiteral, StringerCheck[*ast.ArrayLiteral])
		})
	}
}
