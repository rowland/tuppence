package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/source"
	"github.com/rowland/tuppence/tup/tok"
)

func TestLiteral(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantType ast.NodeType
		wantErr  bool
	}{
		// number
		{"float literal", "1.0", ast.NodeFloatLiteral, false},
		{"integer literal", "1", ast.NodeIntegerLiteral, false},

		// boolean
		{"boolean literal true", "true", ast.NodeBooleanLiteral, false},
		{"boolean literal false", "false", ast.NodeBooleanLiteral, false},

		// string
		{"string literal", "\"hello\"", ast.NodeStringLiteral, false},

		// raw string
		{"raw string literal", "`hello`", ast.NodeRawStringLiteral, false},

		// interpolated string
		// {"interpolated string literal", "`hello ${name}`", ast.NodeInterpolatedStringLiteral, false},

		// multi line string
		// {"multi line string literal", "`hello\nworld`", ast.NodeMultiLineStringLiteral, false},

		// tuple
		// {"tuple literal", "(1, 2, 3)", ast.NodeTupleLiteral, false},

		// array
		// {"array literal", "[1, 2, 3]", ast.NodeArrayLiteral, false},

		// symbol
		// {"symbol literal", ":hello", ast.NodeSymbolLiteral, false},

		// rune
		{"rune literal", "'a'", ast.NodeRuneLiteral, false},

		// fixed size array
		// {"fixed size array literal", "[1, 2, 3]", ast.NodeFixedSizeArrayLiteral, false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			source := source.NewSource([]byte(test.input), "test.tup")
			tokens, err := tok.Tokenize(source.Contents, source.Filename)
			if err != nil {
				t.Errorf("Tokenize(%q) = %v", test.input, err)
			}
			got, _, err := Literal(tokens)
			if test.wantErr {
				if err == nil {
					t.Fatalf("Literal(%q) = %v, want error", test.input, got)
				}
				return
			} else {
				if err != nil {
					t.Fatalf("Literal(%q) = %v, want nil", test.input, err)
				}
			}
			if got == nil {
				t.Fatalf("Literal(%q) = nil, want %v", test.input, test.wantType)
			}
			if got.NodeType() != test.wantType {
				t.Fatalf("Literal(%q) = %v, want %v", test.input, got.NodeType(), test.wantType)
			}
		})
	}
}
