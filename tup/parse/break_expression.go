package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// break_expression = "break" [ expression ] .

func BreakExpression(tokens []tok.Token) (*ast.BreakExpression, []tok.Token, error) {
	remainder := skipTrivia(tokens)
	if peek(remainder).Type != tok.TokKwBreak {
		return nil, tokens, ErrNoMatch
	}
	remainder = remainder[1:]

	// Optional break value is confined to the same statement.
	remainder = skipComments(remainder)
	switch peek(remainder).Type {
	case tok.TokEOF, tok.TokEOL, tok.TokSemiColon, tok.TokCloseBrace:
		return ast.NewBreakExpression(nil), remainder, nil
	}

	expression, remainder, err := Expression(remainder)
	if err == ErrNoMatch {
		return ast.NewBreakExpression(nil), remainder, nil
	} else if err != nil {
		return nil, remainder, err
	}

	return ast.NewBreakExpression(expression), remainder, nil
}
