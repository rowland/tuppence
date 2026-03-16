package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// tuple_update_expression = expression "." labeled_tuple_members .

func TupleUpdateExpression(tokens []tok.Token) (expr *ast.TupleUpdateExpression, remainder []tok.Token, err error) {
	var object ast.Expression
	if object, remainder, err = postfixExpressionWithTails(tokens, postfixBaseExpression, true, postfixTailWithoutTupleUpdate); err != nil {
		return nil, remainder, err
	}

	return tupleUpdateTail(object, remainder)
}

// tuple_update_tail = "." labeled_tuple_members .

func tupleUpdateTail(object ast.Expression, tokens []tok.Token) (expr *ast.TupleUpdateExpression, remainder []tok.Token, err error) {
	remainder = skipTrivia(tokens)

	var found bool
	if remainder, found = Dot(remainder); !found {
		return nil, tokens, ErrNoMatch
	}

	var updateMembers []*ast.TupleMember
	if updateMembers, remainder, err = labeledTupleMembers(remainder); err == ErrNoMatch {
		return nil, remainder, errorExpecting("labeled tuple members", remainder)
	} else if err != nil {
		return nil, remainder, err
	}

	return ast.NewTupleUpdateExpression(object, ast.NewTupleLiteral(true, updateMembers)), remainder, nil
}
