package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// continue_expression = "continue" [ expression ] .

func ContinueExpression(tokens []tok.Token) (*ast.ContinueExpression, []tok.Token, error) {
	remainder, expression, matched, err := keywordOptionalExpression(tokens, tok.TokKwContinue)
	if err != nil {
		return nil, remainder, err
	}
	if !matched {
		return nil, tokens, ErrNoMatch
	}

	return ast.NewContinueExpression(expression), remainder, nil
}
