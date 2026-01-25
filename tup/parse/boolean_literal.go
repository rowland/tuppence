package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// boolean_literal = "true" | "false" .

func BooleanLiteral(tokens []tok.Token) (item *ast.BooleanLiteral, remainder []tok.Token, err error) {
	remainder = skipComments(tokens)
	if peek(remainder).Type != tok.TokBoolLit || peek(remainder).Invalid {
		return nil, nil, errorExpecting(tok.TokenTypes[tok.TokBoolLit], remainder)
	}
	return ast.NewBooleanLiteral(
			remainder[0].Value(),
			remainder[0].Value() == "true",
			remainder[0].File,
			remainder[0].Offset,
			remainder[0].Length),
		remainder[1:],
		nil
}
