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

// array_literal = fixed_size_array array_initializer
//               | type_reference array_initializer
//               | "[" [ array_members ] "]" .

func ArrayLiteral(tokens []tok.Token) (arr *ast.ArrayLiteral, remainder []tok.Token, err error) {
	remainder = skipTrivia(tokens)

	if fixedSizeArray, remainder2, fixedErr := fixedSizeArrayLiteralType(remainder); fixedErr == nil {
		if arr, remainder, err = arrayLiteralWithType(fixedSizeArray, remainder2); err == nil {
			return arr, remainder, nil
		} else if err != ErrNoMatch {
			return nil, remainder, err
		}
	} else if fixedErr != ErrNoMatch {
		return nil, remainder2, fixedErr
	}

	if typeReference, remainder2, typeErr := TypeReference(remainder); typeErr == nil {
		if arr, remainder, err = arrayLiteralWithType(typeReference, remainder2); err == nil {
			return arr, remainder, nil
		} else if err != ErrNoMatch {
			return nil, remainder, err
		}
	} else if typeErr != ErrNoMatch {
		return nil, remainder2, typeErr
	}

	var found bool
	if remainder, found = OpenBracket(remainder); !found {
		return nil, tokens, ErrNoMatch
	}

	var arrayMembers []ast.Expression
	if arrayMembers, remainder, err = ArrayMembers(remainder); err != nil {
		return nil, remainder, err
	}

	remainder = skipTrivia(remainder)
	if remainder, found = CloseBracket(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokCloseBracket, remainder)
	}

	return ast.NewArrayLiteral(nil, arrayMembers, nil), remainder, nil
}

// array_members = expression { "," expression } [ "," ] .

func ArrayMembers(tokens []tok.Token) (members []ast.Expression, remainder []tok.Token, err error) {
	remainder = skipTrivia(tokens)

	for {
		var expression ast.Expression
		if expression, remainder, err = Expression(remainder); err != nil {
			break
		}
		members = append(members, expression)

		var found bool
		if remainder, found = Comma(remainder); !found {
			break
		}
	}

	remainder, _ = Comma(remainder)

	return members, remainder, nil
}

func arrayLiteralWithType(arrayType ast.ArrayLiteralType, tokens []tok.Token) (arr *ast.ArrayLiteral, remainder []tok.Token, err error) {
	if functionBlock, remainder2, blockErr := FunctionBlock(tokens); blockErr == nil {
		return ast.NewArrayLiteral(arrayType, nil, functionBlock), remainder2, nil
	} else if blockErr != ErrNoMatch {
		return nil, remainder2, blockErr
	}

	var found bool
	if remainder, found = OpenBracket(tokens); !found {
		return nil, tokens, ErrNoMatch
	}

	var arrayMembers []ast.Expression
	if arrayMembers, remainder, err = ArrayMembers(remainder); err != nil {
		return nil, remainder, err
	}

	remainder = skipTrivia(remainder)
	if remainder, found = CloseBracket(remainder); !found {
		return nil, remainder, errorExpectingTokenType(tok.TokCloseBracket, remainder)
	}

	return ast.NewArrayLiteral(arrayType, arrayMembers, nil), remainder, nil
}

// fixed_size_array but parsed conservatively so plain array literals like [1, 2, 3]
// do not get consumed as malformed fixed-size array prefixes.

func fixedSizeArrayLiteralType(tokens []tok.Token) (*ast.FixedSizeArrayType, []tok.Token, error) {
	remainder, found := OpenBracket(tokens)
	if !found {
		return nil, tokens, ErrNoMatch
	}

	size, remainder, err := Size(remainder)
	if err != nil {
		return nil, tokens, ErrNoMatch
	}

	if remainder, found = CloseBracket(remainder); !found {
		return nil, tokens, ErrNoMatch
	}

	elementType, remainder, err := ArrayElementType(remainder)
	if err != nil {
		return nil, tokens, ErrNoMatch
	}

	return ast.NewFixedSizeArrayType(elementType, size), remainder, nil
}
