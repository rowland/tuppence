package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// expression = try_expression
//            | binary_expression
//            | unary_expression .

func Expression(tokens []tok.Token) (item ast.Expression, remainder []tok.Token, err error) {
	return nil, nil, nil
}

// try_expression = "try" expression
//                | "try_continue" expression
//                | "try_break" expression .

func TryExpression(tokens []tok.Token) (item *ast.TryExpression, remainder []tok.Token, err error) {
	remainder = skipComments(tokens)

	var variant ast.TryVariant
	tokenType := peek(remainder).Type

	switch tokenType {
	case tok.TokKwTry:
		variant = ast.TryStandard
	case tok.TokKwTryContinue:
		variant = ast.TryContinue
	case tok.TokKwTryBreak:
		variant = ast.TryBreak
	default:
		return nil, tokens, ErrNoMatch
	}

	expression, remainder, err := Expression(remainder)
	if err != nil {
		return nil, remainder, err
	}

	return ast.NewTryExpression(variant, expression), remainder, nil
}

// binary_expression = chained_expression .
// chained_expression = prec1_expression { "|>" function_call } .

// func binaryExpression(tokens []tok.Token) (item ast.BinaryExpression, remainder []tok.Token, err error) {
// 	return nil, nil, nil
// }

// unary_expression = ( "-" | "!" | "~" ) valid_negatable_expression
//                  | primary_expression .

// func unaryExpression(tokens []tok.Token) (item ast.UnaryExpression, remainder []tok.Token, err error) {
// 	return nil, nil, nil
// }
