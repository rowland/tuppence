package main

import (
	"strings"
	"testing"
)

// Helper: tokenize a sequence of tokens and check types.
func testTokenizeSeq(t *testing.T, source string, expected []TokenType) {
	t.Helper()
	tokenizer := NewTokenizer([]byte(source), "test.go")
	for i, exp := range expected {
		token := tokenizer.Next()
		if token.Type != exp {
			t.Errorf("At index %d: expected token type %v, got %v", i, TokenTypes[exp], TokenTypes[token.Type])
		}
		if token.Invalid {
			t.Errorf("At index %d: expected valid token, got invalid for %q", i, token.Value)
		}
	}
	lastToken := tokenizer.Next()
	if lastToken.Type != TokenEOF {
		t.Errorf("Expected EOF token, got %v", TokenTypes[lastToken.Type])
	}

	// For multi-line tokens, the column should be relative to the last line
	lastNewline := strings.LastIndexByte(source, '\n')
	var expectedCol int
	if lastNewline < 0 {
		expectedCol = len(source) + 1
	} else {
		expectedCol = len(source) - lastNewline
	}
	if lastToken.Column != expectedCol {
		t.Errorf("Expected column %d, got %d", expectedCol, lastToken.Column)
	}
}

// Helper for tokens that are expected to be invalid.
func testTokenizeInvalid(t *testing.T, source string, expectedType TokenType) {
	t.Helper()
	tokenizer := NewTokenizer([]byte(source), "test.go")
	token := tokenizer.Next()
	if token.Type != expectedType {
		t.Errorf("Expected token type %v, got %v for source %q", TokenTypes[expectedType], TokenTypes[token.Type], source)
	}
	if !token.Invalid {
		t.Errorf("Expected invalid token for %q", source)
	}
	// For invalid tokens that contain a newline, we only expect the value up to the newline
	lastNewline := strings.IndexByte(source, '\n')
	expectedValue := source
	if lastNewline >= 0 {
		expectedValue = source[:lastNewline]
	}
	if token.Value != expectedValue {
		t.Errorf("Expected token value %q, got %q", expectedValue, token.Value)
	}

	// Get the next token - it could be EOL, EOF, or something else
	nextToken := tokenizer.Next()
	if lastNewline >= 0 && nextToken.Type != TokenEOL {
		t.Errorf("Expected EOL token after newline in invalid token, got %v", TokenTypes[nextToken.Type])
	}

	// For invalid tokens, we need to find where the token became invalid
	var expectedCol int
	if lastNewline < 0 {
		expectedCol = len(source) + 1
	} else {
		// For invalid tokens that contain a newline, the column should be
		// the length up to the newline
		expectedCol = lastNewline + 1
	}
	if nextToken.Column != expectedCol {
		t.Errorf("Expected column %d, got %d", expectedCol, nextToken.Column)
	}
}

func testTokenizeSeqInvalid(t *testing.T, source string, expected []TokenType, invalid bool) {
	t.Helper()
	invalidTokenSeen := false
	tokenizer := NewTokenizer([]byte(source), "test.go")
	for i, exp := range expected {
		token := tokenizer.Next()
		if token.Type != exp {
			t.Errorf("At index %d: expected token type %v, got %v", i, TokenTypes[exp], TokenTypes[token.Type])
		}
		invalidTokenSeen = invalidTokenSeen || token.Invalid
	}
	lastToken := tokenizer.Next()
	if lastToken.Type != TokenEOF {
		t.Errorf("Expected EOF token, got %v", lastToken.Type)
	}
	if invalid && !invalidTokenSeen {
		t.Errorf("Expected to find invalid token")
	}

	// For multi-line tokens, the column should be relative to the last line
	lastNewline := strings.LastIndexByte(source, '\n')
	var expectedCol int
	if lastNewline < 0 {
		expectedCol = len(source) + 1
	} else {
		expectedCol = len(source) - lastNewline
	}
	if lastToken.Column != expectedCol {
		t.Errorf("Expected column %d, got %d", expectedCol, lastToken.Column)
	}
}

func testTokenize(t *testing.T, source string, expectedType TokenType) {
	t.Helper()
	tokenizer := NewTokenizer([]byte(source), "test.go")
	token := tokenizer.Next()
	if token.Type != expectedType {
		t.Errorf("Expected token type %v, got %v for source %q", expectedType, token.Type, source)
	}
	if token.Invalid {
		t.Errorf("Expected valid token for %q", source)
	}
	if token.Value != source {
		t.Errorf("Expected token value %q, got %q", source, token.Value)
	}
	lastToken := tokenizer.Next()
	if lastToken.Type != TokenEOF {
		t.Errorf("Expected EOF token, got %v", lastToken.Type)
	}

	// For multi-line tokens, the column should be relative to the last line
	lastNewline := strings.LastIndexByte(source, '\n')
	var expectedCol int
	if lastNewline < 0 {
		expectedCol = len(source) + 1
	} else {
		expectedCol = len(source) - lastNewline
	}
	if lastToken.Column != expectedCol {
		t.Errorf("Expected column %d, got %d", expectedCol, lastToken.Column)
	}
}

