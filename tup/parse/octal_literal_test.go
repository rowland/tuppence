package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
)

func TestOctalLiteral(t *testing.T) {
	tests := []struct {
		input   string
		want    *ast.IntegerLiteral
		wantErr bool
	}{
		// Single digits
		{input: "0o0", want: ast.NewOctalLiteral("0o0", 0, nil, 0, 3)}, // zero
		{input: "0o1", want: ast.NewOctalLiteral("0o1", 1, nil, 0, 3)}, // one
		{input: "0o2", want: ast.NewOctalLiteral("0o2", 2, nil, 0, 3)}, // two
		{input: "0o3", want: ast.NewOctalLiteral("0o3", 3, nil, 0, 3)}, // three
		{input: "0o4", want: ast.NewOctalLiteral("0o4", 4, nil, 0, 3)}, // four
		{input: "0o5", want: ast.NewOctalLiteral("0o5", 5, nil, 0, 3)}, // five
		{input: "0o6", want: ast.NewOctalLiteral("0o6", 6, nil, 0, 3)}, // six
		{input: "0o7", want: ast.NewOctalLiteral("0o7", 7, nil, 0, 3)}, // seven

		// Invalid digits
		{input: "0o8", want: nil, wantErr: true}, // invalid_8
		{input: "0o9", want: nil, wantErr: true}, // invalid_9
		{input: "0oa", want: nil, wantErr: true}, // invalid_a
		{input: "0ob", want: nil, wantErr: true}, // invalid_b
		{input: "0oc", want: nil, wantErr: true}, // invalid_c
		{input: "0od", want: nil, wantErr: true}, // invalid_d
		{input: "0oe", want: nil, wantErr: true}, // invalid_e
		{input: "0of", want: nil, wantErr: true}, // invalid_f
		{input: "0oz", want: nil, wantErr: true}, // invalid_z

		// Complex numbers
		{input: "0o01234567", want: ast.NewOctalLiteral("0o01234567", 342391, nil, 0, 10)},               // all_octal_digits
		{input: "0o0123_4567", want: ast.NewOctalLiteral("0o0123_4567", 342391, nil, 0, 11)},             // single_underscore
		{input: "0o01_23_45_67", want: ast.NewOctalLiteral("0o01_23_45_67", 342391, nil, 0, 13)},         // multiple_underscores
		{input: "0o0_1_2_3_4_5_6_7", want: ast.NewOctalLiteral("0o0_1_2_3_4_5_6_7", 342391, nil, 0, 17)}, // max_underscores

		// Invalid underscore positions
		{input: "0o_", want: nil, wantErr: true},                             // invalid_leading_underscore
		{input: "0o_0", want: nil, wantErr: true},                            // invalid_underscore_after_prefix
		{input: "0o1_", want: ast.NewOctalLiteral("0o1_", 1, nil, 0, 4)},     // valid_trailing_underscore
		{input: "0o0__1", want: ast.NewOctalLiteral("0o0__1", 1, nil, 0, 6)}, // valid_double_underscore
		{input: "0o0_1_", want: ast.NewOctalLiteral("0o0_1_", 1, nil, 0, 6)}, // valid_middle_underscore

		// Invalid prefix cases
		{input: "0O0", want: nil, wantErr: true}, // invalid_uppercase_prefix
		{input: "0o", want: nil, wantErr: true},  // empty_after_prefix

		// Invalid suffix cases
		{input: "0o1e", want: nil, wantErr: true},  // invalid_e_suffix
		{input: "0o1e0", want: nil, wantErr: true}, // invalid_e_suffix_with_number

		// Sequence cases
		{input: "0o1,", want: ast.NewOctalLiteral("0o1", 1, nil, 0, 3)},   // octal_then_comma
		{input: "0o1_,", want: ast.NewOctalLiteral("0o1_", 1, nil, 0, 4)}, // underscore_then_comma
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			RunParseTest(t, test.input, test.input, test.want, test.wantErr, "OctalLiteral", OctalLiteral,
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
