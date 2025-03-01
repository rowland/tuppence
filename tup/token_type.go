package main

// TokenType represents the various kinds of tokens.
type TokenType int

const (
	// Symbols
	TokenAt           TokenType = iota // @
	TokenCloseBrace                    // }
	TokenCloseBracket                  // ]
	TokenCloseParen                    // )
	TokenColon                         // :
	TokenComma                         // ,
	TokenDot                           // .
	TokenOpenBrace                     // {
	TokenOpenBracket                   // [
	TokenOpenParen                     // (
	TokenQuestionMark                  // ?
	TokenSemiColon                     // ;

	// Operators
	TokenOpCheckedAdd // ?+
	TokenOpCheckedDiv // ?/
	TokenOpCheckedMod // ?%
	TokenOpCheckedMul // ?*
	TokenOpCheckedSub // ?-

	TokenOpDiv        // /
	TokenOpMinus      // -
	TokenOpMod        // %
	TokenOpMul        // *
	TokenOpNot        // !
	TokenOpPlus       // +
	TokenOpPow        // ^
	TokenOpShiftLeft  // <<
	TokenOpShiftRight // >>

	// Bitwise Operators
	TokenOpBitwiseAnd // &
	TokenOpBitwiseOr  // |
	TokenOpBitwiseNot // ~

	// Relational Operators
	TokenOpEqualEqual   // ==
	TokenOpGreaterEqual // >=
	TokenOpGreaterThan  // >
	TokenOpLessEqual    // <=
	TokenOpLessThan     // <
	TokenOpNotEqual     // !=
	TokenOpMatches      // =~
	TokenOpCompareTo    //<=>

	// Logical Operators
	TokenOpLogicalAnd // &&
	TokenOpLogicalOr  // ||

	// Range and Rest Operators
	TokenOpRange
	TokenOpRest

	// Assignment
	TokenOpBitwiseAndEqual // &=
	TokenOpBitwiseOrEqual  // |=
	TokenOpDivEqual        // /=
	TokenOpEqual           // =
	TokenOpLogicalAndEqual // &&=
	TokenOpLogicalOrEqual  // ||=
	TokenOpMinusEqual      // -=
	TokenOpModEqual        // %=
	TokenOpMulEqual        // *=
	TokenOpPlusEqual       // +=
	TokenOpPowEqual        // ^=
	TokenOpShiftLeftEqual  // <<=
	TokenOpShiftRightEqual // >>=

	// Identifiers
	TokenIdentifier     // id
	TokenTypeIdentifier // ID

	// Keywords
	TokenKeywordArray       // array
	TokenKeywordBreak       // break
	TokenKeywordContinue    // continue
	TokenKeywordContract    // contract
	TokenKeywordElse        // else
	TokenKeywordEnum        // enum
	TokenKeywordError       // error
	TokenKeywordFn          // fn
	TokenKeywordFx          // fx
	TokenKeywordFor         // for
	TokenKeywordIf          // if
	TokenKeywordIn          // in
	TokenKeywordIt          // it
	TokenKeywordImport      // import
	TokenKeywordMut         // mut
	TokenKeywordReturn      // return
	TokenKeywordSwitch      // switch
	TokenKeywordTry         // try
	TokenKeywordTryBreak    // try_break
	TokenKeywordTryContinue // try_continue
	TokenKeywordType        // type
	TokenKeywordTypeof      // typeof
	TokenKeywordUnion       // union

	// Literals
	TokenBinaryLiteral             // 10101010
	TokenBooleanLiteral            // false
	TokenCharacterLiteral          // 'A'
	TokenDecimalLiteral            // 1234567890
	TokenFloatLiteral              // 123.456
	TokenHexadecimalLiteral        // 0xDEADBEEF
	TokenInterpolatedStringLiteral // "a\(b)c"
	TokenOctalLiteral              // 0o777
	TokenRawStringLiteral          // `\no\escapes`
	TokenStringLiteral             // "Hello, World"
	TokenMultiLineStringLiteral    // ```\nHello, World\n```
	TokenSymbolLiteral             // :ok

	// Comments
	TokenComment // #

	// Special tokens
	TokenEOL
	TokenEOF
	TokenInvalid
)

func (t TokenType) String() string {
	return TokenTypes[t]
}

