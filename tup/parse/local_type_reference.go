package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// local_type_reference = type_reference | identifier .

func LocalTypeReference(tokens []tok.Token) (item ast.LocalTypeReference, remainder []tok.Token, err error) {
	if typeReference, remainder, err := TypeReference(tokens); err == nil {
		return typeReference, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if identifier, remainder, err := Identifier(tokens); err == nil {
		return identifier, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	return nil, tokens, ErrNoMatch
}
