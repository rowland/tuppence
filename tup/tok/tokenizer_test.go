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
	// Check that error offset is set
	if token.ErrorOffset == 0 {
		t.Errorf("Expected error offset to be set for invalid token %q", source)
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
		TokOpSHL,
		TokOpSHR,
	})
}

func TestBitwiseOperators(t *testing.T) {
	testTokenizeSeq(t, "& | ~", []TokenType{
		TokOpBitAnd,
		TokOpBitOr,
		TokOpBitNot,
	})
}

func TestRelationalOperators(t *testing.T) {
	testTokenizeSeq(t, "== >= > <= < != =~ <=>", []TokenType{
		TokOpEQ,
		TokOpGE,
		TokOpGT,
		TokOpLE,
		TokOpLT,
		TokOpNE,
		TokOpMatch,
		TokOpCompare,
	})
}

func TestLogicalOperators(t *testing.T) {
	testTokenizeSeq(t, "&& ||", []TokenType{
		TokOpLogAnd,
		TokOpLogOr,
	})
}

func TestAssignment(t *testing.T) {
	testTokenizeSeq(t, "&= |= /= = &&= ||= -= %= *= += ^= <<= >>=", []TokenType{
		TokOpBitAndEQ,
		TokOpBitOrEQ,
		TokOpDivEQ,
		TokOpAssign,
		TokOpLogAndEQ,
		TokOpLogOrEQ,
		TokOpMinusEQ,
		TokOpModEQ,
		TokOpMulEQ,
		TokOpPlusEQ,
		TokOpPowEQ,
		TokOpSHL_EQ,
		TokOpSHR_EQ,
	})
}

func TestIdentifiers(t *testing.T) {
	testTokenizeSeq(t, "abc Def", []TokenType{
		TokID,
		TokTypeID,
	})
}

