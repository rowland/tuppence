package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// return_expression = "return" [ expression ] .

func ReturnExpression(tokens []tok.Token) (*ast.ReturnExpression, []tok.Token, error) {
	remainder, expression, matched, err := keywordOptionalExpression(tokens, tok.TokKwReturn)
	if err != nil {
		return nil, remainder, err
	}
	if !matched {
		return nil, tokens, ErrNoMatch
	}

	return ast.NewReturnExpression(expression), remainder, nil
}
