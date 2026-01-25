package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/source"
	"github.com/rowland/tuppence/tup/tok"
)

func TestBinaryLiteral(t *testing.T) {
	tests := []struct {
		input   string
		want    *ast.IntegerLiteral
		wantErr bool
	}{
		// Basic cases
		{input: "0b", want: ast.NewBinaryLiteral("0b", 0, nil, 0, 2), wantErr: true},                     // empty_after_prefix
		{input: "0b0", want: ast.NewBinaryLiteral("0b0", 0, nil, 0, 3), wantErr: false},                  // zero
		{input: "0b1", want: ast.NewBinaryLiteral("0b1", 1, nil, 0, 3), wantErr: false},                  // one
		{input: "0b10101100", want: ast.NewBinaryLiteral("0b10101100", 172, nil, 0, 10), wantErr: false}, // complex_binary

		// Invalid digits
		{input: "0b2", want: ast.NewBinaryLiteral("0b2", 0, nil, 0, 3), wantErr: true}, // invalid_2
		{input: "0b3", want: ast.NewBinaryLiteral("0b3", 0, nil, 0, 3), wantErr: true}, // invalid_3
		{input: "0b4", want: ast.NewBinaryLiteral("0b4", 0, nil, 0, 3), wantErr: true}, // invalid_4
		{input: "0b5", want: ast.NewBinaryLiteral("0b5", 0, nil, 0, 3), wantErr: true}, // invalid_5
		{input: "0b6", want: ast.NewBinaryLiteral("0b6", 0, nil, 0, 3), wantErr: true}, // invalid_6
		{input: "0b7", want: ast.NewBinaryLiteral("0b7", 0, nil, 0, 3), wantErr: true}, // invalid_7
		{input: "0b8", want: ast.NewBinaryLiteral("0b8", 0, nil, 0, 3), wantErr: true}, // invalid_8
		{input: "0b9", want: ast.NewBinaryLiteral("0b9", 0, nil, 0, 3), wantErr: true}, // invalid_9
		{input: "0ba", want: ast.NewBinaryLiteral("0ba", 0, nil, 0, 3), wantErr: true}, // invalid_a
		{input: "0bb", want: ast.NewBinaryLiteral("0bb", 0, nil, 0, 3), wantErr: true}, // invalid_b
		{input: "0bc", want: ast.NewBinaryLiteral("0bc", 0, nil, 0, 3), wantErr: true}, // invalid_c
		{input: "0bd", want: ast.NewBinaryLiteral("0bd", 0, nil, 0, 3), wantErr: true}, // invalid_d
		{input: "0be", want: ast.NewBinaryLiteral("0be", 0, nil, 0, 3), wantErr: true}, // invalid_e
		{input: "0bf", want: ast.NewBinaryLiteral("0bf", 0, nil, 0, 3), wantErr: true}, // invalid_f
		{input: "0bz", want: ast.NewBinaryLiteral("0bz", 0, nil, 0, 3), wantErr: true}, // invalid_z

		// Underscore cases
		{input: "0b_", want: ast.NewBinaryLiteral("0b_", 0, nil, 0, 3), wantErr: true},        // invalid_leading_underscore
		{input: "0b_0", want: ast.NewBinaryLiteral("0b_0", 0, nil, 0, 4), wantErr: true},      // invalid_underscore_after_prefix
		{input: "0b1_", want: ast.NewBinaryLiteral("0b1_", 1, nil, 0, 4), wantErr: false},     // valid_trailing_underscore
		{input: "0b0__1", want: ast.NewBinaryLiteral("0b0__1", 1, nil, 0, 6), wantErr: false}, // valid_double_underscore
		{input: "0b0_1_", want: ast.NewBinaryLiteral("0b0_1_", 1, nil, 0, 6), wantErr: false}, // valid_middle_underscore

		// Other cases
		{input: "0B0", want: ast.NewBinaryLiteral("0B0", 0, nil, 0, 3), wantErr: true},   // invalid_uppercase_prefix
		{input: "0b1e", want: ast.NewBinaryLiteral("0b1e", 0, nil, 0, 4), wantErr: true}, // invalid_e_suffix
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			source := source.NewSource([]byte(test.input), "test.tup")
			tokens, err := tok.Tokenize(source.Contents, source.Filename)
			if err != nil {
				t.Errorf("Tokenize(%q) = %v", test.input, err)
			}
			got, _, err := BinaryLiteral(tokens)
			if test.wantErr {
				if err == nil {
					t.Errorf("BinaryLiteral(%q) = %v, want error", test.input, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("BinaryLiteral(%q) = %v", test.input, err)
			}
			if got.Value != test.want.Value {
				t.Errorf("BinaryLiteral(%q).Value = %v, want %v", test.input, got.Value, test.want.Value)
			}
			if got.IntegerValue != test.want.IntegerValue {
				t.Errorf("BinaryLiteral(%q).IntegerValue = %v, want %v", test.input, got.IntegerValue, test.want.IntegerValue)
			}
			if got.StartOffset != test.want.StartOffset {
				t.Errorf("BinaryLiteral(%q).StartOffset = %v, want %v", test.input, got.StartOffset, test.want.StartOffset)
			}
			if got.Length != test.want.Length {
				t.Errorf("BinaryLiteral(%q).Length = %v, want %v", test.input, got.Length, test.want.Length)
			}
		})
	}
}
