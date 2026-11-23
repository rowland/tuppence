package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// tuple_update_expression = expression tuple_update_tail .

func TupleUpdateExpression(tokens []tok.Token) (expr *ast.TupleUpdateExpression, remainder []tok.Token, err error) {
	var object ast.Expression
	if object, remainder, err = postfixExpressionWithTails(tokens, postfixBaseExpression, true, postfixTailWithoutTupleUpdate); err != nil {
		return nil, remainder, err
	}

	return tupleUpdateTail(object, remainder)
}

// tuple_update_tail = "." "(" tuple_update_member { "," tuple_update_member } [ "," ] ")" .

func tupleUpdateTail(object ast.Expression, tokens []tok.Token) (expr *ast.TupleUpdateExpression, remainder []tok.Token, err error) {
	remainder = skipTrivia(tokens)

	var found bool
	if remainder, found = Dot(remainder); !found {
		return nil, tokens, ErrNoMatch
	}

	if remainder, found = OpenParen(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokOpenParen, remainder)
	}
	var members []ast.TupleUpdateMember
	for {
		selector, rest, err := fieldSelector(remainder)
		if err == ErrNoMatch {
			return nil, rest, errorExpecting("field name or decimal index", rest)
		}
		if err != nil {
			return nil, rest, err
		}
		rest, found = Colon(rest)
		if !found {
			return nil, rest, errorExpectingTokenType(tok.TokColon, rest)
		}
		value, rest, err := Expression(rest)
		if err == ErrNoMatch {
			return nil, rest, errorExpecting("expression", rest)
		}
		if err != nil {
			return nil, rest, err
		}
		members = append(members, ast.TupleUpdateMember{Selector: selector, Value: value})
		remainder, found = Comma(rest)
		if !found || peek(skipTrivia(remainder)).Type == tok.TokCloseParen {
			break
		}
	}
	if remainder, found = CloseParen(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokCloseParen, remainder)
	}
	return ast.NewTupleUpdateExpression(object, members), remainder, nil
}
