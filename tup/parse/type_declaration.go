package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// type_declaration = type_declaration_lhs "=" type_declaration_rhs .

func TypeDeclaration(tokens []tok.Token) (*ast.TypeDeclaration, []tok.Token, error) {
	return nil, nil, nil
}
