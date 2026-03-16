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

// negatable_expression = negatable_postfix_expression .

func NegatableExpression(tokens []tok.Token) (expr ast.Expression, remainder []tok.Token, err error) {
	return postfixExpression(tokens, negatablePostfixBaseExpression, true)
}

// primary_expression = postfix_expression .

func PrimaryExpression(tokens []tok.Token) (expr ast.Expression, remainder []tok.Token, err error) {
	return postfixExpression(tokens, postfixBaseExpression, true)
}

// postfix_expression = postfix_base_expression { postfix_tail }
//                    | type_identifier member_access_tail { postfix_tail } .
//
// The current expression is folded left-to-right through each matched tail, so
// `a.b[0](x)` becomes a chain of AST nodes rooted in the most recently matched tail.

func postfixExpression(tokens []tok.Token, base func([]tok.Token) (ast.Expression, []tok.Token, error), allowTypeMemberAccess bool) (expr ast.Expression, remainder []tok.Token, err error) {
	return postfixExpressionWithTails(tokens, base, allowTypeMemberAccess, postfixTail)
}

func postfixExpressionWithTails(
	tokens []tok.Token,
	base func([]tok.Token) (ast.Expression, []tok.Token, error),
	allowTypeMemberAccess bool,
	nextTail func(ast.Expression, []tok.Token) (ast.Expression, []tok.Token, error),
) (expr ast.Expression, remainder []tok.Token, err error) {
	if allowTypeMemberAccess {
		if expr, remainder, err = typeMemberAccessExpression(tokens); err == nil {
			return continuePostfixExpression(expr, remainder, nextTail)
		} else if err != ErrNoMatch {
			return nil, remainder, err
		}
	}

	if expr, remainder, err = base(tokens); err != nil {
		return nil, remainder, err
	}

	return continuePostfixExpression(expr, remainder, nextTail)
}

// postfix_expression = postfix_base_expression { postfix_tail } .
// continuePostfixExpression parses the repeated postfix_tail portion by
// wrapping the current expression in each matched tail node, left-to-right.

func continuePostfixExpression(
	expr ast.Expression,
	tokens []tok.Token,
	nextTail func(ast.Expression, []tok.Token) (ast.Expression, []tok.Token, error),
) (ast.Expression, []tok.Token, error) {
	remainder := tokens

	for {
		if nextExpr, remainder2, err := nextTail(expr, remainder); err == nil {
			expr = nextExpr
			remainder = remainder2
			continue
		} else if err != ErrNoMatch {
			return nil, remainder2, err
		}

		break
	}

	return expr, remainder, nil
}

// postfix_tail = function_call_tail
//              | member_access_tail
//              | tuple_update_tail
//              | safe_indexed_access_tail
//              | indexed_access_tail .

func postfixTail(expr ast.Expression, tokens []tok.Token) (ast.Expression, []tok.Token, error) {
	return matchPostfixTail(expr, tokens, true, true)
}

// function_call_tail
// | member_access_tail
// | safe_indexed_access_tail
// | indexed_access_tail .
//
// This subset is used while parsing the receiver of tuple_update_expression.

func postfixTailWithoutTupleUpdate(expr ast.Expression, tokens []tok.Token) (ast.Expression, []tok.Token, error) {
	return matchPostfixTail(expr, tokens, true, false)
}

// member_access_tail
// | safe_indexed_access_tail
// | indexed_access_tail .
//
// This subset is used while discovering a function-call receiver, so the first
// function_call_tail is left for FunctionCall(...) itself to consume.

func callableReceiverTail(expr ast.Expression, tokens []tok.Token) (ast.Expression, []tok.Token, error) {
	return matchPostfixTail(expr, tokens, false, false)
}

func matchPostfixTail(expr ast.Expression, tokens []tok.Token, includeFunctionCall bool, includeTupleUpdate bool) (ast.Expression, []tok.Token, error) {
	if includeFunctionCall {
		if functionCall, remainder2, err := functionCallTail(expr, tokens); err == nil {
			return functionCall, remainder2, nil
		} else if err != ErrNoMatch {
			return nil, remainder2, err
		}
	}

	if memberAccess, remainder2, err := memberAccessTail(expr, tokens); err == nil {
		return memberAccess, remainder2, nil
	} else if err != ErrNoMatch {
		return nil, remainder2, err
	}

	if safeIndexedAccess, remainder2, err := safeIndexedAccessTail(expr, tokens); err == nil {
		return safeIndexedAccess, remainder2, nil
	} else if err != ErrNoMatch {
		return nil, remainder2, err
	}

	if indexedAccess, remainder2, err := indexedAccessTail(expr, tokens); err == nil {
		return indexedAccess, remainder2, nil
	} else if err != ErrNoMatch {
		return nil, remainder2, err
	}

	if includeTupleUpdate {
		if tupleUpdateExpression, remainder2, err := tupleUpdateTail(expr, tokens); err == nil {
			return tupleUpdateExpression, remainder2, nil
		} else if err != ErrNoMatch {
			return nil, remainder2, err
		}
	}

	return nil, tokens, ErrNoMatch
}

