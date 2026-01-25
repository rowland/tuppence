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

	return nil, nil, errorExpectingOneOf("number", tokens, errors)

}

// integer_literal = binary_literal
//                 | hexadecimal_literal
//                 | octal_literal
//                 | decimal_literal .

func IntegerLiteral(tokens []tok.Token) (item ast.IntegerLiteral, remainder []tok.Token, err error) {
	var errors []error

	binaryLit, remainder, err := BinaryLiteral(tokens)
	if err == nil {
		return binaryLit, remainder, nil
	}
	errors = append(errors, err)

	hexadecimalLit, remainder, err := HexadecimalLiteral(tokens)
	if err == nil {
		return hexadecimalLit, remainder, nil
	}
	errors = append(errors, err)

	octalLit, remainder, err := OctalLiteral(tokens)
	if err == nil {
		return octalLit, remainder, nil
	}
	errors = append(errors, err)

	decimalLit, remainder, err := DecimalLiteral(tokens)
	if err == nil {
		return decimalLit, remainder, nil
	}
	errors = append(errors, err)

	return nil, nil, errorExpectingOneOf("IntegerLiteral", tokens, errors)
}
