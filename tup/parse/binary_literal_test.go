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
		{input: "0b", want: nil, wantErr: true},                                                          // empty_after_prefix
		{input: "0b0", want: ast.NewBinaryLiteral("0b0", 0, nil, 0, 3), wantErr: false},                  // zero
		{input: "0b1", want: ast.NewBinaryLiteral("0b1", 1, nil, 0, 3), wantErr: false},                  // one
		{input: "0b10101100", want: ast.NewBinaryLiteral("0b10101100", 172, nil, 0, 10), wantErr: false}, // complex_binary

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
		{input: "0b_", want: nil, wantErr: true},                                              // invalid_leading_underscore
		{input: "0b_0", want: nil, wantErr: true},                                             // invalid_underscore_after_prefix
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
				if test.want == nil && got != nil {
					t.Errorf("BinaryLiteral(%q) = %v, want nil", test.input, got)
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
