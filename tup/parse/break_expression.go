package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// break_expression = "break" [ expression ] .

func BreakExpression(tokens []tok.Token) (*ast.BreakExpression, []tok.Token, error) {
	remainder, expression, matched, err := keywordOptionalExpression(tokens, tok.TokKwBreak)
	if err != nil {
		return nil, remainder, err
	}
	if !matched {
		return nil, tokens, ErrNoMatch
	}

	return ast.NewBreakExpression(expression), remainder, nil
}
