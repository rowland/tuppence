package main

// TokenType represents the various kinds of tokens.
type TokenType int

const (
	// Symbols
	TokenAt TokenType = iota
	TokenCloseBrace
	TokenCloseBracket
	TokenCloseParen
	TokenColon
	TokenComma
	TokenDot
	TokenOpenBrace
	TokenOpenBracket
	TokenOpenParen
	TokenQuestionMark
	TokenSemiColon

	// Operators
	TokenOpCheckedAdd
	TokenOpCheckedDiv
	TokenOpCheckedMod
	TokenOpCheckedMul
	TokenOpCheckedSub
	TokenOpDiv
	TokenOpMinus
	TokenOpMod
	TokenOpMul
	TokenOpNot
	TokenOpPlus
	TokenOpPow
	TokenOpShiftLeft
	TokenOpShiftRight

	// Bitwise Operators
	TokenOpBitwiseAnd
	TokenOpBitwiseOr

	// Relational Operators
	TokenOpEqualEqual
	TokenOpGreaterEqual
	TokenOpGreaterThan
	TokenOpLessEqual
	TokenOpLessThan
	TokenOpNotEqual

	// Logical Operators
	TokenOpLogicalAnd
	TokenOpLogicalOr

	// Assignment
	TokenOpBitwiseAndEqual
	TokenOpBitwiseOrEqual
	TokenOpDivEqual
	TokenOpEqual
	TokenOpLogicalAndEqual
	TokenOpLogicalOrEqual
	TokenOpMinusEqual
	TokenOpModEqual
	TokenOpMulEqual
	TokenOpPlusEqual
	TokenOpPowEqual

	// Identifiers
	TokenIdentifier
	TokenTypeIdentifier

	// Keywords
	TokenKeywordBreak
	TokenKeywordElse
	TokenKeywordEnum
	TokenKeywordError
	TokenKeywordFn
	TokenKeywordFx
	TokenKeywordFor
	TokenKeywordIf
	TokenKeywordImport
	TokenKeywordMut
	TokenKeywordReturn
	TokenKeywordSwitch
	TokenKeywordTry
	TokenKeywordType
	TokenKeywordTypeof

	// Literals
	TokenBinaryLiteral
	TokenBooleanLiteral
	TokenCharacterLiteral
	TokenDecimalLiteral
	TokenFloatLiteral
	TokenHexadecimalLiteral
	TokenInterpolatedStringLiteral
	TokenOctalLiteral
	TokenRawStringLiteral
	TokenStringLiteral

	// Comments
	TokenComment

	// Special tokens
	TokenEOL
	TokenEOF
	TokenInvalid
)
