package parse

import (
	"bytes"
	"fmt"
	"unicode/utf8"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// rune_literal = "'" ( byte_escape_sequence | unicode_escape_sequence | escape_sequence | character - eol ) "'" .

func RuneLiteral(tokens []tok.Token) (item *ast.RuneLiteral, remainder []tok.Token, err error) {
	remainder = skipComments(tokens)
	if peek(remainder).Type != tok.TokRuneLit || peek(remainder).Invalid {
		return nil, nil, errorExpecting(tok.TokenTypes[tok.TokRuneLit], remainder)
	}
	value := remainder[0].Value()

	var buf bytes.Buffer
	for i := 1; i < len(value)-1; {
		if value[i] == '\\' {
			if bytes, length, matched := EscapeSequence(value[i:]); matched {
				buf.Write(bytes)
				i += length
			} else if bytes, length, matched, err := ByteEscapeSequence(value[i:]); matched {
				if err != nil {
					return nil, nil, err
				}
				buf.Write(bytes)
				i += length
			} else if bytes, length, matched, err := UnicodeEscapeSequence(value[i:]); matched {
				if err != nil {
					return nil, nil, err
				}
				buf.Write(bytes)
				i += length
			} else {
				return nil, nil, fmt.Errorf("invalid escape sequence: %s", value[i:])
			}
		} else {
			buf.WriteByte(value[i])
			i++
		}
	}
	runeValue, size := utf8.DecodeRune(buf.Bytes())
	if size == 0 || (runeValue == utf8.RuneError && size == 1) {
		return nil, nil, fmt.Errorf("invalid UTF-8 sequence: %s", buf.String())
	}
	if size != buf.Len() {
		return nil, nil, fmt.Errorf("rune literal contains multiple code points: %s", buf.String())
	}
	return ast.NewRuneLiteral(value, runeValue, remainder[0].File, remainder[0].Offset, remainder[0].Length), remainder[1:], nil
}
