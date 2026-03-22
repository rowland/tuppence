package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// scoped_function_identifier = identifier { "." identifier } "." function_identifier
//                            | function_identifier .

func ScopedFunctionIdentifier(tokens []tok.Token) (ast.FunctionCallContextFunction, []tok.Token, error) {
	scope := []*ast.Identifier{}
	remainder := tokens
	for {
		identifier, next, err := Identifier(remainder)
		if err != nil {
			break
		}

		var found bool
		if next, found = Dot(next); !found {
			break
		}

		scope = append(scope, identifier)
		remainder = next

		function, next, err := FunctionIdentifier(remainder)
		if err == nil {
			return ast.NewScopedFunctionIdentifier(scope, function), next, nil
		} else if err != ErrNoMatch {
			return nil, next, err
		}
	}

	if function, remainder, err := FunctionIdentifier(tokens); err == nil {
		return function, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	return nil, tokens, ErrNoMatch
}

// function_call_context = scoped_function_identifier [ "(" [ function_arguments ] ")" ] .

func FunctionCallContext(tokens []tok.Token) (*ast.FunctionCallContext, []tok.Token, error) {
	function, remainder, err := ScopedFunctionIdentifier(tokens)
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
