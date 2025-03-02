package main

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
	source   []byte
	filename string
	index    int
	line     int
	bol      int // beginning-of-line index
	states   []state
}

// NewTokenizer initializes a new Tokenizer.
func NewTokenizer(source []byte, filename string) *Tokenizer {
	idx := 0
	// Skip the UTF-8 BOM if present.
	bom := []byte{0xEF, 0xBB, 0xBF}
	if len(source) >= 3 && string(source[:3]) == string(bom) {
		idx = 3
	}
	return &Tokenizer{
		source:   source,
		filename: filename,
		index:    idx,
		line:     1,
		bol:      idx,
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
		if token.Type == TokenEOF {
			break
		}
	}
	return tokens, nil
}

// Next returns the next token from the input.
func (t *Tokenizer) Next() Token {
	st := stateStart
	start := t.index
	tokenType := TokenEOF
	invalid := false
	escapeDigits := 0
	escapeDigitsExpected := 0
	col := 0 // Initialize column position

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
					tokenType = TokenInvalid
				}
				break outer
			case ' ', '\t', '\r':
				start = t.index + 1
			case '\n':
				tokenType = TokenEOL
				col = t.index - t.bol + 1 // Calculate column before updating line tracking
				t.line++
				t.bol = t.index + 1
				done = true
			case '@':
				tokenType = TokenAt
				done = true
			case '}':
				tokenType = TokenCloseBrace
				done = true
			case ']':
				tokenType = TokenCloseBracket
				done = true
			case ')':
				tokenType = TokenCloseParen
				done = true
			case ':':
				tokenType = TokenColon
				st = stateColon
			case ',':
				tokenType = TokenComma
				done = true
			case '.':
				// Check 3-character operators first
				if t.peek(3) == "..." {
					tokenType = TokenOpRest
					t.index += 2
				} else if t.peek(2) == ".." {
					// Then 2-character operators
					tokenType = TokenOpRange
					t.index++
				} else {
					// Finally, single character operator
					tokenType = TokenDot
				}
				done = true
			case '{':
				tokenType = TokenOpenBrace
				done = true
			case '[':
				tokenType = TokenOpenBracket
				done = true
			case '(':
				tokenType = TokenOpenParen
				done = true
			case '?':
				if t.index+1 < len(t.source) {
					switch t.source[t.index+1] {
					case '+':
						tokenType = TokenOpCheckedAdd
						t.index++
					case '/':
						tokenType = TokenOpCheckedDiv
						t.index++
					case '%':
						tokenType = TokenOpCheckedMod
						t.index++
					case '*':
						tokenType = TokenOpCheckedMul
						t.index++
					case '-':
						tokenType = TokenOpCheckedSub
						t.index++
					default:
						tokenType = TokenQuestionMark
					}
				} else {
					tokenType = TokenQuestionMark
				}
				done = true
			case ';':
				tokenType = TokenSemiColon
				done = true
			case '/':
				if t.peek(2) == "/=" {
					tokenType = TokenOpDivEqual
					t.index++
				} else {
					tokenType = TokenOpDiv
				}
				done = true
			case '-':
				if t.peek(2) == "-=" {
					tokenType = TokenOpMinusEqual
					t.index++
				} else {
					tokenType = TokenOpMinus
				}
				done = true
			case '%':
				if t.peek(2) == "%=" {
					tokenType = TokenOpModEqual
					t.index++
				} else {
					tokenType = TokenOpMod
				}
				done = true
			case '*':
				if t.peek(2) == "*=" {
					tokenType = TokenOpMulEqual
					t.index++
				} else {
					tokenType = TokenOpMul
				}
				done = true
			case '!':
				if t.peek(2) == "!=" {
					tokenType = TokenOpNotEqual
					t.index++
				} else {
					tokenType = TokenOpNot
				}
				done = true
			case '+':
				if t.peek(2) == "+=" {
					tokenType = TokenOpPlusEqual
					t.index++
				} else {
					tokenType = TokenOpPlus
				}
				done = true
			case '^':
				if t.peek(2) == "^=" {
					tokenType = TokenOpPowEqual
					t.index++
				} else {
					tokenType = TokenOpPow
				}
				done = true
			case '>':
				// Check 3-character operators first
				if t.peek(3) == ">>=" {
					tokenType = TokenOpShiftRightEqual
					t.index += 2
				} else if t.peek(2) == ">>" {
					// Then 2-character operators
					tokenType = TokenOpShiftRight
					t.index++
				} else if t.peek(2) == ">=" {
					tokenType = TokenOpGreaterEqual
					t.index++
				} else {
					// Finally, single character operator
					tokenType = TokenOpGreaterThan
				}
				done = true
			case '&':
				if t.peek(3) == "&&=" {
					tokenType = TokenOpLogicalAndEqual
					t.index += 2
				} else if t.peek(2) == "&&" {
					tokenType = TokenOpLogicalAnd
					t.index++
				} else if t.peek(2) == "&=" {
					tokenType = TokenOpBitwiseAndEqual
					t.index++
				} else {
					tokenType = TokenOpBitwiseAnd
				}
				done = true
			case '|':
				if t.peek(3) == "||=" {
					tokenType = TokenOpLogicalOrEqual
					t.index += 2
				} else if t.peek(2) == "||" {
					tokenType = TokenOpLogicalOr
					t.index++
				} else if t.peek(2) == "|=" {
					tokenType = TokenOpBitwiseOrEqual
					t.index++
				} else {
					tokenType = TokenOpBitwiseOr
				}
				done = true
			case '=':
				if t.peek(2) == "==" {
					tokenType = TokenOpEqualEqual
					t.index++
				} else if t.peek(2) == "=~" {
					tokenType = TokenOpMatches
					t.index++
				} else {
					tokenType = TokenOpEqual
				}
				done = true
			case '~':
				tokenType = TokenOpBitwiseNot
				done = true
			case '#':
				tokenType = TokenComment
				st = stateComment
			case '0':
				tokenType = TokenDecimalLiteral
				st = stateNumber
			case '`':
				if t.peek(3) == "```" {
					t.index += 3
					tokenType = TokenMultiLineStringLiteral
					invalid = t.skipMultiLineHeader()
					if invalid {
						break outer
					}
					st = stateMultiLineStringBody
				} else {
					// Regular raw string literal
					tokenType = TokenRawStringLiteral
					st = stateRawStringLiteral
				}
			case '"':
				tokenType = TokenStringLiteral
				st = stateStringLiteral
			case '<':
				// Check 3-character operators first
				if t.peek(3) == "<=>" {
					tokenType = TokenOpCompareTo
					t.index += 2
				} else if t.peek(3) == "<<=" {
					tokenType = TokenOpShiftLeftEqual
					t.index += 2
				} else if t.peek(2) == "<<" {
					// Then 2-character operators
					tokenType = TokenOpShiftLeft
					t.index++
				} else if t.peek(2) == "<=" {
					tokenType = TokenOpLessEqual
					t.index++
				} else {
					// Finally, single character operator
					tokenType = TokenOpLessThan
				}
				done = true
			default:
				// Identifier start: letters or underscore.
				if isIdentifierStart(c) {
					tokenType = TokenIdentifier
					st = stateIdentifier
				} else if isDecimalDigit(c) {
					// Safe to use isDecimalDigit here since '0' is handled in its own case above
					st = stateInt
					tokenType = TokenDecimalLiteral
				} else {
					tokenType = TokenInvalid
					done = true
				}
			}
		case stateColon:
			switch {
			case isIdentifierStart(c):
				tokenType = TokenSymbolLiteral
				st = stateSymbol
			case isDecimalDigit(c):
				tokenType = TokenSymbolLiteral
				st = stateSymbol
				invalid = true
			case c == '"':
				tokenType = TokenSymbolLiteral
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
					tokenType = TokenTypeIdentifier
				}
				break outer
			}
		case stateNumber:
			switch {
			case isDecimalDigit(c) || c == '_':
				st = stateInt
			case c == '.':
				tokenType = TokenFloatLiteral
				st = stateIntDot
			case c == 'b':
				tokenType = TokenBinaryLiteral
				st = stateBinaryFirst
			case c == 'o':
				tokenType = TokenOctalLiteral
				st = stateOctalFirst
			case c == 'x':
				tokenType = TokenHexadecimalLiteral
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
				tokenType = TokenFloatLiteral
				st = stateIntDot
			case c == 'e':
				tokenType = TokenFloatLiteral
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
				tokenType = TokenDecimalLiteral
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
			case c == '\n':
				t.line++
				t.bol = t.index + 1
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
					tokenType = TokenInterpolatedStringLiteral
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
			case c == '\n':
				t.line++
				t.bol = t.index + 1
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
	line := t.line
	if tokenType == TokenEOL {
		line--
	} else {
		col = t.index - t.bol + 1
	}
	value := string(t.source[start:t.index])
	// log.Printf("token: <%v> = %s, invalid: %v", tokenType, value, invalid)
	return Token{
		Type:     tokenType,
		Invalid:  invalid,
		Line:     line,
		Column:   col,
		Value:    value,
		Filename: t.filename,
	}
}

