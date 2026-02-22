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
	return ChainedExpression(tokens)
}

// chained_expression = prec1_expression { "|>" function_call } .

func ChainedExpression(tokens []tok.Token) (item *ast.ChainedExpression, remainder []tok.Token, err error) {
	return nil, nil, nil
}

// prec1_expression = prec2_expression { logical_or_op prec2_expression } .

func Prec1Expression(tokens []tok.Token) (item *ast.Prec1Expression, remainder []tok.Token, err error) {
	return nil, nil, nil
}

// unary_expression = prefixed_unary_expression
//                  | primary_expression .

func UnaryExpression(tokens []tok.Token) (item *ast.UnaryExpression, remainder []tok.Token, err error) {
	remainder = skipComments(tokens)

	expression, remainder, err := PrefixedUnaryExpression(remainder)
	if err == nil {
		return expression, remainder, nil
	} else if err != ErrNoMatch {
		return nil, tokens, err
	}

	primaryExpression, remainder, err := PrimaryExpression(remainder)
	if err == nil {
		return ast.NewUnaryExpression(nil, primaryExpression), remainder, nil
	} else if err != ErrNoMatch {
		return nil, tokens, err
	}

	return nil, tokens, errorExpecting("unary expression", tokens)
}

// prefixed_unary_expression = unary_op negatable_expression .

func PrefixedUnaryExpression(tokens []tok.Token) (item *ast.UnaryExpression, remainder []tok.Token, err error) {
	operator, remainder, err := UnaryOp(remainder)
	if err != nil {
		return nil, tokens, err
	}

	expression, remainder, err := NegatableExpression(remainder)
	if err != nil {
		return nil, remainder, err
	}

	return ast.NewUnaryExpression(&operator, expression), remainder, nil
}

// negatable_expression = "(" expression ")"
//                      | block
//                      | function_call
//                      | member_access
//                      | tuple_update_expression
//                      | safe_indexed_access
//                      | indexed_access
//                      | identifier
//                      | literal .

func NegatableExpression(tokens []tok.Token) (item ast.Expression, remainder []tok.Token, err error) {
	remainder = skipComments(tokens)

	if expression, remainder, err := parenthesizedExpression(remainder); err == nil {
		return expression, remainder, nil
	} else if err != ErrNoMatch {
		return nil, tokens, err
	}

	// block
	// function_call
	// member_access
	// tuple_update_expression
	// safe_indexed_access
	// indexed_access

	if identifier, remainder, err := Identifier(remainder); err == nil {
		return ast.NewIdentifier(identifier.Name, identifier.Source, identifier.StartOffset, identifier.Length), remainder, nil
	} else if err != ErrNoMatch {
		return nil, tokens, err
	}

	if literal, remainder, err := Literal(remainder); err == nil {
		return literal, remainder, nil
	} else if err != ErrNoMatch {
		return nil, tokens, err
	}

	return nil, tokens, ErrNoMatch
}

// primary_expression = "(" expression ")"
//                    | block
//                    | if_expression
//                    | for_expression
//                    | inline_for_expression
//                    | array_function_call
//                    | import_expression
//                    | typeof_expression
//                    | function_call
//                    | type_constructor_call
//                    | return_expression
//                    | break_expression
//                    | continue_expression
//                    | member_access
//                    | tuple_update_expression
//                    | safe_indexed_access
//                    | indexed_access
//                    | range
//                    | identifier
//                    | literal .

func PrimaryExpression(tokens []tok.Token) (item ast.Expression, remainder []tok.Token, err error) {
	if expression, remainder, err := parenthesizedExpression(remainder); err == nil {
		return expression, remainder, nil
	} else if err != ErrNoMatch {
		return nil, tokens, err
	}

	// block
	// if_expression
	// for_expression
	// inline_for_expression
	// array_function_call
	// import_expression
	// typeof_expression
	// function_call
	// type_constructor_call
	// return_expression
	// break_expression
	// continue_expression
	// member_access
	// tuple_update_expression
	// safe_indexed_access
	// indexed_access

	if identifier, remainder, err := Identifier(remainder); err == nil {
		return ast.NewIdentifier(identifier.Name, identifier.Source, identifier.StartOffset, identifier.Length), remainder, nil
	} else if err != ErrNoMatch {
		return nil, tokens, err
	}

	if literal, remainder, err := Literal(remainder); err == nil {
		return literal, remainder, nil
	} else if err != ErrNoMatch {
		return nil, tokens, err
	}

	return nil, tokens, ErrNoMatch
}

// parenthesized_expression = "(" expression ")" .

func parenthesizedExpression(tokens []tok.Token) (item ast.Expression, remainder []tok.Token, err error) {
	remainder = skipComments(tokens)

	remainder, err = OpenParen(remainder)
	if err != nil {
		return nil, tokens, err
	}

	expression, remainder, err := Expression(remainder)
	if err != nil {
		return nil, remainder, err
	}

	remainder, err = CloseParen(remainder)
	if err != nil {
		return nil, remainder, err
	}

	return expression, remainder, nil
}
