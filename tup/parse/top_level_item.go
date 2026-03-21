package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// top_level_item = ( type_qualified_function_declaration
// 	                | type_qualified_declaration
// 	                | type_declaration
// 	                | function_type_declaration
// 	                | function_declaration
// 	                | assignment
// 	                | export_declaration
// 	                ) .

func TopLevelItem(tokens []tok.Token) (item ast.TopLevelItem, remainder []tok.Token, err error) {
	var errors []error
	var tqfd *ast.TypeQualifiedFunctionDeclaration
	if tqfd, remainder, err = TypeQualifiedFunctionDeclaration(tokens); err != nil {
		errors = append(errors, err)
	} else if tqfd != nil {
		return tqfd, remainder, nil
	}

	var tqd *ast.TypeQualifiedDeclaration
	if tqd, remainder, err = TypeQualifiedDeclaration(tokens); err != nil {
		errors = append(errors, err)
	} else if tqd != nil {
		return tqd, remainder, nil
	}

	var td *ast.TypeDeclaration
	if td, remainder, err = TypeDeclaration(tokens); err != nil {
		errors = append(errors, err)
	} else if td != nil {
		return td, remainder, nil
	}

	var ftd *ast.FunctionTypeDeclaration
	if ftd, remainder, err = FunctionTypeDeclaration(tokens); err != nil {
		errors = append(errors, err)
	} else if ftd != nil {
		return ftd, remainder, nil
	}

	var fd *ast.FunctionDeclaration
	if fd, remainder, err = FunctionDeclaration(tokens); err != nil {
		errors = append(errors, err)
	} else if fd != nil {
		return fd, remainder, nil
	}

	var a *ast.Assignment
	if a, remainder, err = Assignment(tokens); err != nil {
		errors = append(errors, err)
	} else if a != nil {
		return a, remainder, nil
	}

	var ed ast.TopLevelItem
	if ed, remainder, err = ExportDeclaration(tokens); err != nil {
		errors = append(errors, err)
	} else if ed != nil {
		return ed, remainder, nil
	}

	return nil, nil, errorExpectingOneOf("top-level item", tokens, errors)
}
