package parse

import (
	"strconv"
	"strings"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// hexadecimal_literal = "0x" hex_digit { hex_digit | "_" } .
// hex_digit = decimal_digit | "a"-"f" | "A"-"F" .

func HexadecimalLiteral(tokens []tok.Token) (item *ast.HexadecimalLiteral, remainder []tok.Token, err error) {
	remainder = skipComments(tokens)
	if peek(remainder).Type != tok.TokHexLit || peek(remainder).Invalid {
		return nil, nil, errorExpecting(tok.TokenTypes[tok.TokHexLit], remainder)
	}
	value := remainder[0].Value()
	integerValue, _ := strconv.ParseInt(strings.ReplaceAll(value, "_", ""), 0, 64)
	return ast.NewHexadecimalLiteral(value, integerValue, remainder[0].File, remainder[0].Offset, remainder[0].Length), remainder[1:], nil
}
