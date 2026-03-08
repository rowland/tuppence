package parse

import (
	"strconv"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

func FloatLiteral(tokens []tok.Token) (lit *ast.FloatLiteral, remainder []tok.Token, err error) {
	remainder = skipComments(tokens)

	t := peek(remainder)
	if t.Type != tok.TokFloatLit {
		return nil, tokens, ErrNoMatch
	} else if t.Invalid {
		return nil, nil, errorExpectingTokenType(tok.TokFloatLit, remainder)
	}

	value := t.Value()
	floatValue, _ := strconv.ParseFloat(value, 64)
	return ast.NewFloatLiteral(value, floatValue, t.File, t.Offset, t.Length), remainder[1:], nil
}
