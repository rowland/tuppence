package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// continue_expression = "continue" [ expression ] .

func ContinueExpression(tokens []tok.Token) (*ast.ContinueExpression, []tok.Token, error) {
	remainder := skipTrivia(tokens)
	if peek(remainder).Type != tok.TokKwContinue {
		return nil, tokens, ErrNoMatch
	}
	remainder = remainder[1:]

	// Optional continue value is confined to the same statement.
	remainder = skipComments(remainder)
	switch peek(remainder).Type {
	case tok.TokEOF, tok.TokEOL, tok.TokSemiColon, tok.TokCloseBrace:
		return ast.NewContinueExpression(nil), remainder, nil
	}

	expression, remainder, err := Expression(remainder)
	if err == ErrNoMatch {
		return ast.NewContinueExpression(nil), remainder, nil
	} else if err != nil {
		return nil, remainder, err
	}

	return ast.NewContinueExpression(expression), remainder, nil
}
