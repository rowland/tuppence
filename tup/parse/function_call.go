package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// function_call = function_identifier [ function_parameter_types ] "(" [ function_arguments ] ")" [ function_block ] .

func FunctionCall(tokens []tok.Token) (expr *ast.FunctionCall, remainder []tok.Token, err error) {
	// fmt.Println("FunctionCall", tok.Types(tokens))
	var functionIdentifier *ast.FunctionIdentifier
	functionIdentifier, remainder, err = FunctionIdentifier(tokens)
	if err != nil {
		return nil, remainder, err
	}

	var functionParameterTypes *ast.FunctionParameterTypes
	functionParameterTypes, remainder, err = FunctionParameterTypes(remainder)
	if err != nil && err != ErrNoMatch {
		return nil, remainder, err
	}

	var found bool
	if remainder, found = OpenParen(remainder); !found {
		return nil, remainder, ErrNoMatch
	}

	var arguments *ast.FunctionArguments
	arguments, remainder, err = FunctionArguments(remainder)
	if err != nil && err != ErrNoMatch {
		return nil, remainder, err
	}

	if remainder, found = CloseParen(remainder); !found {
		return nil, remainder, errorExpecting(tok.TokenTypes[tok.TokCloseParen], remainder)
	}

	var functionBlock *ast.FunctionBlock
	functionBlock, remainder, err = FunctionBlock(remainder)
	if err != nil && err != ErrNoMatch {
		return nil, remainder, err
	}

	return ast.NewFunctionCall(functionIdentifier, functionParameterTypes, arguments, functionBlock), remainder, nil
}

// function_parameter_types = "[" local_type_reference { "," local_type_reference } "]" .

func FunctionParameterTypes(tokens []tok.Token) (expr *ast.FunctionParameterTypes, remainder []tok.Token, err error) {
	// fmt.Println("FunctionParameterTypes", tok.Types(tokens))
	var found bool
	if remainder, found = OpenBracket(tokens); !found {
		return nil, tokens, ErrNoMatch
	}

	var parameters []ast.LocalTypeReference
	if parameters, remainder, err = LocalTypeReferenceList(remainder); err != nil {
		return nil, remainder, err
	}

	if remainder, found = CloseBracket(remainder); !found {
		return nil, remainder, errorExpecting(tok.TokenTypes[tok.TokCloseBracket], remainder)
	}

	return ast.NewFunctionParameterTypes(parameters), remainder, nil
}

func LocalTypeReferenceList(tokens []tok.Token) (parameters []ast.LocalTypeReference, remainder []tok.Token, err error) {
	// fmt.Println("LocalTypeReferenceList", tok.Types(tokens))
	remainder = tokens
	for {
		var parameter ast.LocalTypeReference
		if parameter, remainder, err = LocalTypeReference(remainder); err != nil {
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
	if parameters, remainder, err = BlockParameters(tokens); err != nil && err != ErrNoMatch {
		return nil, remainder, err
	}

	var body *ast.BlockBody
	if body, remainder, err = BlockBody(remainder); err != nil && err != ErrNoMatch {
		return nil, remainder, err
	}

	if remainder, found = CloseBrace(remainder); !found {
		return nil, remainder, errorExpecting(tok.TokenTypes[tok.TokCloseBrace], remainder)
	}

	return ast.NewFunctionBlock(parameters, body), remainder, nil
}

// block_parameters = "|" [ assignment_lhs { "," assignment_lhs } ] "|" .

func BlockParameters(tokens []tok.Token) (expr *ast.BlockParameters, remainder []tok.Token, err error) {
	// fmt.Println("BlockParameters", tok.Types(tokens))
	return nil, nil, ErrNoMatch // TODO: Implement
}

// block_body = { statement } expression .

func BlockBody(tokens []tok.Token) (expr *ast.BlockBody, remainder []tok.Token, err error) {
	// fmt.Println("BlockBody", tok.Types(tokens))
	var statements []ast.Statement
	if statements, remainder, err = Statements(tokens); err != nil && err != ErrNoMatch {
		return nil, remainder, err
	}

	var expression ast.Expression
	if expression, remainder, err = Expression(remainder); err != nil && err != ErrNoMatch {
		return nil, remainder, err
	}

	return ast.NewBlockBody(statements, expression), remainder, nil
}
