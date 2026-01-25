package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/source"
	"github.com/rowland/tuppence/tup/tok"
)

func TestDecimalLiteral(t *testing.T) {
	tests := []struct {
		input   string
		want    *ast.DecimalLiteral
		wantErr bool
	}{
		// Single digits
		{input: "0", want: ast.NewDecimalLiteral("0", 0, nil, 0, 1)}, // zero
		{input: "1", want: ast.NewDecimalLiteral("1", 1, nil, 0, 1)}, // one
		{input: "2", want: ast.NewDecimalLiteral("2", 2, nil, 0, 1)}, // two
		{input: "3", want: ast.NewDecimalLiteral("3", 3, nil, 0, 1)}, // three
		{input: "4", want: ast.NewDecimalLiteral("4", 4, nil, 0, 1)}, // four
		{input: "5", want: ast.NewDecimalLiteral("5", 5, nil, 0, 1)}, // five
		{input: "6", want: ast.NewDecimalLiteral("6", 6, nil, 0, 1)}, // six
		{input: "7", want: ast.NewDecimalLiteral("7", 7, nil, 0, 1)}, // seven
		{input: "8", want: ast.NewDecimalLiteral("8", 8, nil, 0, 1)}, // eight
		{input: "9", want: ast.NewDecimalLiteral("9", 9, nil, 0, 1)}, // nine

		// Leading zeros and underscores
		{input: "0_0", want: ast.NewDecimalLiteral("0_0", 0, nil, 0, 3)},
		{input: "0001", want: ast.NewDecimalLiteral("0001", 1, nil, 0, 4)},

		// Complex numbers
		{input: "01234567890", want: ast.NewDecimalLiteral("01234567890", 1234567890, nil, 0, 11)},
		{input: "012_345_6789_0", want: ast.NewDecimalLiteral("012_345_6789_0", 1234567890, nil, 0, 14)},
		{input: "0_1_2_3_4_5_6_7_8_9_0", want: ast.NewDecimalLiteral("0_1_2_3_4_5_6_7_8_9_0", 1234567890, nil, 0, 21)},

		// Invalid characters in numbers
		{input: "123a", want: ast.NewDecimalLiteral("123a", 0, nil, 0, 4), wantErr: true},
		{input: "123A", want: ast.NewDecimalLiteral("123A", 0, nil, 0, 4), wantErr: true},
		{input: "12a34", want: ast.NewDecimalLiteral("12a34", 0, nil, 0, 5), wantErr: true},

		// Identifiers that look like numbers
		{input: "_123", want: ast.NewDecimalLiteral("_123", 123, nil, 0, 4), wantErr: true},
		{input: "__123", want: ast.NewDecimalLiteral("__123", 123, nil, 0, 5), wantErr: true},

		// Valid number followed by underscore
		{input: "123_", want: ast.NewDecimalLiteral("123_", 123, nil, 0, 4)},
		{input: "123__", want: ast.NewDecimalLiteral("123__", 123, nil, 0, 5)},

		// Sequence cases
		{input: "123,", want: ast.NewDecimalLiteral("123", 123, nil, 0, 3)},   // Should parse as separate tokens
		{input: "123_,", want: ast.NewDecimalLiteral("123_", 123, nil, 0, 4)}, // Should parse as separate tokens
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			source := source.NewSource([]byte(test.input), "test.tup")
			tokens, err := tok.Tokenize(source.Contents, source.Filename)
			if err != nil {
				t.Errorf("Tokenize(%q) = %v", test.input, err)
			}
			got, _, err := DecimalLiteral(tokens)
			if test.wantErr {
				if err == nil {
					t.Errorf("DecimalLiteral(%q) = %v, want error", test.input, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("DecimalLiteral(%q) = %v", test.input, err)
			}
			if got.Value != test.want.Value {
				t.Errorf("DecimalLiteral(%q).Value = %v, want %v", test.input, got.Value, test.want.Value)
			}
			if got.IntegerValue != test.want.IntegerValue {
				t.Errorf("DecimalLiteral(%q).IntegerValue = %v, want %v", test.input, got.IntegerValue, test.want.IntegerValue)
			}
			if got.StartOffset != test.want.StartOffset {
				t.Errorf("DecimalLiteral(%q).StartOffset = %v, want %v", test.input, got.StartOffset, test.want.StartOffset)
			}
			if got.Length != test.want.Length {
				t.Errorf("DecimalLiteral(%q).Length = %v, want %v", test.input, got.Length, test.want.Length)
			}
		})
	}
}
