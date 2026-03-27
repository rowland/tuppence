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
	var tqfd *ast.TypeQualifiedFunctionDeclaration
	if tqfd, remainder, err = TypeQualifiedFunctionDeclaration(tokens); err == nil {
		if tqfd != nil {
			return tqfd, remainder, nil
		}
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	var tqd *ast.TypeQualifiedDeclaration
	if tqd, remainder, err = TypeQualifiedDeclaration(tokens); err == nil {
		if tqd != nil {
			return tqd, remainder, nil
		}
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	var td *ast.TypeDeclaration
	if td, remainder, err = TypeDeclaration(tokens); err == nil {
		if td != nil {
			return td, remainder, nil
		}
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	var ftd *ast.FunctionTypeDeclaration
	if ftd, remainder, err = FunctionTypeDeclaration(tokens); err == nil {
		if ftd != nil {
			return ftd, remainder, nil
		}
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	var fd *ast.FunctionDeclaration
	if fd, remainder, err = FunctionDeclaration(tokens); err == nil {
		if fd != nil {
			return fd, remainder, nil
		}
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	var a *ast.Assignment
	if a, remainder, err = Assignment(tokens); err == nil {
		if a != nil {
			return a, remainder, nil
		}
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	var ed ast.TopLevelItem
	if ed, remainder, err = ExportDeclaration(tokens); err == nil {
		if ed != nil {
			return ed, remainder, nil
		}
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	return nil, nil, errorExpecting("top-level item", tokens)
}
