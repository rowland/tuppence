package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// function_declaration = annotations function_declaration_lhs "=" function_declaration_type block .

func FunctionDeclaration(tokens []tok.Token) (*ast.FunctionDeclaration, []tok.Token, error) {
	annotations, remainder, err := Annotations(tokens)
	if err != nil {
		return nil, remainder, err
	}

	lhs, remainder, err := FunctionDeclarationLHS(remainder)
	if err != nil {
		return nil, remainder, err
	}

	var found bool
	if remainder, found = AssignOp(remainder); !found {
		return nil, tokens, ErrNoMatch
	}

	functionType, remainder, err := FunctionDeclarationType(remainder)
	if err != nil {
		return nil, remainder, err
	}

	body, remainder, err := Block(remainder)
	if err != nil {
		return nil, remainder, err
	}

	return ast.NewFunctionDeclaration(annotations.Annotations, lhs, functionType, body), remainder, nil
}

// function_declaration_lhs = function_identifier [ function_parameter_types ] .

func FunctionDeclarationLHS(tokens []tok.Token) (*ast.FunctionDeclarationLHS, []tok.Token, error) {
	name, remainder, err := FunctionIdentifier(tokens)
	if err != nil {
		return nil, remainder, err
	}

	parameterTypes, remainder, err := FunctionParameterTypes(remainder)
	if err != nil && err != ErrNoMatch {
		return nil, remainder, err
	}

	return ast.NewFunctionDeclarationLHS(name, parameterTypes), remainder, nil
}

// function_declaration_type = ( "fn" "(" [ labeled_parameters | parameters ] ")" ( return_type | "_" ) )
//                           | ( "fx" "(" [ labeled_parameters | parameters ] ")" [ return_type | "_" ] ) .

func FunctionDeclarationType(tokens []tok.Token) (*ast.FunctionDeclarationType, []tok.Token, error) {
	remainder := skipTrivia(tokens)
	hasSideEffects := false
	switch peek(remainder).Type {
	case tok.TokKwFn:
		remainder = remainder[1:]
	case tok.TokKwFx:
		hasSideEffects = true
		remainder = remainder[1:]
	default:
		return nil, tokens, ErrNoMatch
	}

	var found bool
	if remainder, found = OpenParen(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokOpenParen, remainder)
	}

	var parameters []ast.FunctionTypeParameter
	if remainder2, found := CloseParen(remainder); found {
		remainder = remainder2
	} else {
		if parameters, remainder, found = functionTypeParameters(remainder); !found {
			return nil, remainder, errorExpecting("function parameters", remainder)
		}

		if remainder, found = CloseParen(remainder); !found {
			return nil, remainder, errorExpectingTokenType(tok.TokCloseParen, remainder)
		}
	}

	if identifier, remainder2, err := Identifier(remainder); err == nil && identifier.Name == "_" {
		return ast.NewFunctionDeclarationType(hasSideEffects, parameters, nil, true), remainder2, nil
	}

	returnType, remainder, err := ReturnType(remainder)
	if err == nil {
		return ast.NewFunctionDeclarationType(hasSideEffects, parameters, returnType, false), remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if hasSideEffects {
		return ast.NewFunctionDeclarationType(true, parameters, nil, false), remainder, nil
	}

	return nil, remainder, errorExpecting("return type or _", remainder)
}

// block = "{" block_body "}" .

// block_body = { statement } expression .
