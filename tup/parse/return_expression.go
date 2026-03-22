package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// return_expression = "return" [ expression ] .

func ReturnExpression(tokens []tok.Token) (*ast.ReturnExpression, []tok.Token, error) {
	remainder := skipTrivia(tokens)
	if peek(remainder).Type != tok.TokKwReturn {
		return nil, tokens, ErrNoMatch
	}
	remainder = remainder[1:]

	// Optional return value is confined to the same statement.
	remainder = skipComments(remainder)
	switch peek(remainder).Type {
	case tok.TokEOF, tok.TokEOL, tok.TokSemiColon, tok.TokCloseBrace:
		return ast.NewReturnExpression(nil), remainder, nil
	}

	expression, remainder, err := Expression(remainder)
	if err == ErrNoMatch {
		return ast.NewReturnExpression(nil), remainder, nil
	} else if err != nil {
		return nil, remainder, err
	}

	return ast.NewReturnExpression(expression), remainder, nil
}
