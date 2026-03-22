package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// type_qualified_function_declaration = annotations type_identifier "." function_declaration_lhs "=" function_declaration_type block .

func TypeQualifiedFunctionDeclaration(tokens []tok.Token) (*ast.TypeQualifiedFunctionDeclaration, []tok.Token, error) {
	annotations, remainder, err := Annotations(tokens)
	if err != nil {
		return nil, remainder, err
	}

	typeName, remainder, err := TypeIdentifier(remainder)
	if err != nil {
		return nil, remainder, err
	}

	var found bool
	if remainder, found = Dot(remainder); !found {
		return nil, tokens, ErrNoMatch
	}

	lhs, remainder, err := FunctionDeclarationLHS(remainder)
	if err != nil {
		return nil, remainder, err
	}

	if remainder, found = AssignOp(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokOpAssign, remainder)
	}

	functionType, remainder, err := FunctionDeclarationType(remainder)
	if err != nil {
		return nil, remainder, err
	}

	body, remainder, err := Block(remainder)
	if err != nil {
		return nil, remainder, err
	}

	return ast.NewTypeQualifiedFunctionDeclaration(
		typeName,
		ast.NewFunctionDeclaration(annotations.Annotations, lhs, functionType, body),
	), remainder, nil
}
