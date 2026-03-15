package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// tuple_update_expression = expression "." labeled_tuple_members .

func TupleUpdateExpression(tokens []tok.Token) (expr *ast.TupleUpdateExpression, remainder []tok.Token, err error) {
	remainder = skipTrivia(tokens)

	var object ast.Expression
	if object, remainder, err = tupleUpdateReceiver(remainder); err != nil {
		return nil, remainder, err
	}

	var found bool
	if remainder, found = Dot(remainder); !found {
		return nil, tokens, ErrNoMatch
	}

	var updateMembers []*ast.TupleMember
	if updateMembers, remainder, err = labeledTupleMembers(remainder); err == ErrNoMatch {
		return nil, tokens, ErrNoMatch
	} else if err != nil {
		return nil, remainder, err
	}

	return ast.NewTupleUpdateExpression(object, ast.NewTupleLiteral(true, updateMembers)), remainder, nil
}

// temporary helper function to parse the receiver of a tuple update expression
func tupleUpdateReceiver(tokens []tok.Token) (expr ast.Expression, remainder []tok.Token, err error) {
	if expression, remainder, err := parenthesizedExpression(tokens); err == nil {
		return expression, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if block, remainder, err := Block(tokens); err == nil {
		return block, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if functionCall, remainder, err := FunctionCall(tokens); err == nil {
		return functionCall, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if importExpression, remainder, err := ImportExpression(tokens); err == nil {
		return importExpression, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if typeofExpression, remainder, err := TypeofExpression(tokens); err == nil {
		return typeofExpression, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if metaExpression, remainder, err := MetaExpression(tokens); err == nil {
		return metaExpression, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if identifier, remainder, err := Identifier(tokens); err == nil {
		return ast.NewIdentifier(identifier.Name, identifier.Source, identifier.StartOffset, identifier.Length), remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if literal, remainder, err := Literal(tokens); err == nil {
		return literal, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	return nil, tokens, ErrNoMatch
}
