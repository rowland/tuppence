package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// type_qualified_declaration = type_identifier "." identifier "=" expression .

func TypeQualifiedDeclaration(tokens []tok.Token) (*ast.TypeQualifiedDeclaration, []tok.Token, error) {
	typeName, remainder, err := TypeIdentifier(tokens)
	if err != nil {
		return nil, remainder, err
	}

	var found bool
	if remainder, found = Dot(remainder); !found {
		return nil, tokens, ErrNoMatch
	}

	name, remainder, err := Identifier(remainder)
	if err != nil {
		return nil, remainder, err
	}

	if remainder, found = AssignOp(remainder); !found {
		return nil, tokens, ErrNoMatch
	}

	expression, remainder, err := Expression(remainder)
	if err != nil {
		return nil, remainder, err
	}

	return ast.NewTypeQualifiedDeclaration(
		typeName,
		ast.NewAssignment(
			ast.NewOrdinalAssignmentLHS([]*ast.Identifier{name}, nil),
			ast.Immutable,
			expression,
		),
	), remainder, nil
}
