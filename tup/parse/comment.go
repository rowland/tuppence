package parse

import "github.com/rowland/tuppence/tup/tok"

func skipComments(tokens []tok.Token) []tok.Token {
	for len(tokens) > 0 && tokens[0].Type == tok.TokComment {
		tokens = tokens[1:]
	}
	return tokens
}
