package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// expression = try_expression
//            | binary_expression
//            | unary_expression .

func Expression(tokens []tok.Token) (expr ast.Expression, remainder []tok.Token, err error) {
	// fmt.Println("Expression", tok.Types(tokens))
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

	// fmt.Println("UnaryExpression", tokens)
	if unaryExpr, remainder, err := UnaryExpression(tokens); err == nil {
		return unaryExpr, remainder, nil
	} else if err != ErrNoMatch {
		return nil, tokens, err
	}

	return nil, tokens, ErrNoMatch
}

// try_expression = "try" expression
//                | "try_continue" expression
//                | "try_break" expression .

func TryExpression(tokens []tok.Token) (expr *ast.TryExpression, remainder []tok.Token, err error) {
	// fmt.Println("TryExpression", tokens)
	remainder = skipTrivia(tokens)

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

	var expression ast.Expression
	if expression, remainder, err = Expression(remainder); err != nil {
		return nil, remainder, err
	}

	return ast.NewTryExpression(variant, expression), remainder, nil
}

// binary_expression = chained_expression .

func BinaryExpression(tokens []tok.Token) (expr ast.Expression, remainder []tok.Token, err error) {
	// fmt.Println("BinaryExpression", tokens)
	return ChainedExpression(tokens)
}

// chained_expression = logical_or_expression { "|>" function_call } .

