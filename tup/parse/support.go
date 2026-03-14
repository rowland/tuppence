package parse

import (
	"slices"

	"github.com/rowland/tuppence/tup/tok"
)

func peek(tokens []tok.Token) tok.Token {
	if len(tokens) == 0 {
		return tok.Token{Type: tok.TokEOF}
	}
	return tokens[0]
}

func skip(tokens []tok.Token, tokenTypes ...tok.TokenType) []tok.Token {
	for len(tokens) > 0 && slices.Contains(tokenTypes, tokens[0].Type) {
		tokens = tokens[1:]
	}
	return tokens
}

func skipComments(tokens []tok.Token) []tok.Token {
	return skip(tokens, tok.TokComment)
}

func skipTrivia(tokens []tok.Token) []tok.Token {
	return skip(tokens, tok.TokComment, tok.TokEOL)
}
