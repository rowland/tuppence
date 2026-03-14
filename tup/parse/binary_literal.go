package parse

import (
	"strconv"
	"strings"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// binary_literal = "0b" ( "0" | "1" ) { "0" | "1" | "_" } .

func BinaryLiteral(tokens []tok.Token) (lit *ast.IntegerLiteral, remainder []tok.Token, err error) {
	remainder = skipTrivia(tokens)

	t := peek(remainder)
	if t.Type != tok.TokBinLit {
		return nil, tokens, ErrNoMatch
	} else if t.Invalid {
		return nil, remainder, errorExpectingTokenType(tok.TokBinLit, remainder)
	}

	value := t.Value()
	integerValue, _ := strconv.ParseInt(strings.ReplaceAll(value, "_", ""), 0, 64)
	return ast.NewBinaryLiteral(value, integerValue, t.File, t.Offset, t.Length), remainder[1:], nil
}