// Helper for multi-line string tokens that can contain newlines.
func testTokenizeMultiLine(t *testing.T, source string, expectedType TokenType) {
	t.Helper()
	tokenizer := NewTokenizer([]byte(source), "test.go")
	token := tokenizer.Next()
	if token.Type != expectedType {
		t.Errorf("Expected token type %v, got %v for source %q", expectedType, token.Type, source)
	}
	if token.Invalid {
		t.Errorf("Expected valid token for %q, got %q", source, token.Value)
	}
	if token.Value != source {
		t.Errorf("Expected token value %q, got %q", source, token.Value)
	}
	lastToken := tokenizer.Next()
	if lastToken.Type != TokenEOF {
		t.Errorf("Expected EOF token, got %v", lastToken.Type)
	}

	// For multi-line tokens, the column should be relative to the last line
	lastNewline := strings.LastIndexByte(source, '\n')
	var expectedCol int
	if lastNewline < 0 {
		expectedCol = len(source) + 1
	} else {
		expectedCol = len(source) - lastNewline
	}
	if lastToken.Column != expectedCol {
		t.Errorf("Expected column %d, got %d", expectedCol, lastToken.Column)
	}
}

// Helper for invalid multi-line string tokens that can contain newlines.
func testTokenizeMultiLineInvalid(t *testing.T, source string, expectedType TokenType) {
	t.Helper()
	tokenizer := NewTokenizer([]byte(source), "test.go")
	token := tokenizer.Next()
	if token.Type != expectedType {
		t.Errorf("Expected token type %v, got %v for source %q", TokenTypes[expectedType], TokenTypes[token.Type], source)
	}
	if !token.Invalid {
		t.Errorf("Expected invalid token for %q", source)
	}
	if token.Value != source {
		t.Errorf("Expected token value %q, got %q", source, token.Value)
	}
}

func TestSymbols(t *testing.T) {
	testTokenizeSeq(t, "@ } ] ) : , . { [ ( ? ;", []TokenType{
		TokenAt,
		TokenCloseBrace,
		TokenCloseBracket,
		TokenCloseParen,
		TokenColon,
		TokenComma,
		TokenDot,
		TokenOpenBrace,
		TokenOpenBracket,
		TokenOpenParen,
		TokenQuestionMark,
		TokenSemiColon,
	})
}

func TestOperators(t *testing.T) {
	testTokenizeSeq(t, "?+ ?/ ?% ?* ?- / - % * ! + ^ << >>", []TokenType{
		TokenOpCheckedAdd,
		TokenOpCheckedDiv,
		TokenOpCheckedMod,
		TokenOpCheckedMul,
		TokenOpCheckedSub,
		TokenOpDiv,
		TokenOpMinus,
		TokenOpMod,
		TokenOpMul,
		TokenOpNot,
		TokenOpPlus,
		TokenOpPow,
		TokenOpShiftLeft,
		TokenOpShiftRight,
	})
}

func TestBitwiseOperators(t *testing.T) {
	testTokenizeSeq(t, "& | ~", []TokenType{
		TokenOpBitwiseAnd,
		TokenOpBitwiseOr,
		TokenOpBitwiseNot,
	})
}

func TestRelationalOperators(t *testing.T) {
	testTokenizeSeq(t, "== >= > <= < != =~ <=>", []TokenType{
		TokenOpEqualEqual,
		TokenOpGreaterEqual,
		TokenOpGreaterThan,
		TokenOpLessEqual,
		TokenOpLessThan,
		TokenOpNotEqual,
		TokenOpMatches,
		TokenOpCompareTo,
	})
}

func TestLogicalOperators(t *testing.T) {
	testTokenizeSeq(t, "&& ||", []TokenType{
		TokenOpLogicalAnd,
		TokenOpLogicalOr,
	})
}

func TestAssignment(t *testing.T) {
	testTokenizeSeq(t, "&= |= /= = &&= ||= -= %= *= += ^= <<= >>=", []TokenType{
		TokenOpBitwiseAndEqual,
		TokenOpBitwiseOrEqual,
		TokenOpDivEqual,
		TokenOpEqual,
		TokenOpLogicalAndEqual,
		TokenOpLogicalOrEqual,
		TokenOpMinusEqual,
		TokenOpModEqual,
		TokenOpMulEqual,
		TokenOpPlusEqual,
		TokenOpPowEqual,
		TokenOpShiftLeftEqual,
		TokenOpShiftRightEqual,
	})
}

func TestIdentifiers(t *testing.T) {
	testTokenizeSeq(t, "abc Def", []TokenType{
		TokenIdentifier,
		TokenTypeIdentifier,
	})
}

func TestKeywords(t *testing.T) {
	testTokenizeSeq(t,
		"array break continue contract else enum error fn for fx if in it "+
			"import mut return switch try try_break try_continue type typeof union",
		[]TokenType{
			TokenKeywordArray,
			TokenKeywordBreak,
			TokenKeywordContinue,
			TokenKeywordContract,
			TokenKeywordElse,
			TokenKeywordEnum,
			TokenKeywordError,
			TokenKeywordFn,
			TokenKeywordFor,
			TokenKeywordFx,
			TokenKeywordIf,
			TokenKeywordIn,
			TokenKeywordIt,
			TokenKeywordImport,
			TokenKeywordMut,
			TokenKeywordReturn,
			TokenKeywordSwitch,
			TokenKeywordTry,
			TokenKeywordTryBreak,
			TokenKeywordTryContinue,
			TokenKeywordType,
			TokenKeywordTypeof,
			TokenKeywordUnion,
		})
}

