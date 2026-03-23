package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// keywordOptionalExpression parses a control-flow keyword followed by an
// optional same-statement expression payload. Newlines, semicolons, and block
// endings terminate the statement before any payload can begin.
func keywordOptionalExpression(tokens []tok.Token, keyword tok.TokenType) ([]tok.Token, ast.Expression, bool, error) {
	remainder := skipTrivia(tokens)
	if peek(remainder).Type != keyword {
		return tokens, nil, false, ErrNoMatch
	}
	remainder = remainder[1:]

	remainder = skipComments(remainder)
	switch peek(remainder).Type {
	case tok.TokEOF, tok.TokEOL, tok.TokSemiColon, tok.TokCloseBrace:
		return remainder, nil, true, nil
	}

	expression, remainder, err := Expression(remainder)
	if err == ErrNoMatch {
		return remainder, nil, true, nil
	} else if err != nil {
		return remainder, nil, true, err
	}

	return remainder, expression, true, nil
}
