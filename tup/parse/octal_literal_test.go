package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/source"
	"github.com/rowland/tuppence/tup/tok"
)

func TestOctalLiteral(t *testing.T) {
	tests := []struct {
		input   string
		want    *ast.OctalLiteral
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
		{input: "0o8", want: ast.NewOctalLiteral("0o8", 0, nil, 0, 3), wantErr: true}, // invalid_8
		{input: "0o9", want: ast.NewOctalLiteral("0o9", 0, nil, 0, 3), wantErr: true}, // invalid_9
		{input: "0oa", want: ast.NewOctalLiteral("0oa", 0, nil, 0, 3), wantErr: true}, // invalid_a
		{input: "0ob", want: ast.NewOctalLiteral("0ob", 0, nil, 0, 3), wantErr: true}, // invalid_b
		{input: "0oc", want: ast.NewOctalLiteral("0oc", 0, nil, 0, 3), wantErr: true}, // invalid_c
		{input: "0od", want: ast.NewOctalLiteral("0od", 0, nil, 0, 3), wantErr: true}, // invalid_d
		{input: "0oe", want: ast.NewOctalLiteral("0oe", 0, nil, 0, 3), wantErr: true}, // invalid_e
		{input: "0of", want: ast.NewOctalLiteral("0of", 0, nil, 0, 3), wantErr: true}, // invalid_f
		{input: "0oz", want: ast.NewOctalLiteral("0oz", 0, nil, 0, 3), wantErr: true}, // invalid_z

		// Complex numbers
		{input: "0o01234567", want: ast.NewOctalLiteral("0o01234567", 342391, nil, 0, 10)},               // all_octal_digits
		{input: "0o0123_4567", want: ast.NewOctalLiteral("0o0123_4567", 342391, nil, 0, 11)},             // single_underscore
		{input: "0o01_23_45_67", want: ast.NewOctalLiteral("0o01_23_45_67", 342391, nil, 0, 13)},         // multiple_underscores
		{input: "0o0_1_2_3_4_5_6_7", want: ast.NewOctalLiteral("0o0_1_2_3_4_5_6_7", 342391, nil, 0, 17)}, // max_underscores

		// Invalid underscore positions
		{input: "0o_", want: ast.NewOctalLiteral("0o_", 0, nil, 0, 3), wantErr: true},        // invalid_leading_underscore
		{input: "0o_0", want: ast.NewOctalLiteral("0o_0", 0, nil, 0, 4), wantErr: true},      // invalid_underscore_after_prefix
		{input: "0o1_", want: ast.NewOctalLiteral("0o1_", 1, nil, 0, 4), wantErr: false},     // valid_trailing_underscore
		{input: "0o0__1", want: ast.NewOctalLiteral("0o0__1", 1, nil, 0, 6), wantErr: false}, // valid_double_underscore
		{input: "0o0_1_", want: ast.NewOctalLiteral("0o0_1_", 1, nil, 0, 6), wantErr: false}, // valid_middle_underscore

		// Invalid prefix cases
		{input: "0O0", want: ast.NewOctalLiteral("0O0", 0, nil, 0, 3), wantErr: true}, // invalid_uppercase_prefix
		{input: "0o", want: ast.NewOctalLiteral("0o", 0, nil, 0, 2), wantErr: true},   // empty_after_prefix

		// Invalid suffix cases
		{input: "0o1e", want: ast.NewOctalLiteral("0o1e", 0, nil, 0, 4), wantErr: true},   // invalid_e_suffix
		{input: "0o1e0", want: ast.NewOctalLiteral("0o1e0", 0, nil, 0, 5), wantErr: true}, // invalid_e_suffix_with_number

		// Sequence cases
		{input: "0o1,", want: ast.NewOctalLiteral("0o1", 1, nil, 0, 3)},   // octal_then_comma
		{input: "0o1_,", want: ast.NewOctalLiteral("0o1_", 1, nil, 0, 4)}, // underscore_then_comma
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			source := source.NewSource([]byte(test.input), "test.tup")
			tokens, err := tok.Tokenize(source.Contents, source.Filename)
			if err != nil {
				t.Errorf("Tokenize(%q) = %v", test.input, err)
			}
			got, _, err := OctalLiteral(tokens)
			if test.wantErr {
				if err == nil {
					t.Errorf("OctalLiteral(%q) = %v, want error", test.input, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("OctalLiteral(%q) = %v", test.input, err)
			}
			if got.Value != test.want.Value {
				t.Errorf("OctalLiteral(%q).Value = %v, want %v", test.input, got.Value, test.want.Value)
			}
			if got.IntegerValue != test.want.IntegerValue {
				t.Errorf("OctalLiteral(%q).IntegerValue = %v, want %v", test.input, got.IntegerValue, test.want.IntegerValue)
			}
			if got.StartOffset != test.want.StartOffset {
				t.Errorf("OctalLiteral(%q).StartOffset = %v, want %v", test.input, got.StartOffset, test.want.StartOffset)
			}
			if got.Length != test.want.Length {
				t.Errorf("OctalLiteral(%q).Length = %v, want %v", test.input, got.Length, test.want.Length)
			}
		})
	}
}
