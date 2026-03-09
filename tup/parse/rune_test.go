package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
)

func TestRuneLiteral(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.RuneLiteral
		wantErr bool
	}{
		// Basic valid cases
		{"ascii_letter", "'A'", ast.NewRuneLiteral("'A'", 'A', nil, 0, 3), false},
		{"ascii_digit", "'0'", ast.NewRuneLiteral("'0'", '0', nil, 0, 3), false},
		{"ascii_symbol", "'!'", ast.NewRuneLiteral("'!'", '!', nil, 0, 3), false},
		{"space", "' '", ast.NewRuneLiteral("' '", ' ', nil, 0, 3), false},
		{"tab", "'\t'", ast.NewRuneLiteral("'\t'", '\t', nil, 0, 3), false},

		// Unicode characters (multi-byte UTF-8 sequences)
		{"basic_unicode", "'é'", ast.NewRuneLiteral("'é'", 'é', nil, 0, 4), false},     // 2 bytes in UTF-8
		{"greek_letter", "'Ω'", ast.NewRuneLiteral("'Ω'", 'Ω', nil, 0, 4), false},      // 2 bytes in UTF-8
		{"emoji", "'😀'", ast.NewRuneLiteral("'😀'", '😀', nil, 0, 6), false},             // 4 bytes in UTF-8
		{"chinese_character", "'汉'", ast.NewRuneLiteral("'汉'", '汉', nil, 0, 5), false}, // 3 bytes in UTF-8
		{"arabic", "'ش'", ast.NewRuneLiteral("'ش'", 'ش', nil, 0, 4), false},            // 2 bytes in UTF-8
		{"thai", "'ก'", ast.NewRuneLiteral("'ก'", 'ก', nil, 0, 5), false},              // 3 bytes in UTF-8
		{"musical_symbol", "'𝄞'", ast.NewRuneLiteral("'𝄞'", '𝄞', nil, 0, 6), false},    // 4 bytes in UTF-8 (G clef)

		// Simple escape sequences
		{"escape_newline", `'\n'`, ast.NewRuneLiteral(`'\n'`, '\n', nil, 0, 4), false},
		{"escape_tab", `'\t'`, ast.NewRuneLiteral(`'\t'`, '\t', nil, 0, 4), false},
		{"escape_backslash", `'\\'`, ast.NewRuneLiteral(`'\\'`, '\\', nil, 0, 4), false},
		{"escape_single_quote", `'\''`, ast.NewRuneLiteral(`'\''`, '\'', nil, 0, 4), false},
		{"escape_double_quote", `'\"'`, ast.NewRuneLiteral(`'\"'`, '"', nil, 0, 4), false},
		{"escape_carriage_return", `'\r'`, ast.NewRuneLiteral(`'\r'`, '\r', nil, 0, 4), false},
		{"escape_backspace", `'\b'`, ast.NewRuneLiteral(`'\b'`, '\b', nil, 0, 4), false},
		{"escape_form_feed", `'\f'`, ast.NewRuneLiteral(`'\f'`, '\f', nil, 0, 4), false},
		{"escape_vertical_tab", `'\v'`, ast.NewRuneLiteral(`'\v'`, '\v', nil, 0, 4), false},
		{"escape_null", `'\0'`, ast.NewRuneLiteral(`'\0'`, '\x00', nil, 0, 4), false},

		// Hex escapes
		{"hex_escape_lowercase", `'\x61'`, ast.NewRuneLiteral(`'\x61'`, 'a', nil, 0, 6), false}, // 'a'
		{"hex_escape_uppercase", `'\x41'`, ast.NewRuneLiteral(`'\x41'`, 'A', nil, 0, 6), false}, // 'A'
		{"hex_escape_max", `'\xFF'`, nil, true},

		// Unicode escapes
		{"unicode_4digit", `'\u03A9'`, ast.NewRuneLiteral(`'\u03A9'`, 'Ω', nil, 0, 8), false},          // Omega
		{"unicode_8digit", `'\U0001F600'`, ast.NewRuneLiteral(`'\U0001F600'`, '😀', nil, 0, 12), false}, // Grinning face emoji
		{"unicode_mixed_case", `'\u03a9'`, ast.NewRuneLiteral(`'\u03a9'`, 'Ω', nil, 0, 8), false},      // Omega (lowercase hex)

		// Invalid cases
		{"unterminated", `'A`, nil, true},
		{"empty", `''`, nil, true},

		// These should be marked as invalid by the parser, not by the tokenizer
		// The tokenizer just produces tokens, it doesn't validate rune semantics
		{"multi_code_point", `'AB'`, nil, true},

		// These escape sequence errors are still caught by the tokenizer
		{"invalid_escape", `'\z'`, nil, true},
		{"invalid_hex_escape_short", `'\x1'`, nil, true},
		{"invalid_hex_escape_letter", `'\xGG'`, nil, true},
		{"invalid_unicode_escape_short", `'\u123'`, nil, true},
		{"invalid_unicode_escape_letter", `'\uXYZW'`, nil, true},
		{"invalid_unicode_long_escape_short", `'\U1234567'`, nil, true},
		{"newline_in_rune", "'\n'", nil, true},
		{"line_feed_in_rune", `'\n\n'`, nil, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunParseTest(t, test.name, test.input, test.want, test.wantErr, "RuneLiteral", RuneLiteral,
				func(t *testing.T, input, parserName string, got, want *ast.RuneLiteral) {
					if got.Value != want.Value {
						t.Errorf("%s(%q).Value = %v, want %v", parserName, input, got.Value, want.Value)
					}
					if got.RuneValue != want.RuneValue {
						t.Errorf("%s(%q).RuneValue = %v, want %v", parserName, input, got.RuneValue, want.RuneValue)
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
