package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// function_call = function_identifier [ function_parameter_types ] "(" [ function_arguments ] ")" [ function_block ] .

func FunctionCall(tokens []tok.Token) (expr *ast.FunctionCall, remainder []tok.Token, err error) {
	// fmt.Println("FunctionCall", tokens)
	return nil, nil, ErrNoMatch // TODO: Implement
}

// function_parameter_types = "[" [ local_type_reference { "," local_type_reference } ] "]" .

func FunctionParameterTypes(tokens []tok.Token) (expr *ast.FunctionParameterTypes, remainder []tok.Token, err error) {
	// fmt.Println("FunctionParameterTypes", tokens)
	return nil, nil, ErrNoMatch // TODO: Implement
}

// function_arguments = ( labeled_arguments | arguments [ "," labeled_arguments ] ) [ partial_application ] .

func FunctionArguments(tokens []tok.Token) (expr *ast.FunctionArguments, remainder []tok.Token, err error) {
	// fmt.Println("FunctionArguments", tokens)
	return nil, nil, ErrNoMatch // TODO: Implement
}
