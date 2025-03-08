package tok

// We define a state machine for tokenizing.
type state int

const (
	stateStart state = iota
	stateIdentifier
	stateNumber
	stateInt
	stateIntDot
	stateFloat
	stateExponent
	stateExponentSign
	stateExponentInt
	stateBinaryFirst
	stateBinary
	stateHexadecimalFirst
	stateHexadecimal
	stateOctalFirst
	stateOctal
	stateRawStringLiteral
	stateRawStringLiteralEnd
	stateStringLiteral
	stateEscapeSequence
	stateHexEscape
	stateComment
	stateColon
	stateSymbol
	stateQuotedSymbol
	stateMultiLineStringBody
)

// Tokenizer holds the state of the lexer.
type Tokenizer struct {
	file   *Source
	source []byte
	index  int
	// line   int
	bol    int // beginning-of-line index
	states []state
}

// NewTokenizer initializes a new Tokenizer.
func NewTokenizer(source []byte, filename string) *Tokenizer {
	file := NewSource(source, filename)
	idx := 0
	// Skip the UTF-8 BOM if present.
	bom := []byte{0xEF, 0xBB, 0xBF}
	if len(source) >= 3 && string(source[:3]) == string(bom) {
		idx = 3
	}
	return &Tokenizer{
		file:   file,
		source: source,
		index:  idx,
		// line:   0,
		bol: idx,
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
	escapeDigits := 0
	escapeDigitsExpected := 0

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
					tokenType = TokInvalid
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
					tokenType = TokOpDivEqual
					t.index++
				} else {
					tokenType = TokOpDiv
				}
				done = true
			case '-':
				if t.peek(2) == "-=" {
					tokenType = TokOpMinusEqual
					t.index++
				} else {
					tokenType = TokOpMinus
				}
				done = true
			case '%':
				if t.peek(2) == "%=" {
					tokenType = TokOpModEqual
					t.index++
				} else {
					tokenType = TokOpMod
				}
				done = true
			case '*':
				if t.peek(2) == "*=" {
					tokenType = TokOpMulEqual
					t.index++
				} else {
					tokenType = TokOpMul
				}
				done = true
			case '!':
				if t.peek(2) == "!=" {
					tokenType = TokOpNotEqual
					t.index++
				} else {
					tokenType = TokOpNot
				}
				done = true
			case '+':
				if t.peek(2) == "+=" {
					tokenType = TokOpPlusEqual
					t.index++
				} else {
					tokenType = TokOpPlus
				}
				done = true
			case '^':
				if t.peek(2) == "^=" {
					tokenType = TokOpPowEqual
					t.index++
				} else {
					tokenType = TokOpPow
				}
				done = true
			case '>':
				// Check 3-character operators first
				if t.peek(3) == ">>=" {
					tokenType = TokOpShiftRightEqual
					t.index += 2
				} else if t.peek(2) == ">>" {
					// Then 2-character operators
					tokenType = TokOpShiftRight
					t.index++
				} else if t.peek(2) == ">=" {
					tokenType = TokOpGreaterEqual
					t.index++
				} else {
					// Finally, single character operator
					tokenType = TokOpGreaterThan
				}
				done = true
			case '&':
				if t.peek(3) == "&&=" {
					tokenType = TokOpLogicalAndEqual
					t.index += 2
				} else if t.peek(2) == "&&" {
					tokenType = TokOpLogicalAnd
					t.index++
				} else if t.peek(2) == "&=" {
					tokenType = TokOpBitwiseAndEqual
					t.index++
				} else {
					tokenType = TokOpBitwiseAnd
				}
				done = true
			case '|':
				if t.peek(3) == "||=" {
					tokenType = TokOpLogicalOrEqual
					t.index += 2
				} else if t.peek(2) == "||" {
					tokenType = TokOpLogicalOr
					t.index++
				} else if t.peek(2) == "|=" {
					tokenType = TokOpBitwiseOrEqual
					t.index++
				} else {
					tokenType = TokOpBitwiseOr
				}
				done = true
			case '=':
				if t.peek(2) == "==" {
					tokenType = TokOpEqual
					t.index++
				} else if t.peek(2) == "=~" {
					tokenType = TokOpMatches
					t.index++
				} else {
					tokenType = TokOpAssign
				}
				done = true
			case '~':
				tokenType = TokOpBitwiseNot
				done = true
			case '#':
				tokenType = TokComment
				st = stateComment
			case '0':
				tokenType = TokDecimalLit
				st = stateNumber
			case '`':
				if t.peek(3) == "```" {
					t.index += 3
					tokenType = TokMultiLineStringLit
					invalid = t.skipMultiLineHeader()
					if invalid {
						break outer
					}
					st = stateMultiLineStringBody
				} else {
					// Regular raw string literal
					tokenType = TokRawStringLit
					st = stateRawStringLiteral
				}
			case '"':
				tokenType = TokStringLit
				st = stateStringLiteral
			case '<':
				// Check 3-character operators first
				if t.peek(3) == "<=>" {
					tokenType = TokOpCompareTo
					t.index += 2
				} else if t.peek(3) == "<<=" {
					tokenType = TokOpShiftLeftEqual
					t.index += 2
				} else if t.peek(2) == "<<" {
					// Then 2-character operators
					tokenType = TokOpShiftLeft
					t.index++
				} else if t.peek(2) == "<=" {
					tokenType = TokOpLessEqual
					t.index++
				} else {
					// Finally, single character operator
					tokenType = TokOpLessThan
				}
				done = true
			default:
				// Identifier start: letters or underscore.
				if isIdentifierStart(c) {
					tokenType = TokIdentifier
					st = stateIdentifier
				} else if isDecimalDigit(c) {
					// Safe to use isDecimalDigit here since '0' is handled in its own case above
					st = stateInt
					tokenType = TokDecimalLit
				} else {
					tokenType = TokInvalid
					done = true
				}
			}
		case stateColon:
			switch {
			case isIdentifierStart(c):
				tokenType = TokSymbolLit
				st = stateSymbol
			case isDecimalDigit(c):
				tokenType = TokSymbolLit
				st = stateSymbol
				invalid = true
			case c == '"':
				tokenType = TokSymbolLit
				st = stateQuotedSymbol
			default:
				break outer
			}
		case stateSymbol:
			switch {
			case isIdentifierStart(c):
				// continue symbol
			default:
				break outer
			}
		case stateQuotedSymbol:
			switch c {
			case 0:
				invalid = true
				break outer
			case '"':
				done = true
			case '\n':
				invalid = true
				break outer
			default:
				// Just continue consuming characters in the string
			}
		case stateIdentifier:
			switch {
			case isIdentifierStart(c) || isDecimalDigit(c):
				// Continue identifier.
			default:
				lexeme := string(t.source[start:t.index])
				if reserved, ok := GetReserved(lexeme); ok {
					tokenType = reserved
				} else if len(lexeme) > 0 && lexeme[0] >= 'A' && lexeme[0] <= 'Z' {
					tokenType = TokTypeIdentifier
				}
				break outer
			}
		case stateNumber:
			switch {
			case isDecimalDigit(c) || c == '_':
				st = stateInt
			case c == '.':
				tokenType = TokFloatLit
				st = stateIntDot
			case c == 'b':
				tokenType = TokBinaryLit
				st = stateBinaryFirst
			case c == 'o':
				tokenType = TokOctalLit
				st = stateOctalFirst
			case c == 'x':
				tokenType = TokHexLit
				st = stateHexadecimalFirst
			case isInvalidNumberLetter(c):
				invalid = true
			default:
				break outer
			}
		case stateInt:
			switch {
			case isDecimalDigit(c) || c == '_':
				// Continue int.
			case c == '.':
				tokenType = TokFloatLit
				st = stateIntDot
			case c == 'e':
				tokenType = TokFloatLit
				st = stateExponent
			case isInvalidIntegerLetter(c):
				invalid = true
			default:
				break outer
			}
		case stateIntDot:
			switch {
			case isDecimalDigit(c):
				st = stateFloat
			default:
				tokenType = TokDecimalLit
				t.index--
				break outer
			}
		case stateFloat:
			switch {
			case isDecimalDigit(c) || c == '_':
				// Continue float.
			case c == 'e':
				st = stateExponent
			case c == '.':
				break outer
			case isInvalidIntegerLetter(c):
				invalid = true
			default:
				break outer
			}
		case stateExponent:
			switch {
			case c == '+' || c == '-':
				st = stateExponentSign
			case isDecimalDigit(c):
				st = stateExponentInt
			case isInvalidExponentIntChar(c):
				invalid = true
			default:
				invalid = true
				break outer
			}
		case stateExponentSign:
			switch {
			case isDecimalDigit(c):
				st = stateExponentInt
			case isInvalidExponentSignChar(c):
				invalid = true
			default:
				invalid = true
				break outer
			}
		case stateExponentInt:
			switch {
			case isDecimalDigit(c):
				// Continue exponent integer.
			case c == '.':
				break outer
			case isInvalidExponentIntChar(c):
				invalid = true
			default:
				break outer
			}
		case stateBinaryFirst:
			switch {
			case c >= '0' && c <= '1':
				st = stateBinary
			case isInvalidBinaryFirstChar(c):
				st = stateBinary
				invalid = true
			default:
				invalid = true
				break outer
			}
		case stateBinary:
			switch {
			case c == '0' || c == '1' || c == '_':
				// Continue binary.
			case c == '.':
				break outer
			case isInvalidBinaryChar(c):
				invalid = true
			default:
				break outer
			}
		case stateHexadecimalFirst:
			switch {
			case isHexDigit(c):
				st = stateHexadecimal
			case c == 0:
				invalid = true
				break outer
			default:
				st = stateHexadecimal
				invalid = true
			}
		case stateHexadecimal:
			switch {
			case isHexDigit(c) || c == '_':
				// Continue hexadecimal.
			case c == '.':
				break outer
			case isInvalidHexadecimalChar(c):
				invalid = true
			default:
				break outer
			}
		case stateOctalFirst:
			switch {
			case isOctalDigit(c):
				st = stateOctal
			case isInvalidOctalFirstChar(c):
				st = stateOctal
				invalid = true
			default:
				invalid = true
				break outer
			}
		case stateOctal:
			switch {
			case isOctalDigit(c) || c == '_':
				// Continue octal.
			case c == '.':
				break outer
			case isInvalidOctalChar(c):
				invalid = true
			default:
				break outer
			}
		case stateRawStringLiteral:
			switch {
			case c == 0:
				invalid = true
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
		case stateStringLiteral:
			switch {
			case c == 0:
				invalid = true
				break outer
			case c == '\\':
				if t.peek(2) == "\\(" {
					tokenType = TokInterpStringLit
					t.index += 2 // Skip the `\(`
					invalid = t.skipInterpolation()
					if invalid {
						break outer
					}
					t.index--
				} else {
					t.pushState(st)
					st = stateEscapeSequence
				}
			case c == '"':
				done = true
			case c == '\n':
				invalid = true
				break outer
			default:
				// Just continue consuming characters in the string
			}
		case stateMultiLineStringBody:
			switch {
			case c == 0:
				invalid = true
				break outer
			case c == '\\':
				if t.peek(2) == "\\(" {
					t.index += 2 // Skip the `\(`
					invalid = t.skipInterpolation()
					if invalid {
						break outer
					}
					t.index--
				} else {
					t.pushState(st)
					st = stateEscapeSequence
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
						invalid = true
					}
					t.index += 3
					break outer
				}
			default:
				// Just continue consuming characters in the string
			}
		case stateEscapeSequence:
			switch {
			case c == 0:
				invalid = true
				break outer
			case c == 'x':
				st = stateHexEscape
				escapeDigits = 0
				escapeDigitsExpected = 2
			case c == 'u':
				st = stateHexEscape
				escapeDigits = 0
				escapeDigitsExpected = 4
			case c == 'U':
				st = stateHexEscape
				escapeDigits = 0
				escapeDigitsExpected = 8
			case isSimpleEscape(c):
				// Valid single-char escape; return to string literal
				st = t.popState()
			default:
				// Any other char => mark invalid but return to string literal
				st = t.popState()
				invalid = true
			}
		case stateHexEscape:
			switch {
			case c == 0:
				invalid = true
				break outer
			case isHexDigit(c):
				escapeDigits++
				if escapeDigits == escapeDigitsExpected {
					st = t.popState()
				}
			default:
				st = t.popState()
				invalid = true
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
	value := string(t.source[start:t.index])
	// log.Printf("token: <%v> = %s, invalid: %v", tokenType, value, invalid)
	return Token{
		Type:    tokenType,
		Invalid: invalid,
		Value:   value,
		File:    t.file,
		Offset:  start,
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
		}
	}
}
