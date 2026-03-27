package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// function_type_declaration = function_type_declaration_lhs "=" function_type .

func FunctionTypeDeclaration(tokens []tok.Token) (*ast.FunctionTypeDeclaration, []tok.Token, error) {
	name, parameterTypes, remainder, err := FunctionTypeDeclarationLHS(tokens)
	if err != nil {
		return nil, remainder, err
	}

	var found bool
	if remainder, found = AssignOp(remainder); !found {
		return nil, tokens, ErrNoMatch
	}

	functionType, remainder, err := FunctionType(remainder)
	if err == ErrNoMatch {
		return nil, tokens, ErrNoMatch
	} else if err != nil {
		return nil, remainder, err
	}

	return ast.NewFunctionTypeDeclaration(name, parameterTypes, functionType), remainder, nil
}

// function_type_declaration_lhs = function_type_identifier [ function_parameter_types ] .

func FunctionTypeDeclarationLHS(tokens []tok.Token) (*ast.TypeIdentifier, *ast.FunctionParameterTypes, []tok.Token, error) {
	name, remainder, err := FunctionTypeIdentifier(tokens)
	if err != nil {
		return nil, nil, remainder, err
	}

	parameterTypes, remainder, err := FunctionParameterTypes(remainder)
	if err == ErrNoMatch {
		return name, nil, remainder, nil
	} else if err != nil {
		return nil, nil, remainder, err
	}

	return name, parameterTypes, remainder, nil
}

// function_type_identifier = type_identifier .

func FunctionTypeIdentifier(tokens []tok.Token) (*ast.TypeIdentifier, []tok.Token, error) {
	return TypeIdentifier(tokens)
}
