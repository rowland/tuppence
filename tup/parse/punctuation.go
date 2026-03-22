package parse

import (
	"slices"

	"github.com/rowland/tuppence/tup/tok"
)

func expectFunc(tokenTypes ...tok.TokenType) func([]tok.Token) (remainder []tok.Token, found bool) {
	return func(tokens []tok.Token) (remainder []tok.Token, found bool) {
		remainder = skipComments(tokens)

		if slices.Contains(tokenTypes, peek(remainder).Type) {
			return remainder[1:], true
		}

		return tokens, false
	}
}

// assign_op = "=" .

var AssignOp = expectFunc(tok.TokOpAssign)

// at = "@" .

var At = expectFunc(tok.TokAt)

// colon = ":" .

var Colon = expectFunc(tok.TokColon, tok.TokColonNoSpace)

// comma = "," .

var Comma = expectFunc(tok.TokComma)

// dollar = "$" .

var Dollar = expectFunc(tok.TokDollar)

// dot = "." .

var Dot = expectFunc(tok.TokDot)

// eol = "\r\n" | "\r" | "\n" .

var EOL = expectFunc(tok.TokEOL)

// pipe = "|" .

var Pipe = expectFunc(tok.TokOpBitOr)

// pipe_forward = "|>" .

var PipeForward = expectFunc(tok.TokOpPipe)

// semicolon = ";" .

var SemiColon = expectFunc(tok.TokSemiColon)

// rest_op = "..." .

var RestOp = expectFunc(tok.TokOpRest)

// star = "*" .

var Star = expectFunc(tok.TokOpMul)

// open_brace = "{" .

var OpenBrace = expectFunc(tok.TokOpenBrace)

// close_brace = "}" .

var CloseBrace = expectFunc(tok.TokCloseBrace)

// open_bracket = "[" .

var OpenBracket = expectFunc(tok.TokOpenBracket)

// close_bracket = "]" .

var CloseBracket = expectFunc(tok.TokCloseBracket)

// open_paren = "(" .

var OpenParen = expectFunc(tok.TokOpenParen)

// close_paren = ")" .

var CloseParen = expectFunc(tok.TokCloseParen)