func TestBinaryLiterals(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantType TokenType
		wantErr  bool
	}{
		// Basic valid cases
		{"empty_after_prefix", "0b", TokenBinaryLiteral, true},
		{"zero", "0b0", TokenBinaryLiteral, false},
		{"one", "0b1", TokenBinaryLiteral, false},
		{"complex_binary", "0b10101100", TokenBinaryLiteral, false},

		// Invalid digits
		{"invalid_2", "0b2", TokenBinaryLiteral, true},
		{"invalid_3", "0b3", TokenBinaryLiteral, true},
		{"invalid_4", "0b4", TokenBinaryLiteral, true},
		{"invalid_5", "0b5", TokenBinaryLiteral, true},
		{"invalid_6", "0b6", TokenBinaryLiteral, true},
		{"invalid_7", "0b7", TokenBinaryLiteral, true},
		{"invalid_8", "0b8", TokenBinaryLiteral, true},
		{"invalid_9", "0b9", TokenBinaryLiteral, true},
		{"invalid_a", "0ba", TokenBinaryLiteral, true},
		{"invalid_b", "0bb", TokenBinaryLiteral, true},
		{"invalid_c", "0bc", TokenBinaryLiteral, true},
		{"invalid_d", "0bd", TokenBinaryLiteral, true},
		{"invalid_e", "0be", TokenBinaryLiteral, true},
		{"invalid_f", "0bf", TokenBinaryLiteral, true},
		{"invalid_z", "0bz", TokenBinaryLiteral, true},

		// Underscore cases
		{"invalid_leading_underscore", "0b_", TokenBinaryLiteral, true},
		{"invalid_underscore_after_prefix", "0b_0", TokenBinaryLiteral, true},
		{"valid_trailing_underscore", "0b1_", TokenBinaryLiteral, false},
		{"valid_double_underscore", "0b0__1", TokenBinaryLiteral, false},
		{"valid_middle_underscore", "0b0_1_", TokenBinaryLiteral, false},

		// Other cases
		{"invalid_uppercase_prefix", "0B0", TokenDecimalLiteral, true},
		{"invalid_e_suffix", "0b1e", TokenBinaryLiteral, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				testTokenizeInvalid(t, tt.input, tt.wantType)
			} else {
				testTokenize(t, tt.input, tt.wantType)
			}
		})
	}
}

func TestBooleanLiterals(t *testing.T) {
	testTokenize(t, "false", TokenBooleanLiteral)
	testTokenize(t, "true", TokenBooleanLiteral)
}

func TestDecimalLiterals(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantType TokenType
		wantErr  bool
	}{
		// Single digits
		{"zero", "0", TokenDecimalLiteral, false},
		{"one", "1", TokenDecimalLiteral, false},
		{"two", "2", TokenDecimalLiteral, false},
		{"three", "3", TokenDecimalLiteral, false},
		{"four", "4", TokenDecimalLiteral, false},
		{"five", "5", TokenDecimalLiteral, false},
		{"six", "6", TokenDecimalLiteral, false},
		{"seven", "7", TokenDecimalLiteral, false},
		{"eight", "8", TokenDecimalLiteral, false},
		{"nine", "9", TokenDecimalLiteral, false},

		// Leading zeros and underscores
		{"underscore_after_zero", "0_0", TokenDecimalLiteral, false},
		{"leading_zeros", "0001", TokenDecimalLiteral, false},

		// Complex numbers
		{"all_digits", "01234567890", TokenDecimalLiteral, false},
		{"grouped_by_three", "012_345_6789_0", TokenDecimalLiteral, false},
		{"max_underscores", "0_1_2_3_4_5_6_7_8_9_0", TokenDecimalLiteral, false},

		// Invalid characters in numbers
		{"lowercase_letter", "123a", TokenDecimalLiteral, true},
		{"uppercase_letter", "123A", TokenDecimalLiteral, true},
		{"letter_in_middle", "12a34", TokenDecimalLiteral, true},

		// Identifiers that look like numbers
		{"leading_underscore", "_123", TokenIdentifier, false},
		{"multiple_leading_underscores", "__123", TokenIdentifier, false},
		{"only_underscore", "_", TokenIdentifier, false},
		{"multiple_underscores", "___", TokenIdentifier, false},
		{"underscore_and_letter", "_a", TokenIdentifier, false},
		{"letter_and_underscore", "a_", TokenIdentifier, false},

		// Valid number followed by underscore
		{"trailing_underscore", "123_", TokenDecimalLiteral, false},
		{"trailing_multiple_underscores", "123__", TokenDecimalLiteral, false},

		// Sequence cases
		{"number_then_comma", "123,", TokenDecimalLiteral, false},      // Should parse as separate tokens
		{"underscore_then_comma", "123_,", TokenDecimalLiteral, false}, // Should parse as separate tokens
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if strings.HasSuffix(tt.input, ",") {
				// For sequence cases, use testTokenizeSeq
				testTokenizeSeq(t, tt.input, []TokenType{tt.wantType, TokenComma})
			} else if tt.wantErr {
				testTokenizeInvalid(t, tt.input, tt.wantType)
			} else {
				testTokenize(t, tt.input, tt.wantType)
			}
		})
	}
}

