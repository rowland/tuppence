package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// function_call_context = function_identifier [ "(" [ function_arguments ] ")" ] .

func FunctionCallContext(tokens []tok.Token) (*ast.FunctionCallContext, []tok.Token, error) {
	function, remainder, err := FunctionIdentifier(tokens)
	if err != nil {
		return nil, remainder, err
	}

	remainder2, found := OpenParen(remainder)
	if !found {
		return ast.NewFunctionCallContext(function, nil), remainder, nil
	}

	var arguments *ast.FunctionArguments
	if arguments, remainder2, err = FunctionArguments(remainder2); err != nil && err != ErrNoMatch {
		return nil, remainder2, err
	}

	if remainder2, found = CloseParen(remainder2); !found {
		return nil, remainder2, errorExpectingTokenType(tok.TokCloseParen, remainder2)
	}

	return ast.NewFunctionCallContext(function, arguments), remainder2, nil
}