// skipInterpolation skips the interpolation until the closing )
// returns true if the interpolation is invalid, false otherwise
func (t *Tokenizer) skipInterpolation() (invalid bool) {
	parens := 0
	for {
		token := t.Next()
		switch {
		case token.Type == TokenEOF:
			return true
		case token.Type == TokenOpenParen:
			parens++
		case token.Type == TokenCloseParen:
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
		case token.Type == TokenEOF:
			return true
		case token.Type == TokenEOL:
			return false
		}
	}
}

// isHexDigit returns true if c is a valid hexadecimal digit
func isHexDigit(c byte) bool {
	return (c >= '0' && c <= '9') ||
		(c >= 'a' && c <= 'f') ||
		(c >= 'A' && c <= 'F')
}

// isSimpleEscape returns true if c is a valid single-character escape sequence
func isSimpleEscape(c byte) bool {
	return c == 'n' || c == 't' || c == '"' || c == '\'' || c == '\\' ||
		c == 'r' || c == 'b' || c == 'f' || c == 'v' || c == '0'
}

// isOctalDigit returns true if c is a valid octal digit
func isOctalDigit(c byte) bool {
	return c >= '0' && c <= '7'
}

// isDecimalDigit returns true if c is a valid decimal digit
func isDecimalDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

// isLetter returns true if c is a letter (A-Z or a-z)
func isLetter(c byte) bool {
	return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')
}

