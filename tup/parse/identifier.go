package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/tok"
)

// identifier = ( lowercase_letter | "_" ) { letter | decimal_digit | "_" } .

func Identifier(tokens []tok.Token) (item *ast.Identifier, remainder []tok.Token, err error) {
	remainder = skipComments(tokens)
	if peek(remainder).Type != tok.TokID {
		return nil, nil, errorExpecting(tok.TokenTypes[tok.TokID], remainder)
	}
	return ast.NewIdentifier(remainder[0].Value(), remainder[0].File, remainder[0].Offset, remainder[0].Length), remainder[1:], nil
}

// type_identifier = uppercase_letter { letter | decimal_digit | "_" } .

func TypeIdentifier(tokens []tok.Token) (item *ast.TypeIdentifier, remainder []tok.Token, err error) {
	remainder = skipComments(tokens)
	if peek(remainder).Type != tok.TokTypeID {
		return nil, nil, errorExpecting(tok.TokenTypes[tok.TokTypeID], remainder)
	}
	return ast.NewTypeIdentifier(remainder[0].Value(), remainder[0].File, remainder[0].Offset, remainder[0].Length), remainder[1:], nil
}

// rename_identifier = identifier [ ":" identifier ] .

func RenameIdentifier(tokens []tok.Token) (item *ast.RenameIdentifier, remainder []tok.Token, err error) {
	identifier, remainder, err := Identifier(tokens)
	if identifier == nil || err != nil {
		return nil, nil, err
	}
	var original *ast.Identifier
	if peek(remainder).Type == tok.TokColon {
		remainder = remainder[1:]
		original, remainder, err = Identifier(remainder)
		if original == nil || err != nil {
			return nil, nil, err
		}
	}
	return ast.NewRenameIdentifier(identifier, original), remainder, nil
}

// rename_type = type_identifier [ ":" type_identifier ] .

func RenameType(tokens []tok.Token) (item *ast.RenameType, remainder []tok.Token, err error) {
	typeIdentifier, remainder, err := TypeIdentifier(tokens)
	if typeIdentifier == nil || err != nil {
		return nil, nil, err
	}
	var original *ast.TypeIdentifier
	if peek(remainder).Type == tok.TokColon {
		remainder = remainder[1:]
		original, remainder, err = TypeIdentifier(remainder)
		if original == nil || err != nil {
			return nil, nil, err
		}
	}
	return ast.NewRenameType(typeIdentifier, original), remainder, nil
}
