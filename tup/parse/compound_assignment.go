package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// compound_assignment = identifier compound_assignment_op expression .

func CompoundAssignment(tokens []tok.Token) (assignment *ast.CompoundAssignment, remainder []tok.Token, err error) {
	var left *ast.Identifier
	if left, remainder, err = Identifier(tokens); err != nil {
		return nil, remainder, err
	}

	var operator ast.CompoundAssignmentOp
	var found bool
	if operator, remainder, found = CompoundAssignmentOp(remainder); !found {
		return nil, tokens, ErrNoMatch
	}

	var right ast.Expression
	if right, remainder, err = Expression(remainder); err != nil {
		return nil, remainder, err
	}

	return ast.NewCompoundAssignment(left, operator, right), remainder, nil
}
