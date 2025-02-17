package main

import (
	"strings"
	"testing"
)

// Helper: tokenize a sequence of tokens and check types.
func testTokenizeSeq(t *testing.T, source string, expected []TokenType) {
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
	lastToken := tokenizer.Next()
	if lastToken.Type != TokenEOF {
		t.Errorf("Expected EOF token, got %v", lastToken.Type)
	}
	expectedCol := len(source) + 1
	if lastToken.Column != expectedCol {
		t.Errorf("Expected column %d, got %d", expectedCol, lastToken.Column)
	}
}

// Helper for sequences where the first token may be marked invalid.
func testTokenizeSeqInvalid(t *testing.T, source string, expected []TokenType, invalid bool) {
	tokenizer := NewTokenizer([]byte(source), "test.go")
	for i, exp := range expected {
		token := tokenizer.Next()
		if token.Type != exp {
			t.Errorf("At index %d: expected token type %v, got %v", i, exp, token.Type)
		}
		if i == 0 && token.Invalid != invalid {
			t.Errorf("At index %d: expected invalid=%v, got %v", i, invalid, token.Invalid)
		}
	}
	lastToken := tokenizer.Next()
	if lastToken.Type != TokenEOF {
		t.Errorf("Expected EOF token, got %v", lastToken.Type)
	}
	expectedCol := len(source) + 1
	if lastToken.Column != expectedCol {
		t.Errorf("Expected column %d, got %d", expectedCol, lastToken.Column)
	}
}

