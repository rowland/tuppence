package main

// We define a state machine for tokenizing.
type state int

const (
	stateStart state = iota
	stateQuestionMark
	stateOpDiv
	stateOpMinus
	stateOpMod
	stateOpMul
	stateOpNot
	stateOpPlus
	stateOpPow
	stateOpLessThan
	stateOpLessThanEqual
	stateOpGreaterThan
	stateOpBitwiseAnd
	stateOpLogicalAnd
	stateOpBitwiseOr
	stateOpLogicalOr
	stateOpEqual
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
	stateUnicodeEscapeSequence
	stateByteEscapeSequence
	stateComment
)

// Tokenizer holds the state of the lexer.
type Tokenizer struct {
	source              []byte
	filename            string
	index               int
	line                int
	bol                 int // beginning-of-line index
	pendingInvalidToken *Token
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
	if t.pendingInvalidToken != nil {
		token := *t.pendingInvalidToken
		t.pendingInvalidToken = nil
		return token
	}
	st := stateStart
	start := t.index
	tokenType := TokenEOF
	invalid := false
	escapeDigits := 0

	// Use a labeled loop so we can “break out” when a token is complete.
outer:
	for ; t.index <= len(t.source); t.index++ {
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
				t.line++
				start = t.index + 1
				t.bol = start
				tokenType = TokenEOL
				t.index++
				break outer
			case '@':
				tokenType = TokenAt
				t.index++
				break outer
			case '}':
				tokenType = TokenCloseBrace
				t.index++
				break outer
			case ']':
				tokenType = TokenCloseBracket
				t.index++
				break outer
			case ')':
				tokenType = TokenCloseParen
				t.index++
				break outer
			case ':':
				tokenType = TokenColon
				t.index++
				break outer
			case ',':
				tokenType = TokenComma
				t.index++
				break outer
			case '.':
				tokenType = TokenDot
				t.index++
				break outer
			case '{':
				tokenType = TokenOpenBrace
				t.index++
				break outer
			case '[':
				tokenType = TokenOpenBracket
				t.index++
				break outer
			case '(':
				tokenType = TokenOpenParen
				t.index++
				break outer
			case '?':
				tokenType = TokenQuestionMark
				st = stateQuestionMark
			case ';':
				tokenType = TokenSemiColon
				t.index++
				break outer
			case '/':
				tokenType = TokenOpDiv
				st = stateOpDiv
			case '-':
				tokenType = TokenOpMinus
				st = stateOpMinus
			case '%':
				tokenType = TokenOpMod
				st = stateOpMod
			case '*':
				tokenType = TokenOpMul
				st = stateOpMul
			case '!':
				tokenType = TokenOpNot
				st = stateOpNot
			case '+':
				tokenType = TokenOpPlus
				st = stateOpPlus
			case '^':
				tokenType = TokenOpPow
				st = stateOpPow
			case '<':
				tokenType = TokenOpLessThan
				st = stateOpLessThan
			case '>':
				tokenType = TokenOpGreaterThan
				st = stateOpGreaterThan
			case '&':
				tokenType = TokenOpBitwiseAnd
				st = stateOpBitwiseAnd
			case '|':
				tokenType = TokenOpBitwiseOr
				st = stateOpBitwiseOr
			case '=':
				tokenType = TokenOpEqual
				st = stateOpEqual
			case '~':
				tokenType = TokenOpBitwiseNot
				t.index++
				break outer
			case '#':
				tokenType = TokenComment
				st = stateComment
			default:
				// Identifier start: letters or underscore.
				if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || c == '_' {
					tokenType = TokenIdentifier
					st = stateIdentifier
				} else if c == '0' {
					st = stateNumber
					tokenType = TokenDecimalLiteral
				} else if c >= '1' && c <= '9' {
					st = stateInt
					tokenType = TokenDecimalLiteral
				} else if c == '`' {
					st = stateRawStringLiteral
					tokenType = TokenRawStringLiteral
				} else if c == '"' {
					st = stateStringLiteral
					tokenType = TokenStringLiteral
				} else {
					tokenType = TokenInvalid
					t.index++
					break outer
				}
			}
		case stateQuestionMark:
			switch c {
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
			}
			break outer
		case stateOpDiv:
			if c == '=' {
				tokenType = TokenOpDivEqual
				t.index++
			}
			break outer
		case stateOpMinus:
			if c == '=' {
				tokenType = TokenOpMinusEqual
				t.index++
			}
			break outer
		case stateOpMod:
			if c == '=' {
				tokenType = TokenOpModEqual
				t.index++
			}
			break outer
		case stateOpMul:
			if c == '=' {
				tokenType = TokenOpMulEqual
				t.index++
			}
			break outer
		case stateOpNot:
			if c == '=' {
				tokenType = TokenOpNotEqual
				t.index++
			}
			break outer
		case stateOpPlus:
			if c == '=' {
				tokenType = TokenOpPlusEqual
				t.index++
			}
			break outer
		case stateOpPow:
			if c == '=' {
				tokenType = TokenOpPowEqual
				t.index++
			}
			break outer
		case stateOpLessThan:
			if c == '<' {
				tokenType = TokenOpShiftLeft
				t.index++
				break outer
			} else if c == '=' {
				tokenType = TokenOpLessEqual
				st = stateOpLessThanEqual
			} else {
				break outer
			}
		case stateOpLessThanEqual:
			if c == '>' {
				tokenType = TokenOpCompareTo
				t.index++
			}
			break outer
		case stateOpGreaterThan:
			if c == '>' {
				tokenType = TokenOpShiftRight
				t.index++
			} else if c == '=' {
				tokenType = TokenOpGreaterEqual
				t.index++
			}
			break outer
		case stateOpBitwiseAnd:
			if c == '=' {
				tokenType = TokenOpBitwiseAndEqual
				t.index++
				break outer
			} else if c == '&' {
				tokenType = TokenOpLogicalAnd
				st = stateOpLogicalAnd
			} else {
				break outer
			}
		case stateOpLogicalAnd:
			if c == '=' {
				tokenType = TokenOpLogicalAndEqual
				t.index++
			}
			break outer
		case stateOpBitwiseOr:
			if c == '|' {
				tokenType = TokenOpLogicalOr
				st = stateOpLogicalOr
			} else if c == '=' {
				tokenType = TokenOpBitwiseOrEqual
				t.index++
				break outer
			} else {
				break outer
			}
		case stateOpLogicalOr:
			if c == '=' {
				tokenType = TokenOpLogicalOrEqual
				t.index++
			}
			break outer
		case stateOpEqual:
			if c == '=' {
				tokenType = TokenOpEqualEqual
				t.index++
			} else if c == '~' {
				tokenType = TokenOpMatches
				t.index++
			}
			break outer
		case stateIdentifier:
			if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_' || (c >= '0' && c <= '9') {
				// Continue identifier.
			} else {
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
			case (c >= '0' && c <= '9') || c == '_':
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
			case (c >= 'A' && c <= 'Z') || c == 'a' || (c >= 'c' && c <= 'n') || (c >= 'p' && c <= 'w') || (c >= 'y' && c <= 'z'):
				invalid = true
			default:
				break outer
			}
		case stateInt:
			if (c >= '0' && c <= '9') || c == '_' {
				// Continue int.
			} else if c == '.' {
				tokenType = TokenFloatLiteral
				st = stateIntDot
			} else if c == 'e' {
				tokenType = TokenFloatLiteral
				st = stateExponent
			} else if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'd') || (c >= 'f' && c <= 'z') {
				invalid = true
			} else {
				break outer
			}
		case stateIntDot:
			if c >= '0' && c <= '9' {
				st = stateFloat
			} else if (c >= 'A' && c <= 'Z') || c == '_' || (c >= 'a' && c <= 'z') {
				st = stateFloat
				invalid = true
			} else {
				invalid = true
				break outer
			}
		case stateFloat:
			if (c >= '0' && c <= '9') || c == '_' {
				// Continue float.
			} else if c == 'e' {
				st = stateExponent
			} else if (c >= 'A' && c <= 'd') || (c >= 'f' && c <= 'z') {
				invalid = true
			} else {
				break outer
			}
		case stateExponent:
			if c == '+' || c == '-' || (c >= '0' && c <= '9') {
				st = stateExponentSign
			} else if (c >= 'A' && c <= 'Z') || c == '_' || (c >= 'a' && c <= 'z') {
				invalid = true
			} else {
				invalid = true
				break outer
			}
		case stateExponentSign:
			if c >= '0' && c <= '9' {
				st = stateExponentInt
			} else if (c >= 'A' && c <= 'Z') || c == '_' || (c >= 'a' && c <= 'z') {
				invalid = true
			} else {
				invalid = true
				break outer
			}
		case stateExponentInt:
			if c >= '0' && c <= '9' {
				// Continue exponent integer.
			} else if (c >= 'A' && c <= 'Z') || c == '_' || (c >= 'a' && c <= 'z') {
				invalid = true
			} else {
				invalid = true
				break outer
			}
		case stateBinaryFirst:
			if c >= '0' && c <= '1' {
				st = stateBinary
			} else if c == '.' || (c >= '2' && c <= '9') || (c >= 'A' && c <= 'Z') || c == '_' || (c >= 'a' && c <= 'z') {
				st = stateBinary
				invalid = true
			} else {
				invalid = true
				break outer
			}
		case stateBinary:
			if c == '0' || c == '1' || c == '_' {
				// Continue binary.
			} else if c == '.' || (c >= '2' && c <= '9') || (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') {
				invalid = true
			} else {
				break outer
			}
		case stateHexadecimalFirst:
			if (c >= '0' && c <= '9') || (c >= 'A' && c <= 'F') || (c >= 'a' && c <= 'f') {
				st = stateHexadecimal
			} else if c == '.' || (c >= 'G' && c <= 'Z') || c == '_' || (c >= 'g' && c <= 'z') {
				st = stateHexadecimal
				invalid = true
			} else {
				invalid = true
				break outer
			}
		case stateHexadecimal:
			if (c >= '0' && c <= '9') || (c >= 'A' && c <= 'F') || c == '_' || (c >= 'a' && c <= 'f') {
				// Continue hexadecimal.
			} else if c == '.' || (c >= 'G' && c <= 'Z') || (c >= 'g' && c <= 'z') {
				invalid = true
			} else {
				break outer
			}
		case stateOctalFirst:
			if c >= '0' && c <= '7' {
				st = stateOctal
			} else if c == '.' || (c >= '8' && c <= '9') || (c >= 'A' && c <= 'Z') || c == '_' || (c >= 'a' && c <= 'z') {
				st = stateOctal
				invalid = true
			} else {
				invalid = true
				break outer
			}
		case stateOctal:
			if (c >= '0' && c <= '7') || c == '_' {
				// Continue octal.
			} else if c == '.' || (c >= '8' && c <= '9') || (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') {
				invalid = true
			} else {
				break outer
			}
		case stateRawStringLiteral:
			if c == 0 {
				invalid = true
				break outer
			} else if c == '`' {
				st = stateRawStringLiteralEnd
			}
		case stateRawStringLiteralEnd:
			if c == '`' {
				st = stateRawStringLiteral
			} else {
				break outer
			}
		case stateStringLiteral:
			if c == 0 {
				invalid = true
				break outer
			} else if c == '\\' {
				st = stateEscapeSequence
			} else if c == '"' {
				t.index++
				break outer
			}
		case stateEscapeSequence:
			if c == 0 {
				invalid = true
				break outer
			} else if c == 'x' {
				st = stateByteEscapeSequence
				escapeDigits = 0
			} else if c == 'u' {
				st = stateUnicodeEscapeSequence
				escapeDigits = 0
			} else if c == 'n' || c == 't' || c == '"' || c == '\'' || c == '\\' || c == 'r' || c == 'b' || c == 'f' || c == 'v' || c == '0' {
				st = stateStringLiteral
			} else {
				st = stateStringLiteral
				invalid = true
			}
		case stateUnicodeEscapeSequence:
			if c == 0 {
				invalid = true
				break outer
			} else if (c >= '0' && c <= '9') || (c >= 'A' && c <= 'F') || (c >= 'a' && c <= 'f') {
				escapeDigits++
				if escapeDigits == 4 {
					st = stateStringLiteral
				}
			} else {
				st = stateStringLiteral
				invalid = true
			}
		case stateByteEscapeSequence:
			if c == 0 {
				invalid = true
				break outer
			} else if (c >= '0' && c <= '9') || (c >= 'A' && c <= 'F') || (c >= 'a' && c <= 'f') {
				escapeDigits++
				if escapeDigits == 2 {
					st = stateStringLiteral
				}
			} else {
				st = stateStringLiteral
				invalid = true
			}
		case stateComment:
			if c == 0 {
				break outer
			} else if c == '\n' {
				t.index++
				break outer
			}
		}
	}
	col := t.index - t.bol + 1
	return Token{
		Type:     tokenType,
		Invalid:  invalid,
		Line:     t.line,
		Column:   col,
		Value:    string(t.source[start:t.index]),
		Filename: t.filename,
	}
}
