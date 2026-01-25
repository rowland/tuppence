package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

func TestNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantType ast.NodeType
		wantErr  bool
	}{
		// Single digits
		{"zero", "0", ast.NodeIntegerLiteral, false},
		{"one", "1", ast.NodeIntegerLiteral, false},
		{"two", "2", ast.NodeIntegerLiteral, false},
		{"three", "3", ast.NodeIntegerLiteral, false},
		{"four", "4", ast.NodeIntegerLiteral, false},
		{"five", "5", ast.NodeIntegerLiteral, false},
		{"six", "6", ast.NodeIntegerLiteral, false},
		{"seven", "7", ast.NodeIntegerLiteral, false},
		{"eight", "8", ast.NodeIntegerLiteral, false},
		{"nine", "9", ast.NodeIntegerLiteral, false},

		// Multiple-base numbers
		{"binary", "0b1010", ast.NodeIntegerLiteral, false},
		{"hexadecimal", "0x1A", ast.NodeIntegerLiteral, false},
		{"octal", "0o12", ast.NodeIntegerLiteral, false},
		{"decimal", "123", ast.NodeIntegerLiteral, false},

		// Float literals
		{"float", "1.0", ast.NodeFloatLiteral, false},
		{"float_with_exponent", "1.0e10", ast.NodeFloatLiteral, false},
		{"float_with_exponent_and_underscore", "1.0e10", ast.NodeFloatLiteral, false},
		{"float_with_underscore", "1.0_0", ast.NodeFloatLiteral, false},
		{"float_with_underscore_and_exponent", "1.0_0e10", ast.NodeFloatLiteral, false},
		{"float_with_underscore_and_exponent_and_underscore", "1.0_0e100", ast.NodeFloatLiteral, false},

		// Leading zeros and underscores
		{"underscore_after_zero", "0_0", ast.NodeIntegerLiteral, false},
		{"leading_zeros", "0001", ast.NodeIntegerLiteral, false},
		{"all_digits", "01234567890", ast.NodeIntegerLiteral, false},
		{"grouped_by_three", "012_345_6789_0", ast.NodeIntegerLiteral, false},
		{"max_underscores", "0_1_2_3_4_5_6_7_8_9_0", ast.NodeIntegerLiteral, false},

		// Invalid characters in numbers
		{"lowercase_letter", "123a", ast.NodeIntegerLiteral, true},

		// Identifiers that look like numbers
		{"leading_underscore", "_123", ast.NodeIdentifier, true},
		{"multiple_leading_underscores", "__123", ast.NodeIdentifier, true},
		{"trailing_underscore", "123_", ast.NodeIntegerLiteral, false},
		{"trailing_multiple_underscores", "123__", ast.NodeIntegerLiteral, false},
		{"number_then_comma", "123,", ast.NodeIntegerLiteral, false},
		{"underscore_then_comma", "123_,", ast.NodeIntegerLiteral, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := tok.Tokenize([]byte(tt.input), "test.tup")
			if err != nil {
				t.Errorf("Tokenize(%q) = %v", tt.input, err)
			}
			number, _, err := Number(tokens)
			if err != nil && !tt.wantErr {
				t.Errorf("Number(%q) = %v, want nil", tt.input, err)
			}
			if err == nil && number == nil {
				t.Errorf("Number(%q) = nil, want not nil", tt.input)
			}
			if err == nil && number.NodeType().String() != ast.NodeTypes[ast.NodeType(tt.wantType)] {
				t.Errorf("Number(%q) = %v, want %v", tt.input, number.NodeType().String(), ast.NodeTypes[ast.NodeType(tt.wantType)])
			}
		})
	}
}
