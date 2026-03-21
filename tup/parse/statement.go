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

func Statement(tokens []tok.Token) (stmt ast.Statement, remainder []tok.Token, err error) {
	// fmt.Println("Statement", tok.Types(tokens))

	// var typeQualifiedFunctionDeclaration *ast.TypeQualifiedFunctionDeclaration
	// if typeQualifiedFunctionDeclaration, remainder, err = TypeQualifiedFunctionDeclaration(tokens); err == nil {
	// 	return typeQualifiedFunctionDeclaration, remainder, nil
	// } else if err != ErrNoMatch {
	// 	return nil, remainder, err
	// }

	// var typeQualifiedDeclaration *ast.TypeQualifiedDeclaration
	// if typeQualifiedDeclaration, remainder, err = TypeQualifiedDeclaration(tokens); err == nil {
	// 	return typeQualifiedDeclaration, remainder, nil
	// } else if err != ErrNoMatch {
	// 	return nil, remainder, err
	// }

	// var typeDeclaration *ast.TypeDeclaration
	// if typeDeclaration, remainder, err = TypeDeclaration(tokens); err == nil {
	// 	return typeDeclaration, remainder, nil
	// } else if err != ErrNoMatch {
	// 	return nil, remainder, err
	// }

	// var functionDeclaration *ast.FunctionDeclaration
	// if functionDeclaration, remainder, err = FunctionDeclaration(tokens); err == nil {
	// 	return functionDeclaration, remainder, nil
	// } else if err != ErrNoMatch {
	// 	return nil, remainder, err
	// }

	var compoundAssignment *ast.CompoundAssignment
	if compoundAssignment, remainder, err = CompoundAssignment(tokens); err == nil {
		return compoundAssignment, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	var assignment *ast.Assignment
	if assignment, remainder, err = Assignment(tokens); err == nil {
		return assignment, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	var expression ast.Expression
	if expression, remainder, err = Expression(tokens); err == nil {
		statement, ok := expression.(ast.Statement)
		if !ok {
			return nil, remainder, errorExpecting("statement", remainder)
		}
		return statement, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	return nil, tokens, ErrNoMatch
}

// statements = { statement } .

func Statements(tokens []tok.Token) (stmts []ast.Statement, remainder []tok.Token, err error) {
	// fmt.Println("Statements", tok.Types(tokens))

	remainder = tokens
	var statements []ast.Statement
	for {
		var statement ast.Statement
		if statement, remainder, err = Statement(remainder); err == ErrNoMatch {
			break
		} else if err != nil {
			return nil, remainder, err
		}
		statements = append(statements, statement)
		remainder = skipStatmentSeparators(remainder)
	}
	return statements, remainder, nil
}

func skipStatmentSeparators(tokens []tok.Token) []tok.Token {
	return skip(tokens, tok.TokEOL, tok.TokComment, tok.TokSemiColon)
}
