package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// range_bound = postfix_expression .

func RangeBound(tokens []tok.Token) (*ast.RangeBound, []tok.Token, error) {
	expr, remainder, err := postfixExpression(tokens, postfixBaseExpressionWithoutRange, false)
	if err != nil {
		return nil, remainder, err
	}

	return ast.NewRangeBound(expr), remainder, nil
}

// range = range_bound ".." range_bound .

func Range(tokens []tok.Token) (*ast.Range, []tok.Token, error) {
	start, remainder, err := RangeBound(tokens)
	if err != nil {
		return nil, remainder, err
	}

	remainder = skipTrivia(remainder)
	if peek(remainder).Type != tok.TokOpRange {
		return nil, tokens, ErrNoMatch
	}
	remainder = remainder[1:]

	end, remainder, err := RangeBound(remainder)
	if err == ErrNoMatch {
		return nil, remainder, errorExpecting("range bound", remainder)
	} else if err != nil {
		return nil, remainder, err
	}

	return ast.NewRange(start, end), remainder, nil
}
