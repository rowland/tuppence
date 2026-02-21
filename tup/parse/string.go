package parse

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// string_literal = '"' { escape_sequence | byte_escape_sequence | unicode_escape_sequence | character - '"' - eol } '"' .

func StringLiteral(tokens []tok.Token) (item *ast.StringLiteral, remainder []tok.Token, err error) {
	remainder = skipComments(tokens)

	if peek(remainder).Type != tok.TokStrLit {
		return nil, tokens, ErrNoMatch
	} else if peek(remainder).Invalid {
		return nil, remainder, errorExpecting(tok.TokenTypes[tok.TokStrLit], remainder)
	}

	value := remainder[0].Value()

	var sb strings.Builder
	for i := 1; i < len(value)-1; {
		if value[i] == '\\' {
			if bytes, length, matched := EscapeSequence(value[i:]); matched {
				sb.Write(bytes)
				i += length
			} else if bytes, length, matched, err := ByteEscapeSequence(value[i:]); matched {
				if err != nil {
					return nil, nil, err
				}
				sb.Write(bytes)
				i += length
			} else if bytes, length, matched, err := UnicodeEscapeSequence(value[i:]); matched {
				if err != nil {
					return nil, nil, err
				}
				sb.Write(bytes)
				i += length
			} else {
				return nil, nil, fmt.Errorf("invalid escape sequence: %s", value[i:])
			}
		} else {
			sb.WriteByte(value[i])
			i++
		}
	}
	stringValue := sb.String()
	return ast.NewStringLiteral(value, stringValue, remainder[0].File, remainder[0].Offset, remainder[0].Length), remainder[1:], nil
}

// escape_sequence = ( "\\n" | "\\t" | "\\\"" | "\\'" | "\\\\" | "\\r" | "\\b" | "\\f" | "\\v" | "\\0" | "\\`" ) .

func EscapeSequence(value string) (bytes []byte, length int, matched bool) {
	if value[0] == '\\' && value[1] == 'n' {
		return []byte{'\n'}, 2, true
	} else if value[0] == '\\' && value[1] == 't' {
		return []byte{'\t'}, 2, true
	} else if value[0] == '\\' && value[1] == '"' {
		return []byte{'"'}, 2, true
	} else if value[0] == '\\' && value[1] == '\'' {
		return []byte{'\''}, 2, true
	} else if value[0] == '\\' && value[1] == '\\' {
		return []byte{'\\'}, 2, true
	} else if value[0] == '\\' && value[1] == 'r' {
		return []byte{'\r'}, 2, true
	} else if value[0] == '\\' && value[1] == 'b' {
		return []byte{'\b'}, 2, true
	} else if value[0] == '\\' && value[1] == 'f' {
		return []byte{'\f'}, 2, true
	} else if value[0] == '\\' && value[1] == 'v' {
		return []byte{'\v'}, 2, true
	} else if value[0] == '\\' && value[1] == '0' {
		return []byte{0}, 2, true
	} else if value[0] == '\\' && value[1] == '`' {
		return []byte{'`'}, 2, true
	}
	return nil, 0, false
}

// byte_escape_sequence = "\\" "x" hex_digit hex_digit .

func ByteEscapeSequence(value string) (bytes []byte, length int, matched bool, err error) {
	if len(value) >= 4 && value[0] == '\\' && value[1] == 'x' {
		b, err := strconv.ParseUint(value[2:4], 16, 8)
		if err != nil {
			return nil, 0, true, err
		}
		return []byte{byte(b)}, 4, true, nil
	}
	return nil, 0, false, nil
}

// unicode_escape_sequence = "\\" "u" hex_digit hex_digit hex_digit hex_digit
//                         | "\\" "U" hex_digit hex_digit hex_digit hex_digit hex_digit hex_digit hex_digit hex_digit .

func UnicodeEscapeSequence(value string) (bytes []byte, length int, matched bool, err error) {
	if bytes, length, matched, err := Unicode16BitEscapeSequence(value); matched {
		return bytes, length, matched, err
	} else if bytes, length, matched, err := Unicode32BitEscapeSequence(value); matched {
		return bytes, length, matched, err
	}
	return nil, 0, false, nil
}

// "\\" "u" hex_digit hex_digit hex_digit hex_digit

func Unicode16BitEscapeSequence(value string) (bytes []byte, length int, matched bool, err error) {
	const sequenceLength = 6
	if len(value) >= sequenceLength && value[0] == '\\' && value[1] == 'u' {
		codepoint, err := strconv.ParseUint(value[2:sequenceLength], 16, 32)
		if err != nil {
			return nil, 0, true, err
		}
		r := rune(codepoint)
		if !utf8.ValidRune(r) {
			return nil, 0, true, fmt.Errorf("invalid Unicode escape sequence: %s", value[2:sequenceLength])
		}
		var buf [utf8.UTFMax]byte
		n := utf8.EncodeRune(buf[:], r)
		return buf[:n], sequenceLength, true, nil
	}
	return nil, 0, false, nil
}

// "\\" "U" hex_digit hex_digit hex_digit hex_digit hex_digit hex_digit hex_digit hex_digit

func Unicode32BitEscapeSequence(value string) (bytes []byte, length int, matched bool, err error) {
	const sequenceLength = 10
	if len(value) >= sequenceLength && value[0] == '\\' && value[1] == 'U' {
		codepoint, err := strconv.ParseUint(value[2:sequenceLength], 16, 32)
		if err != nil {
			return nil, 0, true, err
		}
		r := rune(codepoint)
		if !utf8.ValidRune(r) {
			return nil, 0, true, fmt.Errorf("invalid Unicode escape sequence: %s", value[2:sequenceLength])
		}
		var buf [utf8.UTFMax]byte
		n := utf8.EncodeRune(buf[:], r)
		return buf[:n], sequenceLength, true, nil
	}
	return nil, 0, false, nil
}