func TestOctalLiterals(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantType TokenType
		wantErr  bool
	}{
		// Single digits
		{"zero", "0o0", TokenOctalLiteral, false},
		{"one", "0o1", TokenOctalLiteral, false},
		{"two", "0o2", TokenOctalLiteral, false},
		{"three", "0o3", TokenOctalLiteral, false},
		{"four", "0o4", TokenOctalLiteral, false},
		{"five", "0o5", TokenOctalLiteral, false},
		{"six", "0o6", TokenOctalLiteral, false},
		{"seven", "0o7", TokenOctalLiteral, false},

		// Invalid digits
		{"invalid_8", "0o8", TokenOctalLiteral, true},
		{"invalid_9", "0o9", TokenOctalLiteral, true},
		{"invalid_a", "0oa", TokenOctalLiteral, true},
		{"invalid_b", "0ob", TokenOctalLiteral, true},
		{"invalid_c", "0oc", TokenOctalLiteral, true},
		{"invalid_d", "0od", TokenOctalLiteral, true},
		{"invalid_e", "0oe", TokenOctalLiteral, true},
		{"invalid_f", "0of", TokenOctalLiteral, true},
		{"invalid_z", "0oz", TokenOctalLiteral, true},

		// Complex numbers
		{"all_octal_digits", "0o01234567", TokenOctalLiteral, false},
		{"single_underscore", "0o0123_4567", TokenOctalLiteral, false},
		{"multiple_underscores", "0o01_23_45_67", TokenOctalLiteral, false},
		{"max_underscores", "0o0_1_2_3_4_5_6_7", TokenOctalLiteral, false},

		// Invalid underscore positions
		{"invalid_leading_underscore", "0o_", TokenOctalLiteral, true},
		{"invalid_underscore_after_prefix", "0o_0", TokenOctalLiteral, true},
		{"valid_trailing_underscore", "0o1_", TokenOctalLiteral, false},
		{"valid_double_underscore", "0o0__1", TokenOctalLiteral, false},
		{"valid_middle_underscore", "0o0_1_", TokenOctalLiteral, false},

		// Invalid prefix cases
		{"invalid_uppercase_prefix", "0O0", TokenDecimalLiteral, true},
		{"empty_after_prefix", "0o", TokenOctalLiteral, true},

		// Invalid suffix cases
		{"invalid_e_suffix", "0o1e", TokenOctalLiteral, true},
		{"invalid_e_suffix_with_number", "0o1e0", TokenOctalLiteral, true},

		// Sequence cases
		{"octal_then_comma", "0o1,", TokenOctalLiteral, false},       // Should parse as separate tokens
		{"underscore_then_comma", "0o1_,", TokenOctalLiteral, false}, // Should parse as separate tokens
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if strings.HasSuffix(tt.input, ",") {
				// For sequence cases, use testTokenizeSeq
				testTokenizeSeq(t, tt.input, []TokenType{tt.wantType, TokenComma})
			} else if tt.wantErr {
				testTokenizeInvalid(t, tt.input, tt.wantType)
			} else {
				testTokenize(t, tt.input, tt.wantType)
			}
		})
	}
}

func TestHexadecimalLiterals(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantType TokenType
		wantErr  bool
	}{
		// Single digits (0-9)
		{"zero", "0x0", TokenHexadecimalLiteral, false},
		{"one", "0x1", TokenHexadecimalLiteral, false},
		{"two", "0x2", TokenHexadecimalLiteral, false},
		{"three", "0x3", TokenHexadecimalLiteral, false},
		{"four", "0x4", TokenHexadecimalLiteral, false},
		{"five", "0x5", TokenHexadecimalLiteral, false},
		{"six", "0x6", TokenHexadecimalLiteral, false},
		{"seven", "0x7", TokenHexadecimalLiteral, false},
		{"eight", "0x8", TokenHexadecimalLiteral, false},
		{"nine", "0x9", TokenHexadecimalLiteral, false},

		// Lowercase hex letters
		{"lowercase_a", "0xa", TokenHexadecimalLiteral, false},
		{"lowercase_b", "0xb", TokenHexadecimalLiteral, false},
		{"lowercase_c", "0xc", TokenHexadecimalLiteral, false},
		{"lowercase_d", "0xd", TokenHexadecimalLiteral, false},
		{"lowercase_e", "0xe", TokenHexadecimalLiteral, false},
		{"lowercase_f", "0xf", TokenHexadecimalLiteral, false},

		// Uppercase hex letters
		{"uppercase_a", "0xA", TokenHexadecimalLiteral, false},
		{"uppercase_b", "0xB", TokenHexadecimalLiteral, false},
		{"uppercase_c", "0xC", TokenHexadecimalLiteral, false},
		{"uppercase_d", "0xD", TokenHexadecimalLiteral, false},
		{"uppercase_e", "0xE", TokenHexadecimalLiteral, false},
		{"uppercase_f", "0xF", TokenHexadecimalLiteral, false},

		// Invalid letters
		{"invalid_g", "0xg", TokenHexadecimalLiteral, true},
		{"invalid_G", "0xG", TokenHexadecimalLiteral, true},
		{"invalid_z", "0xz", TokenHexadecimalLiteral, true},
		{"invalid_Z", "0xZ", TokenHexadecimalLiteral, true},

		// Complex numbers
		{"leading_zeros", "0x0000", TokenHexadecimalLiteral, false},
		{"repeated_letters", "0xAA", TokenHexadecimalLiteral, false},
		{"all_fs", "0xFFFF", TokenHexadecimalLiteral, false},
		{"all_hex_digits", "0x0123456789ABCDEF", TokenHexadecimalLiteral, false},

		// Underscore cases
		{"single_group_underscore", "0x0123_4567_89AB_CDEF", TokenHexadecimalLiteral, false},
		{"multiple_group_underscore", "0x01_23_45_67_89AB_CDE_F", TokenHexadecimalLiteral, false},
		{"max_underscores", "0x0_1_2_3_4_5_6_7_8_9_A_B_C_D_E_F", TokenHexadecimalLiteral, false},
		{"invalid_leading_underscore", "0x_", TokenHexadecimalLiteral, true},
		{"invalid_underscore_after_prefix", "0x_1", TokenHexadecimalLiteral, true},
		{"valid_trailing_underscore", "0x1_", TokenHexadecimalLiteral, false},
		{"valid_double_underscore", "0x0__1", TokenHexadecimalLiteral, false},
		{"valid_middle_underscore", "0x0_1_", TokenHexadecimalLiteral, false},

		// Invalid prefix cases
		{"invalid_uppercase_prefix", "0X0", TokenDecimalLiteral, true},
		{"empty_after_prefix", "0x", TokenHexadecimalLiteral, true},

		// Sequence cases
		{"hex_then_comma", "0x1,", TokenHexadecimalLiteral, false},         // Should parse as separate tokens
		{"underscore_then_comma", "0x1_,", TokenHexadecimalLiteral, false}, // Should parse as separate tokens
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if strings.HasSuffix(tt.input, ",") {
				// For sequence cases, use testTokenizeSeq
				testTokenizeSeq(t, tt.input, []TokenType{tt.wantType, TokenComma})
			} else if tt.wantErr {
				testTokenizeInvalid(t, tt.input, tt.wantType)
			} else {
				testTokenize(t, tt.input, tt.wantType)
			}
		})
	}
}

