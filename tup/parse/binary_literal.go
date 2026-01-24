package parse

import (
	"strconv"
	"strings"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

func BinaryLiteral(tokens []tok.Token) (item *ast.BinaryLiteral, remainder []tok.Token, err error) {
	remainder = skipComments(tokens)
	if peek(remainder).Type != tok.TokBinLit || peek(remainder).Invalid {
		return nil, nil, errorExpecting(tok.TokenTypes[tok.TokBinLit], remainder)
	}
	value := remainder[0].Value()
	integerValue, err := strconv.ParseInt(strings.ReplaceAll(value, "_", ""), 0, 64)
	if err != nil {
		return nil, nil, errorExpecting(tok.TokenTypes[tok.TokBinLit], remainder)
	}
	return ast.NewBinaryLiteral(value, integerValue, remainder[0].File, remainder[0].Offset, remainder[0].Length), remainder[1:], nil
}
