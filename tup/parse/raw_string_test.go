package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
)

func TestRawStringLiteral(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.RawStringLiteral
		wantErr bool
	}{
		// Basic cases
		{"empty", "``", ast.NewRawStringLiteral("``", "", nil, 0, 2), false},
		{"simple", "`abc`", ast.NewRawStringLiteral("`abc`", "abc", nil, 0, 5), false},
		{"spaces", "`   `", ast.NewRawStringLiteral("`   `", "   ", nil, 0, 5), false},
		{"tabs", "`\t\t`", ast.NewRawStringLiteral("`\t\t`", "\t\t", nil, 0, 4), false},

		// Escaped backtick cases
		{"embedded_backtick", "`abc``def`", ast.NewRawStringLiteral("`abc``def`", "abc`def", nil, 0, 10), false},
		{"multiple_escaped_backticks", "`a``b``c`", ast.NewRawStringLiteral("`a``b``c`", "a`b`c", nil, 0, 9), false},
		{"double_backtick", "`a``b`", ast.NewRawStringLiteral("`a``b`", "a`b", nil, 0, 6), false},
		{"escaped_backtick_with_space", "`a`` b`", ast.NewRawStringLiteral("`a`` b`", "a` b", nil, 0, 7), false},
		{"escaped_backtick_with_newline", "`a``\nb`", ast.NewRawStringLiteral("`a``\nb`", "a`\nb", nil, 0, 7), false},

		// Newline cases
		{"single_newline", "`abc\ndef`", ast.NewRawStringLiteral("`abc\ndef`", "abc\ndef", nil, 0, 9), false},
		{"multiple_newlines", "`first\nsecond\nthird`", ast.NewRawStringLiteral("`first\nsecond\nthird`", "first\nsecond\nthird", nil, 0, 20), false},
		{"only_newlines", "`\n\n\n`", ast.NewRawStringLiteral("`\n\n\n`", "\n\n\n", nil, 0, 5), false},
		{"newline_after_backtick", "`abc``\ndef`", ast.NewRawStringLiteral("`abc``\ndef`", "abc`\ndef", nil, 0, 11), false},
		{"newline_before_backtick", "`abc\n``def`", ast.NewRawStringLiteral("`abc\n``def`", "abc\n`def", nil, 0, 11), false},
		{"starts_with_newline", "`\nabc`", ast.NewRawStringLiteral("`\nabc`", "\nabc", nil, 0, 6), false},
		{"ends_with_newline", "`abc\n`", ast.NewRawStringLiteral("`abc\n`", "abc\n", nil, 0, 6), false},

		// Special character cases
		{"special_chars", "`!@#$%^&*()`", ast.NewRawStringLiteral("`!@#$%^&*()`", "!@#$%^&*()", nil, 0, 12), false},
		{"unicode_chars", "`αβγδε`", ast.NewRawStringLiteral("`αβγδε`", "αβγδε", nil, 0, 12), false},
		{"emoji", "`🙂🌟🎉`", ast.NewRawStringLiteral("`🙂🌟🎉`", "🙂🌟🎉", nil, 0, 14), false},
		{"escape_sequences_raw", "`\\n\\t\\r`", ast.NewRawStringLiteral("`\\n\\t\\r`", "\\n\\t\\r", nil, 0, 8), false},
		{"quotes_in_raw", "`\"single\" and 'double'`", ast.NewRawStringLiteral("`\"single\" and 'double'`", "\"single\" and 'double'", nil, 0, 23), false},
		{"brackets_braces", "`[{(<>)}]`", ast.NewRawStringLiteral("`[{(<>)}]`", "[{(<>)}]", nil, 0, 10), false},

		// Whitespace cases
		{"mixed_whitespace", "`space\ttab\nline`", ast.NewRawStringLiteral("`space\ttab\nline`", "space\ttab\nline", nil, 0, 16), false},
		{"carriage_return", "`line1\rline2`", ast.NewRawStringLiteral("`line1\rline2`", "line1\rline2", nil, 0, 13), false},
		{"all_whitespace_types", "`\n\t\r \f\v`", ast.NewRawStringLiteral("`\n\t\r \f\v`", "\n\t\r \f\v", nil, 0, 8), false},

		// Invalid cases
		{"unterminated_after_backtick", "`abc``", nil, true},
		{"unterminated", "`abc", nil, true},
		{"unterminated_with_newlines", "`first\nsecond\nthird", nil, true},
		{"unterminated_after_escape", "`abc``def", nil, true},
		{"unterminated_unicode", "`αβγ", nil, true},
		{"unterminated_emoji", "`🙂🌟", nil, true},
		{"empty_unterminated", "`", nil, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr, "RawStringLiteral", RawStringLiteral,
				func(t *testing.T, input, parserName string, got, want *ast.RawStringLiteral) {
					if got.Value != want.Value {
						t.Errorf("%s(%q).Value = %v, want %v", parserName, input, got.Value, want.Value)
					}
					if got.StringValue != want.StringValue {
						t.Errorf("%s(%q).StringValue = %v, want %v", parserName, input, got.StringValue, want.StringValue)
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
