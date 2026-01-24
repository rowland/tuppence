package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/source"
	"github.com/rowland/tuppence/tup/tok"
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
		{input: "1.2e", want: ast.NewFloatLiteral("1.2e", 0, nil, 0, 4), wantErr: true},       // missing exponent digits
		{input: "1.2ez", want: ast.NewFloatLiteral("1.2ez", 0, nil, 0, 5), wantErr: true},     // non-digit exponent suffix
		{input: "1.2e++3", want: ast.NewFloatLiteral("1.2e++3", 0, nil, 0, 7), wantErr: true}, // double sign
		{input: "123e", want: ast.NewFloatLiteral("123e", 0, nil, 0, 4), wantErr: true},       // no exponent digits
		{input: "12.34e-", want: ast.NewFloatLiteral("12.34e-", 0, nil, 0, 7), wantErr: true}, // minus with no digits
		{input: "0. ", want: ast.NewFloatLiteral("0. ", 0, nil, 0, 2), wantErr: true},         // decimal not followed by a digit or another dot
		{input: "1. ", want: ast.NewFloatLiteral("1. ", 0, nil, 0, 2), wantErr: true},         // decimal not followed by a digit or another dot
		{input: "12. ", want: ast.NewFloatLiteral("12. ", 0, nil, 0, 3), wantErr: true},       // decimal not followed by a digit or another dot
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			source := source.NewSource([]byte(test.input), "test.tup")
			tokens, err := tok.Tokenize(source.Contents, source.Filename)
			if err != nil {
				t.Errorf("Tokenize(%q) = %v", test.input, err)
			}
			got, _, err := FloatLiteral(tokens)
			if test.wantErr {
				if err == nil {
					t.Errorf("FloatLiteral(%q) = %v, want error", test.input, got)
				}
				return
			}
			if err != nil {
				t.Errorf("FloatLiteral(%q) = %v", test.input, err)
			}
			if got.Value != test.want.Value {
				t.Errorf("FloatLiteral(%q).Value = %v, want %v", test.input, got.Value, test.want.Value)
			}
			if got.FloatValue != test.want.FloatValue {
				t.Errorf("FloatLiteral(%q).FloatValue = %v, want %v", test.input, got.FloatValue, test.want.FloatValue)
			}
			if got.StartOffset != test.want.StartOffset {
				t.Errorf("FloatLiteral(%q).StartOffset = %v, want %v", test.input, got.StartOffset, test.want.StartOffset)
			}
			if got.Length != test.want.Length {
				t.Errorf("FloatLiteral(%q).Length = %v, want %v", test.input, got.Length, test.want.Length)
			}
		})
	}
}
