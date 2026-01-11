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
//                | "try_continue" [ expression ]
//                | "try_break" [ expression ] .

// func tryExpression(tokens []tok.Token) (item ast.TryExpression, remainder []tok.Token, err error) {
// 	return nil, nil, nil
// }

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
