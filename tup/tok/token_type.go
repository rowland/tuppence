package tok

// TokenType represents the various kinds of tokens.
type TokenType byte

const (
	// Symbols
	TokAt           TokenType = iota // @
	TokCloseBrace                    // }
	TokCloseBracket                  // ]
	TokCloseParen                    // )
	TokColon                         // :
	TokColonNoSpace                  // : (no space after)
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

	TokOpDiv   // /
	TokOpMinus // -
	TokOpMod   // %
	TokOpMul   // *
	TokOpNot   // !
	TokOpPlus  // +
	TokOpPow   // ^
	TokOpSHL   // <<
	TokOpSHR   // >>

	// Bitwise Operators
	TokOpBitAnd // &
	TokOpBitOr  // |
	TokOpBitNot // ~

	// Relational Operators
	TokOpEQ      // ==
	TokOpGE      // >=
	TokOpGT      // >
	TokOpLE      // <=
	TokOpLT      // <
	TokOpNE      // !=
	TokOpMatch   // =~
	TokOpCompare //<=>

	// Logical Operators
	TokOpLogAnd // &&
	TokOpLogOr  // ||

	// Range and Rest Operators
	TokOpRange
	TokOpRest

	// Assignment
	TokOpBitAndEQ // &=
	TokOpBitOrEQ  // |=
	TokOpDivEQ    // /=
	TokOpAssign   // =
	TokOpLogAndEQ // &&=
	TokOpLogOrEQ  // ||=
	TokOpMinusEQ  // -=
	TokOpModEQ    // %=
	TokOpMulEQ    // *=
	TokOpPlusEQ   // +=
	TokOpPowEQ    // ^=
	TokOpSHL_EQ   // <<=
	TokOpSHR_EQ   // >>=

	// Identifiers
	TokID     // id
	TokFuncID // foo? or foo!
	TokTypeID // ID

	// Keywords
	TokKwArray       // array
	TokKwBreak       // break
	TokKwContinue    // continue
	TokKwContract    // contract
	TokKwElse        // else
	TokKwEnum        // enum
	TokKwError       // error
	TokKwFn          // fn
	TokKwFx          // fx
	TokKwFor         // for
	TokKwIf          // if
	TokKwIn          // in
	TokKwIs          // is
	TokKwIt          // it
	TokKwImport      // import
	TokKwMut         // mut
	TokKwReturn      // return
	TokKwSwitch      // switch
	TokKwTry         // try
	TokKwTryBreak    // try_break
	TokKwTryContinue // try_continue
	TokKwType        // type
	TokKwTypeof      // typeof
	TokKwUnion       // union

	// Literals
	TokBinLit       // 10101010
	TokBoolLit      // false
	TokRuneLit      // 'A'
	TokDecLit       // 1234567890
	TokFloatLit     // 123.456
	TokHexLit       // 0xDEADBEEF
	TokInterpStrLit // "a\(b)c"
	TokOctLit       // 0o777
	TokRawStrLit    // `\no\escapes`
	TokStrLit       // "Hello, World"
	TokMultiStrLit  // ```\nHello, World\n```

	// Comments
	TokComment // #

	// Special tokens
	TokEOL
	TokEOF
	TokINV
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
	TokColonNoSpace: ":",
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

	TokOpDiv:   "/",
	TokOpMinus: "-",
	TokOpMod:   "%",
	TokOpMul:   "*",
	TokOpNot:   "!",
	TokOpPlus:  "+",
	TokOpPow:   "^",
	TokOpSHL:   "<",
	TokOpSHR:   ">",

	// Bitwise Operators
	TokOpBitAnd: "&",
	TokOpBitOr:  "|",
	TokOpBitNot: "~",

	// Relational Operators
	TokOpEQ:      "==",
	TokOpGE:      ">=",
	TokOpGT:      ">",
	TokOpLE:      "<=",
	TokOpLT:      "<",
	TokOpNE:      "!=",
	TokOpMatch:   "=~",
	TokOpCompare: "<=>",

	// Logical Operators
	TokOpLogAnd: "&&",
	TokOpLogOr:  "||",

	// Range and Rest Operators
	TokOpRange: "..",
	TokOpRest:  "...",

	// Assignment
	TokOpBitAndEQ: "&=",
	TokOpBitOrEQ:  "|=",
	TokOpDivEQ:    "/=",
	TokOpAssign:   "=",
	TokOpLogAndEQ: "&&=",
	TokOpLogOrEQ:  "||=",
	TokOpMinusEQ:  "-=",
	TokOpModEQ:    "%=",
	TokOpMulEQ:    "*=",
	TokOpPlusEQ:   "+=",
	TokOpPowEQ:    "^=",
	TokOpSHL_EQ:   "<<=",
	TokOpSHR_EQ:   ">>=",

	// Identifiers
	TokID:     "identifier",
	TokFuncID: "function_identifier",
	TokTypeID: "type_identifier",

	// Keywords
	TokKwArray:       "array",
	TokKwBreak:       "break",
	TokKwContinue:    "continue",
	TokKwContract:    "contract",
	TokKwElse:        "else",
	TokKwEnum:        "enum",
	TokKwError:       "error",
	TokKwFn:          "fn",
	TokKwFx:          "fx",
	TokKwFor:         "for",
	TokKwIf:          "if",
	TokKwIn:          "in",
	TokKwIs:          "is",
	TokKwIt:          "it",
	TokKwImport:      "import",
	TokKwMut:         "mut",
	TokKwReturn:      "return",
	TokKwSwitch:      "switch",
	TokKwTry:         "try",
	TokKwTryBreak:    "try_break",
	TokKwTryContinue: "try_continue",
	TokKwType:        "type",
	TokKwTypeof:      "typeof",
	TokKwUnion:       "union",

	// Literals
	TokBinLit:       "binary_literal",
	TokBoolLit:      "bool_literal",
	TokRuneLit:      "rune_literal",
	TokDecLit:       "decimal_literal",
	TokFloatLit:     "float_literal",
	TokHexLit:       "hex_literal",
	TokInterpStrLit: "interp_string_literal",
	TokOctLit:       "octal_literal",
	TokRawStrLit:    "raw_string_literal",
	TokStrLit:       "string_literal",
	TokMultiStrLit:  "multi_line_string_literal",

	// Comments
	TokComment: "comment",

	// Special tokens
	TokEOL: "EOL",
	TokEOF: "EOF",
	TokINV: "invalid",
}
