package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
)

func TestDecimalLiteral(t *testing.T) {
	tests := []struct {
		input   string
		want    *ast.IntegerLiteral
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
		{input: "123a", want: nil, wantErr: true},
		{input: "123A", want: nil, wantErr: true},
		{input: "12a34", want: nil, wantErr: true},

		// Identifiers that look like numbers
		{input: "_123", want: nil, wantErr: true},
		{input: "__123", want: nil, wantErr: true},

		// Valid number followed by underscore
		{input: "123_", want: ast.NewDecimalLiteral("123_", 123, nil, 0, 4)},
		{input: "123__", want: ast.NewDecimalLiteral("123__", 123, nil, 0, 5)},

		// Sequence cases
		{input: "123,", want: ast.NewDecimalLiteral("123", 123, nil, 0, 3)},   // Should parse as separate tokens
		{input: "123_,", want: ast.NewDecimalLiteral("123_", 123, nil, 0, 4)}, // Should parse as separate tokens
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			RunParseTest(t, test.input, test.input, test.want, test.wantErr, "DecimalLiteral", DecimalLiteral,
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
