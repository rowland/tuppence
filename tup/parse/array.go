package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// size = decimal_literal | identifier .

func Size(tokens []tok.Token) (size ast.Size, remainder []tok.Token, err error) {
	remainder = skipComments(tokens)
	var errors []error

	integerLit, remainder2, err := IntegerLiteral(remainder)
	if err == nil {
		return integerLit, remainder2, nil
	}
	errors = append(errors, err)

	identifier, remainder3, err := Identifier(remainder)
	if err == nil {
		return identifier, remainder3, nil
	}
	errors = append(errors, err)

	return nil, nil, errorExpectingOneOf("size", tokens, errors)
}

// array_literal = "[" [ array_members | array_literal ] "]"
//               | type_identifier "[" [ array_members | array_literal ] "]" .

func ArrayLiteral(tokens []tok.Token) (item *ast.ArrayLiteral, remainder []tok.Token, err error) {
	remainder = skipComments(tokens)

	remainder, err = OpenBracket(remainder)
	if err != nil {
		return nil, nil, err
	}

	arrayMembers, remainder, err := ArrayMembers(remainder)
	if err != nil {
		return nil, nil, err
	}

	remainder, err = CloseBracket(remainder)
	if err != nil {
		return nil, nil, err
	}

	item = ast.NewArrayLiteral(arrayMembers, nil)
	return item, remainder, nil
}

// array_members = expression { "," expression } [ "," ] .

func ArrayMembers(tokens []tok.Token) (members []ast.Expression, remainder []tok.Token, err error) {
	remainder = skipComments(tokens)

	for {
		expression, remainder2, err := Expression(remainder)
		if err != nil {
			break
		}
		members = append(members, expression)
		if peek(remainder2).Type == tok.TokComma {
			remainder = remainder2[1:]
		} else {
			remainder = remainder2
			break
		}
	}

	if peek(remainder).Type == tok.TokComma {
		remainder = remainder[1:]
	}

	return members, remainder, nil
}
