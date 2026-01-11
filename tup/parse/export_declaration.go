package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// export_declaration = ( export_type_qualified_function_declaration
// 	                    | export_type_qualified_declaration
// 	                    | export_type_declaration
// 	                    | export_function_declaration
// 	                    | export_assignment ) .

func ExportDeclaration(tokens []tok.Token) (ast.ExportDeclaration, []tok.Token, error) {
	return nil, nil, nil
}
