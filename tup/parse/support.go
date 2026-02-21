package parse

import (
	"errors"

	"github.com/rowland/tuppence/tup/tok"
)

func peek(tokens []tok.Token) tok.Token {
	if len(tokens) == 0 {
		return tok.Token{Type: tok.TokEOF}
	}
	return tokens[0]
}

var ErrNoMatch = errors.New("no match")
