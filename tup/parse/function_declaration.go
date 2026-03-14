package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// function_declaration = annotations function_declaration_lhs "=" function_declaration_type block .

func FunctionDeclaration(tokens []tok.Token) (*ast.FunctionDeclaration, []tok.Token, error) {
	return nil, tokens, ErrNoMatch // TODO: Implement
}

// function_declaration_lhs = function_identifier [ function_parameter_types ] .

// function_declaration_type = ( "fn" "(" [ labeled_parameters | parameters ] ")" ( return_type | "_" ) )
//                           | ( "fx" "(" [ labeled_parameters | parameters ] ")" [ return_type | "_" ] ) .

// block = "{" block_body "}" .

// block_body = { statement } expression .
