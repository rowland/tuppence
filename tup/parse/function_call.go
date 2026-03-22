package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// function_call = callable_expression function_call_tail { function_call_tail } .

func FunctionCall(tokens []tok.Token) (expr *ast.FunctionCall, remainder []tok.Token, err error) {
	var function ast.Expression
	if function, remainder, err = callableExpression(tokens); err != nil {
		return nil, remainder, err
	}

	if expr, remainder, err = functionCallTail(function, remainder); err != nil {
		return nil, remainder, err
	}

	for {
		var next *ast.FunctionCall
		if next, remainder, err = functionCallTail(expr, remainder); err == ErrNoMatch {
			break
		} else if err != nil {
			return nil, remainder, err
		}
		expr = next
	}

	return expr, remainder, nil
}

// function_call_tail = [ function_parameter_types ] "(" [ function_arguments ] ")" [ function_block ] .

func functionCallTail(function ast.Expression, tokens []tok.Token) (expr *ast.FunctionCall, remainder []tok.Token, err error) {
	var functionParameterTypes *ast.FunctionParameterTypes
	remainder = tokens
	if functionParameterTypes, remainder, err = FunctionParameterTypes(remainder); err != nil && err != ErrNoMatch {
		return nil, remainder, err
	}

	var found bool
	if remainder, found = OpenParen(remainder); !found {
		return nil, tokens, ErrNoMatch
	}

	var arguments *ast.FunctionArguments
	if arguments, remainder, err = FunctionArguments(remainder); err != nil && err != ErrNoMatch {
		return nil, remainder, err
	}

	if remainder, found = CloseParen(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokCloseParen, remainder)
	}

	var functionBlock *ast.FunctionBlock
	if functionBlock, remainder, err = FunctionBlock(remainder); err != nil && err != ErrNoMatch {
		return nil, remainder, err
	}

	return ast.NewFunctionCall(function, functionParameterTypes, arguments, functionBlock), remainder, nil
}

// postfix_base_expression but restricted to forms that may be followed by function_call_tail.

func callableExpression(tokens []tok.Token) (expr ast.Expression, remainder []tok.Token, err error) {
	return postfixExpressionWithTails(tokens, callableBaseExpression, true, callableReceiverTail)
}

// callable_base_expression is a parser helper for the subset of postfix_base_expression
// that is currently supported as a function-call receiver.

func callableBaseExpression(tokens []tok.Token) (expr ast.Expression, remainder []tok.Token, err error) {
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

	if typeConstructorCall, remainder, err := TypeConstructorCall(tokens); err == nil {
		return typeConstructorCall, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	// return_expression
	// break_expression
	// continue_expression
	// range

	if literal, remainder, err := Literal(tokens); err == nil {
		return literal, remainder, nil
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

	return nil, tokens, ErrNoMatch
}

// function_parameter_types = "[" function_parameter_type { "," function_parameter_type } "]" .

func FunctionParameterTypes(tokens []tok.Token) (expr *ast.FunctionParameterTypes, remainder []tok.Token, err error) {
	var found bool
	if remainder, found = OpenBracket(tokens); !found {
		return nil, tokens, ErrNoMatch
	}

	var parameters []ast.FunctionParameterType
	if parameters, remainder, err = FunctionParameterTypeList(remainder); err != nil {
		return nil, remainder, err
	}

	if remainder, found = CloseBracket(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokCloseBracket, remainder)
	}

	return ast.NewFunctionParameterTypes(parameters), remainder, nil
}

// function_parameter_type = local_type_reference
//                         | nilable_type
//                         | fallible_type
//                         | dynamic_array
//                         | fixed_size_array .

func FunctionParameterType(tokens []tok.Token) (ast.FunctionParameterType, []tok.Token, error) {
	if nilableType, remainder, err := NilableType(tokens); err == nil {
		return nilableType, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if fallibleType, remainder, err := FallibleType(tokens); err == nil {
		return fallibleType, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if localTypeReference, remainder, err := LocalTypeReference(tokens); err == nil {
		switch t := localTypeReference.(type) {
		case *ast.TypeReference:
			return t, remainder, nil
		case *ast.Identifier:
			return t, remainder, nil
		default:
			return nil, remainder, errorExpecting("function parameter type", remainder)
		}
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if dynamicArray, remainder, err := DynamicArray(tokens); err == nil {
		return dynamicArray, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if fixedSizeArray, remainder, err := FixedSizeArray(tokens); err == nil {
		return fixedSizeArray, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	return nil, tokens, ErrNoMatch
}

func FunctionParameterTypeList(tokens []tok.Token) (parameters []ast.FunctionParameterType, remainder []tok.Token, err error) {
	remainder = tokens
	for {
		var parameter ast.FunctionParameterType
		if parameter, remainder, err = FunctionParameterType(remainder); err != nil {
			return nil, remainder, err
		}
		parameters = append(parameters, parameter)
		var found bool
		if remainder, found = Comma(remainder); !found {
			break
		}
	}

	if len(parameters) == 0 {
		return nil, remainder, ErrNoMatch
	}
	return parameters, remainder, nil
}

// function_arguments = ( arguments_body [ partial_application ]
//                      | "*"
// 	                    )
// 	                    [ "," ] .

func FunctionArguments(tokens []tok.Token) (expr *ast.FunctionArguments, remainder []tok.Token, err error) {
	// fmt.Println("FunctionArguments", tok.Types(tokens))
	arguments, labeledArgs, remainder, err := ArgumentsBody(tokens)
	if err == nil {
		var partialApplication bool
		remainder, partialApplication = PartialApplication(remainder)
		remainder, _ = Comma(remainder)
		return ast.NewFunctionArguments(arguments, labeledArgs, partialApplication), remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	var partialApplication bool
	remainder, partialApplication = Star(remainder)
	remainder, _ = Comma(remainder)
	return ast.NewFunctionArguments(nil, nil, partialApplication), remainder, nil
}

// function_block = "{" [ block_parameters ] block_body "}" .

func FunctionBlock(tokens []tok.Token) (expr *ast.FunctionBlock, remainder []tok.Token, err error) {
	// fmt.Println("FunctionBlock", tok.Types(tokens))
	var found bool
	if remainder, found = OpenBrace(tokens); !found {
		return nil, tokens, ErrNoMatch
	}

	var parameters *ast.BlockParameters
	if parameters, remainder, err = BlockParameters(remainder); err != nil && err != ErrNoMatch {
		return nil, remainder, err
	}
	// fmt.Println("FunctionBlock parameters", parameters, tok.Types(remainder))

	var body *ast.BlockBody
	if body, remainder, err = BlockBody(remainder); err != nil {
		return nil, remainder, err
	}
	// fmt.Println("FunctionBlock body", body, tok.Types(remainder))

	if remainder, found = CloseBrace(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokCloseBrace, remainder)
	}
	// fmt.Println("FunctionBlock close brace", tok.Types(remainder))

	return ast.NewFunctionBlock(parameters, body), remainder, nil
}

// block_parameters = "|" assignment_lhs "|" .

func BlockParameters(tokens []tok.Token) (expr *ast.BlockParameters, remainder []tok.Token, err error) {
	// fmt.Println("BlockParameters", tok.Types(tokens))

	var found bool
	if remainder, found = Pipe(tokens); !found {
		return nil, tokens, ErrNoMatch
	}

	var parameters ast.AssignmentLHS
	if parameters, remainder, err = AssignmentLHS(remainder); err != nil {
		return nil, remainder, err
	}

	if remainder, found = Pipe(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokOpPipe, remainder)
	}

	return ast.NewBlockParameters(parameters), remainder, nil
}
