package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// size = decimal_literal | identifier .

func Size(tokens []tok.Token) (size ast.Size, remainder []tok.Token, err error) {
	remainder = skipTrivia(tokens)

	integerLit, remainder2, err := IntegerLiteral(remainder)
	if err == nil {
		return integerLit, remainder2, nil
	}

	identifier, remainder3, err := Identifier(remainder)
	if err == nil {
		return identifier, remainder3, nil
	}

	return nil, tokens, ErrNoMatch
}

// array_literal = "[" [ array_members ] "]"
//               | type_identifier "[" [ array_members ] "]" .

func ArrayLiteral(tokens []tok.Token) (arr *ast.ArrayLiteral, remainder []tok.Token, err error) {
	remainder = skipTrivia(tokens)

	var found bool
	if remainder, found = OpenBracket(remainder); !found {
		return nil, nil, ErrNoMatch
	}

	var arrayMembers []ast.Expression
	if arrayMembers, remainder, err = ArrayMembers(remainder); err != nil {
		return nil, nil, err
	}

	if remainder, found = CloseBracket(remainder); !found {
		return nil, nil, ErrNoMatch
	}

	arr = ast.NewArrayLiteral(arrayMembers, nil)
	return arr, remainder, nil
}

// array_members = expression { "," expression } [ "," ] .

func ArrayMembers(tokens []tok.Token) (members []ast.Expression, remainder []tok.Token, err error) {
	remainder = skipTrivia(tokens)

	for {
		var expression ast.Expression
		var remainder2 []tok.Token
		if expression, remainder2, err = Expression(remainder); err != nil {
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
