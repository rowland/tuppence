package parse

import "github.com/rowland/tuppence/tup/tok"

// at = "@" .

func At(tokens []tok.Token) (remainder []tok.Token, err error) {
	remainder = skipComments(tokens)
	if peek(remainder).Type != tok.TokAt {
		return nil, errorExpecting(tok.TokenTypes[tok.TokAt], remainder)
	}
	return remainder[1:], nil
}

// colon = ":" .

func Colon(tokens []tok.Token) (remainder []tok.Token, err error) {
	remainder = skipComments(tokens)
	if peek(remainder).Type != tok.TokColon {
		return nil, errorExpecting(tok.TokenTypes[tok.TokColon], remainder)
	}
	return remainder[1:], nil
}

// dot = "." .

func Dot(tokens []tok.Token) (remainder []tok.Token, err error) {
	remainder = skipComments(tokens)
	if peek(remainder).Type != tok.TokDot {
		return nil, errorExpecting(tok.TokenTypes[tok.TokDot], remainder)
	}
	return remainder[1:], nil
}
