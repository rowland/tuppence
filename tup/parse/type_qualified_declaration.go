package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// type_qualified_declaration = type_identifier "." identifier "=" expression .

func TypeQualifiedDeclaration(tokens []tok.Token) (*ast.TypeQualifiedDeclaration, []tok.Token, error) {
	return nil, nil, nil
}
