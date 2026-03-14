package parse

import (
	"bytes"
	"fmt"
	"unicode/utf8"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// rune_literal = "'" ( byte_escape_sequence | unicode_escape_sequence | escape_sequence | character - eol ) "'" .

func RuneLiteral(tokens []tok.Token) (lit *ast.RuneLiteral, remainder []tok.Token, err error) {
	// fmt.Println("RuneLiteral", tokens)
	remainder = skipTrivia(tokens)
	t := peek(remainder)
	if t.Type != tok.TokRuneLit {
		return nil, tokens, ErrNoMatch
	} else if t.Invalid {
		return nil, remainder, errorExpectingTokenType(tok.TokRuneLit, remainder)
	}
	value := t.Value()

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
	return ast.NewRuneLiteral(value, runeValue, t.File, t.Offset, t.Length), remainder[1:], nil
}