func TestRawStringLiterals(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantType TokenType
		wantErr  bool
	}{
		// Basic cases
		{"empty", "``", TokenRawStringLiteral, false},
		{"simple", "`abc`", TokenRawStringLiteral, false},
		{"spaces", "`   `", TokenRawStringLiteral, false},
		{"tabs", "`\t\t`", TokenRawStringLiteral, false},

		// Escaped backtick cases
		{"embedded_backtick", "`abc``def`", TokenRawStringLiteral, false},
		{"multiple_escaped_backticks", "`a``b``c`", TokenRawStringLiteral, false},
		{"double_backtick", "`a``b`", TokenRawStringLiteral, false},
		{"escaped_backtick_with_space", "`a`` b`", TokenRawStringLiteral, false},
		{"escaped_backtick_with_newline", "`a``\nb`", TokenRawStringLiteral, false},

		// Newline cases
		{"single_newline", "`abc\ndef`", TokenRawStringLiteral, false},
		{"multiple_newlines", "`first\nsecond\nthird`", TokenRawStringLiteral, false},
		{"only_newlines", "`\n\n\n`", TokenRawStringLiteral, false},
		{"newline_after_backtick", "`abc``\ndef`", TokenRawStringLiteral, false},
		{"newline_before_backtick", "`abc\n``def`", TokenRawStringLiteral, false},
		{"starts_with_newline", "`\nabc`", TokenRawStringLiteral, false},
		{"ends_with_newline", "`abc\n`", TokenRawStringLiteral, false},

		// Special character cases
		{"special_chars", "`!@#$%^&*()`", TokenRawStringLiteral, false},
		{"unicode_chars", "`αβγδε`", TokenRawStringLiteral, false},
		{"emoji", "`🙂🌟🎉`", TokenRawStringLiteral, false},
		{"escape_sequences_raw", "`\\n\\t\\r`", TokenRawStringLiteral, false},
		{"quotes_in_raw", "`\"single\" and 'double'`", TokenRawStringLiteral, false},
		{"brackets_braces", "`[{(<>)}]`", TokenRawStringLiteral, false},

		// Whitespace cases
		{"mixed_whitespace", "`space\ttab\nline`", TokenRawStringLiteral, false},
		{"carriage_return", "`line1\rline2`", TokenRawStringLiteral, false},
		{"all_whitespace_types", "`\n\t\r \f\v`", TokenRawStringLiteral, false},

		// Invalid cases
		{"unterminated_after_backtick", "`abc``", TokenRawStringLiteral, true},
		{"unterminated", "`abc", TokenRawStringLiteral, true},
		{"unterminated_with_newlines", "`first\nsecond\nthird", TokenRawStringLiteral, true},
		{"unterminated_after_escape", "`abc``def", TokenRawStringLiteral, true},
		{"unterminated_unicode", "`αβγ", TokenRawStringLiteral, true},
		{"unterminated_emoji", "`🙂🌟", TokenRawStringLiteral, true},
		{"empty_unterminated", "`", TokenRawStringLiteral, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				testTokenizeMultiLineInvalid(t, tt.input, tt.wantType)
			} else {
				testTokenizeMultiLine(t, tt.input, tt.wantType)
			}
		})
	}
}

