package parse

import (
	"testing"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/source"
	"github.com/rowland/tuppence/tup/tok"
)

func TestStringLiteral(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ast.StringLiteral
		wantErr bool
	}{
		// Basic cases
		{"empty", `""`, ast.NewStringLiteral(`""`, "", nil, 0, 2), false},
		{"simple", `"abc"`, ast.NewStringLiteral(`"abc"`, "abc", nil, 0, 5), false},
		{"spaces", `"   "`, ast.NewStringLiteral(`"   "`, "   ", nil, 0, 5), false},
		{"tabs", `"\t\t"`, ast.NewStringLiteral(`"\t\t"`, "\t\t", nil, 0, 6), false},

		// Simple escape sequences
		{"newline", `"\n"`, ast.NewStringLiteral(`"\n"`, "\n", nil, 0, 4), false},
		{"tab", `"\t"`, ast.NewStringLiteral(`"\t"`, "\t", nil, 0, 4), false},
		{"double_quote", `"\""`, ast.NewStringLiteral(`"\""`, "\"", nil, 0, 4), false},
		{"single_quote", `"\'"`, ast.NewStringLiteral(`"\'"`, "'", nil, 0, 4), false},
		{"backslash", `"\\"`, ast.NewStringLiteral(`"\\"`, "\\", nil, 0, 4), false},
		{"carriage_return", `"\r"`, ast.NewStringLiteral(`"\r"`, "\r", nil, 0, 4), false},
		{"backspace", `"\b"`, ast.NewStringLiteral(`"\b"`, "\b", nil, 0, 4), false},
		{"form_feed", `"\f"`, ast.NewStringLiteral(`"\f"`, "\f", nil, 0, 4), false},
		{"vertical_tab", `"\v"`, ast.NewStringLiteral(`"\v"`, "\v", nil, 0, 4), false},
		{"null", `"\0"`, ast.NewStringLiteral(`"\0"`, "\x00", nil, 0, 4), false},
		{"all_escapes", `"\n\t\"\'\\\r\b\f\v\0"`, ast.NewStringLiteral(`"\n\t\"\'\\\r\b\f\v\0"`, "\n\t\"'\\\r\b\f\v\x00", nil, 0, 22), false},

		// Unicode escapes
		{"unicode_4_digit", `"\u1234"`, ast.NewStringLiteral(`"\u1234"`, "\u1234", nil, 0, 8), false},
		{"unicode_8_digit", `"\U12345678"`, ast.NewStringLiteral(`"\U12345678"`, "", nil, 0, 12), true},
		{"unicode_max", `"\uFFFF"`, ast.NewStringLiteral(`"\uFFFF"`, "\uFFFF", nil, 0, 8), false},
		{"unicode_min", `"\u0000"`, ast.NewStringLiteral(`"\u0000"`, "\u0000", nil, 0, 8), false},
		{"unicode_multiple", `"\u1234\u5678"`, ast.NewStringLiteral(`"\u1234\u5678"`, "\u1234\u5678", nil, 0, 14), false},
		{"unicode_mixed_case", `"\uAbCd"`, ast.NewStringLiteral(`"\uAbCd"`, "\uAbCd", nil, 0, 8), false},

		// Hex escapes
		{"hex_byte", `"\xEF"`, ast.NewStringLiteral(`"\xEF"`, "\xEF", nil, 0, 6), false},
		{"hex_multiple", `"\xEF\xBB\xBF"`, ast.NewStringLiteral(`"\xEF\xBB\xBF"`, "\xEF\xBB\xBF", nil, 0, 14), false},
		{"hex_min", `"\x00"`, ast.NewStringLiteral(`"\x00"`, "\x00", nil, 0, 6), false},
		{"hex_max", `"\xFF"`, ast.NewStringLiteral(`"\xFF"`, "\xFF", nil, 0, 6), false},
		{"hex_mixed_case", `"\xaB"`, ast.NewStringLiteral(`"\xaB"`, "\xaB", nil, 0, 6), false},

		// Mixed content
		{"mixed_escapes", `"Hello\n\tWorld\r\n"`, ast.NewStringLiteral(`"Hello\n\tWorld\r\n"`, "Hello\n\tWorld\r\n", nil, 0, 20), false},
		{"mixed_unicode_hex", `"\u1234\xFF"`, ast.NewStringLiteral(`"\u1234\xFF"`, "\u1234\xFF", nil, 0, 12), false},
		{"mixed_all", `"Hello\n\t\u1234\xFF\0World"`, ast.NewStringLiteral(`"Hello\n\t\u1234\xFF\0World"`, "Hello\n\t\u1234\xFF\x00World", nil, 0, 28), false},

		// Invalid cases - Unicode
		{"invalid_unicode_empty", `"\u"`, nil, true},
		{"invalid_unicode_short", `"\u123"`, nil, true},
		{"invalid_unicode_long", `"\U1234"`, nil, true},
		{"invalid_unicode_letters", `"\uXYZ"`, nil, true},
		{"invalid_unicode_partial", `"\u12G4"`, nil, true},
		{"invalid_unicode_space", `"\u 123"`, nil, true},

		// Invalid cases - Hex
		{"invalid_hex_empty", `"\x"`, nil, true},
		{"invalid_hex_short", `"\xF"`, nil, true},
		{"invalid_hex_letters", `"\xXY"`, nil, true},
		{"invalid_hex_space", `"\x F"`, nil, true},

		// Invalid cases - General
		{"unterminated", `"abc`, nil, true},
		{"unterminated_escape", `"abc\`, nil, true},
		{"invalid_escape", `"\k"`, nil, true},
		{"invalid_escape_exclamation", `"\!"`, nil, true},
		{"bare_backslash", `"\"`, nil, true},
		{"newline_in_string", "\"abc\ndef\"", nil, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			source := source.NewSource([]byte(test.input), "test.tup")
			tokens, err := tok.Tokenize(source.Contents, source.Filename)
			if err != nil {
				t.Errorf("Tokenize(%q) = %v", test.input, err)
			}
			got, _, err := StringLiteral(tokens)
			if test.wantErr {
				if err == nil {
					t.Errorf("StringLiteral(%q) = %v, want error", test.input, got)
				}
				if test.want == nil && got != nil {
					t.Errorf("StringLiteral(%q) = %v, want nil", test.input, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("StringLiteral(%q) = %v", test.input, err)
			}
			if got.Value != test.want.Value {
				t.Errorf("StringLiteral(%q).Value = %v, want %v", test.input, got.Value, test.want.Value)
			}
			if got.StringValue != test.want.StringValue {
				t.Errorf("StringLiteral(%q).StringValue = %v, want %v", test.input, got.StringValue, test.want.StringValue)
			}
			if got.StartOffset != test.want.StartOffset {
				t.Errorf("StringLiteral(%q).StartOffset = %v, want %v", test.input, got.StartOffset, test.want.StartOffset)
			}
			if got.Length != test.want.Length {
				t.Errorf("StringLiteral(%q).Length = %v, want %v", test.input, got.Length, test.want.Length)
			}
		})
	}
}
