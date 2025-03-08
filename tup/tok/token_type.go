package tok

// TokenType represents the various kinds of tokens.
type TokenType int

const (
	// Symbols
	TokAt           TokenType = iota // @
	TokCloseBrace                    // }
	TokCloseBracket                  // ]
	TokCloseParen                    // )
	TokColon                         // :
	TokComma                         // ,
	TokDot                           // .
	TokOpenBrace                     // {
	TokOpenBracket                   // [
	TokOpenParen                     // (
	TokQuestionMark                  // ?
	TokSemiColon                     // ;

	// Operators
	TokOpCheckedAdd // ?+
	TokOpCheckedDiv // ?/
	TokOpCheckedMod // ?%
	TokOpCheckedMul // ?*
	TokOpCheckedSub // ?-

	TokOpDiv        // /
	TokOpMinus      // -
	TokOpMod        // %
	TokOpMul        // *
	TokOpNot        // !
	TokOpPlus       // +
	TokOpPow        // ^
	TokOpShiftLeft  // <<
	TokOpShiftRight // >>

	// Bitwise Operators
	TokOpBitwiseAnd // &
	TokOpBitwiseOr  // |
	TokOpBitwiseNot // ~

	// Relational Operators
	TokOpEqual        // ==
	TokOpGreaterEqual // >=
	TokOpGreaterThan  // >
	TokOpLessEqual    // <=
	TokOpLessThan     // <
	TokOpNotEqual     // !=
	TokOpMatches      // =~
	TokOpCompareTo    //<=>

	// Logical Operators
	TokOpLogicalAnd // &&
	TokOpLogicalOr  // ||

	// Range and Rest Operators
	TokOpRange
	TokOpRest

	// Assignment
	TokOpBitwiseAndEqual // &=
	TokOpBitwiseOrEqual  // |=
	TokOpDivEqual        // /=
	TokOpAssign          // =
	TokOpLogicalAndEqual // &&=
	TokOpLogicalOrEqual  // ||=
	TokOpMinusEqual      // -=
	TokOpModEqual        // %=
	TokOpMulEqual        // *=
	TokOpPlusEqual       // +=
	TokOpPowEqual        // ^=
	TokOpShiftLeftEqual  // <<=
	TokOpShiftRightEqual // >>=

	// Identifiers
	TokIdentifier     // id
	TokTypeIdentifier // ID

	// Keywords
	TokKeywordArray       // array
	TokKeywordBreak       // break
	TokKeywordContinue    // continue
	TokKeywordContract    // contract
	TokKeywordElse        // else
	TokKeywordEnum        // enum
	TokKeywordError       // error
	TokKeywordFn          // fn
	TokKeywordFx          // fx
	TokKeywordFor         // for
	TokKeywordIf          // if
	TokKeywordIn          // in
	TokKeywordIt          // it
	TokKeywordImport      // import
	TokKeywordMut         // mut
	TokKeywordReturn      // return
	TokKeywordSwitch      // switch
	TokKeywordTry         // try
	TokKeywordTryBreak    // try_break
	TokKeywordTryContinue // try_continue
	TokKeywordType        // type
	TokKeywordTypeof      // typeof
	TokKeywordUnion       // union

	// Literals
	TokBinaryLit          // 10101010
	TokBoolLit            // false
	TokRuneLit            // 'A'
	TokDecimalLit         // 1234567890
	TokFloatLit           // 123.456
	TokHexLit             // 0xDEADBEEF
	TokInterpStringLit    // "a\(b)c"
	TokOctalLit           // 0o777
	TokRawStringLit       // `\no\escapes`
	TokStringLit          // "Hello, World"
	TokMultiLineStringLit // ```\nHello, World\n```
	TokSymbolLit          // :ok

	// Comments
	TokComment // #

	// Special tokens
	TokEOL
	TokEOF
	TokInvalid
)

func (t TokenType) String() string {
	return TokenTypes[t]
}

