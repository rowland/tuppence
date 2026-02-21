package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// expression = try_expression
//            | binary_expression
//            | unary_expression .

func Expression(tokens []tok.Token) (item ast.Expression, remainder []tok.Token, err error) {
	tryExpr, remainder, err := TryExpression(tokens)
	if err == nil {
		return tryExpr, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	binaryExpr, remainder, err := BinaryExpression(tokens)
	if err == nil {
		return binaryExpr, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	unaryExpr, remainder, err := UnaryExpression(tokens)
	if err == nil {
		return unaryExpr, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	return nil, tokens, errorExpecting("expression", tokens)
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

func BinaryExpression(tokens []tok.Token) (item *ast.BinaryExpression, remainder []tok.Token, err error) {
	return nil, nil, nil
}

// chained_expression = prec1_expression { "|>" function_call } .

func ChainedExpression(tokens []tok.Token) (item *ast.ChainedExpression, remainder []tok.Token, err error) {
	return nil, nil, nil
}

// unary_expression = ( sub_op | logical_not_op | bit_not_op ) valid_negatable_expression
//                  | primary_expression .

func UnaryExpression(tokens []tok.Token) (item *ast.UnaryExpression, remainder []tok.Token, err error) {
	remainder = skipComments(tokens)

	var operator ast.Operator
	switch {
	case peek(remainder).Type == tok.TokOpMinus:
		operator = ast.Operator("-")
	case peek(remainder).Type == tok.TokOpNot:
		operator = ast.Operator("!")
	case peek(remainder).Type == tok.TokOpBitNot:
		operator = ast.Operator("~")
	default:
		return nil, remainder, ErrNoMatch
	}

	expression, remainder, err := Expression(remainder)
	if err != nil {
		return nil, remainder, err
	}

	return ast.NewUnaryExpression(operator, expression), remainder, nil
}