// postfix_base_expression = "(" expression ")"
//                         | block
//                         | if_expression
//                         | for_expression
//                         | inline_for_expression
//                         | array_function_call
//                         | import_expression
//                         | typeof_expression
//                         | meta_expression
//                         | type_constructor_call
//                         | return_expression
//                         | break_expression
//                         | continue_expression
//                         | range
//                         | function_identifier
//                         | identifier
//                         | literal .
//
// The parser currently implements the forms listed below and leaves the others
// as explicit follow-up work in the stacked postfix refactor.

func postfixBaseExpression(tokens []tok.Token) (expr ast.Expression, remainder []tok.Token, err error) {
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

	// if_expression
	// for_expression
	// inline_for_expression

	if arrayFunctionCall, remainder, err := ArrayFunctionCall(tokens); err == nil {
		return arrayFunctionCall, remainder, nil
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

	// type_constructor_call
	// return_expression
	// break_expression
	// continue_expression
	// range

	if functionIdentifier, remainder, err := FunctionIdentifier(tokens); err == nil {
		return functionIdentifier, remainder, nil
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

// negatable_postfix_base_expression = "(" expression ")"
//                                   | block
//                                   | function_identifier
//                                   | identifier
//                                   | literal .

func negatablePostfixBaseExpression(tokens []tok.Token) (expr ast.Expression, remainder []tok.Token, err error) {
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

	if functionIdentifier, remainder, err := FunctionIdentifier(tokens); err == nil {
		return functionIdentifier, remainder, nil
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

// type_identifier member_access_tail .

func typeMemberAccessExpression(tokens []tok.Token) (expr ast.Expression, remainder []tok.Token, err error) {
	var typeIdentifier *ast.TypeIdentifier
	if typeIdentifier, remainder, err = TypeIdentifier(tokens); err != nil {
		return nil, remainder, err
	}

	return memberAccessTail(typeIdentifier, remainder)
}

// member_access_tail = "." ( decimal_literal | identifier ) .

func memberAccessTail(object ast.Node, tokens []tok.Token) (expr *ast.MemberAccess, remainder []tok.Token, err error) {
	remainder = skipTrivia(tokens)

	var found bool
	if remainder, found = Dot(remainder); !found {
		return nil, tokens, ErrNoMatch
	}

	var member ast.MemberAccessMember
	if member, remainder, err = memberAccessMember(remainder); err != nil {
		if err == ErrNoMatch {
			if peek(remainder).Type == tok.TokOpenParen {
				return nil, tokens, ErrNoMatch
			}
			return nil, remainder, errorExpecting("member access member", remainder)
		}
		return nil, remainder, err
	}

	return ast.NewMemberAccess(object, member), remainder, nil
}

// member_access_member = decimal_literal | identifier .

func memberAccessMember(tokens []tok.Token) (member ast.MemberAccessMember, remainder []tok.Token, err error) {
	if decimalLiteral, remainder, err := DecimalLiteral(tokens); err == nil {
		return decimalLiteral, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if identifier, remainder, err := Identifier(tokens); err == nil {
		return identifier, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	return nil, tokens, ErrNoMatch
}

// indexed_access_tail = "[" index "]" .

func indexedAccessTail(object ast.Expression, tokens []tok.Token) (expr *ast.IndexedAccess, remainder []tok.Token, err error) {
	remainder = skipTrivia(tokens)

	var found bool
	if remainder, found = OpenBracket(remainder); !found {
		return nil, tokens, ErrNoMatch
	}

	var index ast.Expression
	if index, remainder, err = Expression(remainder); err != nil {
		return nil, remainder, err
	}

	if remainder, found = CloseBracket(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokCloseBracket, remainder)
	}

	return ast.NewIndexedAccess(object, index), remainder, nil
}

// safe_indexed_access_tail = "[" index "]" "!" .

func safeIndexedAccessTail(object ast.Expression, tokens []tok.Token) (expr *ast.SafeIndexedAccess, remainder []tok.Token, err error) {
	var indexedAccess *ast.IndexedAccess
	if indexedAccess, remainder, err = indexedAccessTail(object, tokens); err != nil {
		return nil, remainder, err
	}

	var found bool
	if remainder, found = expectFunc(tok.TokOpNot)(remainder); !found {
		return nil, tokens, ErrNoMatch
	}

	return ast.NewSafeIndexedAccess(object, indexedAccess.Index), remainder, nil
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