var TokenTypes = map[TokenType]string{
	// Symbols
	TokAt:           "@",
	TokCloseBrace:   "}",
	TokCloseBracket: "]",
	TokCloseParen:   ")",
	TokColon:        ":",
	TokComma:        ",",
	TokDot:          ".",
	TokOpenBrace:    "{",
	TokOpenBracket:  "[",
	TokOpenParen:    "(",
	TokQuestionMark: "?",
	TokSemiColon:    ";",

	// Operators
	TokOpCheckedAdd: "?+",
	TokOpCheckedDiv: "?/",
	TokOpCheckedMod: "?%",
	TokOpCheckedMul: "?*",
	TokOpCheckedSub: "?-",

	TokOpDiv:        "/",
	TokOpMinus:      "-",
	TokOpMod:        "%",
	TokOpMul:        "*",
	TokOpNot:        "!",
	TokOpPlus:       "+",
	TokOpPow:        "^",
	TokOpShiftLeft:  "<",
	TokOpShiftRight: ">",

	// Bitwise Operators
	TokOpBitwiseAnd: "&",
	TokOpBitwiseOr:  "|",
	TokOpBitwiseNot: "~",

	// Relational Operators
	TokOpEqual:        "==",
	TokOpGreaterEqual: ">=",
	TokOpGreaterThan:  ">",
	TokOpLessEqual:    "<=",
	TokOpLessThan:     "<",
	TokOpNotEqual:     "!=",
	TokOpMatches:      "=~",
	TokOpCompareTo:    "<=>",

	// Logical Operators
	TokOpLogicalAnd: "&&",
	TokOpLogicalOr:  "||",

	// Range and Rest Operators
	TokOpRange: "..",
	TokOpRest:  "...",

	// Assignment
	TokOpBitwiseAndEqual: "&=",
	TokOpBitwiseOrEqual:  "|=",
	TokOpDivEqual:        "/=",
	TokOpAssign:          "=",
	TokOpLogicalAndEqual: "&&=",
	TokOpLogicalOrEqual:  "||=",
	TokOpMinusEqual:      "-=",
	TokOpModEqual:        "%=",
	TokOpMulEqual:        "*=",
	TokOpPlusEqual:       "+=",
	TokOpPowEqual:        "^=",
	TokOpShiftLeftEqual:  "<<=",
	TokOpShiftRightEqual: ">>=",

	// Identifiers
	TokIdentifier:     "identifier",
	TokTypeIdentifier: "type_identifier",

	// Keywords
	TokKeywordArray:       "array",
	TokKeywordBreak:       "break",
	TokKeywordContinue:    "continue",
	TokKeywordContract:    "contract",
	TokKeywordElse:        "else",
	TokKeywordEnum:        "enum",
	TokKeywordError:       "error",
	TokKeywordFn:          "fn",
	TokKeywordFx:          "fx",
	TokKeywordFor:         "for",
	TokKeywordIf:          "if",
	TokKeywordIn:          "in",
	TokKeywordIt:          "it",
	TokKeywordImport:      "import",
	TokKeywordMut:         "mut",
	TokKeywordReturn:      "return",
	TokKeywordSwitch:      "switch",
	TokKeywordTry:         "try",
	TokKeywordTryBreak:    "try_break",
	TokKeywordTryContinue: "try_continue",
	TokKeywordType:        "type",
	TokKeywordTypeof:      "typeof",
	TokKeywordUnion:       "union",

	// Literals
	TokBinaryLit:          "binary_literal",
	TokBoolLit:            "bool_literal",
	TokRuneLit:            "rune_literal",
	TokDecimalLit:         "decimal_literal",
	TokFloatLit:           "float_literal",
	TokHexLit:             "hex_literal",
	TokInterpStringLit:    "interp_string_literal",
	TokOctalLit:           "octal_literal",
	TokRawStringLit:       "raw_string_literal",
	TokStringLit:          "string_literal",
	TokMultiLineStringLit: "multi_line_string_literal",
	TokSymbolLit:          "symbol_literal",

	// Comments
	TokComment: "comment",

	// Special tokens
	TokEOL:     "EOL",
	TokEOF:     "EOF",
	TokInvalid: "invalid",
}