func TestStringLiterals(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantType TokenType
		wantErr  bool
	}{
		// Basic cases
		{"empty", `""`, TokenStringLiteral, false},
		{"simple", `"abc"`, TokenStringLiteral, false},
		{"spaces", `"   "`, TokenStringLiteral, false},
		{"tabs", `"\t\t"`, TokenStringLiteral, false},

		// Simple escape sequences
		{"newline", `"\n"`, TokenStringLiteral, false},
		{"tab", `"\t"`, TokenStringLiteral, false},
		{"double_quote", `"\""`, TokenStringLiteral, false},
		{"single_quote", `"\'"`, TokenStringLiteral, false},
		{"backslash", `"\\"`, TokenStringLiteral, false},
		{"carriage_return", `"\r"`, TokenStringLiteral, false},
		{"backspace", `"\b"`, TokenStringLiteral, false},
		{"form_feed", `"\f"`, TokenStringLiteral, false},
		{"vertical_tab", `"\v"`, TokenStringLiteral, false},
		{"null", `"\0"`, TokenStringLiteral, false},
		{"all_escapes", `"\n\t\"\'\\\r\b\f\v\0"`, TokenStringLiteral, false},

		// Unicode escapes
		{"unicode_4_digit", `"\u1234"`, TokenStringLiteral, false},
		{"unicode_8_digit", `"\U12345678"`, TokenStringLiteral, false},
		{"unicode_max", `"\uFFFF"`, TokenStringLiteral, false},
		{"unicode_min", `"\u0000"`, TokenStringLiteral, false},
		{"unicode_multiple", `"\u1234\u5678"`, TokenStringLiteral, false},
		{"unicode_mixed_case", `"\uAbCd"`, TokenStringLiteral, false},

		// Hex escapes
		{"hex_byte", `"\xEF"`, TokenStringLiteral, false},
		{"hex_multiple", `"\xEF\xBB\xBF"`, TokenStringLiteral, false},
		{"hex_min", `"\x00"`, TokenStringLiteral, false},
		{"hex_max", `"\xFF"`, TokenStringLiteral, false},
		{"hex_mixed_case", `"\xaB"`, TokenStringLiteral, false},

		// Mixed content
		{"mixed_escapes", `"Hello\n\tWorld\r\n"`, TokenStringLiteral, false},
		{"mixed_unicode_hex", `"\u1234\xFF"`, TokenStringLiteral, false},
		{"mixed_all", `"Hello\n\t\u1234\xFF\0World"`, TokenStringLiteral, false},

		// Invalid cases - Unicode
		{"invalid_unicode_empty", `"\u"`, TokenStringLiteral, true},
		{"invalid_unicode_short", `"\u123"`, TokenStringLiteral, true},
		{"invalid_unicode_long", `"\U1234"`, TokenStringLiteral, true},
		{"invalid_unicode_letters", `"\uXYZ"`, TokenStringLiteral, true},
		{"invalid_unicode_partial", `"\u12G4"`, TokenStringLiteral, true},
		{"invalid_unicode_space", `"\u 123"`, TokenStringLiteral, true},

		// Invalid cases - Hex
		{"invalid_hex_empty", `"\x"`, TokenStringLiteral, true},
		{"invalid_hex_short", `"\xF"`, TokenStringLiteral, true},
		{"invalid_hex_letters", `"\xXY"`, TokenStringLiteral, true},
		{"invalid_hex_space", `"\x F"`, TokenStringLiteral, true},

		// Invalid cases - General
		{"unterminated", `"abc`, TokenStringLiteral, true},
		{"unterminated_escape", `"abc\`, TokenStringLiteral, true},
		{"invalid_escape", `"\k"`, TokenStringLiteral, true},
		{"invalid_escape_exclamation", `"\!"`, TokenStringLiteral, true},
		{"bare_backslash", `"\"`, TokenStringLiteral, true},
		{"newline_in_string", "\"abc\ndef\"", TokenStringLiteral, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				testTokenizeInvalid(t, tt.input, tt.wantType)
			} else {
				testTokenize(t, tt.input, tt.wantType)
			}
		})
	}
}

func TestInterpolatedStringLiterals(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantType TokenType
		wantErr  bool
	}{
		// Basic interpolation cases
		{"empty", `"Hello \()!"`, TokenInterpolatedStringLiteral, false},
		{"simple_identifier", `"Hello \(name)!"`, TokenInterpolatedStringLiteral, false},
		{"simple_number", `"Count: \(123)!"`, TokenInterpolatedStringLiteral, false},
		{"simple_string", `"Name: \("Alice")!"`, TokenInterpolatedStringLiteral, false},
		{"simple_char", `"Initial: \('A')!"`, TokenInterpolatedStringLiteral, false},

		// Multiple interpolations
		{"two_interpolations", `"Hello \(first) \(last)!"`, TokenInterpolatedStringLiteral, false},
		{"three_interpolations", `"Hello \(title) \(first) \(last)!"`, TokenInterpolatedStringLiteral, false},
		{"adjacent_interpolations", `"\(a)\(b)\(c)"`, TokenInterpolatedStringLiteral, false},

		// Complex expressions
		{"expression_concat", `"Hello \(first + " " + last)!"`, TokenInterpolatedStringLiteral, false},
		{"expression_math", `"Sum: \(a + b + c)!"`, TokenInterpolatedStringLiteral, false},
		{"expression_nested_parens", `"Result: \((a + b) * (c + d))!"`, TokenInterpolatedStringLiteral, false},
		{"expression_function_call", `"Length: \(string.len(name))!"`, TokenInterpolatedStringLiteral, false},

		// Escapes and special characters
		{"escaped_quotes", `"Quote: \("\"nested\"")"`, TokenInterpolatedStringLiteral, false},
		{"escaped_backslash", `"Path: \("C:\\Program Files")"`, TokenInterpolatedStringLiteral, false},
		{"newline_escape", `"Lines: \("first\nsecond")"`, TokenInterpolatedStringLiteral, false},
		{"mixed_escapes", `"Mixed: \("tab\t\"quote\"\nline")"`, TokenInterpolatedStringLiteral, false},

		// Non-interpolation parentheses
		{"non_interpolation_parens", `"Welcome \(opponent1) )))---((( \(opponent2)"`, TokenInterpolatedStringLiteral, false},

		// Invalid cases
		{"unterminated_string", `"Hello \(name`, TokenInterpolatedStringLiteral, true},
		{"unterminated_interpolation", `"Hello \(name"`, TokenInterpolatedStringLiteral, true},
		{"unmatched_parens", `"Hello \((name)!"`, TokenInterpolatedStringLiteral, true},
		{"raw_newline", `"Hello \(name)`, TokenInterpolatedStringLiteral, true},
		{"nested_string_newline", `"Hello \("name`, TokenInterpolatedStringLiteral, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				testTokenizeInvalid(t, tt.input, tt.wantType)
			} else {
				testTokenize(t, tt.input, tt.wantType)
			}
		})
	}
}