var TokenTypes = map[TokenType]string{
	// Symbols
	TokenAt:           "@",
	TokenCloseBrace:   "}",
	TokenCloseBracket: "]",
	TokenCloseParen:   ")",
	TokenColon:        ":",
	TokenComma:        ",",
	TokenDot:          ".",
	TokenOpenBrace:    "{",
	TokenOpenBracket:  "[",
	TokenOpenParen:    "(",
	TokenQuestionMark: "?",
	TokenSemiColon:    ";",

	// Operators
	TokenOpCheckedAdd: "?+",
	TokenOpCheckedDiv: "?/",
	TokenOpCheckedMod: "?%",
	TokenOpCheckedMul: "?*",
	TokenOpCheckedSub: "?-",

	TokenOpDiv:        "/",
	TokenOpMinus:      "-",
	TokenOpMod:        "%",
	TokenOpMul:        "*",
	TokenOpNot:        "!",
	TokenOpPlus:       "+",
	TokenOpPow:        "^",
	TokenOpShiftLeft:  "<",
	TokenOpShiftRight: ">",

	// Bitwise Operators
	TokenOpBitwiseAnd: "&",
	TokenOpBitwiseOr:  "|",
	TokenOpBitwiseNot: "~",

	// Relational Operators
	TokenOpEqualEqual:   "==",
	TokenOpGreaterEqual: ">=",
	TokenOpGreaterThan:  ">",
	TokenOpLessEqual:    "<=",
	TokenOpLessThan:     "<",
	TokenOpNotEqual:     "!=",
	TokenOpMatches:      "=~",
	TokenOpCompareTo:    "<=>",

	// Logical Operators
	TokenOpLogicalAnd: "&&",
	TokenOpLogicalOr:  "||",

	// Range and Rest Operators
	TokenOpRange: "..",
	TokenOpRest:  "...",

	// Assignment
	TokenOpBitwiseAndEqual: "&=",
	TokenOpBitwiseOrEqual:  "|=",
	TokenOpDivEqual:        "/=",
	TokenOpEqual:           "=",
	TokenOpLogicalAndEqual: "&&=",
	TokenOpLogicalOrEqual:  "||=",
	TokenOpMinusEqual:      "-=",
	TokenOpModEqual:        "%=",
	TokenOpMulEqual:        "*=",
	TokenOpPlusEqual:       "+=",
	TokenOpPowEqual:        "^=",
	TokenOpShiftLeftEqual:  "<<=",
	TokenOpShiftRightEqual: ">>=",

	// Identifiers
	TokenIdentifier:     "identifier",
	TokenTypeIdentifier: "TypeIdentifier",

	// Keywords
	TokenKeywordArray:       "array",
	TokenKeywordBreak:       "break",
	TokenKeywordContinue:    "continue",
	TokenKeywordContract:    "contract",
	TokenKeywordElse:        "else",
	TokenKeywordEnum:        "enum",
	TokenKeywordError:       "error",
	TokenKeywordFn:          "fn",
	TokenKeywordFx:          "fx",
	TokenKeywordFor:         "for",
	TokenKeywordIf:          "if",
	TokenKeywordIn:          "in",
	TokenKeywordIt:          "it",
	TokenKeywordImport:      "import",
	TokenKeywordMut:         "mut",
	TokenKeywordReturn:      "return",
	TokenKeywordSwitch:      "switch",
	TokenKeywordTry:         "try",
	TokenKeywordTryBreak:    "try_break",
	TokenKeywordTryContinue: "try_continue",
	TokenKeywordType:        "type",
	TokenKeywordTypeof:      "typeof",
	TokenKeywordUnion:       "union",

	// Literals
	TokenBinaryLiteral:             "binary literal",
	TokenBooleanLiteral:            "boolean literal",
	TokenCharacterLiteral:          "character literal",
	TokenDecimalLiteral:            "decimal literal",
	TokenFloatLiteral:              "float literal",
	TokenHexadecimalLiteral:        "hexadecimal literal",
	TokenInterpolatedStringLiteral: "interpolated string literal",
	TokenOctalLiteral:              "octal literal",
	TokenRawStringLiteral:          "raw string literal",
	TokenStringLiteral:             "string literal",
	TokenMultiLineStringLiteral:    "multi-line string literal",
	TokenSymbolLiteral:             "symbol literal",

	// Comments
	TokenComment: "#",

	// Special tokens
	TokenEOL:     "EOL",
	TokenEOF:     "EOF",
	TokenInvalid: "invalid",
}
