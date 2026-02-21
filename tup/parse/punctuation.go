package parse

import (
	"slices"

	"github.com/rowland/tuppence/tup/tok"
)

func expectFunc(tokenTypes ...tok.TokenType) func([]tok.Token) (remainder []tok.Token, err error) {
	return func(tokens []tok.Token) (remainder []tok.Token, err error) {
		remainder = skipComments(tokens)

		if slices.Contains(tokenTypes, peek(remainder).Type) {
			return remainder[1:], nil
		}

		return nil, ErrNoMatch
	}
}

// at = "@" .

var At = expectFunc(tok.TokAt)

// close_bracket = "]" .

var CloseBracket = expectFunc(tok.TokCloseBracket)

// colon = ":" .

var Colon = expectFunc(tok.TokColon, tok.TokColonNoSpace)

// dot = "." .

var Dot = expectFunc(tok.TokDot)

// eol = "\r\n" | "\r" | "\n" .

var EOL = expectFunc(tok.TokEOL)

// open_bracket = "[" .

var OpenBracket = expectFunc(tok.TokOpenBracket)