func TestKeywords(t *testing.T) {
	testTokenizeSeq(t,
		"array break continue contract else enum error fn for fx if in it "+
			"import mut return switch try try_break try_continue type typeof union",
		[]TokenType{
			TokKwArray,
			TokKwBreak,
			TokKwContinue,
			TokKwContract,
			TokKwElse,
			TokKwEnum,
			TokKwError,
			TokKwFn,
			TokKwFor,
			TokKwFx,
			TokKwIf,
			TokKwIn,
			TokKwIt,
			TokKwImport,
			TokKwMut,
			TokKwReturn,
			TokKwSwitch,
			TokKwTry,
			TokKwTryBreak,
			TokKwTryContinue,
			TokKwType,
			TokKwTypeof,
			TokKwUnion,
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
		{"empty_after_prefix", "0b", TokBinLit, true},
		{"zero", "0b0", TokBinLit, false},
		{"one", "0b1", TokBinLit, false},
		{"complex_binary", "0b10101100", TokBinLit, false},

		// Invalid digits
		{"invalid_2", "0b2", TokBinLit, true},
		{"invalid_3", "0b3", TokBinLit, true},
		{"invalid_4", "0b4", TokBinLit, true},
		{"invalid_5", "0b5", TokBinLit, true},
		{"invalid_6", "0b6", TokBinLit, true},
		{"invalid_7", "0b7", TokBinLit, true},
		{"invalid_8", "0b8", TokBinLit, true},
		{"invalid_9", "0b9", TokBinLit, true},
		{"invalid_a", "0ba", TokBinLit, true},
		{"invalid_b", "0bb", TokBinLit, true},
		{"invalid_c", "0bc", TokBinLit, true},
		{"invalid_d", "0bd", TokBinLit, true},
		{"invalid_e", "0be", TokBinLit, true},
		{"invalid_f", "0bf", TokBinLit, true},
		{"invalid_z", "0bz", TokBinLit, true},

		// Underscore cases
		{"invalid_leading_underscore", "0b_", TokBinLit, true},
		{"invalid_underscore_after_prefix", "0b_0", TokBinLit, true},
		{"valid_trailing_underscore", "0b1_", TokBinLit, false},
		{"valid_double_underscore", "0b0__1", TokBinLit, false},
		{"valid_middle_underscore", "0b0_1_", TokBinLit, false},

		// Other cases
		{"invalid_uppercase_prefix", "0B0", TokDecLit, true},
		{"invalid_e_suffix", "0b1e", TokBinLit, true},
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
		{"zero", "0", TokDecLit, false},
		{"one", "1", TokDecLit, false},
		{"two", "2", TokDecLit, false},
		{"three", "3", TokDecLit, false},
		{"four", "4", TokDecLit, false},
		{"five", "5", TokDecLit, false},
		{"six", "6", TokDecLit, false},
		{"seven", "7", TokDecLit, false},
		{"eight", "8", TokDecLit, false},
		{"nine", "9", TokDecLit, false},

		// Leading zeros and underscores
		{"underscore_after_zero", "0_0", TokDecLit, false},
		{"leading_zeros", "0001", TokDecLit, false},

		// Complex numbers
		{"all_digits", "01234567890", TokDecLit, false},
		{"grouped_by_three", "012_345_6789_0", TokDecLit, false},
		{"max_underscores", "0_1_2_3_4_5_6_7_8_9_0", TokDecLit, false},

		// Invalid characters in numbers
		{"lowercase_letter", "123a", TokDecLit, true},
		{"uppercase_letter", "123A", TokDecLit, true},
		{"letter_in_middle", "12a34", TokDecLit, true},

		// Identifiers that look like numbers
		{"leading_underscore", "_123", TokID, false},
		{"multiple_leading_underscores", "__123", TokID, false},

		// Valid number followed by underscore
		{"trailing_underscore", "123_", TokDecLit, false},
		{"trailing_multiple_underscores", "123__", TokDecLit, false},

		// Sequence cases
		{"number_then_comma", "123,", TokDecLit, false},      // Should parse as separate tokens
		{"underscore_then_comma", "123_,", TokDecLit, false}, // Should parse as separate tokens
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
		{"zero", "0o0", TokOctLit, false},
		{"one", "0o1", TokOctLit, false},
		{"two", "0o2", TokOctLit, false},
		{"three", "0o3", TokOctLit, false},
		{"four", "0o4", TokOctLit, false},
		{"five", "0o5", TokOctLit, false},
		{"six", "0o6", TokOctLit, false},
		{"seven", "0o7", TokOctLit, false},

		// Invalid digits
		{"invalid_8", "0o8", TokOctLit, true},
		{"invalid_9", "0o9", TokOctLit, true},
		{"invalid_a", "0oa", TokOctLit, true},
		{"invalid_b", "0ob", TokOctLit, true},
		{"invalid_c", "0oc", TokOctLit, true},
		{"invalid_d", "0od", TokOctLit, true},
		{"invalid_e", "0oe", TokOctLit, true},
		{"invalid_f", "0of", TokOctLit, true},
		{"invalid_z", "0oz", TokOctLit, true},

		// Complex numbers
		{"all_octal_digits", "0o01234567", TokOctLit, false},
		{"single_underscore", "0o0123_4567", TokOctLit, false},
		{"multiple_underscores", "0o01_23_45_67", TokOctLit, false},
		{"max_underscores", "0o0_1_2_3_4_5_6_7", TokOctLit, false},

		// Invalid underscore positions
		{"invalid_leading_underscore", "0o_", TokOctLit, true},
		{"invalid_underscore_after_prefix", "0o_0", TokOctLit, true},
		{"valid_trailing_underscore", "0o1_", TokOctLit, false},
		{"valid_double_underscore", "0o0__1", TokOctLit, false},
		{"valid_middle_underscore", "0o0_1_", TokOctLit, false},

		// Invalid prefix cases
		{"invalid_uppercase_prefix", "0O0", TokDecLit, true},
		{"empty_after_prefix", "0o", TokOctLit, true},

		// Invalid suffix cases
		{"invalid_e_suffix", "0o1e", TokOctLit, true},
		{"invalid_e_suffix_with_number", "0o1e0", TokOctLit, true},

		// Sequence cases
		{"octal_then_comma", "0o1,", TokOctLit, false},       // Should parse as separate tokens
		{"underscore_then_comma", "0o1_,", TokOctLit, false}, // Should parse as separate tokens
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
		{"invalid_uppercase_prefix", "0X0", TokDecLit, true},
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
		{"empty", "``", TokRawStrLit, false},
		{"simple", "`abc`", TokRawStrLit, false},
		{"spaces", "`   `", TokRawStrLit, false},
		{"tabs", "`\t\t`", TokRawStrLit, false},

		// Escaped backtick cases
		{"embedded_backtick", "`abc``def`", TokRawStrLit, false},
		{"multiple_escaped_backticks", "`a``b``c`", TokRawStrLit, false},
		{"double_backtick", "`a``b`", TokRawStrLit, false},
		{"escaped_backtick_with_space", "`a`` b`", TokRawStrLit, false},
		{"escaped_backtick_with_newline", "`a``\nb`", TokRawStrLit, false},

		// Newline cases
		{"single_newline", "`abc\ndef`", TokRawStrLit, false},
		{"multiple_newlines", "`first\nsecond\nthird`", TokRawStrLit, false},
		{"only_newlines", "`\n\n\n`", TokRawStrLit, false},
		{"newline_after_backtick", "`abc``\ndef`", TokRawStrLit, false},
		{"newline_before_backtick", "`abc\n``def`", TokRawStrLit, false},
		{"starts_with_newline", "`\nabc`", TokRawStrLit, false},
		{"ends_with_newline", "`abc\n`", TokRawStrLit, false},

		// Special character cases
		{"special_chars", "`!@#$%^&*()`", TokRawStrLit, false},
		{"unicode_chars", "`αβγδε`", TokRawStrLit, false},
		{"emoji", "`🙂🌟🎉`", TokRawStrLit, false},
		{"escape_sequences_raw", "`\\n\\t\\r`", TokRawStrLit, false},
		{"quotes_in_raw", "`\"single\" and 'double'`", TokRawStrLit, false},
		{"brackets_braces", "`[{(<>)}]`", TokRawStrLit, false},

		// Whitespace cases
		{"mixed_whitespace", "`space\ttab\nline`", TokRawStrLit, false},
		{"carriage_return", "`line1\rline2`", TokRawStrLit, false},
		{"all_whitespace_types", "`\n\t\r \f\v`", TokRawStrLit, false},

		// Invalid cases
		{"unterminated_after_backtick", "`abc``", TokRawStrLit, true},
		{"unterminated", "`abc", TokRawStrLit, true},
		{"unterminated_with_newlines", "`first\nsecond\nthird", TokRawStrLit, true},
		{"unterminated_after_escape", "`abc``def", TokRawStrLit, true},
		{"unterminated_unicode", "`αβγ", TokRawStrLit, true},
		{"unterminated_emoji", "`🙂🌟", TokRawStrLit, true},
		{"empty_unterminated", "`", TokRawStrLit, true},
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
		{"empty", `""`, TokStrLit, false},
		{"simple", `"abc"`, TokStrLit, false},
		{"spaces", `"   "`, TokStrLit, false},
		{"tabs", `"\t\t"`, TokStrLit, false},

		// Simple escape sequences
		{"newline", `"\n"`, TokStrLit, false},
		{"tab", `"\t"`, TokStrLit, false},
		{"double_quote", `"\""`, TokStrLit, false},
		{"single_quote", `"\'"`, TokStrLit, false},
		{"backslash", `"\\"`, TokStrLit, false},
		{"carriage_return", `"\r"`, TokStrLit, false},
		{"backspace", `"\b"`, TokStrLit, false},
		{"form_feed", `"\f"`, TokStrLit, false},
		{"vertical_tab", `"\v"`, TokStrLit, false},
		{"null", `"\0"`, TokStrLit, false},
		{"all_escapes", `"\n\t\"\'\\\r\b\f\v\0"`, TokStrLit, false},

		// Unicode escapes
		{"unicode_4_digit", `"\u1234"`, TokStrLit, false},
		{"unicode_8_digit", `"\U12345678"`, TokStrLit, false},
		{"unicode_max", `"\uFFFF"`, TokStrLit, false},
		{"unicode_min", `"\u0000"`, TokStrLit, false},
		{"unicode_multiple", `"\u1234\u5678"`, TokStrLit, false},
		{"unicode_mixed_case", `"\uAbCd"`, TokStrLit, false},

		// Hex escapes
		{"hex_byte", `"\xEF"`, TokStrLit, false},
		{"hex_multiple", `"\xEF\xBB\xBF"`, TokStrLit, false},
		{"hex_min", `"\x00"`, TokStrLit, false},
		{"hex_max", `"\xFF"`, TokStrLit, false},
		{"hex_mixed_case", `"\xaB"`, TokStrLit, false},

		// Mixed content
		{"mixed_escapes", `"Hello\n\tWorld\r\n"`, TokStrLit, false},
		{"mixed_unicode_hex", `"\u1234\xFF"`, TokStrLit, false},
		{"mixed_all", `"Hello\n\t\u1234\xFF\0World"`, TokStrLit, false},

		// Invalid cases - Unicode
		{"invalid_unicode_empty", `"\u"`, TokStrLit, true},
		{"invalid_unicode_short", `"\u123"`, TokStrLit, true},
		{"invalid_unicode_long", `"\U1234"`, TokStrLit, true},
		{"invalid_unicode_letters", `"\uXYZ"`, TokStrLit, true},
		{"invalid_unicode_partial", `"\u12G4"`, TokStrLit, true},
		{"invalid_unicode_space", `"\u 123"`, TokStrLit, true},

		// Invalid cases - Hex
		{"invalid_hex_empty", `"\x"`, TokStrLit, true},
		{"invalid_hex_short", `"\xF"`, TokStrLit, true},
		{"invalid_hex_letters", `"\xXY"`, TokStrLit, true},
		{"invalid_hex_space", `"\x F"`, TokStrLit, true},

		// Invalid cases - General
		{"unterminated", `"abc`, TokStrLit, true},
		{"unterminated_escape", `"abc\`, TokStrLit, true},
		{"invalid_escape", `"\k"`, TokStrLit, true},
		{"invalid_escape_exclamation", `"\!"`, TokStrLit, true},
		{"bare_backslash", `"\"`, TokStrLit, true},
		{"newline_in_string", "\"abc\ndef\"", TokStrLit, true},
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
		{"empty", `"Hello \()!"`, TokInterpStrLit, false},
		{"simple_identifier", `"Hello \(name)!"`, TokInterpStrLit, false},
		{"simple_number", `"Count: \(123)!"`, TokInterpStrLit, false},
		{"simple_string", `"Name: \("Alice")!"`, TokInterpStrLit, false},
		{"simple_char", `"Initial: \('A')!"`, TokInterpStrLit, false},

		// Multiple interpolations
		{"two_interpolations", `"Hello \(first) \(last)!"`, TokInterpStrLit, false},
		{"three_interpolations", `"Hello \(title) \(first) \(last)!"`, TokInterpStrLit, false},
		{"adjacent_interpolations", `"\(a)\(b)\(c)"`, TokInterpStrLit, false},

		// Complex expressions
		{"expression_concat", `"Hello \(first + " " + last)!"`, TokInterpStrLit, false},
		{"expression_math", `"Sum: \(a + b + c)!"`, TokInterpStrLit, false},
		{"expression_nested_parens", `"Result: \((a + b) * (c + d))!"`, TokInterpStrLit, false},
		{"expression_function_call", `"Length: \(string.len(name))!"`, TokInterpStrLit, false},

		// Escapes and special characters
		{"escaped_quotes", `"Quote: \("\"nested\"")"`, TokInterpStrLit, false},
		{"escaped_backslash", `"Path: \("C:\\Program Files")"`, TokInterpStrLit, false},
		{"newline_escape", `"Lines: \("first\nsecond")"`, TokInterpStrLit, false},
		{"mixed_escapes", `"Mixed: \("tab\t\"quote\"\nline")"`, TokInterpStrLit, false},

		// Non-interpolation parentheses
		{"non_interpolation_parens", `"Welcome \(opponent1) )))---((( \(opponent2)"`, TokInterpStrLit, false},

		// Invalid cases
		{"unterminated_string", `"Hello \(name`, TokInterpStrLit, true},
		{"unterminated_interpolation", `"Hello \(name"`, TokInterpStrLit, true},
		{"unmatched_parens", `"Hello \((name)!"`, TokInterpStrLit, true},
		{"raw_newline", `"Hello \(name)`, TokInterpStrLit, true},
		{"nested_string_newline", `"Hello \("name`, TokInterpStrLit, true},
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
		{"empty_string", "```\n\n```", TokMultiStrLit, false},
		{"simple_text", "```\nSome text\n```", TokMultiStrLit, false},
		{"with_processor", "```processor\nSome text\n```", TokMultiStrLit, false},
		{"consistent_indentation", "```\n  First line\n  Second line\n  ```", TokMultiStrLit, false},

		// Escape sequences and special characters
		{"escape_sequences", "```\nSpecial: \\n\\t\\\"\\'\n```", TokMultiStrLit, false},
		{"interpolation", "```\nHello \\(name)!\n```", TokMultiStrLit, false},
		{"byte_escapes", "```\nByte: \\x48\\x69\n```", TokMultiStrLit, false},
		{"unicode_escapes", "```\nUnicode: \\u2603 \\U0001F680\n```", TokMultiStrLit, false},
		{"mixed_escapes", "```\nMixed: \\n\\t\\x48\\u2603\n```", TokMultiStrLit, false},

		// Invalid cases
		{"unclosed_string", "```\nUnclosed string", TokMultiStrLit, true},
		{"missing_newline_after_open", "```Some text\n```", TokMultiStrLit, true},
		{"missing_newline_before_close", "```\nSome text```", TokMultiStrLit, true},
		{"invalid_escape", "```\nInvalid escape: \\z\n```", TokMultiStrLit, true},
		{"incomplete_unicode", "```\nBad unicode: \\u26\n```", TokMultiStrLit, true},
		{"unterminated_interpolation", "```\nUnterminated: \\(expr\n```", TokMultiStrLit, true},
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
		TokTypeID,       // List
		TokOpenBracket,  // [
		TokID,           // a
		TokCloseBracket, // ]
		TokOpAssign,     // =
		TokKwUnion,      // union
		TokOpenParen,    // (
		TokEOL,          // \n
		TokTypeID,       // Nil
		TokEOL,          // \n
		TokTypeID,       // Cons
		TokOpenBracket,  // [
		TokID,           // a
		TokCloseBracket, // ]
		TokEOL,          // \n
		TokCloseParen,   // )
		TokEOL,          // \n
		TokEOF,
	})
}

func TestComment(t *testing.T) {
	const source = "a = 5 # assign 5 to a\nb = 10\n"
	testTokenizeSeq(t, source, []TokenType{
		TokID,       // a
		TokOpAssign, // =
		TokDecLit,   // 5
		TokComment,  // #
		TokID,       // b
		TokOpAssign, // =
		TokDecLit,   // 10
		TokEOL,      // \n
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
	testTokenizeInvalid(t, "0. ", TokFloatLit)     // decimal not followed by a digit or another dot
	testTokenizeInvalid(t, "1. ", TokFloatLit)     // decimal not followed by a digit or another dot
	testTokenizeInvalid(t, "12. ", TokFloatLit)    // decimal not followed by a digit or another dot
}

func TestDecimalMemberAccess(t *testing.T) {
	const source = "123.string()"
	testTokenizeSeq(t, source, []TokenType{
		TokDecLit,     // 123
		TokDot,        // .
		TokID,         // string
		TokOpenParen,  // (
		TokCloseParen, // )
		TokEOF,
	})
}

func TestBinaryMemberAccess(t *testing.T) {
	const source = "0b1010.string()"
	testTokenizeSeq(t, source, []TokenType{
		TokBinLit,     // 0b1010
		TokDot,        // .
		TokID,         // string
		TokOpenParen,  // (
		TokCloseParen, // )
		TokEOF,
	})
}

func TestOctalMemberAccess(t *testing.T) {
	const source = "0o123.string()"
	testTokenizeSeq(t, source, []TokenType{
		TokOctLit,     // 0o123
		TokDot,        // .
		TokID,         // string
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
		TokID,         // string
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
		TokID,         // string
		TokOpenParen,  // (
		TokCloseParen, // )
		TokEOF,
	})

	const source2 = "123.456.string()"
	testTokenizeSeq(t, source2, []TokenType{
		TokFloatLit,   // 123.456
		TokDot,        // .
		TokID,         // string
		TokOpenParen,  // (
		TokCloseParen, // )
		TokEOF,
	})
}

func TestRangeOp(t *testing.T) {
	const source = "1..2"
	testTokenizeSeq(t, source, []TokenType{
		TokDecLit,  // 1
		TokOpRange, // ..
		TokDecLit,  // 2
		TokEOF,
	})
}

func TestRestOp(t *testing.T) {
	const source = "...rest"
	testTokenizeSeq(t, source, []TokenType{
		TokOpRest, // ...
		TokID,     // rest
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
		{"uppercase_A", ":A", TokSymLit, false, false, nil},
		{"uppercase_Z", ":Z", TokSymLit, false, false, nil},
		{"lowercase_a", ":a", TokSymLit, false, false, nil},
		{"lowercase_z", ":z", TokSymLit, false, false, nil},
		{"all_uppercase", ":ABCDEFGHIJKLMNOPQRSTUVWXYZ", TokSymLit, false, false, nil},
		{"all_lowercase", ":abcdefghijklmnopqrstuvwxyz", TokSymLit, false, false, nil},
		{"letters and digits", ":a100", TokSymLit, false, false, nil},
		{"letters and digits and underscore", ":a100_", TokSymLit, false, false, nil},

		// Valid quoted symbol literal
		{"quoted_symbol", `:"anything but a newline"`, TokSymLit, false, false, nil},

		// Invalid symbol literals
		{"invalid_digit_0", ":0", TokSymLit, true, false, nil},
		{"invalid_digit_1", ":1", TokSymLit, true, false, nil},
		{"invalid_digit_9", ":9", TokSymLit, true, false, nil},
		{"unterminated_quoted", `:"this symbol does not end`, TokSymLit, true, false, nil},
		{"newline_in_quoted", ":\"no\nnewlines!", TokSymLit, true, false, nil},

		// Symbol in sequence
		{"symbol_in_assignment", "foo = :foo", TokSymLit, false, true, []TokenType{
			TokID,       // foo
			TokOpAssign, // =
			TokSymLit,   // :foo
			TokEOF,
		}},

		{"symbol_in_comparison", `:foo == Symbol("foo")`, TokSymLit, false, true, []TokenType{
			TokSymLit,     // :foo
			TokOpEQ,       // ==
			TokTypeID,     // Symbol
			TokOpenParen,  // (
			TokStrLit,     // "foo"
			TokCloseParen, // )
			TokEOF,
		}},

		// Invalid symbol in sequence
		{"invalid_symbol_in_seq", `foo = :"foo`, TokSymLit, true, true, []TokenType{
			TokID,       // foo
			TokOpAssign, // =
			TokSymLit,   // :"foo
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

// TestErrorOffset specifically tests that the error offset is correctly set
func TestErrorOffset(t *testing.T) {
	tests := []struct {
		name              string
		input             string
		wantType          TokenType
		expectedErrorChar byte // The character where the error is expected
	}{
		{"invalid_hex", "0xG", TokHexLit, 'G'},
		{"invalid_bin", "0b2", TokBinLit, '2'},
		{"invalid_oct", "0o8", TokOctLit, '8'},
		{"invalid_dec_with_letter", "123a", TokDecLit, 'a'},
		{"invalid_exp", "1e+a", TokFloatLit, 'a'},
		{"unterminated_string", "\"abc", TokStrLit, 0}, // 0 represents EOF
		{"invalid_escape", "\"\\k\"", TokStrLit, 'k'},
		{"invalid_hex_escape", "\"\\x1g\"", TokStrLit, 'g'},
		{"invalid_unicode_escape", "\"\\u12zg\"", TokStrLit, 'z'},
		{"symbol_starting_with_digit", ":123", TokSymLit, '1'},
		// These are treated as decimal literals followed by invalid characters
		{"uppercase_hex_prefix", "0X1", TokDecLit, 'X'},
		{"uppercase_bin_prefix", "0B1", TokDecLit, 'B'},
		{"uppercase_oct_prefix", "0O7", TokDecLit, 'O'},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenizer := NewTokenizer([]byte(tt.input), "test.go")
			token := tokenizer.Next()

			// Basic checks
			if token.Type != tt.wantType {
				t.Errorf("Expected token type %v, got %v", TokenTypes[tt.wantType], TokenTypes[token.Type])
			}
			if !token.Invalid {
				t.Errorf("Expected invalid token for %q", tt.input)
			}

			// Error offset checks
			if token.ErrorOffset == 0 {
				t.Errorf("Expected error offset to be set for invalid token %q", tt.input)
			}

			// Check that the error offset points to the correct character
			if tt.expectedErrorChar != 0 {
				// For specific character errors
				index := token.ErrorOffset
				if index >= int32(len(tt.input)) {
					t.Errorf("Error offset %d is out of bounds for input %q", index, tt.input)
				} else if tt.input[index] != tt.expectedErrorChar {
					t.Errorf("Expected error at character %q, but error offset %d points to %q",
						tt.expectedErrorChar, index, tt.input[index])
				}
			} else {
				// For EOF errors, error offset should point to the end
				if token.ErrorOffset != int32(len(tt.input)) {
					t.Errorf("Expected error at EOF, but error offset is %d", token.ErrorOffset)
				}
			}

			// Check the error position methods
			errorLine, errorCol := token.ErrorPosition()
			t.Logf("Error position for %q: line %d, col %d", tt.input, errorLine, errorCol)
		})
	}
}
