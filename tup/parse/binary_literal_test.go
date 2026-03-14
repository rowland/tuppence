package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
)

func TestBinaryLiteral(t *testing.T) {
	tests := []struct {
		input   string
		want    *ast.IntegerLiteral
		wantErr bool
	}{
		// Basic cases
		{input: "0b", want: nil, wantErr: true},                                          // empty_after_prefix
		{input: "0b0", want: ast.NewBinaryLiteral("0b0", 0, nil, 0, 3)},                  // zero
		{input: "0b1", want: ast.NewBinaryLiteral("0b1", 1, nil, 0, 3)},                  // one
		{input: "0b10101100", want: ast.NewBinaryLiteral("0b10101100", 172, nil, 0, 10)}, // complex_binary

		// Invalid digits
		{input: "0b2", want: nil, wantErr: true}, // invalid_2
		{input: "0b3", want: nil, wantErr: true}, // invalid_3
		{input: "0b4", want: nil, wantErr: true}, // invalid_4
		{input: "0b5", want: nil, wantErr: true}, // invalid_5
		{input: "0b6", want: nil, wantErr: true}, // invalid_6
		{input: "0b7", want: nil, wantErr: true}, // invalid_7
		{input: "0b8", want: nil, wantErr: true}, // invalid_8
		{input: "0b9", want: nil, wantErr: true}, // invalid_9
		{input: "0ba", want: nil, wantErr: true}, // invalid_a
		{input: "0bb", want: nil, wantErr: true}, // invalid_b
		{input: "0bc", want: nil, wantErr: true}, // invalid_c
		{input: "0bd", want: nil, wantErr: true}, // invalid_d
		{input: "0be", want: nil, wantErr: true}, // invalid_e
		{input: "0bf", want: nil, wantErr: true}, // invalid_f
		{input: "0bz", want: nil, wantErr: true}, // invalid_z

		// Underscore cases
		{input: "0b_", want: nil, wantErr: true},                              // invalid_leading_underscore
		{input: "0b_0", want: nil, wantErr: true},                             // invalid_underscore_after_prefix
		{input: "0b1_", want: ast.NewBinaryLiteral("0b1_", 1, nil, 0, 4)},     // valid_trailing_underscore
		{input: "0b0__1", want: ast.NewBinaryLiteral("0b0__1", 1, nil, 0, 6)}, // valid_double_underscore
		{input: "0b0_1_", want: ast.NewBinaryLiteral("0b0_1_", 1, nil, 0, 6)}, // valid_middle_underscore

		// Other cases
		{input: "0B0", want: ast.NewBinaryLiteral("0B0", 0, nil, 0, 3), wantErr: true},   // invalid_uppercase_prefix
		{input: "0b1e", want: ast.NewBinaryLiteral("0b1e", 0, nil, 0, 4), wantErr: true}, // invalid_e_suffix
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			RunParseTest(t, test.input, test.input, test.want, test.wantErr, "BinaryLiteral", BinaryLiteral,
				func(t *testing.T, input, parserName string, got, want *ast.IntegerLiteral) {
					if got.Value != want.Value {
						t.Errorf("%s(%q).Value = %v, want %v", parserName, input, got.Value, want.Value)
					}
					if got.IntegerValue != want.IntegerValue {
						t.Errorf("%s(%q).IntegerValue = %v, want %v", parserName, input, got.IntegerValue, want.IntegerValue)
					}
					if got.StartOffset != want.StartOffset {
						t.Errorf("%s(%q).StartOffset = %v, want %v", parserName, input, got.StartOffset, want.StartOffset)
					}
					if got.Length != want.Length {
						t.Errorf("%s(%q).Length = %v, want %v", parserName, input, got.Length, want.Length)
					}
				})
		})
	}
}
