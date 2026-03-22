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
//         | rune_literal .

func Literal(tokens []tok.Token) (item ast.Literal, remainder []tok.Token, err error) {
	// fmt.Println("Literal", tok.Types(tokens))
	remainder = skipTrivia(tokens)

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

	if interpolatedStringLiteral, remainder, err := InterpolatedStringLiteral(remainder); err == nil {
		return interpolatedStringLiteral, remainder, nil
	} else if err != ErrNoMatch {
		return nil, tokens, err
	}

	if rawStringLiteral, remainder, err := RawStringLiteral(remainder); err == nil {
		return rawStringLiteral, remainder, nil
	} else if err != ErrNoMatch {
		return nil, tokens, err
	}

	if multiLineStringLiteral, remainder, err := MultiLineStringLiteral(remainder); err == nil {
		return multiLineStringLiteral, remainder, nil
	} else if err != ErrNoMatch {
		return nil, tokens, err
	}

	tupleLiteral, remainder, err := TupleLiteral(remainder)
	if err == nil {
		return tupleLiteral, remainder, nil
	} else if err != ErrNoMatch {
		return nil, tokens, err
	}

	arrayLiteral, remainder, err := ArrayLiteral(remainder)
	if err == nil {
		return arrayLiteral, remainder, nil
	} else if err != ErrNoMatch {
		return nil, tokens, err
	}

	symbolLiteral, remainder, err := SymbolLiteral(remainder)
	if err == nil {
		return symbolLiteral, remainder, nil
	} else if err != ErrNoMatch {
		return nil, tokens, err
	}

	if runeLiteral, remainder, err := RuneLiteral(remainder); err == nil {
		return runeLiteral, remainder, nil
	} else if err != ErrNoMatch {
		return nil, tokens, err
	}

	return nil, tokens, ErrNoMatch
}

// tuple_literal = empty_tuple | labeled_tuple_members | tuple_members .

func TupleLiteral(tokens []tok.Token) (tupleLiteral *ast.TupleLiteral, remainder []tok.Token, err error) {
	// fmt.Println("TupleLiteral", tok.Types(tokens))

	var members []*ast.TupleMember

	if members, remainder, err = emptyTuple(tokens); err == nil {
		return ast.NewTupleLiteral(false, members), remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if members, remainder, err = labeledTupleMembers(tokens); err == nil {
		return ast.NewTupleLiteral(true, members), remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if members, remainder, err = tupleMembers(tokens); err == nil {
		return ast.NewTupleLiteral(false, members), remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	return nil, tokens, ErrNoMatch
}

// empty_tuple = "(" ")" .

func emptyTuple(tokens []tok.Token) (tupleMembers []*ast.TupleMember, remainder []tok.Token, err error) {
	// fmt.Println("emptyTuple", tok.Types(tokens))

	var found bool
	if remainder, found = OpenParen(tokens); !found {
		return nil, tokens, ErrNoMatch
	}

	if remainder, found = CloseParen(remainder); found {
		return []*ast.TupleMember{}, remainder, nil
	}

	return nil, remainder, ErrNoMatch
}

// labeled_tuple_members = "(" labeled_tuple_member { "," labeled_tuple_member } [ "," ] ")" .

func labeledTupleMembers(tokens []tok.Token) (tupleMembers []*ast.TupleMember, remainder []tok.Token, err error) {
	// fmt.Println("labeledTupleMembers", tok.Types(tokens))

	var found bool
	if remainder, found = OpenParen(tokens); !found {
		return nil, tokens, ErrNoMatch
	}

	for {
		var member *ast.TupleMember
		if member, remainder, err = labeledTupleMember(remainder); err == nil {
			tupleMembers = append(tupleMembers, member)
		} else if err != ErrNoMatch {
			return nil, remainder, err
		}

		if remainder, found = Comma(remainder); !found {
			break
		}
	}

	remainder = skipTrivia(remainder)
	if remainder, found = CloseParen(remainder); !found && len(tupleMembers) > 0 {
		return nil, remainder, errorExpectingTokenType(tok.TokCloseParen, remainder)
	}

	if len(tupleMembers) > 0 {
		return tupleMembers, remainder, nil
	}

	return nil, remainder, ErrNoMatch
}

// labeled_tuple_member = identifier ":" tuple_member .

func labeledTupleMember(tokens []tok.Token) (labeledTupleMember *ast.TupleMember, remainder []tok.Token, err error) {
	// fmt.Println("labeledTupleMember", tok.Types(tokens))

	var identifier *ast.Identifier
	if identifier, remainder, err = Identifier(tokens); err != nil {
		return nil, remainder, err
	}

	var found bool
	if remainder, found = Colon(remainder); !found {
		return nil, tokens, ErrNoMatch
	}

	var expression ast.Expression
	if expression, remainder, err = Expression(remainder); err == nil {
		return ast.NewTupleMember(identifier, expression), remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	return nil, tokens, ErrNoMatch
}

// tuple_members = "(" tuple_member "," { tuple_member "," } [ tuple_member ] ")" .

func tupleMembers(tokens []tok.Token) (tupleMembers []*ast.TupleMember, remainder []tok.Token, err error) {
	// fmt.Println("tupleMembers", tok.Types(tokens))

	var found bool
	if remainder, found = OpenParen(tokens); !found {
		return nil, tokens, ErrNoMatch
	}

	for {
		var member *ast.TupleMember
		if member, remainder, err = tupleMember(remainder); err == nil {
			tupleMembers = append(tupleMembers, member)
		}

		if remainder, found = Comma(remainder); !found {
			if len(tupleMembers) <= 1 {
				return nil, tokens, ErrNoMatch
			}
			break
		}
	}

	remainder = skipTrivia(remainder)
	if remainder, found = CloseParen(remainder); !found && len(tupleMembers) > 0 {
		return nil, remainder, errorExpectingTokenType(tok.TokCloseParen, remainder)
	}

	if len(tupleMembers) > 0 {
		return tupleMembers, remainder, nil
	}

	return nil, remainder, ErrNoMatch
}

// tuple_member = expression .

func tupleMember(tokens []tok.Token) (tupleMember *ast.TupleMember, remainder []tok.Token, err error) {
	// fmt.Println("tupleMember", tok.Types(tokens))

	var expression ast.Expression
	if expression, remainder, err = Expression(tokens); err == nil {
		return ast.NewTupleMember(nil, expression), remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	return nil, remainder, ErrNoMatch
}

// symbol_literal = ":" identifier .

func SymbolLiteral(tokens []tok.Token) (symbolLiteral *ast.SymbolLiteral, remainder []tok.Token, err error) {
	remainder = skipTrivia(tokens)

	if peek(remainder).Type != tok.TokColonNoSpace {
		return nil, tokens, ErrNoMatch
	}

	var identifier *ast.Identifier
	if identifier, remainder, err = Identifier(remainder[1:]); err != nil {
		if err == ErrNoMatch {
			return nil, tokens, ErrNoMatch
		}
		return nil, remainder, err
	}

	return ast.NewSymbolLiteral(
		":"+identifier.String(),
		identifier.Source,
		identifier.StartOffset-1,
		identifier.Length+1,
	), remainder, nil
}
