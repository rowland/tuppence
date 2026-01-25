package parse

import (
	"strconv"
	"strings"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// decimal_literal = decimal_digit { decimal_digit | "_" } .
// decimal_digit = "0"-"9" .

func DecimalLiteral(tokens []tok.Token) (item *ast.IntegerLiteral, remainder []tok.Token, err error) {
	remainder = skipComments(tokens)
	if peek(remainder).Type != tok.TokDecLit || peek(remainder).Invalid {
		return nil, nil, errorExpecting(tok.TokenTypes[tok.TokDecLit], remainder)
	}
	value := remainder[0].Value()
	integerValue, _ := strconv.ParseInt(strings.ReplaceAll(value, "_", ""), 10, 64)
	return ast.NewDecimalLiteral(value, integerValue, remainder[0].File, remainder[0].Offset, remainder[0].Length), remainder[1:], nil
}