func ChainedExpression(tokens []tok.Token) (expr ast.Expression, remainder []tok.Token, err error) {
	// fmt.Println("ChainedExpression", tokens)
	remainder = skipTrivia(tokens)

	var initial ast.Expression
	if initial, remainder, err = LogicalOrExpression(remainder); err == ErrNoMatch {
		return nil, tokens, err
	} else if err != nil {
		return nil, remainder, err
	}

	var functionCalls []*ast.FunctionCall

	for {
		var found bool
		if remainder, found = PipeOp(remainder); !found {
			break
		}

		var functionCall *ast.FunctionCall
		if functionCall, remainder, err = FunctionCall(remainder); err != nil {
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
	// fmt.Println("LogicalOrExpression", tokens)
	var initial ast.Expression
	if initial, remainder, err = LogicalAndExpression(tokens); err == ErrNoMatch {
		return nil, tokens, err
	} else if err != nil {
		return nil, remainder, err
	}

	operands := []ast.Expression{initial}
	for {
		var found bool
		if remainder, found = LogicalOrOp(remainder); !found {
			break
		}

		var operand ast.Expression
		if operand, remainder, err = LogicalAndExpression(remainder); err != nil {
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
	// fmt.Println("LogicalAndExpression", tokens)
	var initial ast.Expression
	if initial, remainder, err = ComparisonExpression(tokens); err == ErrNoMatch {
		return nil, tokens, err
	} else if err != nil {
		return nil, remainder, err
	}

	operands := []ast.Expression{initial}
	for {
		var found bool
		if remainder, found = LogicalAndOp(remainder); !found {
			break
		}

		var operand ast.Expression
		if operand, remainder, err = ComparisonExpression(remainder); err != nil {
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
	// fmt.Println("ComparisonExpression", tokens)
	var left ast.Expression
	if left, remainder, err = AddSubExpression(tokens); err == ErrNoMatch {
		return nil, tokens, err
	} else if err != nil {
		return nil, remainder, err
	}

	var typeComparison ast.Expression
	if typeComparison, remainder, err = TypeComparisonTail(left, remainder); err == nil {
		return typeComparison, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	var relationalComparison ast.Expression
	if relationalComparison, remainder, err = RelationalComparisonTail(left, remainder); err == nil {
		return relationalComparison, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	return left, remainder, nil
}

// add_sub_expression = mul_div_expression { add_sub_op mul_div_expression } .

func AddSubExpression(tokens []tok.Token) (expr ast.Expression, remainder []tok.Token, err error) {
	// fmt.Println("AddSubExpression", tokens)
	var left ast.Expression
	if left, remainder, err = MulDivExpression(tokens); err == ErrNoMatch {
		return nil, tokens, err
	} else if err != nil {
		return nil, remainder, err
	}

	// Parse { add_sub_op mul_div_expression }
	for {
		op, remainder2, match := AddSubOp(remainder)
		if !match {
			return left, remainder, nil // no tail => done
		}

		var right ast.Expression
		if right, remainder, err = MulDivExpression(remainder2); err == ErrNoMatch {
			return nil, remainder, errorExpecting("mul_div expression", remainder2)
		} else if err != nil {
			return nil, remainder, err
		}

		left = ast.NewAddSubExpression(left, op, right)
	}
}

// mul_div_expression = pow_expression { mul_div_op pow_expression } .

func MulDivExpression(tokens []tok.Token) (expr ast.Expression, remainder []tok.Token, err error) {
	// fmt.Println("MulDivExpression", tokens)
	var left ast.Expression
	if left, remainder, err = PowExpression(tokens); err == ErrNoMatch {
		return nil, tokens, err
	} else if err != nil {
		return nil, remainder, err
	}

	// Parse { mul_div_op pow_expression }
	for {
		op, remainder2, match := MulDivOp(remainder)
		if !match {
			return left, remainder, nil // no tail => done
		}

		var right ast.Expression
		if right, remainder, err = PowExpression(remainder2); err == ErrNoMatch {
			return nil, remainder, errorExpecting("pow expression", remainder2)
		} else if err != nil {
			return nil, remainder, err
		}

		left = ast.NewMulDivExpression(left, op, right)
	}
}

// pow_expression = unary_expression { "^" unary_expression } .

func PowExpression(tokens []tok.Token) (expr ast.Expression, remainder []tok.Token, err error) {
	// fmt.Println("PowExpression", tokens)
	var left ast.Expression
	if left, remainder, err = UnaryExpression(tokens); err == ErrNoMatch {
		return nil, tokens, err
	} else if err != nil {
		return nil, remainder, err
	}

	operands := []ast.Expression{left}

	for {
		var found bool
		if remainder, found = PowOp(remainder); !found {
			break
		}

		var operand ast.Expression
		if operand, remainder, err = UnaryExpression(remainder); err != nil {
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
	remainder = skipTrivia(tokens)

	if unaryExpression, remainder, err := PrefixedUnaryExpression(remainder); err == nil {
		return unaryExpression, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if primaryExpression, remainder, err := PrimaryExpression(remainder); err == nil {
		return primaryExpression, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	return nil, tokens, ErrNoMatch
}

// prefixed_unary_expression = unary_op negatable_expression .

func PrefixedUnaryExpression(tokens []tok.Token) (expr *ast.UnaryExpression, remainder []tok.Token, err error) {
	// fmt.Println("PrefixedUnaryExpression", tokens)
	operator, remainder, match := UnaryOp(tokens)
	if !match {
		return nil, tokens, ErrNoMatch
	}

	var expression ast.Expression
	if expression, remainder, err = NegatableExpression(remainder); err != nil {
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
	// fmt.Println("NegatableExpression", tokens)
	remainder = skipTrivia(tokens)

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
		return nil, remainder, err
	}

	if literal, remainder, err := Literal(remainder); err == nil {
		return literal, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
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
//                    | meta_expression
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
	// fmt.Println("PrimaryExpression", tokens)

	if expression, remainder, err := parenthesizedExpression(tokens); err == nil {
		return expression, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	// block
	// if_expression
	// for_expression
	// inline_for_expression

	if arrayFunctionCall, remainder, err := ArrayFunctionCall(tokens); err == nil {
		return arrayFunctionCall, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	// import_expression
	// typeof_expression

	if functionCall, remainder, err := FunctionCall(tokens); err == nil {
		return functionCall, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	// type_constructor_call
	// return_expression
	// break_expression
	// continue_expression
	// member_access
	// tuple_update_expression
	// safe_indexed_access
	// indexed_access
	// range

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

// parenthesized_expression = "(" expression ")" .

func parenthesizedExpression(tokens []tok.Token) (expr ast.Expression, remainder []tok.Token, err error) {
	// fmt.Println("ParenthesizedExpression", tokens)
	remainder = skipTrivia(tokens)

	var found bool
	if remainder, found = OpenParen(remainder); !found {
		return nil, remainder, ErrNoMatch
	}

	expression, remainder, err := Expression(remainder)
	if err != nil {
		return nil, remainder, err
	}

	if remainder, found = CloseParen(remainder); !found {
		return nil, remainder, ErrNoMatch
	}

	return expression, remainder, nil
}

// type_comparison_tail = is_op type_predicate .

func TypeComparisonTail(left ast.Expression, tokens []tok.Token) (expr ast.Expression, remainder []tok.Token, err error) {
	// fmt.Println("TypeComparisonTail", tokens)
	var found bool
	if remainder, found = IsOp(tokens); !found {
		return nil, tokens, ErrNoMatch
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
	// fmt.Println("TypePredicate", tokens)
	return nil, nil, ErrNoMatch // TODO: Implement
}

// relational_comparison_tail = rel_op add_sub_expression .

func RelationalComparisonTail(left ast.Expression, tokens []tok.Token) (expr ast.Expression, remainder []tok.Token, err error) {
	// fmt.Println("RelationalComparisonTail", tokens)
	operator, remainder, match := RelOp(tokens)
	if !match {
		return nil, tokens, ErrNoMatch
	}

	right, remainder, err := AddSubExpression(remainder)
	if err != nil {
		return nil, remainder, err
	}

	return ast.NewRelationalComparison(left, operator, right), remainder, nil
}