func TestMultiLineStringLiteral(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantType TokenType
		wantErr  bool
	}{
		// Basic cases
		// {"empty_string", "```\n\n```", TokenMultiLineStringLiteral, false},
		{"simple_text", "```\nSome text\n```", TokenMultiLineStringLiteral, false},
		{"with_processor", "```processor\nSome text\n```", TokenMultiLineStringLiteral, false},
		{"consistent_indentation", "```\n  First line\n  Second line\n  ```", TokenMultiLineStringLiteral, false},

		// Escape sequences and special characters
		{"escape_sequences", "```\nSpecial: \\n\\t\\\"\\'\n```", TokenMultiLineStringLiteral, false},
		{"interpolation", "```\nHello \\(name)!\n```", TokenMultiLineStringLiteral, false},
		{"byte_escapes", "```\nByte: \\x48\\x69\n```", TokenMultiLineStringLiteral, false},
		{"unicode_escapes", "```\nUnicode: \\u2603 \\U0001F680\n```", TokenMultiLineStringLiteral, false},
		{"mixed_escapes", "```\nMixed: \\n\\t\\x48\\u2603\n```", TokenMultiLineStringLiteral, false},

		// Invalid cases
		{"unclosed_string", "```\nUnclosed string", TokenMultiLineStringLiteral, true},
		{"missing_newline_after_open", "```Some text\n```", TokenMultiLineStringLiteral, true},
		{"missing_newline_before_close", "```\nSome text```", TokenMultiLineStringLiteral, true},
		{"invalid_escape", "```\nInvalid escape: \\z\n```", TokenMultiLineStringLiteral, true},
		{"incomplete_unicode", "```\nBad unicode: \\u26\n```", TokenMultiLineStringLiteral, true},
		{"unterminated_interpolation", "```\nUnterminated: \\(expr\n```", TokenMultiLineStringLiteral, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				testTokenizeMultiLineInvalid(t, tt.input, tt.wantType)
			} else {
				testTokenizeMultiLine(t, tt.input, tt.wantType)
			}
		})
	}
}

func TestUnion(t *testing.T) {
	const source = "List[a] = union(\n  Nil\n  Cons[a]\n)\n"
	testTokenizeSeq(t, source, []TokenType{
		TokenTypeIdentifier, // List
		TokenOpenBracket,    // [
		TokenIdentifier,     // a
		TokenCloseBracket,   // ]
		TokenOpEqual,        // =
		TokenKeywordUnion,   // union
		TokenOpenParen,      // (
		TokenEOL,            // \n
		TokenTypeIdentifier, // Nil
		TokenEOL,            // \n
		TokenTypeIdentifier, // Cons
		TokenOpenBracket,    // [
		TokenIdentifier,     // a
		TokenCloseBracket,   // ]
		TokenEOL,            // \n
		TokenCloseParen,     // )
		TokenEOL,            // \n
		TokenEOF,
	})
}

func TestComment(t *testing.T) {
	const source = "a = 5 # assign 5 to a\nb = 10\n"
	testTokenizeSeq(t, source, []TokenType{
		TokenIdentifier,     // a
		TokenOpEqual,        // =
		TokenDecimalLiteral, // 5
		TokenComment,        // #
		TokenIdentifier,     // b
		TokenOpEqual,        // =
		TokenDecimalLiteral, // 10
		TokenEOL,            // \n
		TokenEOF,
	})
}

func TestFloatLiterals(t *testing.T) {
	// Valid floats
	testTokenize(t, "0.5", TokenFloatLiteral)
	testTokenize(t, "1.23", TokenFloatLiteral)
	testTokenize(t, "12.34e+5", TokenFloatLiteral)
	testTokenize(t, "12.34e-5", TokenFloatLiteral)
	testTokenize(t, "9e9", TokenFloatLiteral)
	testTokenize(t, "10e+10", TokenFloatLiteral)
	testTokenize(t, "10e-10", TokenFloatLiteral)

	// Examples with underscores
	testTokenize(t, "1_2.3_4", TokenFloatLiteral)  // 12.34
	testTokenize(t, "3.14_159", TokenFloatLiteral) // 3.14159

	// Invalid floats
	testTokenizeInvalid(t, "1.2e", TokenFloatLiteral)    // missing exponent digits
	testTokenizeInvalid(t, "1.2ez", TokenFloatLiteral)   // non-digit exponent suffix
	testTokenizeInvalid(t, "1.2e++3", TokenFloatLiteral) // double sign
	// testTokenizeInvalid(t, "1..2", TokenFloatLiteral)    // TODO: range
	testTokenizeInvalid(t, "123e", TokenFloatLiteral)    // no exponent digits
	testTokenizeInvalid(t, "12.34e-", TokenFloatLiteral) // minus with no digits
}

