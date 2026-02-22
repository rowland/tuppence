package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// literal = number
//         | boolean_literal
//         | string_literal
//         | interpolated_string_literal
//         | raw_string_literal
//         | multi_line_string_literal
//         | tuple_literal
//         | array_literal
//         | symbol_literal
//         | rune_literal
//         | fixed_size_array_literal .

func Literal(tokens []tok.Token) (item ast.Literal, remainder []tok.Token, err error) {
	remainder = skipComments(tokens)

	if number, remainder, err := Number(remainder); err == nil {
		return number, remainder, nil
	} else if err != ErrNoMatch {
		return nil, tokens, err
	}

	if booleanLiteral, remainder, err := BooleanLiteral(remainder); err == nil {
		return booleanLiteral, remainder, nil
	} else if err != ErrNoMatch {
		return nil, tokens, err
	}

	if stringLiteral, remainder, err := StringLiteral(remainder); err == nil {
		return stringLiteral, remainder, nil
	} else if err != ErrNoMatch {
		return nil, tokens, err
	}

	// interpolatedStringLiteral, remainder, err := InterpolatedStringLiteral(remainder)
	// if err == nil {
	// 	return interpolatedStringLiteral, remainder, nil
	// } else if err != ErrNoMatch {
	// 	return nil, tokens, err
	// }

	if rawStringLiteral, remainder, err := RawStringLiteral(remainder); err == nil {
		return rawStringLiteral, remainder, nil
	} else if err != ErrNoMatch {
		return nil, tokens, err
	}

	// multiLineStringLiteral, remainder, err := MultiLineStringLiteral(remainder)
	// if err == nil {
	// 	return multiLineStringLiteral, remainder, nil
	// } else if err != ErrNoMatch {
	// 	return nil, tokens, err
	// }

	// tupleLiteral, remainder, err := TupleLiteral(remainder)
	// if err == nil {
	// 	return tupleLiteral, remainder, nil
	// } else if err != ErrNoMatch {
	// 	return nil, tokens, err
	// }

	// arrayLiteral, remainder, err := ArrayLiteral(remainder)
	// if err == nil {
	// 	return arrayLiteral, remainder, nil
	// } else if err != ErrNoMatch {
	// 	return nil, tokens, err
	// }

	// symbolLiteral, remainder, err := SymbolLiteral(remainder)
	// if err == nil {
	// 	return symbolLiteral, remainder, nil
	// } else if err != ErrNoMatch {
	// 	return nil, tokens, err
	// }

	if runeLiteral, remainder, err := RuneLiteral(remainder); err == nil {
		return runeLiteral, remainder, nil
	} else if err != ErrNoMatch {
		return nil, tokens, err
	}

	// fixedSizeArrayLiteral, remainder, err := FixedSizeArrayLiteral(remainder)
	// if err == nil {
	// 	return fixedSizeArrayLiteral, remainder, nil
	// } else if err != ErrNoMatch {
	// 	return nil, tokens, err
	// }

	return nil, tokens, errorExpecting("literal", tokens)
}
