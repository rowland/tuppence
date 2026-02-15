package parse

import "github.com/rowland/tuppence/tup/tok"

func expectFunc(tokenTypes ...tok.TokenType) func([]tok.Token) (remainder []tok.Token, err error) {
	return func(tokens []tok.Token) (remainder []tok.Token, err error) {
		remainder = skipComments(tokens)
		for _, tokenType := range tokenTypes {
			if peek(remainder).Type == tokenType {
				return remainder[1:], nil
			}
			err = errorExpecting(tok.TokenTypes[tokenType], remainder)
		}
		if len(tokenTypes) > 0 {
			return nil, errorExpecting(tok.TokenTypes[tokenTypes[0]], remainder)
		}
		panic("unexpected token type")
	}
}

// at = "@" .

var At = expectFunc(tok.TokAt)

// colon = ":" .

var Colon = expectFunc(tok.TokColon, tok.TokColonNoSpace)

// dot = "." .

var Dot = expectFunc(tok.TokDot)

// eol = "\r\n" | "\r" | "\n" .

var EOL = expectFunc(tok.TokEOL)
