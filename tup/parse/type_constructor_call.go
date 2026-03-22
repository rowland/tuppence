package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// type_constructor_call = type_reference [ function_parameter_types ] "(" [ function_arguments ] ")" [ function_block ] .

func TypeConstructorCall(tokens []tok.Token) (expr *ast.TypeConstructorCall, remainder []tok.Token, err error) {
	var typeReference *ast.TypeReference
	if typeReference, remainder, err = TypeReference(tokens); err != nil {
		return nil, remainder, err
	}

	var parameterTypes *ast.FunctionParameterTypes
	if parameterTypes, remainder, err = FunctionParameterTypes(remainder); err != nil && err != ErrNoMatch {
		return nil, remainder, err
	}

	var found bool
	if remainder, found = OpenParen(remainder); !found {
		return nil, tokens, ErrNoMatch
	}

	var arguments *ast.FunctionArguments
	if arguments, remainder, err = FunctionArguments(remainder); err != nil && err != ErrNoMatch {
		return nil, remainder, err
	}

	if remainder, found = CloseParen(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokCloseParen, remainder)
	}

	var functionBlock *ast.FunctionBlock
	if functionBlock, remainder, err = FunctionBlock(remainder); err != nil && err != ErrNoMatch {
		return nil, remainder, err
	}

	return ast.NewTypeConstructorCall(typeReference, parameterTypes, arguments, functionBlock), remainder, nil
}
