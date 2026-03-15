package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// array_function_call = "array" "(" type_identifier [ "," expression ] ")" .

func ArrayFunctionCall(tokens []tok.Token) (expr *ast.ArrayFunctionCall, remainder []tok.Token, err error) {
	remainder = skipTrivia(tokens)

	if peek(remainder).Type != tok.TokKwArray {
		return nil, tokens, ErrNoMatch
	}
	remainder = remainder[1:]

	var found bool
	if remainder, found = OpenParen(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokOpenParen, remainder)
	}

	var typeArg *ast.TypeIdentifier
	if typeArg, remainder, err = TypeIdentifier(remainder); err != nil {
		return nil, remainder, err
	}

	var sizeArg ast.Expression
	if remainder2, found := Comma(remainder); found {
		if sizeArg, remainder, err = Expression(remainder2); err != nil {
			return nil, remainder, err
		}
	}

	if remainder, found = CloseParen(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokCloseParen, remainder)
	}

	return ast.NewArrayFunctionCall(typeArg, sizeArg), remainder, nil
}
