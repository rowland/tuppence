package parse

import (
	"github.com/rowland/tuppence/tup/ast"
	"github.com/rowland/tuppence/tup/source"
	"github.com/rowland/tuppence/tup/tok"
)

// identifier = ( lowercase_letter | "_" ) { letter | decimal_digit | "_" } .

func Identifier(tokens []tok.Token) (item *ast.Identifier, remainder []tok.Token, err error) {
	remainder = skipComments(tokens)
	if peek(remainder).Type != tok.TokID {
		return nil, tokens, ErrNoMatch
	}
	return ast.NewIdentifier(remainder[0].Value(), remainder[0].File, remainder[0].Offset, remainder[0].Length), remainder[1:], nil
}

// type_identifier = uppercase_letter { letter | decimal_digit | "_" } .

func TypeIdentifier(tokens []tok.Token) (item *ast.TypeIdentifier, remainder []tok.Token, err error) {
	remainder = skipComments(tokens)

	if peek(remainder).Type != tok.TokTypeID {
		return nil, tokens, ErrNoMatch
	} else if peek(remainder).Invalid {
		return nil, remainder, errorExpecting(tok.TokenTypes[tok.TokTypeID], remainder)
	}

	return ast.NewTypeIdentifier(remainder[0].Value(), remainder[0].File, remainder[0].Offset, remainder[0].Length), remainder[1:], nil
}

// rename_identifier = identifier [ ":" identifier ] .

func RenameIdentifier(tokens []tok.Token) (item *ast.RenameIdentifier, remainder []tok.Token, err error) {
	identifier, remainder, err := Identifier(tokens)
	if err == ErrNoMatch {
		return nil, tokens, ErrNoMatch
	} else if err != nil {
		return nil, remainder, err
	}

	remainder2, err2 := Colon(remainder)
	if err2 == ErrNoMatch {
		// no colon found, so no renaming
		return ast.NewRenameIdentifier(identifier, nil), remainder2, nil
	} else if err2 != nil {
		return nil, remainder2, err2
	}

	// colon found, so original identifier expected
	original, remainder3, err3 := Identifier(remainder2)
	if err3 != nil {
		return nil, remainder3, err3
	}

	return ast.NewRenameIdentifier(identifier, original), remainder3, nil
}

// rename_type = type_identifier [ ":" type_identifier ] .

func RenameType(tokens []tok.Token) (item *ast.RenameType, remainder []tok.Token, err error) {
	typeIdentifier, remainder, err := TypeIdentifier(tokens)
	if err == ErrNoMatch {
		return nil, tokens, ErrNoMatch
	} else if err != nil {
		return nil, remainder, err
	}

	remainder2, err2 := Colon(remainder)
	if err2 == ErrNoMatch {
		// no colon found, so no renaming
		return ast.NewRenameType(typeIdentifier, nil), remainder, nil
	} else if err2 != nil {
		return nil, remainder2, err2
	}

	// colon found, so original type identifier expected
	original, remainder3, err3 := TypeIdentifier(remainder2)
	if err3 != nil {
		return nil, remainder3, err3
	}

	return ast.NewRenameType(typeIdentifier, original), remainder3, nil
}

// type_reference = [ identifier { "." identifier } "." ] type_identifier .

func TypeReference(tokens []tok.Token) (item *ast.TypeReference, remainder []tok.Token, err error) {
	remainder = skipComments(tokens)

	var identifiers []*ast.Identifier
	for {
		var identifier *ast.Identifier
		var remainder2 []tok.Token
		identifier, remainder2, err = Identifier(remainder)

		if err == ErrNoMatch {
			break
		} else if err != nil {
			return nil, remainder2, err
		}

		remainder = remainder2
		identifiers = append(identifiers, identifier)
		remainder, err = Dot(remainder)
		if err != nil {
			return nil, nil, errorExpecting(tok.TokenTypes[tok.TokDot], remainder)
		}
	}

	typeIdentifier, remainder, err := TypeIdentifier(remainder)
	if err == ErrNoMatch {
		return nil, tokens, ErrNoMatch
	} else if err != nil {
		return nil, nil, err
	}

	var src *source.Source
	var startOffset, length int32
	if len(identifiers) > 0 {
		src = identifiers[0].Source
		startOffset = identifiers[0].StartOffset
		length = typeIdentifier.StartOffset + typeIdentifier.Length - identifiers[0].StartOffset
	} else {
		src = typeIdentifier.Source
		startOffset = typeIdentifier.StartOffset
		length = typeIdentifier.Length
	}
	return ast.NewTypeReference(identifiers, typeIdentifier, src, startOffset, length), remainder, nil
}

// function_identifier = lowercase_letter { letter | decimal_digit | "_" } [ "?" | "!" ] .

func FunctionIdentifier(tokens []tok.Token) (item *ast.FunctionIdentifier, remainder []tok.Token, err error) {
	remainder = skipComments(tokens)

	if peek(remainder).Type != tok.TokFuncID && peek(remainder).Type != tok.TokID {
		return nil, tokens, ErrNoMatch
	} else if peek(remainder).Invalid {
		return nil, remainder, errorExpecting(tok.TokenTypes[tok.TokFuncID], remainder)
	}

	return ast.NewFunctionIdentifier(
			remainder[0].Value(),
			remainder[0].File,
			remainder[0].Offset,
			remainder[0].Length),
		remainder[1:],
		nil
}
