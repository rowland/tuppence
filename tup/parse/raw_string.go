package parse

import (
	"strings"

	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// raw_string_literal = "`" { "``" | character - "`" } "`" .

func RawStringLiteral(tokens []tok.Token) (lit *ast.RawStringLiteral, remainder []tok.Token, err error) {
	// fmt.Println("RawStringLiteral", tokens)
	remainder = skipTrivia(tokens)

	t := peek(remainder)
	if t.Type != tok.TokRawStrLit {
		return nil, nil, ErrNoMatch
	} else if t.Invalid {
		return nil, remainder, errorExpectingTokenType(tok.TokRawStrLit, remainder)
	}

	value := t.Value()
	stringValue := value[1 : len(value)-1]
	stringValue = strings.ReplaceAll(stringValue, "``", "`")
	return ast.NewRawStringLiteral(value, stringValue, t.File, t.Offset, t.Length), remainder[1:], nil
}
