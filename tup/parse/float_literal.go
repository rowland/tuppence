package parse

import (
	"strconv"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

func FloatLiteral(tokens []tok.Token) (item *ast.FloatLiteral, remainder []tok.Token, err error) {
	remainder = skipComments(tokens)
	if peek(remainder).Type != tok.TokFloatLit || peek(remainder).Invalid {
		return nil, nil, errorExpecting(tok.TokenTypes[tok.TokFloatLit], remainder)
	}
	value := remainder[0].Value()
	floatValue, _ := strconv.ParseFloat(value, 64)
	return ast.NewFloatLiteral(value, floatValue, remainder[0].File, remainder[0].Offset, remainder[0].Length), remainder[1:], nil
}