func testTokenize(t *testing.T, source string, expectedType TokenType) {
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
	expectedCol := len(source) + 1
	if lastToken.Column != expectedCol {
		t.Errorf("Expected column %d, got %d", expectedCol, lastToken.Column)
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
	testTokenizeSeq(t, "&= |= /= = &&= ||= -= %= *= += ^=", []TokenType{
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
	testTokenizeInvalid(t, "0b", TokenBinaryLiteral)
	testTokenize(t, "0b0", TokenBinaryLiteral)
	testTokenize(t, "0b1", TokenBinaryLiteral)
	testTokenize(t, "0b10101100", TokenBinaryLiteral)

	testTokenizeInvalid(t, "0b2", TokenBinaryLiteral)
	testTokenizeInvalid(t, "0b3", TokenBinaryLiteral)
	testTokenizeInvalid(t, "0b4", TokenBinaryLiteral)
	testTokenizeInvalid(t, "0b5", TokenBinaryLiteral)
	testTokenizeInvalid(t, "0b6", TokenBinaryLiteral)
	testTokenizeInvalid(t, "0b7", TokenBinaryLiteral)
	testTokenizeInvalid(t, "0b8", TokenBinaryLiteral)
	testTokenizeInvalid(t, "0b9", TokenBinaryLiteral)
	testTokenizeInvalid(t, "0ba", TokenBinaryLiteral)
	testTokenizeInvalid(t, "0bb", TokenBinaryLiteral)
	testTokenizeInvalid(t, "0bc", TokenBinaryLiteral)
	testTokenizeInvalid(t, "0bd", TokenBinaryLiteral)
	testTokenizeInvalid(t, "0be", TokenBinaryLiteral)
	testTokenizeInvalid(t, "0bf", TokenBinaryLiteral)
	testTokenizeInvalid(t, "0bz", TokenBinaryLiteral)

	// These tests assume that uppercase "0B" is not a valid binary literal.
	testTokenizeInvalid(t, "0B0", TokenDecimalLiteral)
	testTokenizeInvalid(t, "0b_", TokenBinaryLiteral)
	testTokenizeInvalid(t, "0b_0", TokenBinaryLiteral)
	testTokenize(t, "0b1_", TokenBinaryLiteral)
	testTokenize(t, "0b0__1", TokenBinaryLiteral)
	testTokenize(t, "0b0_1_", TokenBinaryLiteral)
	testTokenizeInvalid(t, "0b1e", TokenBinaryLiteral)
}

func TestBooleanLiterals(t *testing.T) {
	testTokenize(t, "false", TokenBooleanLiteral)
	testTokenize(t, "true", TokenBooleanLiteral)
}

func TestDecimalLiterals(t *testing.T) {
	testTokenize(t, "0", TokenDecimalLiteral)
	testTokenize(t, "1", TokenDecimalLiteral)
	testTokenize(t, "2", TokenDecimalLiteral)
	testTokenize(t, "3", TokenDecimalLiteral)
	testTokenize(t, "4", TokenDecimalLiteral)
	testTokenize(t, "5", TokenDecimalLiteral)
	testTokenize(t, "6", TokenDecimalLiteral)
	testTokenize(t, "7", TokenDecimalLiteral)
	testTokenize(t, "8", TokenDecimalLiteral)
	testTokenize(t, "9", TokenDecimalLiteral)

	testTokenize(t, "0_0", TokenDecimalLiteral)
	testTokenize(t, "0001", TokenDecimalLiteral)
	testTokenize(t, "01234567890", TokenDecimalLiteral)
	testTokenize(t, "012_345_6789_0", TokenDecimalLiteral)
	testTokenize(t, "0_1_2_3_4_5_6_7_8_9_0", TokenDecimalLiteral)
}

func TestOctalLiterals(t *testing.T) {
	testTokenize(t, "0o0", TokenOctalLiteral)
	testTokenize(t, "0o1", TokenOctalLiteral)
	testTokenize(t, "0o2", TokenOctalLiteral)
	testTokenize(t, "0o3", TokenOctalLiteral)
	testTokenize(t, "0o4", TokenOctalLiteral)
	testTokenize(t, "0o5", TokenOctalLiteral)
	testTokenize(t, "0o6", TokenOctalLiteral)
	testTokenize(t, "0o7", TokenOctalLiteral)

	testTokenize(t, "0o01234567", TokenOctalLiteral)
	testTokenize(t, "0o0123_4567", TokenOctalLiteral)
	testTokenize(t, "0o01_23_45_67", TokenOctalLiteral)
	testTokenize(t, "0o0_1_2_3_4_5_6_7", TokenOctalLiteral)

	testTokenizeInvalid(t, "0O0", TokenDecimalLiteral)
	testTokenizeInvalid(t, "0o_", TokenOctalLiteral)
	testTokenizeInvalid(t, "0o_0", TokenOctalLiteral)
	testTokenize(t, "0o1_", TokenOctalLiteral)
	testTokenize(t, "0o0__1", TokenOctalLiteral)
	testTokenize(t, "0o0_1_", TokenOctalLiteral)
	testTokenizeInvalid(t, "0o1e", TokenOctalLiteral)
	testTokenizeInvalid(t, "0o1e0", TokenOctalLiteral)
	testTokenizeSeqInvalid(t, "0o_,", []TokenType{TokenOctalLiteral, TokenComma}, true)
}

func TestHexadecimalLiterals(t *testing.T) {
	testTokenize(t, "0x0", TokenHexadecimalLiteral)
	testTokenize(t, "0x1", TokenHexadecimalLiteral)
	testTokenize(t, "0x2", TokenHexadecimalLiteral)
	testTokenize(t, "0x3", TokenHexadecimalLiteral)
	testTokenize(t, "0x4", TokenHexadecimalLiteral)
	testTokenize(t, "0x5", TokenHexadecimalLiteral)
	testTokenize(t, "0x6", TokenHexadecimalLiteral)
	testTokenize(t, "0x7", TokenHexadecimalLiteral)
	testTokenize(t, "0x8", TokenHexadecimalLiteral)
	testTokenize(t, "0x9", TokenHexadecimalLiteral)
	testTokenize(t, "0xa", TokenHexadecimalLiteral)
	testTokenize(t, "0xb", TokenHexadecimalLiteral)
	testTokenize(t, "0xc", TokenHexadecimalLiteral)
	testTokenize(t, "0xd", TokenHexadecimalLiteral)
	testTokenize(t, "0xe", TokenHexadecimalLiteral)
	testTokenize(t, "0xf", TokenHexadecimalLiteral)
	testTokenize(t, "0xA", TokenHexadecimalLiteral)
	testTokenize(t, "0xB", TokenHexadecimalLiteral)
	testTokenize(t, "0xC", TokenHexadecimalLiteral)
	testTokenize(t, "0xD", TokenHexadecimalLiteral)
	testTokenize(t, "0xE", TokenHexadecimalLiteral)
	testTokenize(t, "0xF", TokenHexadecimalLiteral)

	testTokenize(t, "0x0000", TokenHexadecimalLiteral)
	testTokenize(t, "0xAA", TokenHexadecimalLiteral)
	testTokenize(t, "0xFFFF", TokenHexadecimalLiteral)

	testTokenize(t, "0x0123456789ABCDEF", TokenHexadecimalLiteral)
	testTokenize(t, "0x0123_4567_89AB_CDEF", TokenHexadecimalLiteral)
	testTokenize(t, "0x01_23_45_67_89AB_CDE_F", TokenHexadecimalLiteral)
	testTokenize(t, "0x0_1_2_3_4_5_6_7_8_9_A_B_C_D_E_F", TokenHexadecimalLiteral)

	testTokenizeInvalid(t, "0X0", TokenDecimalLiteral)
	testTokenizeInvalid(t, "0x_", TokenHexadecimalLiteral)
	testTokenizeInvalid(t, "0x_1", TokenHexadecimalLiteral)
	testTokenize(t, "0x1_", TokenHexadecimalLiteral)
	testTokenize(t, "0x0__1", TokenHexadecimalLiteral)
	testTokenize(t, "0x0_1_", TokenHexadecimalLiteral)
	testTokenizeSeqInvalid(t, "0x_,", []TokenType{TokenHexadecimalLiteral, TokenComma}, true)
}

func TestRawStringLiterals(t *testing.T) {
	testTokenize(t, "`abc`", TokenRawStringLiteral)
	testTokenize(t, "`abc``def`", TokenRawStringLiteral)
	testTokenizeInvalid(t, "`abc``", TokenRawStringLiteral)
}

func TestStringLiterals(t *testing.T) {
	testTokenize(t, `"abc"`, TokenStringLiteral)
	testTokenize(t, `"\n\t\"\'\\\r\b\f\v\0"`, TokenStringLiteral)
	testTokenize(t, `"\u1234"`, TokenStringLiteral)
	testTokenize(t, `"\xEF\xBB\xBF"`, TokenStringLiteral)

	testTokenizeInvalid(t, `"\u"`, TokenStringLiteral)
	testTokenizeInvalid(t, `"\u123"`, TokenStringLiteral)
	testTokenizeInvalid(t, `"\uXYZ"`, TokenStringLiteral)
	testTokenizeInvalid(t, `"\xXYZ"`, TokenStringLiteral)
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
