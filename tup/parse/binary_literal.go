package parse

import (
	"strconv"
	"strings"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// binary_literal = "0b" ( "0" | "1" ) { "0" | "1" | "_" } .

func BinaryLiteral(tokens []tok.Token) (item *ast.BinaryLiteral, remainder []tok.Token, err error) {
	remainder = skipComments(tokens)
	if peek(remainder).Type != tok.TokBinLit || peek(remainder).Invalid {
		return nil, nil, errorExpecting(tok.TokenTypes[tok.TokBinLit], remainder)
	}
	value := remainder[0].Value()
	integerValue, _ := strconv.ParseInt(strings.ReplaceAll(value, "_", ""), 0, 64)
	return ast.NewBinaryLiteral(value, integerValue, remainder[0].File, remainder[0].Offset, remainder[0].Length), remainder[1:], nil
}
