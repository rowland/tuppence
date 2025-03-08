package tok

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
	if lastToken.Type != TokEOF {
		t.Errorf("Expected EOF token, got %v", TokenTypes[lastToken.Type])
	}

	// For multi-line tokens, the column should be relative to the last line
	lastNewline := strings.LastIndexByte(source, '\n')
	var expectedCol int
	if lastNewline < 0 {
		expectedCol = len(source)
	} else {
		expectedCol = len(source) - (lastNewline + 1)
	}
	if lastToken.Column() != expectedCol {
		t.Errorf("Expected column %d, got %d", expectedCol, lastToken.Column())
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
	if lastNewline >= 0 && nextToken.Type != TokEOL {
		t.Errorf("Expected EOL token after newline in invalid token, got %v", TokenTypes[nextToken.Type])
	}

	// For invalid tokens, we need to find where the token became invalid
	var expectedCol int
	if lastNewline < 0 {
		expectedCol = len(source)
	} else {
		// For invalid tokens that contain a newline, the column should be
		// the length up to the newline
		expectedCol = lastNewline
	}
	if nextToken.Column() != expectedCol {
		t.Errorf("Expected column %d, got %d", expectedCol, nextToken.Column())
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
	if lastToken.Type != TokEOF {
		t.Errorf("Expected EOF token, got %v", lastToken.Type)
	}
	if invalid && !invalidTokenSeen {
		t.Errorf("Expected to find invalid token")
	}

	// For multi-line tokens, the column should be relative to the last line
	lastNewline := strings.LastIndexByte(source, '\n')
	var expectedCol int
	if lastNewline < 0 {
		expectedCol = len(source)
	} else {
		expectedCol = len(source) - lastNewline
	}
	if lastToken.Column() != expectedCol {
		t.Errorf("Expected column %d, got %d", expectedCol, lastToken.Column())
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
	if lastToken.Type != TokEOF {
		t.Errorf("Expected EOF token, got %v", lastToken.Type)
	}

	// For multi-line tokens, the column should be relative to the last line
	lastNewline := strings.LastIndexByte(source, '\n')
	var expectedCol int
	if lastNewline < 0 {
		expectedCol = len(source)
	} else {
		expectedCol = len(source) - lastNewline
	}
	if lastToken.Column() != expectedCol {
		t.Errorf("Expected column %d, got %d", expectedCol, lastToken.Column())
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
	if lastToken.Type != TokEOF {
		t.Errorf("Expected EOF token, got %v", lastToken.Type)
	}

	// For multi-line tokens, the column should be relative to the last line
	lastNewline := strings.LastIndexByte(source, '\n')
	var expectedCol int
	if lastNewline < 0 {
		expectedCol = len(source)
	} else {
		expectedCol = len(source) - (lastNewline + 1)
	}
	if lastToken.Column() != expectedCol {
		t.Errorf("Expected column %d, got %d", expectedCol, lastToken.Column())
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
		TokAt,
		TokCloseBrace,
		TokCloseBracket,
		TokCloseParen,
		TokColon,
		TokComma,
		TokDot,
		TokOpenBrace,
		TokOpenBracket,
		TokOpenParen,
		TokQuestionMark,
		TokSemiColon,
	})
}

func TestOperators(t *testing.T) {
	testTokenizeSeq(t, "?+ ?/ ?% ?* ?- / - % * ! + ^ << >>", []TokenType{
		TokOpCheckedAdd,
		TokOpCheckedDiv,
		TokOpCheckedMod,
		TokOpCheckedMul,
		TokOpCheckedSub,
		TokOpDiv,
		TokOpMinus,
		TokOpMod,
		TokOpMul,
		TokOpNot,
		TokOpPlus,
		TokOpPow,
		TokOpShiftLeft,
		TokOpShiftRight,
	})
}

func TestBitwiseOperators(t *testing.T) {
	testTokenizeSeq(t, "& | ~", []TokenType{
		TokOpBitwiseAnd,
		TokOpBitwiseOr,
		TokOpBitwiseNot,
	})
}

func TestRelationalOperators(t *testing.T) {
	testTokenizeSeq(t, "== >= > <= < != =~ <=>", []TokenType{
		TokOpEqual,
		TokOpGreaterEqual,
		TokOpGreaterThan,
		TokOpLessEqual,
		TokOpLessThan,
		TokOpNotEqual,
		TokOpMatches,
		TokOpCompareTo,
	})
}

func TestLogicalOperators(t *testing.T) {
	testTokenizeSeq(t, "&& ||", []TokenType{
		TokOpLogicalAnd,
		TokOpLogicalOr,
	})
}

func TestAssignment(t *testing.T) {
	testTokenizeSeq(t, "&= |= /= = &&= ||= -= %= *= += ^= <<= >>=", []TokenType{
		TokOpBitwiseAndEqual,
		TokOpBitwiseOrEqual,
		TokOpDivEqual,
		TokOpAssign,
		TokOpLogicalAndEqual,
		TokOpLogicalOrEqual,
		TokOpMinusEqual,
		TokOpModEqual,
		TokOpMulEqual,
		TokOpPlusEqual,
		TokOpPowEqual,
		TokOpShiftLeftEqual,
		TokOpShiftRightEqual,
	})
}

func TestIdentifiers(t *testing.T) {
	testTokenizeSeq(t, "abc Def", []TokenType{
		TokIdentifier,
		TokTypeIdentifier,
	})
}

func TestKeywords(t *testing.T) {
	testTokenizeSeq(t,
		"array break continue contract else enum error fn for fx if in it "+
			"import mut return switch try try_break try_continue type typeof union",
		[]TokenType{
			TokKeywordArray,
			TokKeywordBreak,
			TokKeywordContinue,
			TokKeywordContract,
			TokKeywordElse,
			TokKeywordEnum,
			TokKeywordError,
			TokKeywordFn,
			TokKeywordFor,
			TokKeywordFx,
			TokKeywordIf,
			TokKeywordIn,
			TokKeywordIt,
			TokKeywordImport,
			TokKeywordMut,
			TokKeywordReturn,
			TokKeywordSwitch,
			TokKeywordTry,
			TokKeywordTryBreak,
			TokKeywordTryContinue,
			TokKeywordType,
			TokKeywordTypeof,
			TokKeywordUnion,
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
		{"empty_after_prefix", "0b", TokBinaryLit, true},
		{"zero", "0b0", TokBinaryLit, false},
		{"one", "0b1", TokBinaryLit, false},
		{"complex_binary", "0b10101100", TokBinaryLit, false},

		// Invalid digits
		{"invalid_2", "0b2", TokBinaryLit, true},
		{"invalid_3", "0b3", TokBinaryLit, true},
		{"invalid_4", "0b4", TokBinaryLit, true},
		{"invalid_5", "0b5", TokBinaryLit, true},
		{"invalid_6", "0b6", TokBinaryLit, true},
		{"invalid_7", "0b7", TokBinaryLit, true},
		{"invalid_8", "0b8", TokBinaryLit, true},
		{"invalid_9", "0b9", TokBinaryLit, true},
		{"invalid_a", "0ba", TokBinaryLit, true},
		{"invalid_b", "0bb", TokBinaryLit, true},
		{"invalid_c", "0bc", TokBinaryLit, true},
		{"invalid_d", "0bd", TokBinaryLit, true},
		{"invalid_e", "0be", TokBinaryLit, true},
		{"invalid_f", "0bf", TokBinaryLit, true},
		{"invalid_z", "0bz", TokBinaryLit, true},

		// Underscore cases
		{"invalid_leading_underscore", "0b_", TokBinaryLit, true},
		{"invalid_underscore_after_prefix", "0b_0", TokBinaryLit, true},
		{"valid_trailing_underscore", "0b1_", TokBinaryLit, false},
		{"valid_double_underscore", "0b0__1", TokBinaryLit, false},
		{"valid_middle_underscore", "0b0_1_", TokBinaryLit, false},

		// Other cases
		{"invalid_uppercase_prefix", "0B0", TokDecimalLit, true},
		{"invalid_e_suffix", "0b1e", TokBinaryLit, true},
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
	testTokenize(t, "false", TokBoolLit)
	testTokenize(t, "true", TokBoolLit)
}

func TestDecimalLiterals(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantType TokenType
		wantErr  bool
	}{
		// Single digits
		{"zero", "0", TokDecimalLit, false},
		{"one", "1", TokDecimalLit, false},
		{"two", "2", TokDecimalLit, false},
		{"three", "3", TokDecimalLit, false},
		{"four", "4", TokDecimalLit, false},
		{"five", "5", TokDecimalLit, false},
		{"six", "6", TokDecimalLit, false},
		{"seven", "7", TokDecimalLit, false},
		{"eight", "8", TokDecimalLit, false},
		{"nine", "9", TokDecimalLit, false},

		// Leading zeros and underscores
		{"underscore_after_zero", "0_0", TokDecimalLit, false},
		{"leading_zeros", "0001", TokDecimalLit, false},

		// Complex numbers
		{"all_digits", "01234567890", TokDecimalLit, false},
		{"grouped_by_three", "012_345_6789_0", TokDecimalLit, false},
		{"max_underscores", "0_1_2_3_4_5_6_7_8_9_0", TokDecimalLit, false},

		// Invalid characters in numbers
		{"lowercase_letter", "123a", TokDecimalLit, true},
		{"uppercase_letter", "123A", TokDecimalLit, true},
		{"letter_in_middle", "12a34", TokDecimalLit, true},

		// Identifiers that look like numbers
		{"leading_underscore", "_123", TokIdentifier, false},
		{"multiple_leading_underscores", "__123", TokIdentifier, false},

		// Valid number followed by underscore
		{"trailing_underscore", "123_", TokDecimalLit, false},
		{"trailing_multiple_underscores", "123__", TokDecimalLit, false},

		// Sequence cases
		{"number_then_comma", "123,", TokDecimalLit, false},      // Should parse as separate tokens
		{"underscore_then_comma", "123_,", TokDecimalLit, false}, // Should parse as separate tokens
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if strings.HasSuffix(tt.input, ",") {
				// For sequence cases, use testTokenizeSeq
				testTokenizeSeq(t, tt.input, []TokenType{tt.wantType, TokComma})
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
		{"zero", "0o0", TokOctalLit, false},
		{"one", "0o1", TokOctalLit, false},
		{"two", "0o2", TokOctalLit, false},
		{"three", "0o3", TokOctalLit, false},
		{"four", "0o4", TokOctalLit, false},
		{"five", "0o5", TokOctalLit, false},
		{"six", "0o6", TokOctalLit, false},
		{"seven", "0o7", TokOctalLit, false},

		// Invalid digits
		{"invalid_8", "0o8", TokOctalLit, true},
		{"invalid_9", "0o9", TokOctalLit, true},
		{"invalid_a", "0oa", TokOctalLit, true},
		{"invalid_b", "0ob", TokOctalLit, true},
		{"invalid_c", "0oc", TokOctalLit, true},
		{"invalid_d", "0od", TokOctalLit, true},
		{"invalid_e", "0oe", TokOctalLit, true},
		{"invalid_f", "0of", TokOctalLit, true},
		{"invalid_z", "0oz", TokOctalLit, true},

		// Complex numbers
		{"all_octal_digits", "0o01234567", TokOctalLit, false},
		{"single_underscore", "0o0123_4567", TokOctalLit, false},
		{"multiple_underscores", "0o01_23_45_67", TokOctalLit, false},
		{"max_underscores", "0o0_1_2_3_4_5_6_7", TokOctalLit, false},

		// Invalid underscore positions
		{"invalid_leading_underscore", "0o_", TokOctalLit, true},
		{"invalid_underscore_after_prefix", "0o_0", TokOctalLit, true},
		{"valid_trailing_underscore", "0o1_", TokOctalLit, false},
		{"valid_double_underscore", "0o0__1", TokOctalLit, false},
		{"valid_middle_underscore", "0o0_1_", TokOctalLit, false},

		// Invalid prefix cases
		{"invalid_uppercase_prefix", "0O0", TokDecimalLit, true},
		{"empty_after_prefix", "0o", TokOctalLit, true},

		// Invalid suffix cases
		{"invalid_e_suffix", "0o1e", TokOctalLit, true},
		{"invalid_e_suffix_with_number", "0o1e0", TokOctalLit, true},

		// Sequence cases
		{"octal_then_comma", "0o1,", TokOctalLit, false},       // Should parse as separate tokens
		{"underscore_then_comma", "0o1_,", TokOctalLit, false}, // Should parse as separate tokens
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if strings.HasSuffix(tt.input, ",") {
				// For sequence cases, use testTokenizeSeq
				testTokenizeSeq(t, tt.input, []TokenType{tt.wantType, TokComma})
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
		{"zero", "0x0", TokHexLit, false},
		{"one", "0x1", TokHexLit, false},
		{"two", "0x2", TokHexLit, false},
		{"three", "0x3", TokHexLit, false},
		{"four", "0x4", TokHexLit, false},
		{"five", "0x5", TokHexLit, false},
		{"six", "0x6", TokHexLit, false},
		{"seven", "0x7", TokHexLit, false},
		{"eight", "0x8", TokHexLit, false},
		{"nine", "0x9", TokHexLit, false},

		// Lowercase hex letters
		{"lowercase_a", "0xa", TokHexLit, false},
		{"lowercase_b", "0xb", TokHexLit, false},
		{"lowercase_c", "0xc", TokHexLit, false},
		{"lowercase_d", "0xd", TokHexLit, false},
		{"lowercase_e", "0xe", TokHexLit, false},
		{"lowercase_f", "0xf", TokHexLit, false},

		// Uppercase hex letters
		{"uppercase_a", "0xA", TokHexLit, false},
		{"uppercase_b", "0xB", TokHexLit, false},
		{"uppercase_c", "0xC", TokHexLit, false},
		{"uppercase_d", "0xD", TokHexLit, false},
		{"uppercase_e", "0xE", TokHexLit, false},
		{"uppercase_f", "0xF", TokHexLit, false},

		// Invalid letters
		{"invalid_g", "0xg", TokHexLit, true},
		{"invalid_G", "0xG", TokHexLit, true},
		{"invalid_z", "0xz", TokHexLit, true},
		{"invalid_Z", "0xZ", TokHexLit, true},

		// Complex numbers
		{"leading_zeros", "0x0000", TokHexLit, false},
		{"repeated_letters", "0xAA", TokHexLit, false},
		{"all_fs", "0xFFFF", TokHexLit, false},
		{"all_hex_digits", "0x0123456789ABCDEF", TokHexLit, false},

		// Underscore cases
		{"single_group_underscore", "0x0123_4567_89AB_CDEF", TokHexLit, false},
		{"multiple_group_underscore", "0x01_23_45_67_89AB_CDE_F", TokHexLit, false},
		{"max_underscores", "0x0_1_2_3_4_5_6_7_8_9_A_B_C_D_E_F", TokHexLit, false},
		{"invalid_leading_underscore", "0x_", TokHexLit, true},
		{"invalid_underscore_after_prefix", "0x_1", TokHexLit, true},
		{"valid_trailing_underscore", "0x1_", TokHexLit, false},
		{"valid_double_underscore", "0x0__1", TokHexLit, false},
		{"valid_middle_underscore", "0x0_1_", TokHexLit, false},

		// Invalid prefix cases
		{"invalid_uppercase_prefix", "0X0", TokDecimalLit, true},
		{"empty_after_prefix", "0x", TokHexLit, true},

		// Sequence cases
		{"hex_then_comma", "0x1,", TokHexLit, false},         // Should parse as separate tokens
		{"underscore_then_comma", "0x1_,", TokHexLit, false}, // Should parse as separate tokens
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if strings.HasSuffix(tt.input, ",") {
				// For sequence cases, use testTokenizeSeq
				testTokenizeSeq(t, tt.input, []TokenType{tt.wantType, TokComma})
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
		{"empty", "``", TokRawStringLit, false},
		{"simple", "`abc`", TokRawStringLit, false},
		{"spaces", "`   `", TokRawStringLit, false},
		{"tabs", "`\t\t`", TokRawStringLit, false},

		// Escaped backtick cases
		{"embedded_backtick", "`abc``def`", TokRawStringLit, false},
		{"multiple_escaped_backticks", "`a``b``c`", TokRawStringLit, false},
		{"double_backtick", "`a``b`", TokRawStringLit, false},
		{"escaped_backtick_with_space", "`a`` b`", TokRawStringLit, false},
		{"escaped_backtick_with_newline", "`a``\nb`", TokRawStringLit, false},

		// Newline cases
		{"single_newline", "`abc\ndef`", TokRawStringLit, false},
		{"multiple_newlines", "`first\nsecond\nthird`", TokRawStringLit, false},
		{"only_newlines", "`\n\n\n`", TokRawStringLit, false},
		{"newline_after_backtick", "`abc``\ndef`", TokRawStringLit, false},
		{"newline_before_backtick", "`abc\n``def`", TokRawStringLit, false},
		{"starts_with_newline", "`\nabc`", TokRawStringLit, false},
		{"ends_with_newline", "`abc\n`", TokRawStringLit, false},

		// Special character cases
		{"special_chars", "`!@#$%^&*()`", TokRawStringLit, false},
		{"unicode_chars", "`αβγδε`", TokRawStringLit, false},
		{"emoji", "`🙂🌟🎉`", TokRawStringLit, false},
		{"escape_sequences_raw", "`\\n\\t\\r`", TokRawStringLit, false},
		{"quotes_in_raw", "`\"single\" and 'double'`", TokRawStringLit, false},
		{"brackets_braces", "`[{(<>)}]`", TokRawStringLit, false},

		// Whitespace cases
		{"mixed_whitespace", "`space\ttab\nline`", TokRawStringLit, false},
		{"carriage_return", "`line1\rline2`", TokRawStringLit, false},
		{"all_whitespace_types", "`\n\t\r \f\v`", TokRawStringLit, false},

		// Invalid cases
		{"unterminated_after_backtick", "`abc``", TokRawStringLit, true},
		{"unterminated", "`abc", TokRawStringLit, true},
		{"unterminated_with_newlines", "`first\nsecond\nthird", TokRawStringLit, true},
		{"unterminated_after_escape", "`abc``def", TokRawStringLit, true},
		{"unterminated_unicode", "`αβγ", TokRawStringLit, true},
		{"unterminated_emoji", "`🙂🌟", TokRawStringLit, true},
		{"empty_unterminated", "`", TokRawStringLit, true},
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
		{"empty", `""`, TokStringLit, false},
		{"simple", `"abc"`, TokStringLit, false},
		{"spaces", `"   "`, TokStringLit, false},
		{"tabs", `"\t\t"`, TokStringLit, false},

		// Simple escape sequences
		{"newline", `"\n"`, TokStringLit, false},
		{"tab", `"\t"`, TokStringLit, false},
		{"double_quote", `"\""`, TokStringLit, false},
		{"single_quote", `"\'"`, TokStringLit, false},
		{"backslash", `"\\"`, TokStringLit, false},
		{"carriage_return", `"\r"`, TokStringLit, false},
		{"backspace", `"\b"`, TokStringLit, false},
		{"form_feed", `"\f"`, TokStringLit, false},
		{"vertical_tab", `"\v"`, TokStringLit, false},
		{"null", `"\0"`, TokStringLit, false},
		{"all_escapes", `"\n\t\"\'\\\r\b\f\v\0"`, TokStringLit, false},

		// Unicode escapes
		{"unicode_4_digit", `"\u1234"`, TokStringLit, false},
		{"unicode_8_digit", `"\U12345678"`, TokStringLit, false},
		{"unicode_max", `"\uFFFF"`, TokStringLit, false},
		{"unicode_min", `"\u0000"`, TokStringLit, false},
		{"unicode_multiple", `"\u1234\u5678"`, TokStringLit, false},
		{"unicode_mixed_case", `"\uAbCd"`, TokStringLit, false},

		// Hex escapes
		{"hex_byte", `"\xEF"`, TokStringLit, false},
		{"hex_multiple", `"\xEF\xBB\xBF"`, TokStringLit, false},
		{"hex_min", `"\x00"`, TokStringLit, false},
		{"hex_max", `"\xFF"`, TokStringLit, false},
		{"hex_mixed_case", `"\xaB"`, TokStringLit, false},

		// Mixed content
		{"mixed_escapes", `"Hello\n\tWorld\r\n"`, TokStringLit, false},
		{"mixed_unicode_hex", `"\u1234\xFF"`, TokStringLit, false},
		{"mixed_all", `"Hello\n\t\u1234\xFF\0World"`, TokStringLit, false},

		// Invalid cases - Unicode
		{"invalid_unicode_empty", `"\u"`, TokStringLit, true},
		{"invalid_unicode_short", `"\u123"`, TokStringLit, true},
		{"invalid_unicode_long", `"\U1234"`, TokStringLit, true},
		{"invalid_unicode_letters", `"\uXYZ"`, TokStringLit, true},
		{"invalid_unicode_partial", `"\u12G4"`, TokStringLit, true},
		{"invalid_unicode_space", `"\u 123"`, TokStringLit, true},

		// Invalid cases - Hex
		{"invalid_hex_empty", `"\x"`, TokStringLit, true},
		{"invalid_hex_short", `"\xF"`, TokStringLit, true},
		{"invalid_hex_letters", `"\xXY"`, TokStringLit, true},
		{"invalid_hex_space", `"\x F"`, TokStringLit, true},

		// Invalid cases - General
		{"unterminated", `"abc`, TokStringLit, true},
		{"unterminated_escape", `"abc\`, TokStringLit, true},
		{"invalid_escape", `"\k"`, TokStringLit, true},
		{"invalid_escape_exclamation", `"\!"`, TokStringLit, true},
		{"bare_backslash", `"\"`, TokStringLit, true},
		{"newline_in_string", "\"abc\ndef\"", TokStringLit, true},
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
		{"empty", `"Hello \()!"`, TokInterpStringLit, false},
		{"simple_identifier", `"Hello \(name)!"`, TokInterpStringLit, false},
		{"simple_number", `"Count: \(123)!"`, TokInterpStringLit, false},
		{"simple_string", `"Name: \("Alice")!"`, TokInterpStringLit, false},
		{"simple_char", `"Initial: \('A')!"`, TokInterpStringLit, false},

		// Multiple interpolations
		{"two_interpolations", `"Hello \(first) \(last)!"`, TokInterpStringLit, false},
		{"three_interpolations", `"Hello \(title) \(first) \(last)!"`, TokInterpStringLit, false},
		{"adjacent_interpolations", `"\(a)\(b)\(c)"`, TokInterpStringLit, false},

		// Complex expressions
		{"expression_concat", `"Hello \(first + " " + last)!"`, TokInterpStringLit, false},
		{"expression_math", `"Sum: \(a + b + c)!"`, TokInterpStringLit, false},
		{"expression_nested_parens", `"Result: \((a + b) * (c + d))!"`, TokInterpStringLit, false},
		{"expression_function_call", `"Length: \(string.len(name))!"`, TokInterpStringLit, false},

		// Escapes and special characters
		{"escaped_quotes", `"Quote: \("\"nested\"")"`, TokInterpStringLit, false},
		{"escaped_backslash", `"Path: \("C:\\Program Files")"`, TokInterpStringLit, false},
		{"newline_escape", `"Lines: \("first\nsecond")"`, TokInterpStringLit, false},
		{"mixed_escapes", `"Mixed: \("tab\t\"quote\"\nline")"`, TokInterpStringLit, false},

		// Non-interpolation parentheses
		{"non_interpolation_parens", `"Welcome \(opponent1) )))---((( \(opponent2)"`, TokInterpStringLit, false},

		// Invalid cases
		{"unterminated_string", `"Hello \(name`, TokInterpStringLit, true},
		{"unterminated_interpolation", `"Hello \(name"`, TokInterpStringLit, true},
		{"unmatched_parens", `"Hello \((name)!"`, TokInterpStringLit, true},
		{"raw_newline", `"Hello \(name)`, TokInterpStringLit, true},
		{"nested_string_newline", `"Hello \("name`, TokInterpStringLit, true},
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
		{"empty_string", "```\n\n```", TokMultiLineStringLit, false},
		{"simple_text", "```\nSome text\n```", TokMultiLineStringLit, false},
		{"with_processor", "```processor\nSome text\n```", TokMultiLineStringLit, false},
		{"consistent_indentation", "```\n  First line\n  Second line\n  ```", TokMultiLineStringLit, false},

		// Escape sequences and special characters
		{"escape_sequences", "```\nSpecial: \\n\\t\\\"\\'\n```", TokMultiLineStringLit, false},
		{"interpolation", "```\nHello \\(name)!\n```", TokMultiLineStringLit, false},
		{"byte_escapes", "```\nByte: \\x48\\x69\n```", TokMultiLineStringLit, false},
		{"unicode_escapes", "```\nUnicode: \\u2603 \\U0001F680\n```", TokMultiLineStringLit, false},
		{"mixed_escapes", "```\nMixed: \\n\\t\\x48\\u2603\n```", TokMultiLineStringLit, false},

		// Invalid cases
		{"unclosed_string", "```\nUnclosed string", TokMultiLineStringLit, true},
		{"missing_newline_after_open", "```Some text\n```", TokMultiLineStringLit, true},
		{"missing_newline_before_close", "```\nSome text```", TokMultiLineStringLit, true},
		{"invalid_escape", "```\nInvalid escape: \\z\n```", TokMultiLineStringLit, true},
		{"incomplete_unicode", "```\nBad unicode: \\u26\n```", TokMultiLineStringLit, true},
		{"unterminated_interpolation", "```\nUnterminated: \\(expr\n```", TokMultiLineStringLit, true},
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
		TokTypeIdentifier, // List
		TokOpenBracket,    // [
		TokIdentifier,     // a
		TokCloseBracket,   // ]
		TokOpAssign,       // =
		TokKeywordUnion,   // union
		TokOpenParen,      // (
		TokEOL,            // \n
		TokTypeIdentifier, // Nil
		TokEOL,            // \n
		TokTypeIdentifier, // Cons
		TokOpenBracket,    // [
		TokIdentifier,     // a
		TokCloseBracket,   // ]
		TokEOL,            // \n
		TokCloseParen,     // )
		TokEOL,            // \n
		TokEOF,
	})
}

func TestComment(t *testing.T) {
	const source = "a = 5 # assign 5 to a\nb = 10\n"
	testTokenizeSeq(t, source, []TokenType{
		TokIdentifier, // a
		TokOpAssign,   // =
		TokDecimalLit, // 5
		TokComment,    // #
		TokIdentifier, // b
		TokOpAssign,   // =
		TokDecimalLit, // 10
		TokEOL,        // \n
		TokEOF,
	})
}

func TestFloatLiterals(t *testing.T) {
	// Valid floats
	testTokenize(t, "0.5", TokFloatLit)
	testTokenize(t, "1.23", TokFloatLit)
	testTokenize(t, "12.34e+5", TokFloatLit)
	testTokenize(t, "12.34e-5", TokFloatLit)
	testTokenize(t, "9e9", TokFloatLit)
	testTokenize(t, "10e+10", TokFloatLit)
	testTokenize(t, "10e-10", TokFloatLit)

	// Examples with underscores
	testTokenize(t, "1_2.3_4", TokFloatLit)  // 12.34
	testTokenize(t, "3.14_159", TokFloatLit) // 3.14159

	// Invalid floats
	testTokenizeInvalid(t, "1.2e", TokFloatLit)    // missing exponent digits
	testTokenizeInvalid(t, "1.2ez", TokFloatLit)   // non-digit exponent suffix
	testTokenizeInvalid(t, "1.2e++3", TokFloatLit) // double sign
	// testTokenizeInvalid(t, "1..2", TokenFloatLiteral)    // TODO: range
	testTokenizeInvalid(t, "123e", TokFloatLit)    // no exponent digits
	testTokenizeInvalid(t, "12.34e-", TokFloatLit) // minus with no digits
}

func TestDecimalMemberAccess(t *testing.T) {
	const source = "123.string()"
	testTokenizeSeq(t, source, []TokenType{
		TokDecimalLit, // 123
		TokDot,        // .
		TokIdentifier, // string
		TokOpenParen,  // (
		TokCloseParen, // )
		TokEOF,
	})
}

func TestBinaryMemberAccess(t *testing.T) {
	const source = "0b1010.string()"
	testTokenizeSeq(t, source, []TokenType{
		TokBinaryLit,  // 0b1010
		TokDot,        // .
		TokIdentifier, // string
		TokOpenParen,  // (
		TokCloseParen, // )
		TokEOF,
	})
}

func TestOctalMemberAccess(t *testing.T) {
	const source = "0o123.string()"
	testTokenizeSeq(t, source, []TokenType{
		TokOctalLit,   // 0o123
		TokDot,        // .
		TokIdentifier, // string
		TokOpenParen,  // (
		TokCloseParen, // )
		TokEOF,
	})
}

func TestHexadecimalMemberAccess(t *testing.T) {
	const source = "0xDEADBEEF.string()"
	testTokenizeSeq(t, source, []TokenType{
		TokHexLit,     // 0xDEADBEEF
		TokDot,        // .
		TokIdentifier, // string
		TokOpenParen,  // (
		TokCloseParen, // )
		TokEOF,
	})
}

func TestFloatMemberAccess(t *testing.T) {
	const source1 = "123e456.string()"
	testTokenizeSeq(t, source1, []TokenType{
		TokFloatLit,   // 123e456
		TokDot,        // .
		TokIdentifier, // string
		TokOpenParen,  // (
		TokCloseParen, // )
		TokEOF,
	})

	const source2 = "123.456.string()"
	testTokenizeSeq(t, source2, []TokenType{
		TokFloatLit,   // 123.456
		TokDot,        // .
		TokIdentifier, // string
		TokOpenParen,  // (
		TokCloseParen, // )
		TokEOF,
	})
}

func TestRangeOp(t *testing.T) {
	const source = "1..2"
	testTokenizeSeq(t, source, []TokenType{
		TokDecimalLit, // 1
		TokOpRange,    // ..
		TokDecimalLit, // 2
		TokEOF,
	})
}

func TestRestOp(t *testing.T) {
	const source = "...rest"
	testTokenizeSeq(t, source, []TokenType{
		TokOpRest,     // ...
		TokIdentifier, // rest
		TokEOF,
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
		{"uppercase_A", ":A", TokSymbolLit, false, false, nil},
		{"uppercase_Z", ":Z", TokSymbolLit, false, false, nil},
		{"lowercase_a", ":a", TokSymbolLit, false, false, nil},
		{"lowercase_z", ":z", TokSymbolLit, false, false, nil},
		{"all_uppercase", ":ABCDEFGHIJKLMNOPQRSTUVWXYZ", TokSymbolLit, false, false, nil},
		{"all_lowercase", ":abcdefghijklmnopqrstuvwxyz", TokSymbolLit, false, false, nil},

		// Valid quoted symbol literal
		{"quoted_symbol", `:"anything but a newline"`, TokSymbolLit, false, false, nil},

		// Invalid symbol literals
		{"invalid_digit_0", ":0", TokSymbolLit, true, false, nil},
		{"invalid_digit_1", ":1", TokSymbolLit, true, false, nil},
		{"invalid_digit_9", ":9", TokSymbolLit, true, false, nil},
		{"unterminated_quoted", `:"this symbol does not end`, TokSymbolLit, true, false, nil},
		{"newline_in_quoted", ":\"no\nnewlines!", TokSymbolLit, true, false, nil},

		// Symbol in sequence
		{"symbol_in_assignment", "foo = :foo", TokSymbolLit, false, true, []TokenType{
			TokIdentifier, // foo
			TokOpAssign,   // =
			TokSymbolLit,  // :foo
			TokEOF,
		}},

		{"symbol_in_comparison", `:foo == Symbol("foo")`, TokSymbolLit, false, true, []TokenType{
			TokSymbolLit,      // :foo
			TokOpEqual,        // ==
			TokTypeIdentifier, // Symbol
			TokOpenParen,      // (
			TokStringLit,      // "foo"
			TokCloseParen,     // )
			TokEOF,
		}},

		// Invalid symbol in sequence
		{"invalid_symbol_in_seq", `foo = :"foo`, TokSymbolLit, true, true, []TokenType{
			TokIdentifier, // foo
			TokOpAssign,   // =
			TokSymbolLit,  // :"foo
			TokEOF,
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
