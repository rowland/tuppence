package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// number = float_literal | integer_literal .

func Number(tokens []tok.Token) (number ast.Number, remainder []tok.Token, err error) {
	// fmt.Println("Number", tokens)

	var floatLit *ast.FloatLiteral
	if floatLit, remainder, err = FloatLiteral(tokens); err == nil {
		return floatLit, remainder, nil
	}

	var integerLit *ast.IntegerLiteral
	if integerLit, remainder, err = IntegerLiteral(tokens); err == nil {
		return integerLit, remainder, nil
	}

	return nil, nil, ErrNoMatch
}
