package parse

import (
	"strconv"
	"strings"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// octal_literal = "0o" octal_digit { octal_digit } .
// octal_digit = "0"-"7" .

func OctalLiteral(tokens []tok.Token) (lit *ast.IntegerLiteral, remainder []tok.Token, err error) {
	remainder = skipComments(tokens)

	t := peek(remainder)
	if t.Type != tok.TokOctLit {
		return nil, tokens, ErrNoMatch
	} else if t.Invalid {
		return nil, remainder, errorExpectingTokenType(tok.TokOctLit, remainder)
	}

	value := t.Value()
	integerValue, _ := strconv.ParseInt(strings.ReplaceAll(value, "_", ""), 0, 64)
	return ast.NewOctalLiteral(value, integerValue, t.File, t.Offset, t.Length), remainder[1:], nil
}
