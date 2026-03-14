package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// boolean_literal = "true" | "false" .

func BooleanLiteral(tokens []tok.Token) (lit *ast.BooleanLiteral, remainder []tok.Token, err error) {
	// fmt.Println("BooleanLiteral", tokens)
	remainder = skipTrivia(tokens)

	t := peek(remainder)
	if t.Type != tok.TokBoolLit {
		return nil, tokens, ErrNoMatch
	} else if t.Invalid {
		return nil, remainder, errorExpectingTokenType(tok.TokBoolLit, remainder)
	}

	return ast.NewBooleanLiteral(t.Value(), t.Value() == "true", t.File, t.Offset, t.Length), remainder[1:], nil
}
