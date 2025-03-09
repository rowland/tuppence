package tok

import "bytes"

// We define a state machine for tokenizing.
type state int

const (
	stateStart state = iota
	stateID
	stateNum
	stateInt
	stateIntDot
	stateFloat
	stateExp
	stateExpSign
	stateExpInt
	stateBinFirst
	stateBin
	stateHexFirst
	stateHex
	stateOctFirst
	stateOct
	stateRawStrLit
	stateStrLit
	stateEscSeq
	stateHexEsc
	stateComment
	stateColon
	stateSym
	stateQuotedSym
	stateMultiStrBody
)

// Tokenizer holds the state of the lexer.
type Tokenizer struct {
	file   *Source
	source []byte
	index  int
	states []state
}

var bom = []byte{0xEF, 0xBB, 0xBF} // UTF-8 BOM

// NewTokenizer initializes a new Tokenizer.
func NewTokenizer(source []byte, filename string) *Tokenizer {
	file := NewSource(source, filename)
	idx := 0
	// Skip the UTF-8 BOM if present.
	if bytes.Equal(source[:3], bom) {
		idx = 3
	}
	return &Tokenizer{
		file:   file,
		source: source,
		index:  idx,
	}
}

// pushState pushes the current state onto the stack.
func (t *Tokenizer) pushState(s state) {
	t.states = append(t.states, s)
}

// popState pops the last state from the stack.
func (t *Tokenizer) popState() state {
	s := t.states[len(t.states)-1]
	t.states = t.states[:len(t.states)-1]
	return s
}

// peek returns the next n bytes from the current position without advancing the index.
// If there aren't enough bytes remaining, it returns an empty string.
func (t *Tokenizer) peek(n int) string {
	if t.index+n > len(t.source) {
		return ""
	}
	return string(t.source[t.index : t.index+n])
}

// Tokenize converts the source code into a slice of tokens.
func Tokenize(source []byte, filename string) ([]Token, error) {
	tokens := []Token{}
	tokenizer := NewTokenizer(source, filename)
	for {
		token := tokenizer.Next()
		tokens = append(tokens, token)
		if token.Type == TokEOF {
			break
		}
	}
	return tokens, nil
}

