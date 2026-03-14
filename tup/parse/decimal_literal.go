package parse

import (
	"strconv"
	"strings"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// decimal_literal = decimal_digit { decimal_digit | "_" } .
// decimal_digit = "0"-"9" .

func DecimalLiteral(tokens []tok.Token) (lit *ast.IntegerLiteral, remainder []tok.Token, err error) {
	remainder = skipTrivia(tokens)

	t := peek(remainder)
	if t.Type != tok.TokDecLit {
		return nil, tokens, ErrNoMatch
	} else if t.Invalid {
		return nil, remainder, errorExpectingTokenType(tok.TokDecLit, remainder)
	}

	value := t.Value()
	integerValue, _ := strconv.ParseInt(strings.ReplaceAll(value, "_", ""), 10, 64)
	return ast.NewDecimalLiteral(value, integerValue, t.File, t.Offset, t.Length), remainder[1:], nil
}
