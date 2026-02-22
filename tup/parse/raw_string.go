package parse

import (
	"strings"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// raw_string_literal = "`" { "``" | character - "`" } "`" .

func RawStringLiteral(tokens []tok.Token) (item *ast.RawStringLiteral, remainder []tok.Token, err error) {
	remainder = skipComments(tokens)

	if peek(remainder).Type != tok.TokRawStrLit {
		return nil, nil, ErrNoMatch
	} else if peek(remainder).Invalid {
		return nil, remainder, errorExpecting(tok.TokenTypes[tok.TokRawStrLit], remainder)
	}

	value := remainder[0].Value()
	stringValue := value[1 : len(remainder[0].Value())-1]
	stringValue = strings.ReplaceAll(stringValue, "``", "`")
	return ast.NewRawStringLiteral(value, stringValue, remainder[0].File, remainder[0].Offset, remainder[0].Length), remainder[1:], nil
}
