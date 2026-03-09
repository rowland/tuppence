package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
)

func TestFloatLiteral(t *testing.T) {
	tests := []struct {
		input   string
		want    *ast.FloatLiteral
		wantErr bool
	}{
		// Valid floats
		{input: "0.5", want: ast.NewFloatLiteral("0.5", 0.5, nil, 0, 3)},
		{input: "1.23", want: ast.NewFloatLiteral("1.23", 1.23, nil, 0, 4)},
		{input: "12.34e+5", want: ast.NewFloatLiteral("12.34e+5", 12.34e+5, nil, 0, 8)},
		{input: "12.34e-5", want: ast.NewFloatLiteral("12.34e-5", 12.34e-5, nil, 0, 8)},
		{input: "9e9", want: ast.NewFloatLiteral("9e9", 9e9, nil, 0, 3)},
		{input: "10e+10", want: ast.NewFloatLiteral("10e+10", 10e+10, nil, 0, 6)},
		{input: "10e-10", want: ast.NewFloatLiteral("10e-10", 10e-10, nil, 0, 6)},
		{input: "1_2.3_4", want: ast.NewFloatLiteral("1_2.3_4", 12.34, nil, 0, 7)},
		{input: "3.14_159", want: ast.NewFloatLiteral("3.14_159", 3.14159, nil, 0, 8)},

		// Examples with underscores
		{input: "1_2.3_4", want: ast.NewFloatLiteral("1_2.3_4", 12.34, nil, 0, 7)},
		{input: "3.14_159", want: ast.NewFloatLiteral("3.14_159", 3.14159, nil, 0, 8)},

		// Invalid floats
		{input: "1.2e", want: nil, wantErr: true},    // missing exponent digits
		{input: "1.2ez", want: nil, wantErr: true},   // non-digit exponent suffix
		{input: "1.2e++3", want: nil, wantErr: true}, // double sign
		{input: "123e", want: nil, wantErr: true},    // no exponent digits
		{input: "12.34e-", want: nil, wantErr: true}, // minus with no digits
		{input: "0. ", want: nil, wantErr: true},     // decimal not followed by a digit or another dot
		{input: "1. ", want: nil, wantErr: true},     // decimal not followed by a digit or another dot
		{input: "12. ", want: nil, wantErr: true},    // decimal not followed by a digit or another dot
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			RunParseTest(t, test.input, test.input, test.want, test.wantErr, "FloatLiteral", FloatLiteral,
				func(t *testing.T, input, parserName string, got, want *ast.FloatLiteral) {
					if got.Value != want.Value {
						t.Errorf("%s(%q).Value = %v, want %v", parserName, input, got.Value, want.Value)
					}
					if got.FloatValue != want.FloatValue {
						t.Errorf("%s(%q).FloatValue = %v, want %v", parserName, input, got.FloatValue, want.FloatValue)
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
