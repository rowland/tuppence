package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
)

func TestHexadecimalLiteral(t *testing.T) {
	tests := []struct {
		input   string
		want    *ast.IntegerLiteral
		wantErr bool
	}{
		// Single digits (0-9)
		{input: "0x0", want: ast.NewHexadecimalLiteral("0x0", 0, nil, 0, 3)}, // zero
		{input: "0x1", want: ast.NewHexadecimalLiteral("0x1", 1, nil, 0, 3)}, // one
		{input: "0x2", want: ast.NewHexadecimalLiteral("0x2", 2, nil, 0, 3)}, // two
		{input: "0x3", want: ast.NewHexadecimalLiteral("0x3", 3, nil, 0, 3)}, // three
		{input: "0x4", want: ast.NewHexadecimalLiteral("0x4", 4, nil, 0, 3)}, // four
		{input: "0x5", want: ast.NewHexadecimalLiteral("0x5", 5, nil, 0, 3)}, // five
		{input: "0x6", want: ast.NewHexadecimalLiteral("0x6", 6, nil, 0, 3)}, // six
		{input: "0x7", want: ast.NewHexadecimalLiteral("0x7", 7, nil, 0, 3)}, // seven
		{input: "0x8", want: ast.NewHexadecimalLiteral("0x8", 8, nil, 0, 3)}, // eight
		{input: "0x9", want: ast.NewHexadecimalLiteral("0x9", 9, nil, 0, 3)}, // nine

		// Lowercase hex letters
		{input: "0xa", want: ast.NewHexadecimalLiteral("0xa", 10, nil, 0, 3)}, // lowercase_a
		{input: "0xb", want: ast.NewHexadecimalLiteral("0xb", 11, nil, 0, 3)}, // lowercase_b
		{input: "0xc", want: ast.NewHexadecimalLiteral("0xc", 12, nil, 0, 3)}, // lowercase_c
		{input: "0xd", want: ast.NewHexadecimalLiteral("0xd", 13, nil, 0, 3)}, // lowercase_d
		{input: "0xe", want: ast.NewHexadecimalLiteral("0xe", 14, nil, 0, 3)}, // lowercase_e
		{input: "0xf", want: ast.NewHexadecimalLiteral("0xf", 15, nil, 0, 3)}, // lowercase_f

		// Uppercase hex letters
		{input: "0xA", want: ast.NewHexadecimalLiteral("0xA", 10, nil, 0, 3)}, // uppercase_a
		{input: "0xB", want: ast.NewHexadecimalLiteral("0xB", 11, nil, 0, 3)}, // uppercase_b
		{input: "0xC", want: ast.NewHexadecimalLiteral("0xC", 12, nil, 0, 3)}, // uppercase_c
		{input: "0xD", want: ast.NewHexadecimalLiteral("0xD", 13, nil, 0, 3)}, // uppercase_d
		{input: "0xE", want: ast.NewHexadecimalLiteral("0xE", 14, nil, 0, 3)}, // uppercase_e
		{input: "0xF", want: ast.NewHexadecimalLiteral("0xF", 15, nil, 0, 3)}, // uppercase_f

		// Invalid letters
		{input: "0xg", want: nil, wantErr: true}, // invalid_g
		{input: "0xG", want: nil, wantErr: true}, // invalid_G
		{input: "0xz", want: nil, wantErr: true}, // invalid_z
		{input: "0xZ", want: nil, wantErr: true}, // invalid_Z

		// Complex numbers
		{input: "0x0000", want: ast.NewHexadecimalLiteral("0x0000", 0, nil, 0, 6)},                                          // leading_zeros
		{input: "0xAA", want: ast.NewHexadecimalLiteral("0xAA", 170, nil, 0, 4)},                                            // repeated_letters
		{input: "0xFFFF", want: ast.NewHexadecimalLiteral("0xFFFF", 65535, nil, 0, 6)},                                      // all_fs
		{input: "0x0123456789ABCDEF", want: ast.NewHexadecimalLiteral("0x0123456789ABCDEF", 81985529216486895, nil, 0, 18)}, // all_hex_digits

		// Underscore cases
		{input: "0x0123_4567_89AB_CDEF", want: ast.NewHexadecimalLiteral("0x0123_4567_89AB_CDEF", 81985529216486895, nil, 0, 21), wantErr: false},                         // single_group_underscore
		{input: "0x01_23_45_67_89AB_CDE_F", want: ast.NewHexadecimalLiteral("0x01_23_45_67_89AB_CDE_F", 81985529216486895, nil, 0, 24), wantErr: false},                   // multiple_group_underscore
		{input: "0x0_1_2_3_4_5_6_7_8_9_A_B_C_D_E_F", want: ast.NewHexadecimalLiteral("0x0_1_2_3_4_5_6_7_8_9_A_B_C_D_E_F", 81985529216486895, nil, 0, 33), wantErr: false}, // max_underscores
		{input: "0x_", want: nil, wantErr: true},                                                   // invalid_leading_underscore
		{input: "0x_1", want: nil, wantErr: true},                                                  // invalid_underscore_after_prefix
		{input: "0x1_", want: ast.NewHexadecimalLiteral("0x1_", 1, nil, 0, 4), wantErr: false},     // valid_trailing_underscore
		{input: "0x0__1", want: ast.NewHexadecimalLiteral("0x0__1", 1, nil, 0, 6), wantErr: false}, // valid_double_underscore
		{input: "0x0_1_", want: ast.NewHexadecimalLiteral("0x0_1_", 1, nil, 0, 6), wantErr: false}, // valid_middle_underscore

		// Invalid prefix cases
		{input: "0X0", want: nil, wantErr: true}, // invalid_uppercase_prefix
		{input: "0x", want: nil, wantErr: true},  // empty_after_prefix

		// Sequence cases
		{input: "0x1,", want: ast.NewHexadecimalLiteral("0x1", 1, nil, 0, 3)},   // hex_then_comma
		{input: "0x1_,", want: ast.NewHexadecimalLiteral("0x1_", 1, nil, 0, 4)}, // underscore_then_comma
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			RunParseTest(t, test.input, test.input, test.want, test.wantErr, "HexadecimalLiteral", HexadecimalLiteral,
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
