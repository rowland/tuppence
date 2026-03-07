package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// statement = ( type_qualified_function_declaration
// 	           | type_qualified_declaration
// 	           | type_declaration
// 	           | function_declaration
// 	           | compound_assignment
// 	           | assignment
// 	           | expression
// 	           ) .

func Statement(tokens []tok.Token) (expr *ast.Statement, remainder []tok.Token, err error) {
	return nil, nil, ErrNoMatch // TODO: Implement
}

func Statements(tokens []tok.Token) (expr []ast.Statement, remainder []tok.Token, err error) {
	return nil, nil, ErrNoMatch // TODO: Implement
}
