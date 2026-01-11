package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// function_declaration = annotations function_declaration_lhs "=" function_declaration_type block .

func FunctionDeclaration(tokens []tok.Token) (*ast.FunctionDeclaration, []tok.Token, error) {
	return nil, nil, nil
}
