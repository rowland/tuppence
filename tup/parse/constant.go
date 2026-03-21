package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// scoped_identifier = identifier { "." identifier } .

func ScopedIdentifier(tokens []tok.Token) (*ast.ScopedIdentifier, []tok.Token, error) {
	first, remainder, err := Identifier(tokens)
	if err != nil {
		return nil, remainder, err
	}

	identifiers := []*ast.Identifier{first}
	for {
		remainder2, found := Dot(remainder)
		if !found {
			break
		}

		next, remainder3, err := Identifier(remainder2)
		if err == ErrNoMatch {
			return nil, remainder2, errorExpecting("identifier", remainder2)
		} else if err != nil {
			return nil, remainder3, err
		}

		identifiers = append(identifiers, next)
		remainder = remainder3
	}

	return ast.NewScopedIdentifier(identifiers), remainder, nil
}

// constant = literal
//          | scoped_identifier .

func Constant(tokens []tok.Token) (*ast.Constant, []tok.Token, error) {
	if literal, remainder, err := Literal(tokens); err == nil {
		switch value := literal.(type) {
		case *ast.FloatLiteral:
			return ast.NewConstant(value), remainder, nil
		case *ast.IntegerLiteral:
			return ast.NewConstant(value), remainder, nil
		case *ast.BooleanLiteral:
			return ast.NewConstant(value), remainder, nil
		case *ast.StringLiteral:
			return ast.NewConstant(value), remainder, nil
		case *ast.InterpolatedStringLiteral:
			return ast.NewConstant(value), remainder, nil
		case *ast.RawStringLiteral:
			return ast.NewConstant(value), remainder, nil
		case *ast.MultiLineStringLiteral:
			return ast.NewConstant(value), remainder, nil
		case *ast.TupleLiteral:
			return ast.NewConstant(value), remainder, nil
		case *ast.ArrayLiteral:
			return ast.NewConstant(value), remainder, nil
		case *ast.SymbolLiteral:
			return ast.NewConstant(value), remainder, nil
		case *ast.RuneLiteral:
			return ast.NewConstant(value), remainder, nil
		default:
			return nil, remainder, errorExpecting("constant value", remainder)
		}
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	if scopedIdentifier, remainder, err := ScopedIdentifier(tokens); err == nil {
		return ast.NewConstant(scopedIdentifier), remainder, nil
	} else if err != ErrNoMatch {
		return nil, remainder, err
	}

	return nil, tokens, ErrNoMatch
}
