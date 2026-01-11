package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// export_type_qualified_function_declaration = annotations type_identifier "." function_declaration_lhs ":" function_declaration_type block .

func TypeQualifiedFunctionDeclaration(tokens []tok.Token) (*ast.TypeQualifiedFunctionDeclaration, []tok.Token, error) {
	return nil, nil, nil
}
