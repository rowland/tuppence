package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// expression = try_expression
//            | binary_expression
//            | unary_expression .

func Expression(tokens []tok.Token) (expr ast.Expression, remainder []tok.Token, err error) {
	if tryExpr, remainder, err := TryExpression(tokens); err == nil {
		return tryExpr, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if binaryExpr, remainder, err := BinaryExpression(tokens); err == nil {
		return binaryExpr, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if unaryExpr, remainder, err := UnaryExpression(tokens); err == nil {
		return unaryExpr, remainder, nil
	} else if err != ErrNoMatch {
		return nil, tokens, err
	}

	return nil, tokens, errorExpecting("expression", tokens)
}

// try_expression = "try" expression
//                | "try_continue" expression
//                | "try_break" expression .

func TryExpression(tokens []tok.Token) (expr *ast.TryExpression, remainder []tok.Token, err error) {
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

func BinaryExpression(tokens []tok.Token) (expr ast.Expression, remainder []tok.Token, err error) {
	return ChainedExpression(tokens)
}

// chained_expression = logical_or_expression { "|>" function_call } .

func ChainedExpression(tokens []tok.Token) (expr ast.Expression, remainder []tok.Token, err error) {
	remainder = skipComments(tokens)

	initial, remainder, err := LogicalOrExpression(remainder)
	if err == ErrNoMatch {
		return nil, tokens, err
	} else if err != nil {
		return nil, remainder, err
	}

	var functionCalls []*ast.FunctionCall

	for {
		remainder, err = PipeOp(remainder)
		if err == ErrNoMatch {
			break
		} else if err != nil {
			return nil, remainder, err
		}

		var functionCall *ast.FunctionCall
		functionCall, remainder, err = FunctionCall(remainder)
		if err != nil {
			return nil, remainder, err
		}
		functionCalls = append(functionCalls, functionCall)
	}

	if len(functionCalls) == 0 {
		return initial, remainder, nil
	}

	return ast.NewChainedExpression(initial, functionCalls), remainder, nil
}

// logical_or_expression = logical_and_expression { logical_or_op logical_and_expression } .

func LogicalOrExpression(tokens []tok.Token) (expr ast.Expression, remainder []tok.Token, err error) {
	initial, remainder, err := LogicalAndExpression(remainder)
	if err == ErrNoMatch {
		return nil, tokens, err
	} else if err != nil {
		return nil, remainder, err
	}

	operands := []ast.Expression{initial}
	for {
		remainder, err = LogicalOrOp(remainder)
		if err == ErrNoMatch {
			break
		} else if err != nil {
			return nil, remainder, err
		}

		var operand ast.Expression
		operand, remainder, err = LogicalAndExpression(remainder)
		if err != nil {
			return nil, remainder, err
		}

		operands = append(operands, operand)
	}

	if len(operands) == 1 {
		return operands[0], remainder, nil
	}

	return ast.NewLogicalOrExpression(operands), remainder, nil
}

// logical_and_expression = comparison_expression { logical_and_op comparison_expression } .

func LogicalAndExpression(tokens []tok.Token) (expr ast.Expression, remainder []tok.Token, err error) {
	initial, remainder, err := ComparisonExpression(remainder)
	if err == ErrNoMatch {
		return nil, tokens, err
	} else if err != nil {
		return nil, remainder, err
	}

	operands := []ast.Expression{initial}
	for {
		remainder, err = LogicalAndOp(remainder)
		if err == ErrNoMatch {
			break
		} else if err != nil {
			return nil, remainder, err
		}

		var operand ast.Expression
		operand, remainder, err = ComparisonExpression(remainder)
		if err != nil {
			return nil, remainder, err
		}

		operands = append(operands, operand)
	}

	if len(operands) == 1 {
		return operands[0], remainder, nil
	}

	return ast.NewLogicalAndExpression(operands), remainder, nil
}

// comparison_expression = add_sub_expression [ type_comparison_tail | relational_comparison_tail ] .

func ComparisonExpression(tokens []tok.Token) (expr ast.Expression, remainder []tok.Token, err error) {
	left, remainder, err := AddSubExpression(remainder)
	if err == ErrNoMatch {
		return nil, tokens, err
	} else if err != nil {
		return nil, remainder, err
	}

	typeComparison, remainder, err := TypeComparisonTail(left, remainder)
	if err == nil {
		return typeComparison, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	relationalComparison, remainder, err := RelationalComparisonTail(left, remainder)
	if err == nil {
		return relationalComparison, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	return left, remainder, nil
}

// add_sub_expression = mul_div_expression { add_sub_op mul_div_expression } .

func AddSubExpression(tokens []tok.Token) (expr ast.Expression, remainder []tok.Token, err error) {
	left, remainder, err := MulDivExpression(tokens)
	if err == ErrNoMatch {
		return nil, tokens, err
	} else if err != nil {
		return nil, remainder, err
	}

	// Parse { add_sub_op mul_div_expression }
	for {
		op, remainder2, err := AddSubOp(remainder)
		if err == ErrNoMatch {
			return left, remainder, nil // no tail => done
		} else if err != nil {
			return nil, remainder2, err
		}

		var right ast.Expression
		right, remainder, err = MulDivExpression(remainder2)
		if err == ErrNoMatch {
			return nil, remainder, errorExpecting("mul_div expression", remainder2)
		} else if err != nil {
			return nil, remainder, err
		}

		left = ast.NewAddSubExpression(left, op, right)
	}
}

// mul_div_expression = pow_expression { mul_div_op pow_expression } .

func MulDivExpression(tokens []tok.Token) (expr ast.Expression, remainder []tok.Token, err error) {
	left, remainder, err := PowExpression(remainder)
	if err == ErrNoMatch {
		return nil, tokens, err
	} else if err != nil {
		return nil, remainder, err
	}

	// Parse { mul_div_op pow_expression }
	for {
		op, remainder2, err := MulDivOp(remainder)
		if err == ErrNoMatch {
			return left, remainder, nil // no tail => done
		} else if err != nil {
			return nil, remainder2, err
		}

		var right ast.Expression
		right, remainder, err = PowExpression(remainder2)
		if err == ErrNoMatch {
			return nil, remainder, errorExpecting("pow expression", remainder2)
		} else if err != nil {
			return nil, remainder, err
		}

		left = ast.NewMulDivExpression(left, op, right)
	}
}

// pow_expression = unary_expression { "^" unary_expression } .

func PowExpression(tokens []tok.Token) (expr ast.Expression, remainder []tok.Token, err error) {
	left, remainder, err := UnaryExpression(remainder)
	if err == ErrNoMatch {
		return nil, tokens, err
	} else if err != nil {
		return nil, remainder, err
	}

	operands := []ast.Expression{left}

	for {
		remainder, err = PowOp(remainder)
		if err == ErrNoMatch {
			break
		} else if err != nil {
			return nil, remainder, err
		}

		var operand ast.Expression
		operand, remainder, err = UnaryExpression(remainder)
		if err != nil {
			return nil, remainder, err
		}

		operands = append(operands, operand)
	}

	if len(operands) == 1 {
		return operands[0], remainder, nil
	}

	return ast.NewPowExpression(operands), remainder, nil
}

// unary_expression = prefixed_unary_expression
//                  | primary_expression .

func UnaryExpression(tokens []tok.Token) (expr ast.Expression, remainder []tok.Token, err error) {
	remainder = skipComments(tokens)

	if unaryExpression, remainder, err := PrefixedUnaryExpression(remainder); err == nil {
		return unaryExpression, remainder, nil
	} else if err != ErrNoMatch {
		return nil, tokens, err
	}

	if primaryExpression, remainder, err := PrimaryExpression(remainder); err == nil {
		return primaryExpression, remainder, nil
	} else if err != ErrNoMatch {
		return nil, tokens, err
	}

	return nil, tokens, errorExpecting("unary expression", tokens)
}

// prefixed_unary_expression = unary_op negatable_expression .

func PrefixedUnaryExpression(tokens []tok.Token) (expr *ast.UnaryExpression, remainder []tok.Token, err error) {
	operator, remainder, err := UnaryOp(remainder)
	if err != nil {
		return nil, tokens, err
	}

	expression, remainder, err := NegatableExpression(remainder)
	if err != nil {
		return nil, remainder, err
	}

	return ast.NewUnaryExpression(operator, expression), remainder, nil
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

func NegatableExpression(tokens []tok.Token) (expr ast.Expression, remainder []tok.Token, err error) {
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

func PrimaryExpression(tokens []tok.Token) (expr ast.Expression, remainder []tok.Token, err error) {
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

func parenthesizedExpression(tokens []tok.Token) (expr ast.Expression, remainder []tok.Token, err error) {
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

// function_call = function_identifier [ function_parameter_types ] "(" [ function_arguments ] ")" [ function_block ] .

func FunctionCall(tokens []tok.Token) (expr *ast.FunctionCall, remainder []tok.Token, err error) {
	return nil, nil, ErrNoMatch // TODO: Implement
}

// type_comparison_tail = is_op type_predicate .

func TypeComparisonTail(left ast.Expression, tokens []tok.Token) (expr ast.Expression, remainder []tok.Token, err error) {
	remainder, err = IsOp(tokens)
	if err == ErrNoMatch {
		return nil, tokens, err
	} else if err != nil {
		return nil, remainder, err
	}

	right, remainder, err := TypePredicate(remainder)
	if err == ErrNoMatch {
		return nil, remainder, errorExpecting("type predicate", remainder)
	} else if err != nil {
		return nil, remainder, err
	}

	return ast.NewTypeComparison(left, right), remainder, nil
}

// type_predicate = type_reference | inline_union .

func TypePredicate(tokens []tok.Token) (expr ast.TypePredicate, remainder []tok.Token, err error) {
	return nil, nil, ErrNoMatch // TODO: Implement
}

// relational_comparison_tail = rel_op add_sub_expression .

func RelationalComparisonTail(left ast.Expression, tokens []tok.Token) (expr ast.Expression, remainder []tok.Token, err error) {
	operator, remainder, err := RelOp(tokens)
	if err == ErrNoMatch {
		return nil, tokens, err
	} else if err != nil {
		return nil, remainder, err
	}

	right, remainder, err := AddSubExpression(remainder)
	if err != nil {
		return nil, remainder, err
	}

	return ast.NewRelationalComparison(left, operator, right), remainder, nil
}
