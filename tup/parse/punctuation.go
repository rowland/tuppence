package parse

import "github.com/rowland/tuppence/tup/tok"

func expectFunc(tokenType tok.TokenType) func([]tok.Token) (remainder []tok.Token, err error) {
	return func(tokens []tok.Token) (remainder []tok.Token, err error) {
		remainder = skipComments(tokens)
		if peek(remainder).Type != tokenType {
			return remainder, errorExpecting(tok.TokenTypes[tokenType], remainder)
		}
		return remainder[1:], nil
	}
}

// at = "@" .

var At = expectFunc(tok.TokAt)

// colon = ":" .

var Colon = expectFunc(tok.TokColon)

// dot = "." .

var Dot = expectFunc(tok.TokDot)

// eol = "\r\n" | "\r" | "\n" .

var EOL = expectFunc(tok.TokEOL)
