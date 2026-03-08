package parse

import (
	"strconv"
	"strings"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// hexadecimal_literal = "0x" hex_digit { hex_digit | "_" } .
// hex_digit = decimal_digit | "a"-"f" | "A"-"F" .

func HexadecimalLiteral(tokens []tok.Token) (lit *ast.IntegerLiteral, remainder []tok.Token, err error) {
	remainder = skipComments(tokens)

	t := peek(remainder)
	if t.Type != tok.TokHexLit {
		return nil, tokens, ErrNoMatch
	} else if t.Invalid {
		return nil, remainder, errorExpectingTokenType(tok.TokHexLit, remainder)
	}

	value := t.Value()
	integerValue, _ := strconv.ParseInt(strings.ReplaceAll(value, "_", ""), 0, 64)
	return ast.NewHexadecimalLiteral(value, integerValue, t.File, t.Offset, t.Length), remainder[1:], nil
}