// Next returns the next token from the input.
func (t *Tokenizer) Next() Token {
	st := stateStart
	start := t.index
	tokenType := TokEOF
	invalid := false
	errorIndex := 0 // Track the position where error is first encountered
	escDigits := 0
	escDigitsExpected := 0

	// Helper function to mark a token as invalid and set the error index if not already set
	markInvalid := func() {
		invalid = true
		if errorIndex == 0 {
			errorIndex = t.index
		}
	}

	// Use a labeled loop so we can "break out" when a token is complete.
outer:
	for done := false; t.index <= len(t.source) && !done; t.index++ {
		var c byte
		if t.index < len(t.source) {
			c = t.source[t.index]
		} else {
			c = 0
		}
		switch st {
		case stateStart:
			switch c {
			case 0:
				if t.index != len(t.source) {
					tokenType = TokINV
					markInvalid()
				}
				break outer
			case ' ', '\t', '\r':
				start = t.index + 1
			case '\n':
				tokenType = TokEOL
				done = true
			case '@':
				tokenType = TokAt
				done = true
			case '}':
				tokenType = TokCloseBrace
				done = true
			case ']':
				tokenType = TokCloseBracket
				done = true
			case ')':
				tokenType = TokCloseParen
				done = true
			case ':':
				tokenType = TokColon
				st = stateColon
			case ',':
				tokenType = TokComma
				done = true
			case '.':
				// Check 3-character operators first
				if t.peek(3) == "..." {
					tokenType = TokOpRest
					t.index += 2
				} else if t.peek(2) == ".." {
					// Then 2-character operators
					tokenType = TokOpRange
					t.index++
				} else {
					// Finally, single character operator
					tokenType = TokDot
				}
				done = true
			case '{':
				tokenType = TokOpenBrace
				done = true
			case '[':
				tokenType = TokOpenBracket
				done = true
			case '(':
				tokenType = TokOpenParen
				done = true
			case '?':
				if t.index+1 < len(t.source) {
					switch t.source[t.index+1] {
					case '+':
						tokenType = TokOpCheckedAdd
						t.index++
					case '/':
						tokenType = TokOpCheckedDiv
						t.index++
					case '%':
						tokenType = TokOpCheckedMod
						t.index++
					case '*':
						tokenType = TokOpCheckedMul
						t.index++
					case '-':
						tokenType = TokOpCheckedSub
						t.index++
					default:
						tokenType = TokQuestionMark
					}
				} else {
					tokenType = TokQuestionMark
				}
				done = true
			case ';':
				tokenType = TokSemiColon
				done = true
			case '/':
				if t.peek(2) == "/=" {
					tokenType = TokOpDivEQ
					t.index++
				} else {
					tokenType = TokOpDiv
				}
				done = true
			case '-':
				if t.peek(2) == "-=" {
					tokenType = TokOpMinusEQ
					t.index++
				} else {
					tokenType = TokOpMinus
				}
				done = true
			case '%':
				if t.peek(2) == "%=" {
					tokenType = TokOpModEQ
					t.index++
				} else {
					tokenType = TokOpMod
				}
				done = true
			case '*':
				if t.peek(2) == "*=" {
					tokenType = TokOpMulEQ
					t.index++
				} else {
					tokenType = TokOpMul
				}
				done = true
			case '!':
				if t.peek(2) == "!=" {
					tokenType = TokOpNE
					t.index++
				} else {
					tokenType = TokOpNot
				}
				done = true
			case '+':
				if t.peek(2) == "+=" {
					tokenType = TokOpPlusEQ
					t.index++
				} else {
					tokenType = TokOpPlus
				}
				done = true
			case '^':
				if t.peek(2) == "^=" {
					tokenType = TokOpPowEQ
					t.index++
				} else {
					tokenType = TokOpPow
				}
				done = true
			case '>':
				// Check 3-character operators first
				if t.peek(3) == ">>=" {
					tokenType = TokOpSHR_EQ
					t.index += 2
				} else if t.peek(2) == ">>" {
					// Then 2-character operators
					tokenType = TokOpSHR
					t.index++
				} else if t.peek(2) == ">=" {
					tokenType = TokOpGE
					t.index++
				} else {
					// Finally, single character operator
					tokenType = TokOpGT
				}
				done = true
			case '&':
				if t.peek(3) == "&&=" {
					tokenType = TokOpLogAndEQ
					t.index += 2
				} else if t.peek(2) == "&&" {
					tokenType = TokOpLogAnd
					t.index++
				} else if t.peek(2) == "&=" {
					tokenType = TokOpBitAndEQ
					t.index++
				} else {
					tokenType = TokOpBitAnd
				}
				done = true
			case '|':
				if t.peek(3) == "||=" {
					tokenType = TokOpLogOrEQ
					t.index += 2
				} else if t.peek(2) == "||" {
					tokenType = TokOpLogOr
					t.index++
				} else if t.peek(2) == "|=" {
					tokenType = TokOpBitOrEQ
					t.index++
				} else {
					tokenType = TokOpBitOr
				}
				done = true
			case '=':
				if t.peek(2) == "==" {
					tokenType = TokOpEQ
					t.index++
				} else if t.peek(2) == "=~" {
					tokenType = TokOpMatch
					t.index++
				} else {
					tokenType = TokOpAssign
				}
				done = true
			case '~':
				tokenType = TokOpBitNot
				done = true
			case '#':
				tokenType = TokComment
				st = stateComment
			case '0':
				tokenType = TokDecLit
				st = stateNum
			case '`':
				if t.peek(3) == "```" {
					t.index += 3
					tokenType = TokMultiStrLit
					invalid = t.skipMultiLineHeader()
					if invalid {
						break outer
					}
					st = stateMultiStrBody
				} else {
					// Regular raw string literal
					tokenType = TokRawStrLit
					st = stateRawStrLit
				}
			case '"':
				tokenType = TokStrLit
				st = stateStrLit
			case '<':
				// Check 3-character operators first
				if t.peek(3) == "<=>" {
					tokenType = TokOpCompare
					t.index += 2
				} else if t.peek(3) == "<<=" {
					tokenType = TokOpSHL_EQ
					t.index += 2
				} else if t.peek(2) == "<<" {
					// Then 2-character operators
					tokenType = TokOpSHL
					t.index++
				} else if t.peek(2) == "<=" {
					tokenType = TokOpLE
					t.index++
				} else {
					// Finally, single character operator
					tokenType = TokOpLT
				}
				done = true
			default:
				// Identifier start: letters or underscore.
				if isIdentifierStart(c) {
					tokenType = TokID
					st = stateID
				} else if isDecDigit(c) {
					// Safe to use isDecimalDigit here since '0' is handled in its own case above
					st = stateInt
					tokenType = TokDecLit
				} else {
					tokenType = TokINV
					markInvalid()
					done = true
				}
			}
		case stateColon:
			switch {
			case isIdentifierStart(c):
				tokenType = TokSymLit
				st = stateSym
			case isDecDigit(c):
				tokenType = TokSymLit
				st = stateSym
				markInvalid()
			case c == '"':
				tokenType = TokSymLit
				st = stateQuotedSym
			default:
				break outer
			}
		case stateSym:
			switch {
			case isIdentifierChar(c):
				// continue symbol
			default:
				break outer
			}
		case stateQuotedSym:
			switch c {
			case 0:
				markInvalid()
				break outer
			case '"':
				done = true
			case '\n':
				markInvalid()
				break outer
			default:
				// Just continue consuming characters in the string
			}
		case stateID:
			switch {
			case isIdentifierStart(c) || isDecDigit(c):
				// Continue identifier.
			default:
				lexeme := string(t.source[start:t.index])
				if reserved, ok := GetReserved(lexeme); ok {
					tokenType = reserved
				} else if len(lexeme) > 0 && lexeme[0] >= 'A' && lexeme[0] <= 'Z' {
					tokenType = TokTypeID
				}
				break outer
			}
		case stateNum:
			switch {
			case isDecDigit(c) || c == '_':
				st = stateInt
			case c == '.':
				tokenType = TokFloatLit
				st = stateIntDot
			case c == 'b':
				tokenType = TokBinLit
				st = stateBinFirst
			case c == 'o':
				tokenType = TokOctLit
				st = stateOctFirst
			case c == 'x':
				tokenType = TokHexLit
				st = stateHexFirst
			case isInvNumLetter(c):
				markInvalid()
			default:
				break outer
			}
		case stateInt:
			switch {
			case isDecDigit(c) || c == '_':
				// Continue int.
			case c == '.':
				tokenType = TokFloatLit
				st = stateIntDot
			case c == 'e':
				tokenType = TokFloatLit
				st = stateExp
			case isInvIntLetter(c):
				markInvalid()
			default:
				break outer
			}
		case stateIntDot:
			switch {
			case isDecDigit(c):
				st = stateFloat
			case isIdentifierStart(c) || c == '.':
				// This is not an error - it's likely UFCS (uniform function call syntax) or a range operator
				tokenType = TokDecLit
				t.index-- // Back up to process the next character in next token
				break outer
			default:
				markInvalid()
				done = true
			}
		case stateFloat:
			switch {
			case isDecDigit(c) || c == '_':
				// Continue float.
			case c == 'e':
				st = stateExp
			case c == '.':
				break outer
			case isInvIntLetter(c):
				markInvalid()
			default:
				break outer
			}
		case stateExp:
			switch {
			case c == '+' || c == '-':
				st = stateExpSign
			case isDecDigit(c):
				st = stateExpInt
			case isInvExpIntChar(c):
				markInvalid()
			default:
				markInvalid()
				break outer
			}
		case stateExpSign:
			switch {
			case isDecDigit(c):
				st = stateExpInt
			case isInvExpSignChar(c):
				markInvalid()
			default:
				markInvalid()
				break outer
			}
		case stateExpInt:
			switch {
			case isDecDigit(c):
				// Continue exponent integer.
			case c == '.':
				break outer
			case isInvExpIntChar(c):
				markInvalid()
			default:
				break outer
			}
		case stateBinFirst:
			switch {
			case c >= '0' && c <= '1':
				st = stateBin
			case isInvBinFirstChar(c):
				st = stateBin
				markInvalid()
			default:
				markInvalid()
				break outer
			}
		case stateBin:
			switch {
			case isBinaryDigit(c) || c == '_':
				// Continue binary.
			case c == '.':
				break outer
			case isInvBinChar(c):
				markInvalid()
			default:
				break outer
			}
		case stateHexFirst:
			switch {
			case isHexDigit(c):
				st = stateHex
			case c == 0:
				markInvalid()
				break outer
			default:
				st = stateHex
				markInvalid()
			}
		case stateHex:
			switch {
			case isHexDigit(c) || c == '_':
				// Continue hexadecimal.
			case c == '.':
				break outer
			case isInvHexChar(c):
				markInvalid()
			default:
				break outer
			}
		case stateOctFirst:
			switch {
			case isOctDigit(c):
				st = stateOct
			case isInvOctFirstChar(c):
				st = stateOct
				markInvalid()
			default:
				markInvalid()
				break outer
			}
		case stateOct:
			switch {
			case isOctDigit(c) || c == '_':
				// Continue octal.
			case c == '.':
				break outer
			case isInvOctChar(c):
				markInvalid()
			default:
				break outer
			}
		case stateRawStrLit:
			switch {
			case c == 0:
				markInvalid()
				break outer
			case c == '`':
				// Check if it's a double backtick (escaped backtick)
				if t.peek(2) == "``" {
					t.index++ // Skip the second backtick
				} else {
					// Single backtick - end of string
					done = true
				}
			default:
				// Continue reading characters
			}
		case stateStrLit:
			switch {
			case c == 0:
				markInvalid()
				break outer
			case c == '\\':
				if t.peek(2) == "\\(" {
					tokenType = TokInterpStrLit
					t.index += 2 // Skip the `\(`
					if t.skipInterpolation() {
						markInvalid()
						break outer
					}
					t.index--
				} else {
					t.pushState(st)
					st = stateEscSeq
				}
			case c == '"':
				done = true
			case c == '\n':
				markInvalid()
				break outer
			default:
				// Just continue consuming characters in the string
			}
		case stateMultiStrBody:
			switch {
			case c == 0:
				markInvalid()
				break outer
			case c == '\\':
				if t.peek(2) == "\\(" {
					t.index += 2 // Skip the `\(`
					if t.skipInterpolation() {
						markInvalid()
						break outer
					}
					t.index--
				} else {
					t.pushState(st)
					st = stateEscSeq
				}
			case c == '`':
				if t.peek(3) == "```" {
					// Check if we have a newline before the closing sequence
					// Look backwards for a newline, skipping any whitespace
					i := t.index - 1
					for i >= 0 && (t.source[i] == ' ' || t.source[i] == '\t') {
						i--
					}
					if i >= 0 && t.source[i] != '\n' {
						markInvalid()
					}
					t.index += 3
					break outer
				}
			default:
				// Just continue consuming characters in the string.
			}
		case stateEscSeq:
			switch {
			case c == 0:
				markInvalid()
				break outer
			case c == 'x':
				st = stateHexEsc
				escDigits = 0
				escDigitsExpected = 2
			case c == 'u':
				st = stateHexEsc
				escDigits = 0
				escDigitsExpected = 4
			case c == 'U':
				st = stateHexEsc
				escDigits = 0
				escDigitsExpected = 8
			case isSimpleEsc(c):
				// Valid single-char escape; return to string literal
				st = t.popState()
			default:
				// Any other char => mark invalid but return to string literal
				st = t.popState()
				markInvalid()
			}
		case stateHexEsc:
			switch {
			case c == 0:
				markInvalid()
				break outer
			case isHexDigit(c):
				escDigits++
				if escDigits == escDigitsExpected {
					st = t.popState()
				}
			default:
				st = t.popState()
				markInvalid()
			}
		case stateComment:
			switch {
			case c == 0:
				break outer
			case c == '\n':
				done = true
			default:
				// Continue consuming comment characters.
			}
		}
	}
	// log.Printf("token: <%v> = %s, invalid: %v", tokenType, value, invalid)
	return Token{
		Type:        tokenType,
		Invalid:     invalid,
		File:        t.file,
		Offset:      int32(start),
		Length:      int32(t.index - start),
		ErrorOffset: int32(errorIndex),
	}
}

// skipInterpolation skips the interpolation until the closing )
// returns true if the interpolation is invalid, false otherwise
func (t *Tokenizer) skipInterpolation() (invalid bool) {
	parens := 0
	for {
		token := t.Next()
		switch {
		case token.Type == TokEOF:
			return true
		case token.Type == TokOpenParen:
			parens++
		case token.Type == TokCloseParen:
			parens--
		case token.Invalid:
			return true
		}
		if parens < 0 {
			return false
		}
	}
}

// skipMultiLineHeader skips the multi-line string header until it finds a newline
// returns true if the header is invalid, false otherwise
func (t *Tokenizer) skipMultiLineHeader() (invalid bool) {
	for {
		token := t.Next()
		switch {
		case token.Type == TokEOF:
			return true
		case token.Type == TokEOL:
			return false
		case token.Invalid:
			return true
		}
	}
}
