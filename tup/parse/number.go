package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// number = float_literal | integer_literal .

func Number(tokens []tok.Token) (item ast.Number, remainder []tok.Token, err error) {
	var errors []error
	floatLit, remainder, err := FloatLiteral(tokens)
	if err == nil {
		return floatLit, remainder, nil
	}
	errors = append(errors, err)

	integerLit, remainder, err := IntegerLiteral(tokens)
	if err == nil {
		return integerLit, remainder, nil
	}
	errors = append(errors, err)

	return nil, nil, ErrNoMatch
}
