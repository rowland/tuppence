package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// export_assignment = assignment_lhs ":" expression .

func ExportAssignment(tokens []tok.Token) (*ast.ExportAssignment, []tok.Token, error) {
	lhs, remainder, err := AssignmentLHS(tokens)
	if err != nil {
		return nil, remainder, err
	}

	var found bool
	if remainder, found = Colon(remainder); !found {
		return nil, tokens, ErrNoMatch
	}

	expression, remainder, err := Expression(remainder)
	if err != nil {
		return nil, remainder, err
	}

	return ast.NewExportAssignment(*ast.NewAssignment(lhs, ast.Immutable, expression)), remainder, nil
}

// export_function_declaration = annotations function_declaration_lhs ":" function_declaration_type block .

func ExportFunctionDeclaration(tokens []tok.Token) (*ast.ExportFunctionDeclaration, []tok.Token, error) {
	annotations, remainder, err := Annotations(tokens)
	if err != nil {
		return nil, remainder, err
	}

	lhs, remainder, err := FunctionDeclarationLHS(remainder)
	if err != nil {
		return nil, remainder, err
	}

	var found bool
	if remainder, found = Colon(remainder); !found {
		return nil, tokens, ErrNoMatch
	}

	functionType, remainder, err := FunctionDeclarationType(remainder)
	if err != nil {
		return nil, remainder, err
	}

	body, remainder, err := Block(remainder)
	if err != nil {
		return nil, remainder, err
	}

	return ast.NewExportFunctionDeclaration(
		ast.NewFunctionDeclaration(annotations.Annotations, lhs, functionType, body),
	), remainder, nil
}

// export_function_type_declaration = function_type_declaration_lhs ":" function_type .

func ExportFunctionTypeDeclaration(tokens []tok.Token) (*ast.ExportFunctionTypeDeclaration, []tok.Token, error) {
	name, parameterTypes, remainder, err := FunctionTypeDeclarationLHS(tokens)
	if err != nil {
		return nil, remainder, err
	}

	var found bool
	if remainder, found = Colon(remainder); !found {
		return nil, tokens, ErrNoMatch
	}

	functionType, remainder, err := FunctionType(remainder)
	if err == ErrNoMatch {
		return nil, remainder, errorExpecting("function type", remainder)
	} else if err != nil {
		return nil, remainder, err
	}

	return ast.NewExportFunctionTypeDeclaration(
		ast.NewFunctionTypeDeclaration(name, parameterTypes, functionType),
	), remainder, nil
}

// export_type_declaration = type_declaration_lhs ":" type_declaration_rhs .

func ExportTypeDeclaration(tokens []tok.Token) (*ast.ExportTypeDeclaration, []tok.Token, error) {
	lhs, remainder, err := TypeDeclarationLHS(tokens)
	if err != nil {
		return nil, remainder, err
	}

	var found bool
	if remainder, found = Colon(remainder); !found {
		return nil, tokens, ErrNoMatch
	}

	rhs, remainder, err := TypeDeclarationRHS(remainder)
	if err == ErrNoMatch {
		return nil, remainder, errorExpecting("type declaration right-hand side", remainder)
	} else if err != nil {
		return nil, remainder, err
	}

	return ast.NewExportTypeDeclaration(*ast.NewTypeDeclaration(lhs, rhs)), remainder, nil
}

// export_declaration = ( export_type_qualified_function_declaration
// 	                    | export_type_qualified_declaration
// 	                    | export_function_type_declaration
// 	                    | export_type_declaration
// 	                    | export_function_declaration
// 	                    | export_assignment ) .
//
// This parser currently implements the first four non-qualified alternatives plus
// export_assignment. The type-qualified export forms remain deferred.

func ExportDeclaration(tokens []tok.Token) (ast.ExportDeclaration, []tok.Token, error) {
	var errors []error

	if declaration, remainder, err := ExportFunctionTypeDeclaration(tokens); err == nil {
		return declaration, remainder, nil
	} else {
		errors = append(errors, err)
	}

	if declaration, remainder, err := ExportTypeDeclaration(tokens); err == nil {
		return declaration, remainder, nil
	} else {
		errors = append(errors, err)
	}

	if declaration, remainder, err := ExportFunctionDeclaration(tokens); err == nil {
		return declaration, remainder, nil
	} else {
		errors = append(errors, err)
	}

	if declaration, remainder, err := ExportAssignment(tokens); err == nil {
		return declaration, remainder, nil
	} else {
		errors = append(errors, err)
	}

	return nil, nil, errorExpectingOneOf("export declaration", tokens, errors)
}