// isIdentifierStart returns true if c is a letter or underscore
func isIdentifierStart(c byte) bool {
	return isLetter(c) || c == '_'
}

// isInvalidNumberLetter returns true if c is a letter that would make a number invalid
// This excludes 'b', 'o', and 'x' which are handled separately as valid number prefixes
func isInvalidNumberLetter(c byte) bool {
	return isLetter(c) && c != 'b' && c != 'o' && c != 'x'
}

// isInvalidIntegerLetter returns true if c is a letter that would make an integer invalid
// This excludes 'e' which is handled separately as an exponent marker
func isInvalidIntegerLetter(c byte) bool {
	return isLetter(c) && c != 'e'
}

// isInvalidExponentSignChar returns true if c is an invalid character after an exponent sign
// This includes letters, underscore, and signs (+ or -) since we've already handled the sign
func isInvalidExponentSignChar(c byte) bool {
	return isLetter(c) || c == '_' || c == '+' || c == '-'
}

// isInvalidExponentIntChar returns true if c is an invalid character in an exponent's integer part
// This includes letters and underscore since only digits are allowed
func isInvalidExponentIntChar(c byte) bool {
	return isLetter(c) || c == '_'
}

// isInvalidBinaryFirstChar returns true if c is an invalid character for the first position of a binary number
// This includes any character that is not 0 or 1
func isInvalidBinaryFirstChar(c byte) bool {
	return (c >= '2' && c <= '9') || isLetter(c) || c == '_' || c == '.'
}

// isInvalidBinaryChar returns true if c is an invalid character for a binary number
// This includes any character that is not 0, 1, or underscore
func isInvalidBinaryChar(c byte) bool {
	return (c >= '2' && c <= '9') || isLetter(c)
}

// isInvalidOctalChar returns true if c is an invalid character for an octal number
// This includes any character that is not an octal digit (0-7) or underscore
func isInvalidOctalChar(c byte) bool {
	return (c >= '8' && c <= '9') || isLetter(c)
}

// isInvalidOctalFirstChar returns true if c is an invalid character for the first position of an octal number
// This includes any character that is not an octal digit (0-7)
func isInvalidOctalFirstChar(c byte) bool {
	return (c >= '8' && c <= '9') || isLetter(c) || c == '_' || c == '.'
}

// isInvalidHexadecimalChar returns true if c is a letter that would make a hexadecimal number invalid
// This includes letters G-Z and g-z, since A-F and a-f are valid hex digits
func isInvalidHexadecimalChar(c byte) bool {
	return (c >= 'G' && c <= 'Z') || (c >= 'g' && c <= 'z')
}