func TestDecimalMemberAccess(t *testing.T) {
	const source = "123.string()"
	testTokenizeSeq(t, source, []TokenType{
		TokenDecimalLiteral, // 123
		TokenDot,            // .
		TokenIdentifier,     // string
		TokenOpenParen,      // (
		TokenCloseParen,     // )
		TokenEOF,
	})
}

func TestBinaryMemberAccess(t *testing.T) {
	const source = "0b1010.string()"
	testTokenizeSeq(t, source, []TokenType{
		TokenBinaryLiteral, // 0b1010
		TokenDot,           // .
		TokenIdentifier,    // string
		TokenOpenParen,     // (
		TokenCloseParen,    // )
		TokenEOF,
	})
}

func TestOctalMemberAccess(t *testing.T) {
	const source = "0o123.string()"
	testTokenizeSeq(t, source, []TokenType{
		TokenOctalLiteral, // 0o123
		TokenDot,          // .
		TokenIdentifier,   // string
		TokenOpenParen,    // (
		TokenCloseParen,   // )
		TokenEOF,
	})
}

func TestHexadecimalMemberAccess(t *testing.T) {
	const source = "0xDEADBEEF.string()"
	testTokenizeSeq(t, source, []TokenType{
		TokenHexadecimalLiteral, // 0xDEADBEEF
		TokenDot,                // .
		TokenIdentifier,         // string
		TokenOpenParen,          // (
		TokenCloseParen,         // )
		TokenEOF,
	})
}

func TestFloatMemberAccess(t *testing.T) {
	const source1 = "123e456.string()"
	testTokenizeSeq(t, source1, []TokenType{
		TokenFloatLiteral, // 123e456
		TokenDot,          // .
		TokenIdentifier,   // string
		TokenOpenParen,    // (
		TokenCloseParen,   // )
		TokenEOF,
	})

	const source2 = "123.456.string()"
	testTokenizeSeq(t, source2, []TokenType{
		TokenFloatLiteral, // 123.456
		TokenDot,          // .
		TokenIdentifier,   // string
		TokenOpenParen,    // (
		TokenCloseParen,   // )
		TokenEOF,
	})
}

func TestRangeOp(t *testing.T) {
	const source = "1..2"
	testTokenizeSeq(t, source, []TokenType{
		TokenDecimalLiteral, // 1
		TokenOpRange,        // ..
		TokenDecimalLiteral, // 2
		TokenEOF,
	})
}

func TestRestOp(t *testing.T) {
	const source = "...rest"
	testTokenizeSeq(t, source, []TokenType{
		TokenOpRest,     // ...
		TokenIdentifier, // rest
		TokenEOF,
	})
}

func TestSymbolLiterals(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantType TokenType
		wantErr  bool
		isSeq    bool
		seqTypes []TokenType
	}{
		// Valid simple symbol literals
		{"uppercase_A", ":A", TokenSymbolLiteral, false, false, nil},
		{"uppercase_Z", ":Z", TokenSymbolLiteral, false, false, nil},
		{"lowercase_a", ":a", TokenSymbolLiteral, false, false, nil},
		{"lowercase_z", ":z", TokenSymbolLiteral, false, false, nil},
		{"all_uppercase", ":ABCDEFGHIJKLMNOPQRSTUVWXYZ", TokenSymbolLiteral, false, false, nil},
		{"all_lowercase", ":abcdefghijklmnopqrstuvwxyz", TokenSymbolLiteral, false, false, nil},

		// Valid quoted symbol literal
		{"quoted_symbol", `:"anything but a newline"`, TokenSymbolLiteral, false, false, nil},

		// Invalid symbol literals
		{"invalid_digit_0", ":0", TokenSymbolLiteral, true, false, nil},
		{"invalid_digit_1", ":1", TokenSymbolLiteral, true, false, nil},
		{"invalid_digit_9", ":9", TokenSymbolLiteral, true, false, nil},
		{"unterminated_quoted", `:"this symbol does not end`, TokenSymbolLiteral, true, false, nil},
		{"newline_in_quoted", ":\"no\nnewlines!", TokenSymbolLiteral, true, false, nil},

		// Symbol in sequence
		{"symbol_in_assignment", "foo = :foo", TokenSymbolLiteral, false, true, []TokenType{
			TokenIdentifier,    // foo
			TokenOpEqual,       // =
			TokenSymbolLiteral, // :foo
			TokenEOF,
		}},

		{"symbol_in_comparison", `:foo == Symbol("foo")`, TokenSymbolLiteral, false, true, []TokenType{
			TokenSymbolLiteral,  // :foo
			TokenOpEqualEqual,   // ==
			TokenTypeIdentifier, // Symbol
			TokenOpenParen,      // (
			TokenStringLiteral,  // "foo"
			TokenCloseParen,     // )
			TokenEOF,
		}},

		// Invalid symbol in sequence
		{"invalid_symbol_in_seq", `foo = :"foo`, TokenSymbolLiteral, true, true, []TokenType{
			TokenIdentifier,    // foo
			TokenOpEqual,       // =
			TokenSymbolLiteral, // :"foo
			TokenEOF,
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.isSeq {
				if tt.wantErr {
					testTokenizeSeqInvalid(t, tt.input, tt.seqTypes, true)
				} else {
					testTokenizeSeq(t, tt.input, tt.seqTypes)
				}
			} else {
				if tt.wantErr {
					testTokenizeInvalid(t, tt.input, tt.wantType)
				} else {
					testTokenize(t, tt.input, tt.wantType)
				}
			}
		})
	}
}
