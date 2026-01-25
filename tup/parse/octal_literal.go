package parse

import (
	"strconv"
	"strings"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// octal_literal = "0o" octal_digit { octal_digit } .
// octal_digit = "0"-"7" .

func OctalLiteral(tokens []tok.Token) (item *ast.IntegerLiteral, remainder []tok.Token, err error) {
	remainder = skipComments(tokens)
	if peek(remainder).Type != tok.TokOctLit || peek(remainder).Invalid {
		return nil, nil, errorExpecting(tok.TokenTypes[tok.TokOctLit], remainder)
	}
	value := remainder[0].Value()
	integerValue, _ := strconv.ParseInt(strings.ReplaceAll(value, "_", ""), 0, 64)
	return ast.NewOctalLiteral(value, integerValue, remainder[0].File, remainder[0].Offset, remainder[0].Length), remainder[1:], nil
}
