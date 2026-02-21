package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// integer_literal = binary_literal
//                 | hexadecimal_literal
//                 | octal_literal
//                 | decimal_literal .

func IntegerLiteral(tokens []tok.Token) (item *ast.IntegerLiteral, remainder []tok.Token, err error) {
	binaryLit, remainder, err := BinaryLiteral(tokens)
	if err == nil {
		return binaryLit, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	hexadecimalLit, remainder, err := HexadecimalLiteral(tokens)
	if err == nil {
		return hexadecimalLit, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	octalLit, remainder, err := OctalLiteral(tokens)
	if err == nil {
		return octalLit, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	decimalLit, remainder, err := DecimalLiteral(tokens)
	if err == nil {
		return decimalLit, remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	return nil, tokens, ErrNoMatch
}
