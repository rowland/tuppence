package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// typeof_expression = "typeof" "(" expression ")" .

func TypeofExpression(tokens []tok.Token) (expr *ast.TypeofExpression, remainder []tok.Token, err error) {
	remainder = skipTrivia(tokens)

	if peek(remainder).Type != tok.TokKwTypeof {
		return nil, tokens, ErrNoMatch
	}
	remainder = remainder[1:]

	var found bool
	if remainder, found = OpenParen(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokOpenParen, remainder)
	}

	var expression ast.Expression
	if expression, remainder, err = Expression(remainder); err != nil {
		return nil, remainder, err
	}

	if remainder, found = CloseParen(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokCloseParen, remainder)
	}

	return ast.NewTypeofExpression(expression), remainder, nil
}
